package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// SubmitMsg is emitted once a TurnView's action button is triggered and its
// local validation (action) passes. TurnView only knows how to collect and
// locally validate raw text; it has no notion of grpc, addresses, or game
// state, so it hands the raw values up to whichever parent model owns
// those (cli.model) to turn into an actual network call. Kept in this
// package (not cli) so TurnView stays self-contained and cli can depend on
// actions without a cycle.
type SubmitMsg struct {
	Action GameAction
	Values []string
}

// TurnView collects the argument(s) one GameAction needs (as free-form
// text, one field per input) and validates them locally before handing
// them off via SubmitMsg.
type TurnView struct {
	inputs       []textinput.Model
	focusedIndex int
	focused      bool
	actionText   string
	gameAction   GameAction
	error        error
	// action validates the current input values, read fresh at submit time
	// (not closed over at construction — TurnView's Update has a value
	// receiver, so a closure captured over the construction-time inputs
	// would never see what the user actually typed). A nil action always
	// passes.
	action func(values []string) error
}

// NewTurnView builds a TurnView for gameAction/actionText, backed by
// inputs, validated at submit time by action (nil means "always valid").
func NewTurnView(gameAction GameAction, actionText string, inputs []textinput.Model, action func(values []string) error) TurnView {
	v := TurnView{
		gameAction: gameAction,
		actionText: actionText,
		inputs:     inputs,
		action:     action,
	}
	v.Focus()
	if len(v.inputs) > 0 {
		v.inputs[0].Focus()
	}
	return v
}

// isGemLetters reports whether s is empty or made up only of gem-color
// letters (y/g/b/p, case-insensitive) — the wire format PlayGemCard's
// discards, GetGemCard's payment, and Rest's discards all share.
func isGemLetters(s string) bool {
	for _, c := range strings.ToLower(s) {
		switch c {
		case 'y', 'g', 'b', 'p':
		default:
			return false
		}
	}
	return true
}

func (v TurnView) currentValues() []string {
	values := make([]string, len(v.inputs))
	for i, input := range v.inputs {
		values[i] = input.Value()
	}
	return values
}

func (v TurnView) Init() tea.Cmd {
	return nil
}

func (v TurnView) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(v.inputs))

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
				values := v.currentValues()
				if v.action != nil {
					if err := v.action(values); err != nil {
						v.error = err
						return v, nil
					}
				}
				v.error = nil
				return v, func() tea.Msg {
					return SubmitMsg{Action: v.gameAction, Values: values}
				}
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

	return v, tea.Batch(cmd, tea.Batch(cmds...))
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
		b.WriteString(blurredStyle.Render(buttonText))
	}

	if v.error != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(v.error.Error()))
	}

	return b.String()
}

func (v TurnView) Reset() TurnView {
	// reset all input values
	for i := 0; i < len(v.inputs); i++ {
		v.inputs[i].Reset()
	}

	v.focusedIndex = 0
	v.error = nil

	return v
}

func (v *TurnView) Focus() {
	v.focused = true
}

// errValueRequired is returned by field validators when a required value
// is blank.
func errValueRequired(field string) error {
	return errors.New(field + " is required")
}
