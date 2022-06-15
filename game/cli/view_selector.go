package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type ViewSelector struct {
	selector  Selector
	views     []tea.Model
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
	if !s.focused {
		return s, nil
	}

	if s.viewFocus {
		s.views[s.selected()], cmd = s.selectedView().Update(msg)
	} else {
		var selector tea.Model
		selector, cmd = s.selector.Update(msg)
		s.selector = selector.(Selector)
	}

	return s, cmd
}

func (s ViewSelector) selected() int {
	return s.selector.selected
}

func (s ViewSelector) selectedView() tea.Model {
	return s.views[s.selected()]
}
