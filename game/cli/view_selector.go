package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli/turns"
	"strings"
)

type ViewSelector struct {
	selector  Selector
	views     []turns.TurnView
	viewFocus bool
	focused   bool
}

func (s ViewSelector) Init() tea.Cmd {
	return nil
}

func (s ViewSelector) View() string {
	var b strings.Builder

	if s.viewFocus {
		b.WriteString("placeholder")
	} else {

	}
	// displays the selector and then the corresponding tea model
	// focus is displayed on the selected element
	return fmt.Sprintf("%v\n%v", s.selector.View(), s.selectedView().View())
}

func (s ViewSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "tab" {
			s.selector.focused = !s.selector.focused
			s.views[s.selected()].Focus()
			s.viewFocus = !s.viewFocus

			return s, nil
		}
	}

	if s.viewFocus {
		model, cmd = s.selectedView().Update(msg)
		s.views[s.selected()] = model.(turns.TurnView)
	} else {
		var selector tea.Model
		selector, cmd = s.selector.Update(msg)
		s.selector = selector.(Selector)
	}

	return s, cmd
}

func (s ViewSelector) Reset() ViewSelector {
	// reset focus, reset views
	s.selector.Reset()

	for i := 0; i < len(s.views); i++ {
		s.views[i] = s.views[i].Reset()
	}

	s.viewFocus = false

	return s
}

func (s ViewSelector) selected() int {
	return s.selector.selected
}

func (s ViewSelector) selectedView() tea.Model {
	return s.views[s.selected()]
}
