package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli/actions"
)

type model struct {
  game Game
}

var actionsChoices = []tea.Model{turns.PlayGemCard, turns.GetGemCard, turns.Rest, turns.GetGolemCard}

func InitialModel() model {
	return model {
    game: Game{},
	}
}

func (m model) Init() tea.Cmd {
	return nil //createGame
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

  case Game:
    m.game = msg
	}

	// pass the message along (and cast to view selector)
	return m, nil
}

func CardToString(card GemCard) string {
  return fmt.Sprintf("%vy%vg%vb%vp", card.In.Yellow, card.In.Green, card.In.Blue, card.In.Pink)
}

func (m model) View() string {
  return fmt.Sprintf("%v\n", m.game)
}
