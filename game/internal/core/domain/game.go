package domain

import (
	"errors"
	"fmt"
  pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"

)

type Game struct {
  ID            string
	Players       []*Player
	Golems        LineUp[GolemCard]
	GemCards      LineUp[GemCard]
	LooseGems     []GemValues
	SilverCoins   int
	GoldCoins     int
	CurrentPlayer int
	GolemLimit    int
}

func NewGame(id string, players int) *Game {
	var golemLimit int
	if players > 3 {
		golemLimit = 5
	} else {
		golemLimit = 6
	}
	game := Game{
    ID:            id,
		Golems:        NewGolemLinup(),
		GemCards:      NewGemLineUp(),
		LooseGems:     make([]GemValues, 6),
		SilverCoins:   2 * players,
		GoldCoins:     2 * players,
		CurrentPlayer: 0,
		GolemLimit:    golemLimit,
	}

	for i := 1; i <= players; i++ {
		game.Players = append(game.Players, NewPlayer(i))
	}

	return &game
}

// play gem card
func (g *Game) PlayGemCard(gem_card GemCard, discards GemValues) error {
	player := g.Players[g.CurrentPlayer]
  

	err := player.PlayGemCard(gem_card, discards)

	g.Players[g.CurrentPlayer] = player
	return err
}

// GetGemCard picks up the gem card at index from the face-up lineup for the
// current player, paying one gem (per parsed character of payment) onto
// each bypassed lineup slot.
//
// This is the mutation half of the old prompt-and-mutate GetGemCard: the
// prompting ("which card", "which gems to pay") is a CLI concern and lives
// client-side (see cli.PromptGetGemCard); this method only applies an
// already-decided index/payment, so it can run inside a server process that
// has no terminal to prompt on.
//
// Pre: index is in range for g.GemCards.Stack; if index > 0, payment has
// exactly `index` gem-letter characters (y/g/b/p) and the current player
// holds those gems.
// Post: on success, the drawn card joins the current player's hand, the
// loose gems sitting on index move to the player, and payment gems land one
// per bypassed pile. On error, no state is mutated for cases caught before
// the draw; errors from GemCards.Draw itself may still have consumed cards
// per LineUp's own contract.
func (g *Game) GetGemCard(index int, payment string) error {
	player := g.Players[g.CurrentPlayer]

	if index > 0 {
		if len(payment) != index {
			return errors.New("You must pay one gem per card you bypass")
		}
		gemCost, err := parseGemInput(payment)
		if err != nil {
			return err
		}

		if err := player.Gems.remove(gemCost); err != nil {
			return err
		}
	}

	// Draw card
	card, err := g.GemCards.Draw(index)
	if err != nil {
		return err
	}
	player.AddCard(card)
	player.Gems.add(g.LooseGems[index])
	g.LooseGems[index] = GemValues{}

	// Add gems to piles
	for i, gem := range payment {
		pile := g.LooseGems[i]

		switch string(gem) {
		case "y", "Y":
			pile.Yellow += 1
		case "g", "G":
			pile.Green += 1
		case "b", "B":
			pile.Blue += 1
		case "p", "P":
			pile.Pink += 1
		}

		g.LooseGems[i] = pile
	}

	g.Players[g.CurrentPlayer] = player
	return nil
}

// pick up cards
func (g *Game) PickUpGemCards() error {
	player := g.Players[g.CurrentPlayer]
	player.PickupCards()
	g.Players[g.CurrentPlayer] = player

	return nil
}

// BuyGolem buys the golem card at index from the face-up lineup for the
// current player.
//
// Mutation half of the old prompt-and-mutate BuyGolem: prompting ("which
// golem") is a CLI concern (see cli.PromptBuyGolem); this method only
// applies an already-decided index, so it can run server-side.
//
// Pre: 0 <= index < len(g.Golems.Stack).
// Post: on success, the golem joins the current player's Golems, its cost
// is removed from the player's gems, and a bonus coin is granted per the
// index-based coin rule below. On error (out of range, or the player can't
// afford it), no state is mutated.
func (g *Game) BuyGolem(index int) error {
	player := g.Players[g.CurrentPlayer]

	// check index in bounds
	if index >= len(g.Golems.Stack) || index < 0 {
		return errors.New("There is no golem at that index")
	}

	// check player has the gems
	golem := g.Golems.Stack[index]
	if !player.CanAfford(golem) {
		return errors.New("You cannot afford this card")
	}
	// Remove gems and add card. CanAfford already confirmed the player holds
	// golem.Cost(), so remove() failing here would mean CanAfford and remove
	// disagree on affordability — a design bug, not a caller error.
	if err := player.Gems.remove(golem.Cost()); err != nil {
		return fmt.Errorf("CanAfford said yes but remove failed: %w", err)
	}
	player.AddGolemCard(golem)

	// Add coins if applicable
	switch index {
	case 1:
		if g.GoldCoins > 0 {
			player.GoldCoins += 1
			g.GoldCoins -= 1
		}
	case 2:
		if g.SilverCoins > 0 {
			player.SilverCoins += 1
			g.SilverCoins -= 1
		}
	}

	g.Players[g.CurrentPlayer] = player
	return nil
}

