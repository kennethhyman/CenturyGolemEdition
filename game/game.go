package main

import "fmt"
import . "github.com/kennethhyman/CenturyGolemEdition/models"

func main() {
	//gem_stack := GemValues{
	//	Yellow: 0,
	//	Green:  0,
	//	Blue:   0,
	//	Pink:   0,
	//}
	game := NewGame(2)
	for !game.Finished() {
		// print gamestate
		var err error
		fmt.Printf("%v\n", game.String())

		// give options
		fmt.Println("What would you like to do for your turn?")
		fmt.Println("1. Buy a gem card")
		fmt.Println("2. Buy a golem card")
		fmt.Println("3. Play a Gem Card")
		fmt.Println("4. Pickup your discard pile")
		// get user input
		var turn string
		fmt.Scan(&turn)
		// take turn
		switch string(turn) {
		case "1":
			// buy gem card
			err = game.GetGemCard()
		case "2":
			// buy golem card
			err = game.BuyGolem()
		case "3":
			// play gem card
			game.PlayGemCard()
		case "4":
			// pickup discard pile
			game.PickUpGemCards()
		}

		fmt.Printf("%v\n", err)
		// change player
		if err == nil {
			game.NextTurn()
		}
	}
}
