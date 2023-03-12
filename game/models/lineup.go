package models

import "errors"
import "strings"
import "fmt"

type card interface {
	String() string
	Cost() GemValues
}

type LineUp[CardType card] struct {
	Stack []CardType
	deck  *Deck[CardType] // the supplying deck
}

func NewGemLineUp() LineUp[GemCard] {
	d := NewGemDeck()
	d.Shuffle()
	cards, _ := d.DrawCards(6)

	lineup := LineUp[GemCard]{
		deck:  d,
		Stack: cards,
	}

	return lineup
}

func NewGolemLinup() LineUp[GolemCard] {
	d := NewGolemDeck()
	d.Shuffle()

	cards, _ := d.DrawCards(6)
	return LineUp[GolemCard]{
		deck:  d,
		Stack: cards,
	}
}

// Can pick a card from anywhere
func (l LineUp[CardType]) String() string {
	var cards []string
	for i, card := range l.Stack {
		cards = append(cards, fmt.Sprintf("%v:%v", i, card.String()))
	}

	return strings.Join(cards[:], " ") + fmt.Sprintf("\t[%v]", len(l.deck.stack))
}

func (l *LineUp[CardType]) Draw(index int) (CardType, error) {
	var card CardType
	var err error
	if index >= len(l.Stack) {
		return card, errors.New("index is out of range")
	}

	card = l.Stack[index]

	//
	for index < (len(l.Stack) - 1) {
		l.Stack[index] = l.Stack[index+1]
		index++
	}

	l.Stack[len(l.Stack)-1], err = l.deck.DrawCard()
	if err != nil {
		fmt.Println("no more cards to draw")
		// reduce the lineup length by one
		l.Stack = l.Stack[0 : len(l.Stack)-1]
	}

	return card, nil
}
