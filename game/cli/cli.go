package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli/turns"
	. "github.com/kennethhyman/CenturyGolemEdition/models"
)

type model struct {
	actionSelector ViewSelector
}

var game = NewGame(2)

var restView = turns.NewRestView(game)
var actionsChoices = []Viewable{turns.PlayGemCard, turns.GetGemCard, turns.Rest, turns.GetGolemCard}

func InitialModel() model {
	turnSelector := Selector{options: actionsChoices, focused: true}
	views := []turns.TurnView{turns.NewPlayCardView(game), turns.NewGemCardView(game), turns.NewRestView(game), turns.NewGolemCardView(game)}
	return model{
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

	// pass the message along (and cast to view selector)
	selector, _ := m.actionSelector.Update(msg)
	m.actionSelector = selector.(ViewSelector)
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%v\n%v\n", game.String(), m.actionSelector.View())
}
