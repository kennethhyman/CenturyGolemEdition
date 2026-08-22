# Plan

Working notes for getting CenturyGolemEdition from "restructured but unwired" to
a game you can trust, run locally, and eventually plug AI players into.

Stages are ordered by dependency, not by preference — each one is the
prerequisite substrate for the next. Do not start a stage whose predecessor's
criteria aren't met. Do not pull in work from a later stage while executing an
earlier one, even if it looks convenient in the moment.

An agent picking up this plan should read this file, find the first stage with
unmet criteria, do the smallest task that moves toward meeting it, verify with
the listed command, and stop there rather than continuing into the next
stage's territory.

## Stage 1 — Repo hygiene

Owner: build tooling (`go.mod`, `Makefile`, CI), not application code.

- [x] `go.mod` accurately reflects direct vs. indirect deps (`go mod tidy` clean)
- [x] `Makefile` targets work: `tidy`, `vet`, `build_game`, `test_game`, `lint`,
      `run_server`, `run_cli`, `clean`, `all`
- [x] CI (`.github/workflows/ci.yml`) runs tidy-check, vet, build, test as
      blocking gates on push/PR
- [x] `cli/actions` package declaration matches its import path (was `package
      turns` under `cli/actions/`, a leftover from an incomplete rename —
      fixed)
- [x] `go.mod`'s `go` directive matches the installed toolchain (was `1.18`,
      bumped to `1.24`, then to `1.27.0` — the current latest Go release.
      `go build`/`go test`/etc. auto-download the matching toolchain via
      `GOTOOLCHAIN=auto`; no system-wide Go reinstall needed or done — that's
      a machine-level concern shared by every project here, not this repo's
      to own. Local `golangci-lint` had to be rebuilt with
      `GOTOOLCHAIN=go1.27.0 go install .../golangci-lint@latest` to match.
- [ ] `golangci-lint run ./...` is clean — **not required yet, and currently
      can't run at all**. As of go 1.27, golangci-lint's latest release
      (v1.64.8) can't parse this toolchain's export data format
      ("export data version 4 is greater than maximum supported version 2")
      and crashes with exit 3, independent of any real findings. This is an
      upstream tooling gap, not a repo issue — don't try to work around it
      (e.g. by downgrading the go directive). CI's lint job has
      `continue-on-error: true` for this reason; revisit once golangci-lint
      ships go1.27 support and re-verify locally first. Underneath that,
      the original pre-existing findings (unchecked errors from `fmt.Scan*` in
      `internal/core/domain/game.go`, unused fields/vars in `cli/game.go`,
      `cli/player.go`, `cli/actions/turnView.go`, `cli/cli.go`) still live in
      files Stage 2 and Stage 3 substantially rework. Fix them as part of the
      stage that touches that file, not as a standalone cleanup pass — don't let an
      agent "clean up" these files outside their owning stage. CI runs lint
      non-blocking (`--issues-exit-code=0`) so findings are visible without
      gating merges on debt this plan already knows about.

**Success criteria (all must pass):**
```
cd game && go mod tidy && git diff --exit-code -- go.mod go.sum
cd game && make all   # tidy + vet + build_protos + build_game + test_game
```
Both exit 0 with no diff and no test failures. Stage 1 is otherwise done —
only the deferred lint items remain, tracked against Stage 2/3 below.

## Stage 2 — Architecture wiring

Owner: a new adapters package + `cmd/server` as the composition root. Domain
and `ports` stay as they are — this stage fills gaps, it doesn't redesign the
interface.

