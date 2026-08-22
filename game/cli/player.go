package cli

import (
	"fmt"

	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

// Player is the CLI's own display/view model for a player, unmarshalled
// from the pb.Player embedded in the last GameState the server sent —
// never mutated locally; every field here reflects what the server last
// reported, not local guesswork.
type Player struct {
	hand         []GemCard
	golems       []GolemCard
	discards     []GemCard
	gems         GemValue
	gold_coins   int
	silver_coins int
}

func (p Player) CoinString() string {
	return fmt.Sprintf("gold: %v\t silver: %v", p.gold_coins, p.silver_coins)
}

func (p Player) String() string {
	// name : coins | golems
	// gems [discard_count]
	// gem cards |
	return fmt.Sprintf("%v | %v\nGems: %v\nHand: %v [discard: %v]",
		p.CoinString(), p.golems, p.gems, p.hand, len(p.discards))
}

// UnmarshallPlayer is MarshallPlayer's inverse on the domain side
// (domain.MarshallPlayer): converts the wire Player into the CLI's display
// type. A nil p (e.g. a not-yet-populated GameState) yields the zero
// Player rather than panicking.
func UnmarshallPlayer(p *pb.Player) Player {
	if p == nil {
		return Player{}
	}

	var hand []GemCard
	var golems []GolemCard
	var discards []GemCard

	for _, gc := range p.GetHand() {
		hand = append(hand, UnmarshallGemCard(gc))
	}

	for _, golem := range p.GetGolems() {
		golems = append(golems, UnmarshallGolemCard(golem))
	}

	for _, gc := range p.GetDiscardPile() {
		discards = append(discards, UnmarshallGemCard(gc))
	}

	return Player{
		hand:     hand,
		golems:   golems,
		discards: discards,
		gems: GemValue{
			Yellow: int(p.GetGems().GetYellow()),
			Green:  int(p.GetGems().GetGreen()),
			Blue:   int(p.GetGems().GetBlue()),
			Pink:   int(p.GetGems().GetPink()),
		},
		gold_coins:   int(p.GetGoldCoins()),
		silver_coins: int(p.GetSilverCoins()),
	}
}
