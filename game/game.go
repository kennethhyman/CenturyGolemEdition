package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/kennethhyman/CenturyGolemEdition/models"
	"os"
)

type model struct {
	game   Game
	action GameAction
}

func initialModel() model {
	return model{game: *NewGame(2)}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "1":
			m.action = PlayGemCard
		case "2":
			m.action = GetGemCard
		case "3":
			m.action = Rest
		case "4":
			m.action = GetGolemCard
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%v\n%v", m.game.String(), m.action.view())
}

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

// func main() {
// 	game := NewGame(2)
// 	for !game.Finished() {
// 		// print gamestate
// 		var err error
// 		fmt.Printf("%v\n", game.String())

// 		// give options
// 		fmt.Println("What would you like to do for your turn?")
// 		fmt.Println("1. Buy a gem card")
// 		fmt.Println("2. Buy a golem card")
// 		fmt.Println("3. Play a Gem Card")
// 		fmt.Println("4. Pickup your discard pile")
// 		// get user input
// 		var turn string
// 		fmt.Scan(&turn)
// 		// take turn
// 		switch string(turn) {
// 		case "1":
// 			// buy gem card
// 			err = game.GetGemCard()
// 		case "2":
// 			// buy golem card
// 			err = game.BuyGolem()
// 		case "3":
// 			// play gem card
// 			game.PlayGemCard()
// 		case "4":
// 			// pickup discard pile
// 			game.PickUpGemCards()
// 		}

// 		fmt.Printf("%v\n", err)
// 		// change player
// 		if err == nil {
// 			game.NextTurn()
// 		}
// 	}
// }
