package cli

type GameAction string

const (
	PlayGemCard  GameAction = "PlayGemCard"
	GetGemCard   GameAction = "GetGemCard"
	Rest         GameAction = "Rest"
	GetGolemCard GameAction = "GetGolemCard"
)

func (ga GameAction) View() string {
	return string(ga)
}
