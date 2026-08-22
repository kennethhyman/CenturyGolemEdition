package actions

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

// NewPlayCardView builds the TurnView for the PlayGemCard action: which
// card (by hand index) to play. The typed index is only checked for being
// a well-formed non-negative integer here — bounds-checking it against the
// player's actual hand length happens where the hand is known (cli.model,
// on SubmitMsg), not here.
func NewPlayCardView() TurnView {
	handIndex := textinput.New()
	handIndex.Placeholder = "hand index to play"

	return NewTurnView(PlayGemCard, "Play Card", []textinput.Model{handIndex}, validatePlayGemCard)
}

func validatePlayGemCard(values []string) error {
	if len(values) < 1 || values[0] == "" {
		return errValueRequired("hand index")
	}
	if _, err := strconv.Atoi(values[0]); err != nil {
		return fmt.Errorf("hand index %q is not a number", values[0])
	}
	return nil
}
