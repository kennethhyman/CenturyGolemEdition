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
  // Rest ends player's turn. discards is a string of gem letters (y/g/b/p)
  // to discard first, needed (and validated) only when player holds more
  // than 10 gems; pass "" otherwise. Added here (the original signature had
  // no way to carry a discard choice) to mirror GetGemCard's `payment`
  // convention — without it Rest could never succeed once a player held too
  // many gems, since there was no way for a caller to ever supply the
  // discard the domain rule requires.
  Rest(id string, player int, discards string) (domain.Game, error)
}
