package turns

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

func NewGolemCardView(game *domain.Game) TurnView {
	cardSelector := textinput.New()
	cardSelector.Placeholder = "Desired card"
	rv := TurnView{focusedIndex: 0, focused: false, actionText: "Get Golem Card", action: game.GetGemCard, inputs: []textinput.Model{cardSelector}}
	return rv
}
