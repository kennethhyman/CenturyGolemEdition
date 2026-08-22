package actions

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

// NewGolemCardView builds the TurnView for the GetGolemCard (buy golem)
// action: which lineup golem to buy. Only checked for being a well-formed
// non-negative integer here — bounds-checking against the current golem
// lineup, and affordability, happen where that state is known (cli.model
// on SubmitMsg, and the server).
func NewGolemCardView() TurnView {
	index := textinput.New()
	index.Placeholder = "lineup index of golem to buy"

	return NewTurnView(GetGolemCard, "Buy Golem", []textinput.Model{index}, validateGetGolemCard)
}

func validateGetGolemCard(values []string) error {
	if len(values) < 1 || values[0] == "" {
		return errValueRequired("lineup index")
	}
	if _, err := strconv.Atoi(values[0]); err != nil {
		return fmt.Errorf("lineup index %q is not a number", values[0])
	}
	return nil
}
