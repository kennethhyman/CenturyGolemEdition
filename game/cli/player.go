package cli


import (
  "fmt"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
)

type Player struct {
  hand []GemCard
  golems []GolemCard
  discards []GemCard
  gems GemValue
  gold_coins int8
  silver_coins int8
}

func (p Player) CoinString() string {
  return fmt.Sprintf("gold: %v\t silver: %v", p.gold_coins, p.silver_coins)
}

func (p Player) String() string {
  // name : coins | golems 
  // gems [discard_count]
  // gem cards | 
  return fmt.Sprintf("%v | %v", p.CoinString(), p.golems)
}

func UnmarshallPlayer(p *pb.Player) Player {
  if p == nil {
    return Player{}
  }

  var hand []GemCard
  var golems []GolemCard
  var discards []GemCard

 // for _, gc := range(p.Hand) {
 //   hand = append(hand, UnmarshallGemCard(gc))
 // }

 // for _, golem := range(p.Golems) {
 //   golems = append(golems, UnmarshallGolemCard(golem))
 // }

 // for _, gc := range(p.DiscardPile) {
 //   discards = append(discards, UnmarshallGemCard(gc))
 // }
 fmt.Printf("%v\n", p)

  return Player{
    hand: hand,
    golems: golems,
    discards: discards,
    gems: GemValue {
      Yellow: int(p.Gems.Yellow),
      Green: int(p.Gems.Green),
      Blue: int(p.Gems.Blue),
      Pink: int(p.Gems.Pink),
    },
    gold_coins: int8(p.GoldCoins),
    silver_coins: int8(p.SilverCoins),
  }
}
