package main

import "fmt"

type GameAction int64

const (
	PlayGemCard GameAction = iota
	GetGemCard
	Rest
	GetGolemCard
)

func (ga GameAction) view() string {
	return fmt.Sprintf("Game action: %v", ga)
}