func (g Game) String() string {
	// print coins available
	available_coins := fmt.Sprintf("Remaining Coins | gold: %v\tsilver: %v", g.GoldCoins, g.SilverCoins)
	golem_string := g.Golems.String()
	gemcard_string := g.GemCards.String()
	player_string := g.Players[g.CurrentPlayer].String()

	return fmt.Sprintf("%v\n%v\n%v\n\n%v\n",
		available_coins, golem_string, gemcard_string, player_string)
}

func (g Game) Finished() bool {
	var golem_limit_reached = false

	for _, player := range g.Players {
		if len(player.Golems) >= g.GolemLimit {
			golem_limit_reached = true
			break
		}
	}

	// game is over when first player is up and someone has X golems
	return (g.CurrentPlayer == 0) && golem_limit_reached
}

// NextTurn ends the current player's turn and advances to the next one.
//
// If the current player holds more than 10 gems, discards (one gem letter
// y/g/b/p per character) must bring them down to 10 or fewer — this is the
// mutation half of the old prompt-in-a-loop NextTurn; deciding what to
// discard is a CLI concern (see cli.PromptRest), which loops client-side
// until it has a valid answer and then calls this once. Pass "" when the
// player already holds 10 or fewer gems.
//
// Pre: if the current player holds more than 10 gems, discards must be
// non-empty and, once removed, bring their count to 10 or fewer.
// Post: on success, discards are removed from the current player's gems and
// g.CurrentPlayer advances (wrapping) to the next player. On error, no
// state is mutated and the turn does not advance.
func (g *Game) NextTurn(discards string) error {
	current_player := g.Players[g.CurrentPlayer]

	if current_player.Gems.count() > 10 {
		discardGems, err := parseGemInput(discards)
		if err != nil {
			return err
		}
		if err := current_player.Gems.remove(discardGems); err != nil {
			return err
		}
		if current_player.Gems.count() > 10 {
			return errors.New("You have too many gems, you must discard down to 10")
		}
	}

	g.Players[g.CurrentPlayer] = current_player
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	return nil
}

// MarshallGameState builds the wire GameState for game from the current
// player's point of view: Player is the hand/gems/coins of whoever's turn
// it is, and CurrentPlayer is that player's index — callers must echo it
// back as the `player` field on their next request (PlayGemCard/GetCard/
// BuyGolem/Rest all validate the caller is actually up).
func MarshallGameState(game *Game) *pb.GameState {
  var gem_lineup []*pb.GemCard
  var golem_lineup []*pb.GolemCard

  for _, card := range game.GemCards.Stack {
    gem_lineup = append(gem_lineup, MarshallGemCard(&card))
  }

  for _, card := range game.Golems.Stack {
    golem_lineup = append(golem_lineup, MarshallGolemCard(&card))
  }

  return &pb.GameState{
    GameId: game.ID,
    CurrentPlayer: int32(game.CurrentPlayer),
    GemLineup: gem_lineup,
    GemDeckSize: int32(game.GemCards.Remaining()),
    GolemLineup: golem_lineup,
    GolemDeckSize: int32(game.Golems.Remaining()),
    GoldCoins: int32(game.GoldCoins),
    SilverCoins: int32(game.SilverCoins),
    Player: MarshallPlayer(game.Players[game.CurrentPlayer]),
  }
}

// MarshallGame wraps MarshallGameState in the NewGame RPC's response shape.
func MarshallGame(game *Game) *pb.CreateGameResponse {
  return &pb.CreateGameResponse{GameState: MarshallGameState(game)}
}
