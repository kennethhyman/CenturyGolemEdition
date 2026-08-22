package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

// TestCreateGameRoundTripsOverGRPC is the Stage 2 smoke test: it starts the
// real composition root (sqlite repository + game_service.Service + grpc
// adapter) on a real TCP listener, then drives it with a real grpc client —
// proving NewGame round-trips grpc -> service -> sqlite -> back, not just
// that the pieces compile against each other.
func TestCreateGameRoundTripsOverGRPC(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "games.db")

	lis, srv, repo, err := build("127.0.0.1:0", dbPath)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	t.Cleanup(func() { repo.Close() })

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(lis) }()
	t.Cleanup(srv.GracefulStop)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("DialContext(%s): %v", lis.Addr().String(), err)
	}
	t.Cleanup(func() { conn.Close() })

	client := pb.NewGameClient(conn)

	resp, err := client.NewGame(ctx, &pb.CreateGameMessage{PlayerCount: 3})
	if err != nil {
		t.Fatalf("NewGame: %v", err)
	}

	state := resp.GetGameState()
	if state == nil {
		t.Fatal("NewGame response has nil GameState")
	}
	if got, want := len(state.GetGemLineup()), 6; got != want {
		t.Errorf("len(GemLineup) = %d, want %d", got, want)
	}
	if got, want := len(state.GetGolemLineup()), 5; got != want {
		t.Errorf("len(GolemLineup) = %d, want %d", got, want)
	}
	// 2*players gold/silver coins per domain.NewGame, for 3 players.
	if got, want := state.GetGoldCoins(), int32(6); got != want {
		t.Errorf("GoldCoins = %d, want %d", got, want)
	}
	if got, want := state.GetSilverCoins(), int32(6); got != want {
		t.Errorf("SilverCoins = %d, want %d", got, want)
	}
	if state.GetPlayer() == nil {
		t.Error("GameState.Player is nil")
	}

	select {
	case err := <-serveErr:
		t.Fatalf("srv.Serve returned early: %v", err)
	default:
	}
}

// TestNewGameRejectsNonPositivePlayerCount exercises the adapter's own
// validation, independent of game_service.
func TestNewGameRejectsNonPositivePlayerCount(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "games.db")

	lis, srv, repo, err := build("127.0.0.1:0", dbPath)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	t.Cleanup(func() { repo.Close() })

	go srv.Serve(lis)
	t.Cleanup(srv.GracefulStop)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("DialContext(%s): %v", lis.Addr().String(), err)
	}
	t.Cleanup(func() { conn.Close() })

	client := pb.NewGameClient(conn)

	if _, err := client.NewGame(ctx, &pb.CreateGameMessage{PlayerCount: 0}); err == nil {
		t.Fatal("NewGame with PlayerCount: 0, want error, got nil")
	}
}
