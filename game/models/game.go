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
	CurrentPlayer int
	GolemLimit    int
}

func NewGame(players int) *Game {
	var golemLimit int
	if players > 3 {
		golemLimit = 5
	} else {
		golemLimit = 6
	}
	game := Game{
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
func (g *Game) PlayGemCard() error {
	player := g.Players[g.CurrentPlayer]

	fmt.Printf("which card would you like to play?\n%v\n", player.Hand)
	var index int
	fmt.Scan(&index)

	err := player.PlayGemCard(player.Hand[index])
	fmt.Printf("%v\n", err)

	g.Players[g.CurrentPlayer] = player
	return err
}

// get gem card
func (g *Game) GetGemCard() error {
	player := g.Players[g.CurrentPlayer]
	var index int
	var err error

	fmt.Printf("Which gemcard would you like to pick up?\n%v\n", g.GemCards.String())
	for _, err := fmt.Scan(&index); err != nil; {
		fmt.Printf("invalid input...")
		fmt.Printf("Which gemcard would you like to pick up?\n%v\n", g.GemCards.String())
	}

	var cost string
	if index > 0 {
		fmt.Printf("Which gems would you like to place on the previous cards? (select %v gem)\n", index)
		fmt.Scanln(&cost)

		if len(cost) != index {
			return errors.New("You must pay one gem per card you bypass")
		}
		gemCost, err := parseGemInput(cost)

		if err != nil {
			return err
		}

		err = player.Gems.remove(gemCost)
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	}

	// Draw card
	card, err := g.GemCards.Draw(index)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	player.AddCard(card)
	player.Gems.add(g.LooseGems[index])
	g.LooseGems[index] = GemValues{}
	fmt.Printf("Hand addr: %p\n", &player.Hand)
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

	g.Players[g.CurrentPlayer] = player
	return err
}

// pick up cards
func (g *Game) PickUpGemCards() {
	player := g.Players[g.CurrentPlayer]
	player.PickupCards()
	g.Players[g.CurrentPlayer] = player
}

// buy golems
func (g *Game) BuyGolem() error {
	player := g.Players[g.CurrentPlayer]
	var index int

	fmt.Printf("Which golem would you like to buy?\n%v\n", g.Golems)
	fmt.Scan(&index)
	// check index in bounds
	if index >= len(g.Golems.stack) || index < 0 {
		return errors.New("There is no golem at that index")
	}

	// check player has the gems
	golem := g.Golems.stack[index]
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

func (g *Game) NextTurn() {
	current_player := g.Players[g.CurrentPlayer]

	// Trigger gem discards
	for current_player.Gems.count() > 10 {
		var user_input string
		fmt.Printf("You have too many gems, you must discard down to 10\n%v\n", current_player.Gems)
		fmt.Scan(&user_input)
		discards, err := parseGemInput(user_input)
		if err == nil {
			current_player.Gems.remove(discards)
		}
	}
	// set next Players

	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
}
