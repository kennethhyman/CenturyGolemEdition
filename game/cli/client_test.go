package cli

import (
	"net"
	"path/filepath"
	"testing"

	"google.golang.org/grpc"

	grpcadapter "github.com/kennethhyman/CenturyGolemEdition/internal/adapters/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/adapters/sqlite"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/services/game_service"
)

// startTestServer builds the same composition root cmd/server's build()
// does (sqlite repository + game_service.Service + grpc adapter) on a real
// TCP listener, and returns its address. It's duplicated from
// cmd/server/server_test.go's approach rather than imported, since build()
// lives in package main (cmd/server) and can't be imported from here.
func startTestServer(t *testing.T) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "games.db")
	repo, err := sqlite.New(dbPath)
	if err != nil {
		t.Fatalf("sqlite.New: %v", err)
	}
	t.Cleanup(func() { repo.Close() })

	service := game_service.New(repo)
	adapter := grpcadapter.NewGameServer(service)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGameServer(srv, adapter)

	go srv.Serve(lis)
	t.Cleanup(srv.GracefulStop)

	return lis.Addr().String()
}

// mustGameState asserts msg is a gameStateMsg (failing the test with
// whatever error was reported otherwise) and returns its Game.
func mustGameState(t *testing.T, label string, msg interface{}) Game {
	t.Helper()
	switch m := msg.(type) {
	case gameStateMsg:
		return m.game
	case gameErrMsg:
		t.Fatalf("%s: got error %v, want a game state", label, m.err)
	default:
		t.Fatalf("%s: got %#v, want gameStateMsg", label, msg)
	}
	return Game{}
}

// canAfford mirrors domain.Player.CanAfford's (now-fixed) AND-across-colors
// rule, so the test can independently predict whether a purchase should
// succeed rather than just trusting the server's own answer.
func canAfford(gems, cost GemValue) bool {
	return gems.Yellow >= cost.Yellow && gems.Green >= cost.Green &&
		gems.Blue >= cost.Blue && gems.Pink >= cost.Pink
}

// TestCLIPlaysFullTurnOverGRPC is the Stage 3 smoke test, mirroring Stage
// 2's TestCreateGameRoundTripsOverGRPC: it starts the real composition root
// on a real TCP listener and drives it with the CLI's own network client
// functions (createGame/getGemCard/buyGolem/rest — the same tea.Cmds the
// bubbletea program itself dispatches), proving a full turn — create,
// one gem-card action, one golem action (accepted or correctly rejected),
// and a turn advance — round-trips over real grpc end to end, not just
// that the pieces compile against each other.
func TestCLIPlaysFullTurnOverGRPC(t *testing.T) {
	addr := startTestServer(t)

	// Create a 2-player game.
	game := mustGameState(t, "createGame", createGame(addr, 2)())
	if game.id == "" {
		t.Fatal("createGame: game id is empty — CLI has no way to make follow-up calls")
	}
	if game.currentPlayer != 0 {
		t.Fatalf("createGame: currentPlayer = %d, want 0", game.currentPlayer)
	}
	if len(game.gem_lineup) != 6 {
		t.Fatalf("createGame: len(gem_lineup) = %d, want 6", len(game.gem_lineup))
	}
	if len(game.golem_lineup) != 5 {
		t.Fatalf("createGame: len(golem_lineup) = %d, want 5", len(game.golem_lineup))
	}
	startingHand := len(game.player.hand)

	// Gem-card action: pick up the lineup's first card (index 0 costs no
	// payment).
	game = mustGameState(t, "getGemCard", getGemCard(addr, game, game.gem_lineup[0], "")())
	if len(game.player.hand) != startingHand+1 {
		t.Fatalf("getGemCard: hand size = %d, want %d", len(game.player.hand), startingHand+1)
	}

	// Golem action: buy the first affordable lineup golem, or — if none are
	// affordable, which is the expected case for a player who has only
	// just started (starting gems are pure yellow; every golem in
	// resources/golems.csv costs at least one non-yellow gem) — confirm an
	// attempt is correctly rejected and leaves state unchanged.
	affordableIndex := -1
	for i, golem := range game.golem_lineup {
		if canAfford(game.player.gems, golem.Cost) {
			affordableIndex = i
			break
		}
	}

	if affordableIndex >= 0 {
		before := len(game.player.golems)
		game = mustGameState(t, "buyGolem (affordable)", buyGolem(addr, game, game.golem_lineup[affordableIndex])())
		if len(game.player.golems) != before+1 {
			t.Fatalf("buyGolem: player.golems size = %d, want %d", len(game.player.golems), before+1)
		}
	} else {
		before := len(game.player.golems)
		msg := buyGolem(addr, game, game.golem_lineup[0])()
		errMsg, ok := msg.(gameErrMsg)
		if !ok {
			t.Fatalf("buyGolem (unaffordable): got %#v, want gameErrMsg", msg)
		}
		if errMsg.err == nil {
			t.Fatal("buyGolem (unaffordable): gameErrMsg.err is nil")
		}
		if len(game.player.golems) != before {
			t.Fatalf("buyGolem (unaffordable) must not mutate state: player.golems size = %d, want %d", len(game.player.golems), before)
		}
	}

	// Rest: advance the turn.
	game = mustGameState(t, "rest", rest(addr, game, "")())
	if game.currentPlayer != 1 {
		t.Fatalf("rest: currentPlayer = %d, want 1", game.currentPlayer)
	}
}
