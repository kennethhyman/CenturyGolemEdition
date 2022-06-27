package turns

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type TurnView struct {
	inputs       []textinput.Model
	focusedIndex int
	focused      bool
	actionText   string
	error        error
	action       func() error
}

func (v TurnView) Init() tea.Cmd {
	return nil
}

func (v TurnView) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(v.inputs))
	if v.focused && v.focusedIndex < len(v.inputs) {
		v.inputs[v.focusedIndex].Focus()
	}

	for i := range v.inputs {
		v.inputs[i], cmds[i] = v.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (v TurnView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// skip when not focused
	if !v.focused {
		return v, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if v.focusedIndex == len(v.inputs) {
				// send game action message
				v.action()

				return v, nil
			}
			v.focusedIndex += 1
		case "down":
			v.focusedIndex += 1

		case "up":
			v.focusedIndex -= 1
		}
	}

	if v.focusedIndex > len(v.inputs) {
		v.focusedIndex = 0
	}
	if v.focusedIndex < 0 {
		v.focusedIndex = len(v.inputs)
	}

	cmds := make([]tea.Cmd, len(v.inputs))
	for i := 0; i < len(v.inputs); i++ {
		if i == v.focusedIndex {
			cmds[i] = v.inputs[i].Focus()
		} else {
			v.inputs[i].Blur()
		}
	}

	cmd := v.updateInputs(msg)

	return v, cmd
}

func (v TurnView) View() string {
	var b strings.Builder
	buttonText := fmt.Sprintf("[%v]", v.actionText)
	// draws all of the inputs
	for i := range v.inputs {
		b.WriteString(v.inputs[i].View())
		b.WriteString("\n")
	}
	// draws the action button

	if v.focusedIndex == len(v.inputs) {
		b.WriteString(focusedStyle.Render(buttonText))
	} else {
		b.WriteString(buttonText)
	}

	return b.String()
}

func (v TurnView) Reset() TurnView {
	// reset all input values
	for i := 0; i < len(v.inputs); i++ {
		v.inputs[i].Reset()
	}

	v.focusedIndex = 0

	return v
}

func (v *TurnView) Focus() {
	v.focused = true
}
