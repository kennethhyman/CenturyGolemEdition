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

Owner: `cmd/cli` + `game/cli` package. Depends on Stage 2 existing — don't
start wiring a client to an endpoint that doesn't do anything yet.

Tasks:
- [ ] Add networking from the CLI to a running server over grpc (currently
      missing/commented out — CLI runs entirely against local/offline logic).
- [ ] Support pointing the CLI at either a local server or a deployed one
      (flag/env for the address).
- [ ] Clean up the dead/unused fields flagged by lint in `cli/game.go`,
      `cli/player.go`, `cli/actions/turnView.go`, `cli/cli.go` as part of
      rewiring these files for real network calls — don't do this as a
      separate pass; it's only cheap to fix now because you're already
      rewriting the file for the connectivity change.
- [ ] Extract the `fmt.Scan`/`fmt.Scanln` terminal-I/O prompts currently
      living in `internal/core/domain/game.go` (`GetGemCard`, `BuyGolem`,
      `NextTurn`) out of the domain layer and into the CLI. Flagged during
      Stage 2 (see that stage's task list) rather than done there: Stage 2
      didn't need to touch `game.go`, and once the CLI actually talks to a
      running server (this stage's whole point), those domain methods either
      stop being the CLI's code path entirely or need their I/O and their
      game-state mutation pulled apart so the state-mutation half can run
      server-side while the prompt-for-input half stays client-side. Also
      resolves the errcheck lint findings on those unchecked `fmt.Scan*`
      calls (see Stage 1's note).

**Success criteria:**
```
cd game && make run_server &
cd game && ./bin/cli   # can create/join a game and play a full turn against
                        # the running server, not local-only logic
```

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
