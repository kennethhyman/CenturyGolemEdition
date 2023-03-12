package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
)

func UnmarshallPlayer(player *pb.Player) *models.Player {
  var hand = []models.GemCard{}
  var discard = []models.GemCard{}
  var golems = []models.GolemCard{}

  for _, g := range(player.Hand) {
    hand = append(hand, *UnmarshallGemCard(g))
  }

  for _, g := range(player.DiscardPile) {
    discard = append(discard, *UnmarshallGemCard(g))
  }

  for _, g := range(player.Golems) {
    golems = append(golems, *UnmarshallGolemCard(g))
  }

  return &models.Player{
    Hand: hand,
    DiscardPile: discard,
    Gems: *UnmarshallGems(player.Gems),
    GoldCoins: int(player.GoldCoins),
    SilverCoins: int(player.SilverCoins),
    Golems: golems,
  }
}
