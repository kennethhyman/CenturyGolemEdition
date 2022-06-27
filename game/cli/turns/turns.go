package turns

import (
	tea "github.com/charmbracelet/bubbletea"
)

type GameAction string

const (
	PlayGemCard  GameAction = "PlayGemCard"
	GetGemCard   GameAction = "GetGemCard"
	Rest         GameAction = "Rest"
	GetGolemCard GameAction = "GetGolemCard"
)

func (ga GameAction) Init() tea.Cmd {
	return nil
}

func (ga GameAction) View() string {
	return string(ga)
}

func (ga GameAction) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return ga, nil
}
