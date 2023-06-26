package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
  "fmt"
)

func MarshallGame(game *models.Game) *pb.CreateGameResponse {
  //var otherPlayers []*pb.Player
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
      Player: MarshallPlayer(game.Players[0]),
    },
  }
}
