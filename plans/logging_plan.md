# Logging Plan — Century: Golem Edition Server

## Goals

- Provide operational visibility into the running gRPC server (startup, requests, errors)
- Capture enough game-action context to debug game state issues and support analysis
- Avoid logging sensitive or redundant information that adds noise
- logging to stdout is the desired outcome; we can wire this up with an observability solution down the line

---

## Logging Library

Use Go's `log/slog` (stdlib, Go 1.21+). It supports structured key-value logging with
JSON or text output, zero external dependencies, and integrates cleanly with gRPC interceptors.

No third-party library is needed at this stage. If log aggregation (e.g. Loki, Datadog)
is added later, slog's JSON handler makes that straightforward.

---

## Log Levels

| Level | Purpose |
|-------|---------|
| `DEBUG` | Verbose game-state details (hand contents, gem counts) — off by default |
| `INFO`  | Normal operational events (server start, game created, player actions) |
| `WARN`  | Recoverable issues (game not found, invalid argument from client) |
| `ERROR` | Unexpected failures (repository errors, internal panics) |

---

## What to Log — By Layer

### 1. Server Startup (`cmd/server/server.go`)

| Event | Level | Fields |
|-------|-------|--------|
| Listener bound successfully | INFO | `port=50051` |
| Failed to bind listener | ERROR (fatal) | `err` |
| gRPC server ready to serve | INFO | — |
| Server stopped unexpectedly | ERROR (fatal) | `err` |

### 2. gRPC Handler (`adapters/handlers/grpc_handler.go`)

One log entry per RPC call — at the boundary where the server receives and
responds to a request. This is the most important layer for operational visibility.

#### `NewGame`
| Event | Level | Fields |
|-------|-------|--------|
| Request received | INFO | `rpc=NewGame`, `player_count` |
| Game created successfully | INFO | `rpc=NewGame`, `game_id` (int32 external), `uuid` (internal) |
| Service returned error | ERROR | `rpc=NewGame`, `err` |

#### `PlayGemCard`
| Event | Level | Fields |
|-------|-------|--------|
| Request received | INFO | `rpc=PlayGemCard`, `game_id`, `card_in`, `card_out`, `card_upgrades` |
| Game ID not found in id map | WARN | `rpc=PlayGemCard`, `game_id` |
| Missing card in request | WARN | `rpc=PlayGemCard`, `game_id` |
| Card played successfully | INFO | `rpc=PlayGemCard`, `game_id`, `player` |
| Service returned error | ERROR | `rpc=PlayGemCard`, `game_id`, `err` |

#### `GetCard` (currently unimplemented)
| Event | Level | Fields |
|-------|-------|--------|
| Called while unimplemented | WARN | `rpc=GetCard` |

#### Future RPCs (`GetGemCard`, `GetGolemCard`, `Rest`)
Follow the same pattern: log request received (INFO), log success (INFO), log
validation failures (WARN), log service/repo errors (ERROR).

### 3. Game Service (`internal/core/services/game_service/service.go`)

Log at the service layer only for errors and notable business events — not for
every call (the handler layer already covers that).

| Event | Level | Fields |
|-------|-------|--------|
| `Create`: new game initialized | INFO | `uuid`, `player_count` |
| `Create`: repository save failed | ERROR | `uuid`, `err` |
| `Get`: game not found in repository | WARN | `uuid`, `err` |
| `PlayGemCard`: domain rule violation | WARN | `uuid`, `player`, `err` |
| `PlayGemCard`: repository save failed | ERROR | `uuid`, `player`, `err` |
| `GetGemCard` / `GetGolemCard` / `Rest`: stub called (not yet implemented) | WARN | `uuid`, `player`, `rpc` |

### 4. Repository (`adapters/repository/in_memory.go`)

The in-memory repository is simple enough that its errors are already propagated
upward. Add DEBUG-level logs only:

