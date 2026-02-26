package turns

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

func NewPlayCardView(game *domain.Game) TurnView {
	cardSelector := textinput.New()
	cardSelector.Placeholder = "Which Card would you like to play?"
	rv := TurnView{focusedIndex: 0, focused: false, actionText: "Play Card", action: game.GetGemCard, inputs: []textinput.Model{cardSelector}}
	return rv
}
