package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

type Viewable interface {
	View() string
}

type Selector struct {
	focused  bool
	options  []Viewable
	selected Viewable
}

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) View() string {
	var b strings.Builder

	var str string
	for i, option := range s.options {
		str = fmt.Sprintf("%v: %v\t", i+1, option.View())
		if s.selected == option {
			b.WriteString(focusedStyle.Render(str))
		} else {
			b.WriteString(str)
		}
	}

	b.WriteString("\n")

	return b.String()
}

func (s Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// is the input a number?
		if i, err := strconv.Atoi(msg.String()); err == nil {
			// is that number in our list?
			if i > 0 && i <= len(s.options) {
				s.selected = s.options[i-1]
			}
		}
	}

	return s, nil
}
