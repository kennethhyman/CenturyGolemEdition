package game_service

import (
	"errors"
	"fmt"

  "github.com/google/uuid"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/ports"
)

type Service struct{
  gamesRepository ports.GameRepository
}

var _ ports.GameService = (*Service)(nil)

// New builds a Service backed by gamesRepository.
//
// Pre: gamesRepository is non-nil — which repository implementation (sqlite
// today, e.g. Postgres later) is the composition root's (cmd/server's)
// decision, never this package's.
func New(gamesRepository ports.GameRepository) *Service {
  return &Service{gamesRepository: gamesRepository}
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

// PlayGemCard plays card for player in game id, discarding no gems (this
// entry point carries no discard choice — see ports.GameService.PlayGemCard).
//
// Pre: player == the game's CurrentPlayer (else "not this player's turn");
// player holds card in hand and can afford its Inputs.
// Post: on success, the updated game is saved and returned.
func (srv *Service) PlayGemCard(id string, player int, card domain.GemCard) (domain.Game, error) {
  game, err := srv.loadForPlayer(id, player)
  if err != nil {
    return domain.Game{}, err
  }

  if err := game.PlayGemCard(card, domain.GemValues{}); err != nil {
    return domain.Game{}, fmt.Errorf("play gem card: %w", err)
  }

  return srv.save(game)
}

// GetGemCard picks up card from game id's gem lineup for player, paying
// payment (one gem letter y/g/b/p per lineup card bypassed) for cards
// bypassed to reach it.
//
// Pre: player == the game's CurrentPlayer; card matches (by value) an entry
// in the lineup.
// Post: on success, the updated game is saved and returned.
func (srv *Service) GetGemCard(id string, player int, card domain.GemCard, payment string) (domain.Game, error) {
  game, err := srv.loadForPlayer(id, player)
  if err != nil {
    return domain.Game{}, err
  }

  index, err := indexOfGemCard(game.GemCards.Stack, card)
  if err != nil {
    return domain.Game{}, err
  }

  if err := game.GetGemCard(index, payment); err != nil {
    return domain.Game{}, fmt.Errorf("get gem card: %w", err)
  }

  return srv.save(game)
}

// GetGolemCard buys golem from game id's golem lineup for player.
//
// Pre: player == the game's CurrentPlayer; golem matches (by value) an
// entry in the lineup; player can afford golem.Cost().
// Post: on success, the updated game is saved and returned.
func (srv *Service) GetGolemCard(id string, player int, golem domain.GolemCard) (domain.Game, error) {
  game, err := srv.loadForPlayer(id, player)
  if err != nil {
    return domain.Game{}, err
  }

  index, err := indexOfGolemCard(game.Golems.Stack, golem)
  if err != nil {
    return domain.Game{}, err
  }

  if err := game.BuyGolem(index); err != nil {
    return domain.Game{}, fmt.Errorf("buy golem: %w", err)
  }

  return srv.save(game)
}

// Rest ends player's turn, applying discards first if player holds more
// than 10 gems (see ports.GameService.Rest for discards' format).
//
// Pre: player == the game's CurrentPlayer; if player holds more than 10
// gems, discards must bring them to 10 or fewer.
// Post: on success, the updated game (CurrentPlayer advanced) is saved and
// returned.
func (srv *Service) Rest(id string, player int, discards string) (domain.Game, error) {
  game, err := srv.loadForPlayer(id, player)
  if err != nil {
    return domain.Game{}, err
  }

  if err := game.NextTurn(discards); err != nil {
    return domain.Game{}, fmt.Errorf("rest: %w", err)
  }

  return srv.save(game)
}

// loadForPlayer fetches game id and checks player is the one currently up.
//
// Pre: id names a game in the repository.
// Post: on success, the returned game's CurrentPlayer == player. Violating
// that (a stale or wrong player index from the caller) is a caller bug, not
// this service's — it always reports it rather than silently acting on
// whichever player is actually current.
func (srv *Service) loadForPlayer(id string, player int) (domain.Game, error) {
  game, err := srv.gamesRepository.Get(id)
  if err != nil {
    return domain.Game{}, fmt.Errorf("get game %q: %w", id, err)
  }

  if player != game.CurrentPlayer {
    return domain.Game{}, fmt.Errorf("player %d is not up (current player is %d)", player, game.CurrentPlayer)
  }

  return game, nil
}

// save persists game and returns it, wrapping any repository error.
func (srv *Service) save(game domain.Game) (domain.Game, error) {
  if err := srv.gamesRepository.Save(game); err != nil {
    return domain.Game{}, fmt.Errorf("save game %q: %w", game.ID, err)
  }
  return game, nil
}

// indexOfGemCard returns the index of the first card in stack equal (by
// value) to card, or an error if none matches.
func indexOfGemCard(stack []domain.GemCard, card domain.GemCard) (int, error) {
  for i, c := range stack {
    if c == card {
      return i, nil
    }
  }
  return -1, errors.New("no matching gem card in the lineup")
}

// indexOfGolemCard returns the index of the first card in stack equal (by
// value) to golem, or an error if none matches.
func indexOfGolemCard(stack []domain.GolemCard, golem domain.GolemCard) (int, error) {
  for i, c := range stack {
    if c == golem {
      return i, nil
    }
  }
  return -1, errors.New("no matching golem card in the lineup")
}

func GetUID() string {
  return uuid.NewString()
}
