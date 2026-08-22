package actions

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

// NewRestView builds the TurnView for the Rest action: which gems (if any)
// to discard to end the turn. Blank is valid whenever the player holds 10
// or fewer gems — that check (and the "must bring it to 10 or fewer")
// happens server-side, since only the server knows the player's current
// gem count.
func NewRestView() TurnView {
	discards := textinput.New()
	discards.Placeholder = "discard letters (blank if 10 or fewer gems)"

	return NewTurnView(Rest, "Rest", []textinput.Model{discards}, validateRest)
}

func validateRest(values []string) error {
	if len(values) < 1 || values[0] == "" {
		return nil
	}
	if !isGemLetters(values[0]) {
		return fmt.Errorf("discards must be gem letters (y/g/b/p), got %q", values[0])
	}
	return nil
}
