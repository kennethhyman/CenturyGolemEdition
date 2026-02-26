package turns

import (
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
)

func NewRestView(game *domain.Game) TurnView {
	rv := TurnView{focusedIndex: 0, focused: false, actionText: "Rest", action: game.PickUpGemCards}
	return rv
}
