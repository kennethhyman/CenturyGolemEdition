package models

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

//go:embed resources/*
var resources embed.FS

type Deck[CardType any] struct {
	stack []CardType // array of card IDs
}

func GetStartingDeck() []GemCard {

	f, err := resources.Open("resources/starter_cards.csv")
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(f)

	newDeck := []GemCard{}
	for scanner.Scan() {
		line := scanner.Text()
		card := *parseGemCard(line)
		newDeck = append(newDeck, card)
	}

	f.Close()
	return newDeck
}

func NewGemDeck() *Deck[GemCard] {
	// Read in values of cards in order
	f, err := resources.Open("resources/cards.csv")
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(f)

	newDeck := []GemCard{}
	for scanner.Scan() {
		line := scanner.Text()
		card := *parseGemCard(line)
		newDeck = append(newDeck, card)
	}

	f.Close()

	return &Deck[GemCard]{
		stack: newDeck,
	}
}

func NewGolemDeck() *Deck[GolemCard] {
	// Read in values of cards in order
	f, err := resources.Open("resources/golems.csv")
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(f)

	newDeck := []GolemCard{}
	for scanner.Scan() {
		line := scanner.Text()
		card := *parseGolemCard(line)
		newDeck = append(newDeck, card)
	}

	f.Close()

	return &Deck[GolemCard]{
		stack: newDeck,
	}
}

func (d *Deck[_]) Shuffle() {
	//shuffle deck
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.stack), func(i, j int) { d.stack[i], d.stack[j] = d.stack[j], d.stack[i] })
}

func (d *Deck[CardType]) DrawCard() (CardType, error) {
	var empty_card CardType
	if len(d.stack) <= 0 {
		return empty_card, errors.New("This deck is empty")
	}

	firstCard := d.stack[0]
	d.stack = d.stack[1:len(d.stack)]

	return firstCard, nil
}

func (d *Deck[CardType]) DrawCards(count int) ([]CardType, error) {
	var cards []CardType

	if count >= len(d.stack) {
		return cards, errors.New("Not enough cards to draw from")
	}

	cards = d.stack[0:count]
	d.stack = d.stack[count:len(d.stack)]

	return cards, nil
}
