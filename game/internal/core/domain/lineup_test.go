package domain

import "testing"
import "fmt"
import "github.com/stretchr/testify/assert"

func TestLineup(t *testing.T) {
	gem_deck := NewGemDeck()
	start_deck_size := len(gem_deck.stack)
	lineup := NewGemLineUp()
	fmt.Printf("%v\n", lineup)

	// test intial deck setup
	assert.Equal(t, start_deck_size-6, lineup.Remaining(), "Creating a lineup should draw 6 cards from the deck")
	assert.Equal(t, len(lineup.Stack), 6, "New lineups should be of length 6")
}

func TestLineUpDraw(t *testing.T) {
	lineup := NewGemLineUp()
	fmt.Printf("%v\n", lineup)

	first_card := lineup.Stack[0]
	last_card := lineup.Stack[len(lineup.Stack)-1]

	fmt.Println("Draw the first card")
	drawn_card, _ := lineup.Draw(0)
	fmt.Printf("%v\n", lineup)
	// test intial deck setup
	assert.Equal(t, first_card, drawn_card, "The card returned should match the first listed card")
	assert.NotEqual(t, drawn_card, lineup.Stack[0], "Cards should shift down one slot when a card is drawn")
	assert.Equal(t, last_card, lineup.Stack[len(lineup.Stack)-2])
}
