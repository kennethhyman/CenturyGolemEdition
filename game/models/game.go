package models

import (
	"errors"
	"fmt"
)

type Game struct {
	Players       []Player
	Golems        LineUp[GolemCard]
	GemCards      LineUp[GemCard]
	LooseGems     []GemValues
	SilverCoins   int
	GoldCoins     int
	LastPlayer    *Player
	CurrentPlayer int
}

func NewGame(players int) *Game {
	game := Game{
		Golems:        NewGolemLinup(),
		GemCards:      NewGemLineUp(),
		SilverCoins:   2 * players,
		GoldCoins:     2 * players,
		LastPlayer:    nil,
		CurrentPlayer: 0,
	}

	for i := 1; i <= players; i++ {
		game.Players = append(game.Players, NewPlayer(i))
	}

	return &game
}

// play gem card
func (g *Game) PlayGemCard(card GemCard) error {
	player := g.Players[g.CurrentPlayer]
	err := player.PlayGemCard(card)
	return err
}

// get gem card
func (g *Game) GetGemCard(index int) error {
	player := g.Players[g.CurrentPlayer]
	var cost string
	if index > 0 {
		fmt.Scanln(&cost)

		if len(cost) != index {
			return errors.New("You must pay one gem per card you bypass")
		}
	}

	gemCost, err := parseGemInput(cost)

	if err != nil {
		return err
	}

	// Remove gems first
	err = player.Gems.remove(gemCost)
	if err != nil {
		return err
	}
	// Draw card
	card, err := g.GemCards.Draw(index)
	if err != nil {
		return err
	}
	player.AddCard(card)
	player.Gems.add(g.LooseGems[index-1])
	g.LooseGems[index-1] = GemValues{}

	// Add gems to piles
	for i, gem := range cost {
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
	}
	return err
}

// pick up cards
func (g *Game) PickUpGemCards() {
	player := g.Players[g.CurrentPlayer]
	player.PickupCards()
}

// buy golems
func (g *Game) BuyGolem(index int) error {
	player := g.Players[g.CurrentPlayer]
	// check index in bounds
	if index >= len(g.Golems.stack) || index < 1 {
		return errors.New("There is no golem at that index")
	}

	// check player has the gems
	golem := g.Golems.stack[index-1]
	if !player.CanAfford(golem) {
		return errors.New("You cannot afford this card")
	}
	// Remove gems and add card
	player.Gems.remove(golem.Cost())
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
