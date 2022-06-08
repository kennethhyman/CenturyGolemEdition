package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"strings"
)

type GameAction string

const (
	PlayGemCard  GameAction = "PlayGemCard"
	GetGemCard   GameAction = "GetGemCard"
	Rest         GameAction = "Rest"
	GetGolemCard GameAction = "GetGolemCard"
)

type ActionMenu struct {
	actionType GameAction
	focusIndex int
	inputs     []textinput.Model
	error      error
}

func NewActionMenu(ga GameAction) ActionMenu {
	var actionMenu = ActionMenu{actionType: ga}
	switch ga {
	case PlayGemCard:
		// which card would you like to play?
		var t = textinput.New()
		t.CharLimit = 32
		t.Placeholder = "Which Card would you like to play?"
		t.Focus()

		actionMenu.inputs = []textinput.Model{t}
	case GetGemCard:

	case Rest:

	case GetGolemCard:

	}
	return actionMenu
}

func (am ActionMenu) View() string {
	var b strings.Builder

	for i := range am.inputs {
		b.WriteString(am.inputs[i].View())
		if i < len(am.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

var options = fmt.Sprintf("1. %v\t2. %v\t 3.%v\t4. %v\n", PlayGemCard, GetGemCard, Rest, GetGolemCard)

func (ga GameAction) View() string {
	var b strings.Builder

	b.WriteString("What action would you like to do?\n")
	b.WriteString(options)
	b.WriteString(fmt.Sprintf("Game action: %v", ga))

	switch ga {
	case PlayGemCard:

	case GetGemCard:

	case Rest:

	case GetGolemCard:

	}
	return b.String()
}

func PlayGemCardView() string {

	return ""
}
