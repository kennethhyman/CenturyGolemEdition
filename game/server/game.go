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
  var lineup []*pb.GemCard

  for _, card := range(game.GemCards.Stack) {
    marshalled_card := MarshallGemCard(&card)
    lineup = append(lineup, marshalled_card)
    fmt.Printf("%v\n", marshalled_card)
  }

  return &pb.CreateGameResponse {
    GameState: &pb.GameState{
      GemLineup: lineup,
    },
  }
}
