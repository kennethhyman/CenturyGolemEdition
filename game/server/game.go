package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
  "fmt"
)

func UnmarshallGame(game *pb.GameState) *models.Game {
  var players = []*models.Player{}
  var gem_cards = []models.GemCard{}
  var golems = []models.GolemCard{}

  for _, player := range(game.Players) {
    players = append(players, UnmarshallPlayer(player))
  }

  for _, card := range(game.GemLineup) {
    gem_cards = append(gem_cards, *UnmarshallGemCard(card))
  }

  for _, golem := range(game.GolemLineup) {
    golems = append(golems, *UnmarshallGolemCard(golem))
  }

  return &models.Game{
    Players: players,
    GolemLimit: 7,
    CurrentPlayer: 1,
    SilverCoins: 8,
    GoldCoins: 8,
  }
}

func MarshallGame(game *models.Game) *pb.CreateGameResponse {
  var gem_lineup []*pb.GemCard
  var golem_lineup []*pb.GolemCard

  for _, card := range(game.GemCards.Stack) {
    marshalled_card := MarshallGemCard(&card)
    gem_lineup = append(gem_lineup, marshalled_card)
    fmt.Printf("%v\n", marshalled_card)
  }

  for _, card := range(game.Golems.Stack) {
    marshalled_card := MarshallGolemCard(&card)
    golem_lineup = append(golem_lineup, marshalled_card)
    fmt.Printf("%v\n", marshalled_card)
  }
  fmt.Printf("GemCard length: %v", len(game.GemCards.Stack))

  return &pb.CreateGameResponse {
    GameState: &pb.GameState{
      GemLineup: gem_lineup,
      GemDeckSize: int32(game.GemCards.Remaining()),
      GolemLineup: golem_lineup,
      GolemDeckSize: int32(game.Golems.Remaining()),
      GoldCoins: int32(game.GoldCoins),
      SilverCoins: int32(game.SilverCoins),
    },
  }
}
