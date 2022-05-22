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

	deck := NewGemDeck()
	deck.Shuffle()
	fmt.Printf("Deck: %v\n", deck)
	card, err := deck.DrawCard()
	if err == nil {
		fmt.Printf("first card: %v\n", card)
	}
}

func playGemCard(gems *GemValues, card *GemCard) {
	fmt.Println("required inputs: ", card.Inputs)
	fmt.Println("current gem stack: ", *gems)

	if gems.HasInputs(card.Inputs) {
		fmt.Println("You can play this card!")
		gems.Yellow += card.Outputs.Yellow
		gems.Green += card.Outputs.Green
		gems.Blue += card.Outputs.Blue
		gems.Pink += card.Outputs.Pink
		fmt.Println("New Gem stack: ", *gems)
		// TODO: Handle case where you end up with too many gems
	}
}
