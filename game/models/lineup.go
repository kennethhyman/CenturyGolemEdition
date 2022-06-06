package models

import "errors"
import "strings"
import "fmt"

type card interface {
	String() string
	Cost() GemValues
}

type LineUp[CardType card] struct {
	stack []CardType
	deck  *Deck[CardType] // the supplying deck
}

func NewGemLineUp() LineUp[GemCard] {
	d := NewGemDeck()
	d.Shuffle()
	cards, _ := d.DrawCards(6)

	lineup := LineUp[GemCard]{
		deck:  d,
		stack: cards,
	}

	return lineup
}

func NewGolemLinup() LineUp[GolemCard] {
	d := NewGolemDeck()
	d.Shuffle()

	cards, _ := d.DrawCards(6)
	return LineUp[GolemCard]{
		deck:  d,
		stack: cards,
	}
}

// Can pick a card from anywhere
func (l LineUp[CardType]) String() string {
	var cards []string
	for i, card := range l.stack {
		cards = append(cards, fmt.Sprintf("%v:%v", i, card.String()))
	}

	return strings.Join(cards[:], " ") + fmt.Sprintf("\t[%v]", len(l.deck.stack))
}

func (l *LineUp[CardType]) Draw(index int) (CardType, error) {
	var card CardType
	var err error
	if index >= len(l.stack) {
		return card, errors.New("index is out of range")
	}

	card = l.stack[index]

	//
	for index < (len(l.stack) - 1) {
		l.stack[index] = l.stack[index+1]
		index++
	}

	l.stack[len(l.stack)-1], err = l.deck.DrawCard()
	if err != nil {
		fmt.Println("no more cards to draw")
		// reduce the lineup length by one
		l.stack = l.stack[0 : len(l.stack)-1]
	}

	return card, nil
}
