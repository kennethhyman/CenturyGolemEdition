package models

import "testing"
import "github.com/stretchr/testify/assert"

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

func TestPlayerPlayCard(t *testing.T) {
	player_1 := NewPlayer(0)
	starter_gems := player_1.Gems
	yellowCard := GemCard{
		Outputs: GemValues{
			Yellow: 2,
		},
	}

	err := player_1.PlayCard(yellowCard)

	assert.Equal(t, err, nil, ">2y should be a legal first move")
	assert.NotEqual(t, player_1.Gems, starter_gems, "playing the first card adds gems")
	assert.False(t, player_1.HasCardAvailable(yellowCard), "Playing a card should remove it from hand")

	starter_gems = player_1.Gems
	err = player_1.PlayCard(GemCard{
		Outputs: GemValues{
			Pink: 2,
		},
	})

	assert.NotEqual(t, err, nil, "Playing a card you don't have should throw an error")
	assert.Equal(t, starter_gems, player_1.Gems, "If you cant play a card, gems don't change")
}
