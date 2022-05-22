package models

import "errors"
import "strings"
import "fmt"

type LineUp struct {
	stack []GemCard
  deck *Deck[GemCard] // the supplying deck
}

// Can Create from deck
func (l *LineUp) PickCard(index int) (GemCard, error) {
  var empty_card GemCard
  if index >= len(l.stack) {
    return empty_card, errors.New("index is out of range")
  }
  
  return empty_card, nil
}

func NewGemLineUp(d *Deck[GemCard]) (LineUp) {
  cards, _ := d.DrawCards(6)

  lineup := LineUp{
    deck:  d,
    stack: cards,
  }
  
  return lineup
}

// Can pick a card from anywhere
func (l LineUp) String() string {
  var cards []string
  for i, card  := range(l.stack) {
    cards = append(cards, fmt.Sprintf("%v:%v", i, card.String()))
  }

  return strings.Join(cards[:], " ") + fmt.Sprintf("\t[%v]", len(l.deck.stack))
}

func (l *LineUp) Draw(index int) (GemCard, error) {
  var card GemCard
  var err error
  if (index >= len(l.stack)) {
    return card, errors.New("index is out of range")
  }
  
  card = l.stack[index]

  // 
  for index < (len(l.stack) - 1) {
    l.stack[index] = l.stack[index + 1]
    index++
  }
  l.stack[len(l.stack)-1], err = l.deck.DrawCard()
  if (err != nil) {
    fmt.Println("no more cards to draw")
    // reduce the lineup length by one
    l.stack = l.stack[0:len(l.stack) - 1]
  }

  return card, nil
} 
