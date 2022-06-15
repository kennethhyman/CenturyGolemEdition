package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var unfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
var focusedStyle = lipgloss.NewStyle().Bold(true)

type Viewable interface {
	View() string
}

type Selector struct {
	focused  bool
	options  []Viewable
	selected int
}

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) selectedOption() Viewable {
	return s.options[s.selected]
}

func (s Selector) View() string {
	var b strings.Builder
	var style = unfocusedStyle

	if s.focused {
		style = focusedStyle
	}

	var str string
	for i, option := range s.options {
		str = style.Render(fmt.Sprintf("%v: %v\t", i+1, option.View()))
		if s.selectedOption() == option {
			b.WriteString(selectedStyle.Render(str))
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
				s.selected = i - 1
			}
		}
	}

	return s, nil
}