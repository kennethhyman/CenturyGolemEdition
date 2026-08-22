package actions

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

// NewGemCardView builds the TurnView for the GetGemCard action: which
// lineup card to pick up, and the payment (one gem letter per card
// bypassed) for cards left behind. Only format is checked here (index is
// numeric, payment length matches the typed index, payment is gem
// letters) — whether the index is actually in range for the current
// lineup, and whether the player holds those gems, is checked where the
// lineup/hand are known (cli.model, on SubmitMsg, and ultimately the
// server).
func NewGemCardView() TurnView {
	index := textinput.New()
	index.Placeholder = "lineup index"

	payment := textinput.New()
	payment.Placeholder = "payment letters (y/g/b/p), one per card bypassed"

	return NewTurnView(GetGemCard, "Get Gem Card", []textinput.Model{index, payment}, validateGetGemCard)
}

func validateGetGemCard(values []string) error {
	if len(values) < 1 || values[0] == "" {
		return errValueRequired("lineup index")
	}
	index, err := strconv.Atoi(values[0])
	if err != nil {
		return fmt.Errorf("lineup index %q is not a number", values[0])
	}

	payment := ""
	if len(values) > 1 {
		payment = values[1]
	}
	if len(payment) != index {
		return fmt.Errorf("payment must have exactly %d letter(s) (one per card bypassed), got %q", index, payment)
	}
	if !isGemLetters(payment) {
		return fmt.Errorf("payment must be gem letters (y/g/b/p), got %q", payment)
	}
	return nil
}
