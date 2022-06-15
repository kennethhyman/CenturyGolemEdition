package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/kennethhyman/CenturyGolemEdition/models"
)

type model struct {
	game           Game
	actionSelector tea.Model
}

var actionsChoices = []Viewable{PlayGemCard, GetGemCard, Rest, GetGolemCard}

func InitialModel() model {
	turnSelector := Selector{options: actionsChoices, focused: true}
	views := []tea.Model{PlayGemCard, GetGemCard, Rest, GetGolemCard}
	return model{
		game:           *NewGame(2),
		actionSelector: ViewSelector{turnSelector, views, false, true},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// pass the message along
	m.actionSelector, _ = m.actionSelector.Update(msg)

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%v\n%v\n", m.game.String(), m.actionSelector.View())
}
