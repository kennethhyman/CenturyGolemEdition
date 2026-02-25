package repository

import (
	"fmt"
	"sync"

	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

type InMemoryRepository struct {
	mu    sync.RWMutex
	games map[string]domain.Game
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		games: make(map[string]domain.Game),
	}
}

func (r *InMemoryRepository) Get(id string) (domain.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	game, ok := r.games[id]
	if !ok {
		return domain.Game{}, fmt.Errorf("game %v not found", id)
	}

	return game, nil
}

func (r *InMemoryRepository) Save(game domain.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.games[game.ID] = game
	return nil
}
