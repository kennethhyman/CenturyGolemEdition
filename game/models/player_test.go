package models

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func TestNewPlayer(t *testing.T) {
	player_1 := NewPlayer(0)
	player_2 := NewPlayer(1)
	var emptyStack []GemCard

	assert.NotEqual(t, player_1.Gems, player_2.Gems,
		"players should start with different gems counts based on turn order")

	assert.Equal(t, player_1.Hand, player_2.Hand, "players should start with the same cards")
	assert.NotEqual(t, len(player_1.Hand), 0, "Players should start with cards")
	assert.Equal(t, player_1.DiscardPile, emptyStack, "A player should start with an empty discard pile")
}

func TestPlayerPlayGemCard(t *testing.T) {
	player_1 := NewPlayer(0)
	player_gems := player_1.Gems
	yellowCard := GemCard{
		Outputs: GemValues{
			Yellow: 2,
		},
	}

	err := player_1.PlayGemCard(yellowCard)
	yellow_to_green := GemCard{
		Inputs: GemValues{
			Yellow: 2,
		},
		Outputs: GemValues{
			Green: 2,
		},
	}
	player_1.AddCard(yellow_to_green)

	assert.Equal(t, err, nil, ">2y should be a legal first move")
	assert.NotEqual(t, player_1.Gems, player_gems, "playing the first card adds gems")
	available, _ := player_1.HasCardAvailable(yellowCard)
	assert.False(t, available, "Playing a card should remove it from hand")

	player_gems = player_1.Gems
	err = player_1.PlayGemCard(GemCard{
		Outputs: GemValues{
			Pink: 2,
		},
	})

	assert.NotEqual(t, nil, err, "Playing a card you don't have should throw an error")
	assert.Equal(t, player_gems, player_1.Gems, "If you cant play a card, gems don't change")

	err = player_1.PlayGemCard(yellow_to_green)
	assert.Equal(t, nil, err, "Player should legally be allowed to play gems")
	expected_gems := GemValues{
		Yellow: 3,
		Green:  2,
	}
	assert.Equal(t, expected_gems, player_1.Gems, fmt.Sprintf("The player should end with %v", expected_gems))
}

func TestPlayerPlayUpgradeCard(t *testing.T) {
	player_1 := NewPlayer(0)
	starter_gems := player_1.Gems
	upgradeCard := GemCard{
		Upgrades: 2,
	}
	in := GemValues{Yellow: 2}
	out := GemValues{Green: 2}
	err := player_1.PlayUpgradeCard(upgradeCard, in, out)

	assert.Equal(t, nil, err, fmt.Sprintf("%v should be a legal first move", upgradeCard))
	assert.NotEqual(t, player_1.Gems, starter_gems, "playing the first card adds gems")
}
