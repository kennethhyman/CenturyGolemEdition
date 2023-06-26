package ports

import (
  "github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

type GameRepository interface {
  Get(id string) (domain.Game, error)
  Save(game domain.Game) error
}

type GameService interface {
  Create(players int) (domain.Game, error)
  Get(id string) (domain.Game, error)
  PlayGemCard(id string, player int, card domain.GemCard) (domain.Game, error)
  GetGemCard(id string, player int, card domain.GemCard, payment string) (domain.Game, error)
  GetGolemCard(id string, player int, golem domain.GolemCard) (domain.Game, error)
  Rest(id string, player int) (domain.Game, error)
}
