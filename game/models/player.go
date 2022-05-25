package models

import "errors"

var starting_gems = []GemValues{
	GemValues{
		Yellow: 3,
	},
	GemValues{
		Yellow: 4,
	},
	GemValues{
		Yellow: 4,
	},
	GemValues{
		Yellow: 3,
		Green:  1,
	},
	GemValues{
		Yellow: 3,
		Green:  1,
	},
}

// 2nd player gets 4 yellow crystals
// 3rd player gets 4 yellow crystals
// 4th player gets 3 yellow crystals and 1 green crystal
// 5th player gets 3 yellow crystals and 1 green crystal

type Player struct {
	Hand        []GemCard
	DiscardPile []GemCard
	Gems        GemValues
	GoldCoins   int
	SilverCoins int
	//Golems []GolemCard
}

func NewPlayer(turnOrder int) Player {
	return Player{
		Hand: GetStartingDeck(),
		Gems: starting_gems[turnOrder],
	}
}

func (p *Player) PlayCard(card GemCard) error {
	// Check that the player has the gems
	if !p.HasInputs(card) {
		return errors.New("Player does not have the gems to play this card")
	}
	// Check that player has the card
	if !p.HasCardAvailable(card) {
		return errors.New("Player does not have the card available for play")
	}

	// play / remove the card

	// add / remove the gems
	return nil
}

func (p Player) HasInputs(card GemCard) bool {
	inputs := card.Inputs
	gems := p.Gems
	return (gems.Yellow >= inputs.Yellow && gems.Green >= inputs.Green && gems.Blue >= inputs.Blue && gems.Pink >= inputs.Pink)
}

func (p Player) HasCardAvailable(card GemCard) bool {
	// iterate through deck
	for _, ownedCard := range p.Hand {
		if ownedCard == card {
			return true
		}
	}

	return false
}

func (p *Player) GetCard(index int, gems string) error {

	return nil
}
