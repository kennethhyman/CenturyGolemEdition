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

func MarshallPlayer(player *models.Player) *pb.Player {
  var golems []*pb.GolemCard
  var hand []*pb.GemCard
  var discards []*pb.GemCard

  for _, card := range(player.Hand) {
    hand = append(hand, MarshallGemCard(&card))
  }

  for _, card := range(player.Golems) {
    golems = append(golems, MarshallGolemCard(&card))
  }

  for _, card := range(player.DiscardPile) {
    discards = append(discards, MarshallGemCard(&card))
  }

  return &pb.Player{
    GoldCoins: int32(player.GoldCoins),
    SilverCoins: int32(player.SilverCoins),
    Gems: MarshallGems(&player.Gems),
    Golems: golems,
    Hand: hand,
    DiscardPile: discards,
  }
}
