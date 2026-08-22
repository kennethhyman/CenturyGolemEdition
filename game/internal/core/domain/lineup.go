package domain

import "errors"
import "strings"
import "fmt"
import "encoding/json"

type card interface {
	String() string
	Cost() GemValues
}

type LineUp[CardType card] struct {
	Stack []CardType
	deck  *Deck[CardType] // the supplying deck
}

// lineupJSON is LineUp's JSON wire shape. It exists because LineUp's `deck`
// field is unexported: encoding/json silently skips unexported fields, so
// marshalling a LineUp directly would drop every card still left to draw.
// That matters for persistence (Stage 2 stores domain.Game as a JSON blob)
// where a round-tripped LineUp must still support Draw/Remaining/String
// exactly like a freshly built one.
type lineupJSON[CardType card] struct {
	Stack []CardType
	Deck  []CardType
}

// MarshalJSON encodes both the face-up Stack and the supplying deck's
// remaining cards.
func (l LineUp[CardType]) MarshalJSON() ([]byte, error) {
	var deckCards []CardType
	if l.deck != nil {
		deckCards = l.deck.stack
	}

	return json.Marshal(lineupJSON[CardType]{Stack: l.Stack, Deck: deckCards})
}

// UnmarshalJSON is MarshalJSON's inverse.
//
// Post: l.deck is non-nil, so Remaining/Draw/String behave the same as on a
// LineUp built via NewGemLineUp/NewGolemLinup.
func (l *LineUp[CardType]) UnmarshalJSON(data []byte) error {
	var aux lineupJSON[CardType]
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	l.Stack = aux.Stack
	l.deck = &Deck[CardType]{stack: aux.Deck}
	return nil
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

	cards, _ := d.DrawCards(5)
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

func (l LineUp[CardType]) Remaining() int {
  return l.deck.Length()
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