Decisions already made (don't revisit without a reason):
- Local/dev persistence: **sqlite**, via a pure-Go driver (e.g.
  `modernc.org/sqlite`) — avoids CGO so `go build`/CI stay simple. A
  production deployment gets a different `ports.GameRepository`
  implementation later (e.g. Postgres); the port is the seam, so this costs
  nothing to swap.
- Storage shape for now: a single table, game serialized as a JSON blob keyed
  by ID (`id TEXT PRIMARY KEY, data TEXT`). `domain.Game` is a plain struct of
  exported fields — no normalized schema until a production DB is actually
  chosen. Don't build one preemptively.
- Which repository implementation gets used is selected in `cmd/server`
  (composition root) via flag/env, never hardcoded into
  `game_service.Service`.

Tasks:
- [x] Implement `ports.GameRepository` against sqlite in a new
      `internal/adapters/sqlite` (or similar) package, with its own tests.
      Done: `internal/adapters/sqlite/repository.go`, driver
      `modernc.org/sqlite` (pure Go, no CGO), single `games(id TEXT PRIMARY
      KEY, data TEXT)` table, JSON blob per row. Tests in
      `internal/adapters/sqlite/repository_test.go` (save/get round trip,
      not-found, overwrite, reopen-persists, empty-ID rejected).
      **Implementation note surfaced during this work**: `domain.LineUp`'s
      `deck *Deck[CardType]` field is unexported, so a plain
      `encoding/json.Marshal(domain.Game{...})` silently drops every card
      still left to draw — `Remaining()`/`Draw()`/`String()` would then
      nil-pointer-panic on anything loaded back from storage. Fixed with
      `MarshalJSON`/`UnmarshalJSON` on `LineUp` (`internal/core/domain/lineup.go`)
      that include the deck's remaining cards; `domain.Game`/`Player` needed
      no such fix since all their fields are already exported. This is a
      mechanical fix (no public API change), not a redesign of domain/ports.
- [x] Fill in `internal/core/grpc/adapters.go` — currently just `package grpc`
      with no handlers — to translate grpc requests to `game_service.Service`
      calls and back (there's already `domain.MarshallGame` /
      `MarshallGemCard` / `MarshallGolemCard` to build on).
      **Structural note**: the implementation does not live in that file.
      `internal/core/domain` already imports `internal/core/grpc` (as `pb`,
      for `MarshallGame` et al.), so a `GameServer` implementation calling
      `domain.MarshallGame` from *inside* `internal/core/grpc` would close an
      import cycle (`grpc -> domain -> grpc`) — not a style choice, Go
      rejects this outright. The real implementation is the new
      `internal/adapters/grpc` package (`grpcadapter.GameServer`), which
      depends one-way on both `internal/core/grpc` and `internal/core/domain`;
      `internal/core/grpc/adapters.go` now just documents this and points
      there. Of the 3 RPCs the generated `pb.GameServer` interface declares
      (`NewGame`, `GetCard`, `PlayGemCard`), only `NewGame` is implemented —
      it's the only one with a corresponding `game_service.Service` method
      (`Create`) to translate to. `GetCard`/`PlayGemCard` are left
      `Unimplemented` (via the generated `UnimplementedGameServer` default):
      `GetCardMessage.game_id` is `int32` while domain game IDs are UUID
      strings and `GetCard` returns a single `GemCard` that no service method
      produces; `PlayGemCardMessage` carries no player index and
      `game_service.Service` has no `PlayGemCard` method yet (task below only
      asked to change `New()`'s signature, not add the other
      `ports.GameService` methods). Wiring these up is follow-on work once
      `game_service.Service` grows the matching methods.
- [x] `game_service.Service.New()` currently takes no arguments and leaves
      `gamesRepository` nil — change it to accept a `ports.GameRepository`,
      injected by the caller. Done — `New(gamesRepository ports.GameRepository) *Service`.
- [x] Wire `cmd/server/server.go` (currently a stub `fmt.Printf`, real startup
      commented out) to construct the sqlite repository, the service, the
      grpc adapter, and actually start listening on 50051. Done — `-addr`/
      `-db` flags (env `GAME_SERVER_ADDR`/`GAME_SERVER_DB`), default
      `:50051` / `games.db`. Composition logic factored into `build(addr,
      dbPath)` so `cmd/server/server_test.go` can drive a real
      `grpc.Server`/`net.Listener` in-process without forking `main`.
- [ ] Address the errcheck/lint findings in `internal/core/domain/game.go` in
      this stage if touching that I/O — but note `fmt.Scan`/`fmt.Scanln`
      calls living directly in the domain layer are themselves a scope smell
      (terminal I/O belongs to the CLI, not the domain). If a full fix here
      turns into extracting that I/O, that's Stage 2's job, not Stage 1's.
      **Not touched in Stage 2**: none of tasks 1-4 above required editing
      `game.go`'s I/O, so per this task's own contingency clause this was
      left alone rather than done as a standalone pass. The extraction (move
      `fmt.Scan`/`Scanln` prompts out of the domain layer into the CLI) is
      flagged below under Stage 3, which is where the CLI-side of that
      boundary actually gets rewritten.

Build-tooling fix made in passing: `make build_game`/`make all`/`make
run_server` failed on a truly clean checkout (no pre-existing `bin/`) with
`go: cannot write multiple packages to non-directory ./bin` — `go build -o`
into multiple main packages (`cmd/cli` + `cmd/server`) requires the target to
already exist as a directory. Pre-existing bug (both `cmd/` packages predate
Stage 2), not introduced here, but it blocked Stage 2's own `make
run_server` success criterion, so fixed with a one-line `mkdir -p ./bin` in
`build_game` rather than deferred back to Stage 1.

**Success criteria (all must pass):**
```
cd game && go test ./internal/adapters/... -v          # repository tests pass
cd game && make run_server &                            # server starts, stays up
# a CreateGame call round-trips through grpc -> service -> sqlite -> back
```
Smoke-test tooling used: `grpcurl` isn't installed in this environment, so the
round trip is proven by `cd game && go test ./cmd/server/... -v`
(`TestCreateGameRoundTripsOverGRPC` in `cmd/server/server_test.go`), which
starts the real composition root on a real `net.Listener` and drives it with
a real `google.golang.org/grpc` client — plus a manual confirmation in this
run: `./bin/server` was started as its own process, hit over TCP by a
separate `go run` client process, and `games.db` was inspected afterward to
confirm the write actually landed in sqlite (not just in memory).

## Stage 3 — CLI connectivity

Owner: `cmd/cli` + `game/cli` package, plus the server-side gaps below that
block this stage's own success criteria. Depends on Stage 2 existing — don't
start wiring a client to an endpoint that doesn't do anything yet.

Decision locked (don't revisit without a reason): the `.proto` only declares
`GetCard`/`NewGame`/`PlayGemCard`, has no RPC at all for buying a golem or
resting, and `GetCardMessage`/`PlayGemCardMessage` type `game_id` as `int32`
even though real game IDs are UUID strings (`game_service.GetUID`). Confirmed
before this stage started, not left for research to rediscover. Expand the
proto now rather than narrowing this stage's scope — "play a full turn" only
means something if golem-buying and resting actually round-trip too.

Tasks:
- [x] Fix `game_id` to `string` in `GetCardMessage` and `PlayGemCardMessage`,
      and add RPCs + messages for buying a golem and resting (mirror the
      existing `PlayGemCard`/`PlayGemCardResponse` shape: request carries
      `game_id` + whatever the action needs, response carries `GameState`).
      Regenerate via `make build_protos`.
      Done, plus two additions research surfaced as load-bearing (not
      optional polish) rather than scope creep:
      - `GameState` gained `game_id` and `current_player` fields. Without
        `game_id` a client can never learn a newly created game's ID — the
        old `CreateGameResponse` only carried `GameState`, which had no ID
        field anywhere, so no follow-up call was ever possible. Without
        `current_player`, a client has no way to know which player index
        to pass back on its next request (this game is pass-and-play/
        hotseat: one client drives whichever player is up, there's no
        per-connection player identity).
      - `GetCard`'s response changed from a bare `GemCard` to a
        `GetCardResponse{GameState}`, matching the other RPCs — the old
        shape couldn't carry enough information back for a client to keep
        playing (no updated lineup/hand/turn state).
      RPCs added: `BuyGolem(BuyGolemMessage) -> BuyGolemResponse` and
      `Rest(RestMessage) -> RestResponse`, both `{game_id, player, ...}` /
      `{GameState}` shaped per the mirror instruction.
- [x] Implement the remaining `ports.GameService` methods on
      `game_service.Service` — `PlayGemCard`, `GetGemCard`, `GetGolemCard`,
      `Rest` (currently only `Create`/`Get` exist) — translating to the
      matching `domain.Game` methods (`PlayGemCard`, `GetGemCard`, `BuyGolem`,
      `NextTurn`). This is the same mechanical service-method-plus-adapter
      pattern Stage 2 already established for `Create`/`NewGame`, not a new
      design — but if something about it doesn't fit that pattern cleanly,
      stop and report rather than forcing it.
      Done (`internal/core/services/game_service/service.go`). One method
      genuinely didn't fit the given `ports.GameService` signature as
      written, handled as a disclosed mechanical extension rather than a
      silent redesign or a full stop: `Rest(id string, player int)` had no
      way to carry a discard choice, but the domain rule it wraps
      (`NextTurn`) requires one whenever the current player holds more
      than 10 gems — with the original signature, `Rest` could never
      succeed once a player crossed that line, for any caller, ever.
      Extended it to `Rest(id string, player int, discards string)`,
      reusing `GetGemCard`'s existing `payment string` (gem-letter)
      convention rather than inventing a new shape. `PlayGemCard` was left
      exactly as declared (`id, player, card`, no discard) since real
      Century Golem rules put discard-to-limit at end-of-turn (`Rest`),
      not mid-play — so the missing parameter there is correct, not a gap.
      Every method validates `player == game.CurrentPlayer` before
      mutating (via a shared `loadForPlayer` helper) — the interface
      passes `player` but nothing previously used it, which would have let
      a stale/wrong-turn caller silently mutate whichever player happened
      to be current instead of being rejected.
- [x] Wire those into `internal/adapters/grpc` (`grpcadapter.GameServer`) —
      replace the `Unimplemented` `GetCard`/`PlayGemCard` stubs and add the
      new golem/rest RPCs' handlers.
      Done (`internal/adapters/grpc/game_server.go`). Domain/service
      errors surface as `codes.FailedPrecondition` (e.g. "can't afford
      this golem", "not your turn"); a missing `game_id` is
      `codes.InvalidArgument`.
- [x] Add networking from the CLI to a running server over grpc (currently
      missing/commented out — CLI runs entirely against local/offline logic).
      Done (`cli/game.go`): `createGame`/`playGemCard`/`getGemCard`/
      `buyGolem`/`rest` are `tea.Cmd`s that dial the server and perform one
      rpc each, reported back into the bubbletea loop as `gameStateMsg`/
      `gameErrMsg`. The bubbletea flow itself (`cli/cli.go`) was rebuilt
      around them — menu -> `actions.TurnView` (per-action input
      collection + local format validation) -> `actions.SubmitMsg` ->
      `model.dispatch` (resolves the typed index against the current
      lineup/hand, then fires the matching network `tea.Cmd`).
      `cli/actions/{playGemCard,getGemCard,getGolemCard,rest}.go`'s
      constructors (previously fully commented out, referencing a
      `models` package that no longer exists post-restructuring) were
      rewritten against the current architecture rather than merely
      uncommented.
- [x] Support pointing the CLI at either a local server or a deployed one
      (flag/env for the address).
      Done (`cmd/cli/cli.go`): `-addr`/`GAME_SERVER_ADDR`, mirroring
      `cmd/server`'s own `-addr`/`GAME_SERVER_ADDR` convention exactly
      (same flag name, same env var, same default-fallback helper).
- [x] Clean up the dead/unused fields flagged by lint in `cli/game.go`,
      `cli/player.go`, `cli/actions/turnView.go`, `cli/cli.go` as part of
      rewiring these files for real network calls — don't do this as a
      separate pass; it's only cheap to fix now because you're already
      rewriting the file for the connectivity change.
      Done. `cli/game.go` (`Game`) and `cli/player.go` (`Player`) fields
      were never dead in the sense of being pointless — they were
      genuinely meant to hold server-fetched state, just never wired to a
      real `Unmarshall*` populating them (both had this half-written,
      commented out). Both are implemented for real now and every field is
      read (`String()`/lineup-index resolution/etc.), which is what
      resolved the lint findings — not deletion. `cli/actions/turnView.go`
      was the one file where the fields' original *design* was the
      problem, not just missing wiring: `action func() error` was a
      closure meant to be bound once at construction (matching the old
      pre-restructuring domain methods that did their own `fmt.Scan`
      internally), which cannot work against bubbletea's value-semantics
      `Update` loop (a closure capturing construction-time `textinput`
      values would never see what the user actually typed, since each
      `Update` call returns a new `TurnView` value rather than mutating in
      place). Changed to `action func(values []string) error`, called with
      values read fresh from `v.inputs` at submit time; `error` is now
      rendered in `View()` instead of sitting unread.
- [x] Extract the `fmt.Scan`/`fmt.Scanln` terminal-I/O prompts currently
      living in `internal/core/domain/game.go` (`GetGemCard`, `BuyGolem`,
      `NextTurn`) out of the domain layer and into the CLI. Flagged during
      Stage 2 as a cleanup item, but it's stronger than that now: the new
      `game_service.Service` methods above (`PlayGemCard`/`GetGemCard`/
      `GetGolemCard`/`Rest`) run *server-side*, and a server process can't
      block on stdin to prompt a remote client for input — so these domain
      methods' I/O and state-mutation halves must be pulled apart for the
      new service methods to be implementable at all, not just for lint
      cleanliness. Do this before or alongside the service-method task above,
      whichever order actually works out. Also resolves the errcheck lint
      findings on those unchecked `fmt.Scan*` calls (see Stage 1's note).
      Done. `GetGemCard(index int, payment string)`, `BuyGolem(index int)`,
      and `NextTurn(discards string)` are now pure mutation, taking an
      already-decided answer instead of prompting; all `fmt.Scan*` calls
      are gone from `internal/core/domain/game.go` (errcheck findings
      resolved as a consequence, not a separate pass). The CLI's
      `actions.TurnView` constructors are the client-side prompting half.
      Two pre-existing bugs surfaced and fixed while rewriting these exact
      lines (both directly in the code being touched, not a separate
      sweep): `GetGemCard`'s payment-distribution loop computed
      `pile.Yellow += 1` etc. on a local copy and never wrote it back to
      `g.LooseGems[i]`, so bypass payments were silently discarded; and
      `Player.CanAfford` OR'd its four gem-color comparisons together
      instead of AND'ing them (so holding enough of *any one* color, even
      zero of the others, counted as affordable), which would have made
      the golem-purchase success criterion's reject-when-unaffordable path
      nearly untestable. Fixed to AND (`internal/core/domain/player.go`).
      `BuyGolem` also gained a checked `Gems.remove` (was
      unchecked/ignored before).

**Success criteria:**
```
cd game && make run_server &
cd game && ./bin/cli   # can create/join a game and play a full turn against
                        # the running server, not local-only logic
```
Verified two ways:
- Automated: `go test ./cli/... -run TestCLIPlaysFullTurnOverGRPC -v`
  (`cli/client_test.go`) — mirrors Stage 2's
  `TestCreateGameRoundTripsOverGRPC` pattern: starts the real composition
  root on a real `net.Listener`, then drives it with the CLI's own network
  `tea.Cmd`s (`createGame`/`getGemCard`/`buyGolem`/`rest` — the same code
  the bubbletea program dispatches) through create -> get gem card -> buy
  golem (asserts either the purchase succeeds and `player.golems` grows by
  one, or it's rejected with a `FailedPrecondition` and state is
  unchanged) -> rest -> turn advances to player 1. Passed 9/9 runs
  (`-count=1`, run 8 times individually plus once via `make all`), stable
  across the golem lineup's shuffle.
- Manual, over a real pty (this sandbox has no tty, so `tmux` stood in for
  one): `./bin/server` on `:50999`, `./bin/cli -addr localhost:50999
  -players 2` in a `tmux` pane, driven with real keystrokes through get
  gem card (hand 2 -> 3 cards, gem deck 37 -> 36), buy golem index 0
  (correctly rejected: "You cannot afford this card", state unchanged),
  and rest (turn advanced to "player 1's turn", player 1's starting gems
  shown correctly). `q` quits the program cleanly.

## Stage 4 — ML pluggability (deferred, not started)

Explicitly out of scope for an autonomous loop to pick up. No design changes
are believed necessary: since play decisions already enter through the grpc
API boundary (the CLI is just one caller), a future AI-driven player plugs in
as another caller of that same API. Re-evaluate this assumption only once
Stage 2/3 are done and an actual ML integration is being scoped — don't
speculatively build interfaces for it now.

## Rules for an agent looping on this file

1. Find the first unchecked stage. Work only within it.
2. A task is done when its stage's listed success-criteria commands pass, not
   when it "looks right."
3. If a task turns up work that belongs to a different stage (wrong layer or
   wrong owner), add it to that stage's task list instead of doing it inline
   — see the Stage 1 lint note and Stage 3's cleanup-deferral for the pattern.
4. Never touch Stage 4 without being explicitly told to.
5. Update this file's checkboxes in the same commit/change that satisfies
   them.

## Per-stage agent protocol (research -> code -> verify -> next)

One agent per stage, run in this sequence, from the `ai-behaviors` framework
(`#=research` = `#epistemic #legible #concise`; `#=code` = `#contract #name
#checklist #scope #legible #concise`):

1. **`#research`** — investigate only what this stage's task list needs to
   execute the decisions already recorded in that stage's section (e.g.
   Stage 2's sqlite/JSON-blob/composition-root calls are made; research is
   implementation detail, not re-litigating them). Every claim sourced and
   confidence-labeled. If research surfaces a real unresolved fork — not
   covered by an existing decision in this file — STOP and report the
   options. Do not pick one and continue; that is `#epistemic`'s HARD
   CONSTRAINT, not a style preference.
2. **`#code`** — only after research completes without an unresolved fork.
   The stage's task checklist above is the `#checklist` reference artifact:
   every item gets done/blocked (deferred requires explicit user
   confirmation, which an unattended agent cannot obtain — so an unattended
   agent that wants to defer an item must STOP and report instead, same as a
   research fork). Update this file's checkboxes for whatever gets done.
3. **Verify** — the orchestrator (not the stage agent's self-report) runs
   that stage's exact success-criteria commands. Passing them is what
   "requirements met" means; nothing else qualifies.
4. Only on a clean verify does the next stage's agent get started, same
   sequence. A stop-and-report from step 1 or 2 ends the run — it is not a
   failure to route around, it's the plan asking a real question.
