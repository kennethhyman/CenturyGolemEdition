package sqlite

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "games.db")
	repo, err := New(dbPath)
	if err != nil {
		t.Fatalf("New(%q): %v", dbPath, err)
	}
	t.Cleanup(func() {
		if err := repo.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})

	return repo
}

func TestSaveThenGetRoundTrips(t *testing.T) {
	repo := newTestRepository(t)
	game := *domain.NewGame("game-1", 3)

	if err := repo.Save(game); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Get("game-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.ID != game.ID {
		t.Errorf("ID = %q, want %q", got.ID, game.ID)
	}
	if len(got.Players) != len(game.Players) {
		t.Errorf("len(Players) = %d, want %d", len(got.Players), len(game.Players))
	}
	if got.GoldCoins != game.GoldCoins || got.SilverCoins != game.SilverCoins {
		t.Errorf("coins = (%d,%d), want (%d,%d)", got.GoldCoins, got.SilverCoins, game.GoldCoins, game.SilverCoins)
	}
	if got.GolemLimit != game.GolemLimit {
		t.Errorf("GolemLimit = %d, want %d", got.GolemLimit, game.GolemLimit)
	}

	// LineUp's `deck` field is unexported; this exercises LineUp's custom
	// MarshalJSON/UnmarshalJSON, which is what keeps Remaining/Draw/String
	// working on a round-tripped game instead of nil-pointer-panicking.
	if got.GemCards.Remaining() != game.GemCards.Remaining() {
		t.Errorf("GemCards.Remaining() = %d, want %d", got.GemCards.Remaining(), game.GemCards.Remaining())
	}
	if got.Golems.Remaining() != game.Golems.Remaining() {
		t.Errorf("Golems.Remaining() = %d, want %d", got.Golems.Remaining(), game.Golems.Remaining())
	}
	if !reflect.DeepEqual(got.GemCards.Stack, game.GemCards.Stack) {
		t.Errorf("GemCards.Stack mismatch:\ngot  %+v\nwant %+v", got.GemCards.Stack, game.GemCards.Stack)
	}
	if !reflect.DeepEqual(got.Golems.Stack, game.Golems.Stack) {
		t.Errorf("Golems.Stack mismatch:\ngot  %+v\nwant %+v", got.Golems.Stack, game.Golems.Stack)
	}

	// A round-tripped LineUp must still be able to draw without panicking.
	if _, err := got.GemCards.Draw(0); err != nil {
		t.Errorf("GemCards.Draw(0) on round-tripped game: %v", err)
	}
}

func TestGetMissingReturnsNotFound(t *testing.T) {
	repo := newTestRepository(t)

	_, err := repo.Get("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get: err = %v, want it to wrap ErrNotFound", err)
	}
}

func TestSaveOverwritesExisting(t *testing.T) {
	repo := newTestRepository(t)
	game := *domain.NewGame("game-2", 2)

	if err := repo.Save(game); err != nil {
		t.Fatalf("Save: %v", err)
	}

	game.GoldCoins = 99
	if err := repo.Save(game); err != nil {
		t.Fatalf("Save (overwrite): %v", err)
	}

	got, err := repo.Get("game-2")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.GoldCoins != 99 {
		t.Errorf("GoldCoins = %d, want 99", got.GoldCoins)
	}
}

func TestSaveRejectsEmptyID(t *testing.T) {
	repo := newTestRepository(t)

	err := repo.Save(domain.Game{})
	if err == nil {
		t.Fatal("Save with empty ID: want error, got nil")
	}
}

func TestPersistsAcrossReopen(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "games.db")

	repo, err := New(dbPath)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	game := *domain.NewGame("game-3", 4)
	if err := repo.Save(game); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := repo.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	reopened, err := New(dbPath)
	if err != nil {
		t.Fatalf("New (reopen): %v", err)
	}
	t.Cleanup(func() { reopened.Close() })

	got, err := reopened.Get("game-3")
	if err != nil {
		t.Fatalf("Get after reopen: %v", err)
	}
	if got.ID != game.ID {
		t.Errorf("ID = %q, want %q", got.ID, game.ID)
	}
}
