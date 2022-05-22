package models

import "testing"

// Test that draw card decreases card count appropriately
func TestDrawCard(t *testing.T) {
	deck := NewGemDeck()
	start_deck_size := len(deck.stack)

	// DrawCard
	_, err := deck.DrawCard()
	if err != nil {
		t.Errorf(err.Error())
	}

	// ensure deck decreases in size
	if len(deck.stack) >= start_deck_size {
		t.Errorf("expected deck size to decrease: Start deck size was %v, current deck size is %v",
			start_deck_size, len(deck.stack))
	}
}
