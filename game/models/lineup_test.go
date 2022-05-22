package models

import "testing"
import "fmt"

func TestLineup(t *testing.T) {
	gem_deck := NewGemDeck()
	gem_deck.Shuffle()
	start_deck_size := len(gem_deck.stack)
	lineup := NewGemLineUp(gem_deck)

	fmt.Printf("%v\n", lineup)
	// test drawing
	if len(gem_deck.stack) >= start_deck_size {
		t.Errorf("expected deck size to decrease: Start deck size was %v, current deck size is %v",
			start_deck_size, len(gem_deck.stack))
	}
}