| Event | Level | Fields |
|-------|-------|--------|
| `Get`: cache hit | DEBUG | `uuid` |
| `Get`: cache miss | DEBUG | `uuid` |
| `Save`: game stored | DEBUG | `uuid`, `total_games` |

---

## What NOT to Log

- Full game state (hand, discard pile, lineup) at INFO — that's DEBUG territory and
  very verbose. Log game state only on explicit debug flag or when reproducing a bug.
- Player gem counts on every action — too noisy; relevant only for DEBUG.
- The internal UUID on every handler log at INFO — use the external `game_id` (int32)
  at the handler boundary; the UUID is for service/repo logs.

---

## Implementation Approach

### Structured Logger Setup

Initialize a single `*slog.Logger` in `main()` and pass it down via dependency injection
(constructor arguments), not a global variable. This keeps things testable.

```
// cmd/server/server.go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo, // override with env var or flag
}))
```

Pass `logger` to:
- `handlers.New(svc, logger)`
- `game_service.New(repo, logger)`
- (optionally) `repository.NewInMemoryRepository(logger)` for DEBUG repo logs

### Log Level from Environment

Read a `LOG_LEVEL` environment variable at startup and set the handler's `Level`
accordingly. This avoids recompiling to change verbosity.

### gRPC Interceptor

A unary server interceptor in `cmd/server/server.go` will automatically log every
RPC's method name, duration, and final gRPC status code. This gives a single
consistent entry point for all request logging rather than duplicating "request
received" boilerplate in every handler method.

```
grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(logger)))
```

The interceptor will be added from the start. Distributed tracing (e.g. OpenTelemetry)
will be layered on top of this interceptor later — the interceptor is the natural
place to inject and propagate trace spans.

---

## Player Action Audit Log

Player actions require a distinct logging path from operational logs. The purpose is
to produce a clean, structured record of every game decision — suitable for later
use in training models to play the game.

### What it is

A separate `*slog.Logger` (or `slog.Handler`) writing to its own output stream
(e.g. a dedicated file or stdout with a parseable prefix), so the audit records
can be extracted independently from operational noise.

### What each record captures

Every entry represents one player action and should include enough context to
reconstruct the game state before and after the action:

| Field | Description |
|-------|-------------|
| `event` | Action type: `play_gem_card`, `get_gem_card`, `buy_golem`, `rest` |
| `game_id` | External game identifier |
| `player` | Player index (0-based) |
| `turn` | Turn number within the game |
| `card` | The card acted upon (inputs, outputs, upgrades) |
| `gems_before` | Player's gem counts before the action |
| `gems_after` | Player's gem counts after the action |
| `hand_size` | Number of cards in hand after action |
| `discard_size` | Number of cards in discard pile after action |
| `golem_lineup` | Current golem lineup (for buy decisions) |
| `gem_lineup` | Current gem card lineup (for get-card decisions) |

Game-end state (winner, scores, total turns) should also be emitted as a final
`game_over` event per game.

### What it is NOT

The audit log is not a debug log. It should not include error traces, warnings,
or server internals — only successful, committed game actions.

### Implementation note

Emit audit events from the service layer (`game_service`), after the domain action
has succeeded and the updated game has been saved to the repository. The handler
layer should not write audit events — it doesn't have the before/after game state.

---

## Log Format Example (JSON)

```json
{"time":"2026-02-26T12:00:00Z","level":"INFO","msg":"game created","game_id":1,"uuid":"f47ac10b-...","player_count":3}
{"time":"2026-02-26T12:00:01Z","level":"INFO","msg":"gem card played","game_id":1,"player":0,"card_in":{"yellow":2},"card_out":{"green":1}}
{"time":"2026-02-26T12:00:02Z","level":"WARN","msg":"game not found","game_id":99}
{"time":"2026-02-26T12:00:03Z","level":"ERROR","msg":"repository save failed","uuid":"f47ac10b-...","err":"..."}
```

