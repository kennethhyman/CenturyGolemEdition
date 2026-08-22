// Package sqlite implements ports.GameRepository against a local sqlite
// database, via the pure-Go modernc.org/sqlite driver (no CGO, so build/CI
// stay simple). Each domain.Game is stored as a JSON blob keyed by its ID —
// there is no normalized schema, matching Stage 2's decision to defer that
// until a production database is actually chosen.
package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/ports"
)

// ErrNotFound is returned by Get when no game is stored under the given ID.
var ErrNotFound = errors.New("sqlite: game not found")

const schema = `
CREATE TABLE IF NOT EXISTS games (
	id   TEXT PRIMARY KEY,
	data TEXT NOT NULL
);`

// Repository is a ports.GameRepository backed by sqlite.
type Repository struct {
	db *sql.DB
}

var _ ports.GameRepository = (*Repository)(nil)

// New opens the sqlite database at dsn (a file path, "file::memory:", or
// ":memory:" — passed straight to modernc.org/sqlite) and ensures the games
// table exists.
//
// Post: on success, the returned *Repository's Get/Save are ready to use.
// Caller owns the returned Repository and must Close it.
func New(dsn string) (*Repository, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open %q: %w", dsn, err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: create games table: %w", err)
	}

	return &Repository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *Repository) Close() error {
	return r.db.Close()
}

// Get returns the game stored under id.
//
// Pre: id is the ID a game was previously Saved under.
// Post: returns a domain.Game equivalent (by value) to what was Saved,
// including LineUp/Deck state, via LineUp's custom (Un)MarshalJSON. If no
// game is stored under id, returns an error wrapping ErrNotFound.
func (r *Repository) Get(id string) (domain.Game, error) {
	var data string
	err := r.db.QueryRow(`SELECT data FROM games WHERE id = ?`, id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Game{}, fmt.Errorf("sqlite: game %q: %w", id, ErrNotFound)
		}
		return domain.Game{}, fmt.Errorf("sqlite: get game %q: %w", id, err)
	}

	var game domain.Game
	if err := json.Unmarshal([]byte(data), &game); err != nil {
		return domain.Game{}, fmt.Errorf("sqlite: decode game %q: %w", id, err)
	}

	return game, nil
}

// Save upserts game, keyed by game.ID.
//
// Pre: game.ID is non-empty.
// Post: a subsequent Get(game.ID) returns a game equivalent to this one.
// Invariant: at most one row exists per ID — a second Save with the same ID
// replaces the first rather than erroring or duplicating.
func (r *Repository) Save(game domain.Game) error {
	if game.ID == "" {
		return errors.New("sqlite: save game: empty ID")
	}

	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("sqlite: encode game %q: %w", game.ID, err)
	}

	_, err = r.db.Exec(
		`INSERT INTO games (id, data) VALUES (?, ?)
		 ON CONFLICT(id) DO UPDATE SET data = excluded.data`,
		game.ID, string(data),
	)
	if err != nil {
		return fmt.Errorf("sqlite: save game %q: %w", game.ID, err)
	}

	return nil
}
