package game_service

import (
	"fmt"

  "github.com/google/uuid"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/ports"
)

type Service struct{
  gamesRepository ports.GameRepository
}

func New() *Service {
  return &Service{}
}

func (srv *Service) Get(id string) (domain.Game, error) {
  game, err := srv.gamesRepository.Get(id)

  if err != nil {
    return domain.Game{}, fmt.Errorf("Game Repository failed to get game, %v\n", err)
  }

  return game, err
}

func (srv *Service) Create(players int) (domain.Game, error) {
  game := *domain.NewGame(GetUID(), players)

  if err := srv.gamesRepository.Save(game); err != nil {
    return domain.Game{}, fmt.Errorf("create game into repository has failed: %v\n", err)
	}

  return game, nil
}

func (srv *Service) PlayGemCard(id string, player int, card domain.GemCard, discards domain.GemValues) (domain.Game, error) {
  game, err := srv.Get(id)
  if err != nil {
    return game, err
  }

  err = game.PlayGemCard(card, discards)
  if err != nil {
    return game, err
  }

  err = srv.gamesRepository.Save(game)
  if err != nil {
    return game, err
  }

  return game, nil
}

func (srv *Service) GetGemCard(id string, player int, card domain.GemCard, payment string) (domain.Game, error) {
	return domain.Game{}, nil
}

func (srv *Service) GetGolemCard(id string, player int, golem domain.GolemCard) (domain.Game, error) {
	return domain.Game{}, nil
}

func (srv *Service) Rest(id string, player int) (domain.Game, error) {
	return domain.Game{}, nil
}

func GetUID() string {
  return uuid.NewString()
}
