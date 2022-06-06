package models

import (
	"errors"
	"fmt"
	"strings"
)

var starting_gems = []GemValues{
	{
		Yellow: 3,
	},
	{
		Yellow: 4,
	},
	{
		Yellow: 4,
	},
	{
		Yellow: 3,
		Green:  1,
	},
	{
		Yellow: 3,
		Green:  1,
	},
}

type Player struct {
	Hand        []GemCard
	DiscardPile []GemCard
	Gems        GemValues
	GoldCoins   int
	SilverCoins int
	Golems      []GolemCard
}

func NewPlayer(turnOrder int) Player {
	return Player{
		Hand: GetStartingDeck(),
		Gems: starting_gems[turnOrder-1],
	}
}

func (p *Player) PlayGemCard(card GemCard) error {
	// Check that player has the card
	available, index := p.HasCardAvailable(card)
	if !available {
		return errors.New("Player does not have the card available for play")
	}

	// add / remove the gems
	err := p.Gems.remove(card.Inputs)
	if err != nil {
		return err
	}
	p.Gems.add(card.Outputs)
	fmt.Printf("New gem count: %v\n", p.Gems)

	// play / remove the card
	p.DiscardPile = append(p.DiscardPile, card)
	p.Hand = append(p.Hand[:index], p.Hand[index+1:]...)

	return nil
}

func (p *Player) PlayUpgradeCard(card GemCard, inputs GemValues, outputs GemValues) error {
	// Check that player has the card
	available, index := p.HasCardAvailable(card)
	if !available {
		return errors.New("Player does not have the card available for play")
	}

	// Check legality of the upgrade
	upgradeValue := outputs.effectiveValue() - inputs.effectiveValue()
	legalUpgrade := upgradeValue >= 0 && upgradeValue <= card.Upgrades &&
		outputs.strictlyGreater(inputs) && inputs.count() == outputs.count() &&
		inputs.Pink == 0

	if !legalUpgrade {
		return errors.New("Invalid input / output for this card")
	}

	// add / remove the gems
	err := p.Gems.remove(inputs)
	if err != nil {
		return err
	}
	p.Gems.add(outputs)

	// play / remove the card
	p.DiscardPile = append(p.DiscardPile, card)
	p.Hand = append(p.Hand[:index], p.Hand[index+1:]...)
	return nil
}

func (p Player) HasCardAvailable(card GemCard) (bool, int) {
	// iterate through deck
	for i, ownedCard := range p.Hand {
		if ownedCard == card {
			return true, i
		}
	}

	return false, -1
}

func (p *Player) AddCard(card GemCard) {
	fmt.Printf("Adding card %v\n", card)
	p.Hand = append(p.Hand, card)
	fmt.Printf("Hand addr: %p\n", &p.Hand)
}

func (p *Player) PickupCards() {
	p.Hand = append(p.Hand, p.DiscardPile...)
	p.DiscardPile = []GemCard{}
}

func (p *Player) AddGolemCard(card GolemCard) {
	p.Golems = append(p.Golems, card)
}

func (p Player) CanAfford(card card) bool {
	player_gems := p.Gems
	cost := card.Cost()

	return player_gems.Yellow >= cost.Yellow || player_gems.Green >= cost.Green ||
		player_gems.Blue >= cost.Blue || player_gems.Pink > cost.Pink
}

func (p Player) CardString() string {
	var hand []string
	for i, card := range p.Hand {
		hand = append(hand, fmt.Sprintf("%v : %v", i, card))
	}
	return fmt.Sprintf("%v\tdiscard: %v", strings.Join(hand[:], " "), len(p.DiscardPile))
}

func (p Player) GolemString() string {
	var golems []string
	for _, card := range p.Golems {
		golems = append(golems, card.String())
	}

	return fmt.Sprintf("Golems: %v", strings.Join(golems[:], " "))
}

func (p Player) String() string {
	hand := p.CardString()
	golems := p.GolemString()
	coins := fmt.Sprintf("gold: %v\tsilver: %v", p.GoldCoins, p.SilverCoins)

	return fmt.Sprintf("%v\n%v\n%v\n%v\n", coins, golems, hand, p.Gems.String())
}
