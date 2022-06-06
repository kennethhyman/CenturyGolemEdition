package models

import "strconv"
import "strings"
import "fmt"

//import "errors"

type GolemCard struct {
	Inputs GemValues
	Points int
}

func (g GolemCard) String() string {
	return fmt.Sprintf("[ %v\033[0m | %v ]", g.Inputs, g.Points)
}

func parseGolemCard(card string) *GolemCard {
	// convert str to int
	strArr := strings.Split(card, ",")
	convertedArr := make([]int, len(strArr))

	for i, val := range strArr {
		convertedArr[i], _ = strconv.Atoi(val)
	}

	// Craft our card
	cost := convertedArr[0:4]
	points := convertedArr[4]
	return &GolemCard{
		Inputs: GemValues{
			Yellow: cost[0],
			Green:  cost[1],
			Blue:   cost[2],
			Pink:   cost[3],
		},
		Points: points,
	}
}

func (g GolemCard) Cost() GemValues {
	return g.Inputs
}
