# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go implementation of the Century: Golem Edition board game. Goals:
1. Answer strategic/optimality questions about the game via simulation
2. Create a playable online version of the game

The `game/` directory contains all Go source code. The root also contains Python/Ruby scripts (`parse_cards.rb`, `golem_analysis.py`, etc.) for data parsing and analysis.

## Commands

All commands are run from within `game/`:

```bash
cd game

# Build everything (protos + game binaries)
make

# Build just the game binaries
make build_game

# Build protobuf/gRPC files (requires protoc with go and go-grpc plugins)
make build_protos

# Run tests (also builds first)
make test_game

# Run tests directly with verbose output
go test ./internal/core/domain/... -v

# Run a specific test
go test ./internal/core/domain/... -run TestName -v
```

The `make test_game` target currently runs `go test ./models/ -v` (an outdated path); use the `go test` commands above instead for the restructured layout.

Binaries are output to `game/bin/`: `cli` and `server`.

## Architecture

The project follows hexagonal (ports and adapters) architecture:

```
game/
├── cmd/
│   ├── cli/        # Entry point: launches Bubbletea TUI
│   └── server/     # Entry point: gRPC server (stub)
├── internal/core/
│   ├── domain/     # Pure game logic (no I/O dependencies)
│   ├── ports/      # Interfaces: GameRepository, GameService
│   ├── services/game_service/  # GameService implementation
│   └── grpc/       # Protobuf generated code + marshal/unmarshal adapters
├── adapters/       # Intended for adapter implementations (currently empty)
└── cli/            # Bubbletea TUI views and actions
```

### Domain Layer (`internal/core/domain/`)

The core game model:
- **`Game`** — top-level state: players, gem card lineup, golem lineup, coins, current player
- **`Player`** — hand, discard pile, gems (`GemValues`), coins, acquired golems
- **`GemCard`** — has `Inputs`, `Outputs` (GemValues), and `Upgrades` count
- **`GolemCard`** — victory point cards purchased with gems
- **`GemValues`** — `{Yellow, Green, Blue, Pink int}` representing gem counts (gem value hierarchy: Yellow=1, Green=2, Blue=3, Pink=4)
- **`LineUp[T]`** — a visible row of cards backed by a `Deck[T]`; drawing from a position slides cards and refills from deck
- **`Deck[T]`** — generic shuffleable deck loaded from embedded CSV resources

Card data is embedded via `//go:embed resources/*` — CSV files in `internal/core/domain/resources/` (`cards.csv`, `golems.csv`, `starter_cards.csv`).

### Player Actions (implemented on `Game`)

- `PlayGemCard(card, discards)` — play a conversion card, enforcing 10-gem hand limit
- `GetGemCard()` — pick up a card from the lineup (paying gems to bypass earlier slots)
- `PickUpGemCards()` / Rest — return discard pile to hand
- `BuyGolem()` — purchase a golem card with gems

### Service Layer (`internal/core/services/game_service/`)

`Service` implements `ports.GameService` and `ports.GameRepository`. Currently `New()` creates a service with a nil repository — the repository must be wired in. Methods `GetGemCard`, `GetGolemCard`, and `Rest` are stubs.

### Transport Layer

- **gRPC**: Proto defined at `internal/core/grpc/CenturyGolemEdition.proto`. Marshall/Unmarshall functions live on domain types (e.g., `MarshallGame`, `MarshallPlayer`, `UnmarshallGemCard`).
- **CLI**: Bubbletea TUI in `cli/` package, launched from `cmd/cli/`.

### Current State

The codebase is mid-restructuring to hexagonal architecture. Key gaps:
- `adapters/` directory exists but is empty (no repository implementations)
- `game_service.New()` doesn't inject a repository — service will panic on use
- Several `GameService` interface methods are unimplemented stubs
- `make test_game` has a stale test path
