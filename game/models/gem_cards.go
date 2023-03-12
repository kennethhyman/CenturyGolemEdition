package models

import "strconv"
import "strings"
import "errors"
import "regexp"

type GemValues struct {
	Yellow int
	Green  int
	Blue   int
	Pink   int
}

const yellowString = "\033[1;33m○\033[0m"
const greenString = "\033[2;32m○"
const blueString = "\033[2;34m○"
const pinkString = "\033[38;5;206m○"
const upgradeString = "🌟"
const yieldsString = "\033[0m->"
const gemRE = "[y|g|b|p]+"

type GemCard struct {
	Inputs   GemValues
	Outputs  GemValues
	Upgrades int
}

func NewGemCard() GemCard {
	return GemCard{
		Inputs: GemValues{
			Yellow: 0,
			Green:  0,
			Blue:   0,
			Pink:   0,
		},
		Outputs: GemValues{
			Yellow: 0,
			Green:  0,
			Blue:   0,
			Pink:   0,
		},
		Upgrades: 0,
	}
}

func (g GemValues) String() string {
	output := ""
	output += strings.Repeat(yellowString, g.Yellow)
	output += strings.Repeat(greenString, g.Green)
	output += strings.Repeat(blueString, g.Blue)
	output += strings.Repeat(pinkString, g.Pink)

	return output
}

func (g GemValues) count() int {
	return g.Yellow + g.Green + g.Blue + g.Pink
}

func (g GemValues) effectiveValue() int {
	return (g.Yellow) + (g.Green * 2) + (g.Blue * 3) + (g.Pink * 4)
}

func (g *GemValues) add(newGems GemValues) {
	g.Yellow = g.Yellow + newGems.Yellow
	g.Green = g.Green + newGems.Green
	g.Blue = g.Blue + newGems.Blue
	g.Pink = g.Pink + newGems.Pink
}

func (g *GemValues) remove(cost GemValues) error {
	if cost.Yellow > g.Yellow || cost.Green > g.Green ||
		cost.Blue > g.Blue || cost.Pink > g.Pink {
		return errors.New("You do not have the gems required")
	}

	g.Yellow = g.Yellow - cost.Yellow
	g.Green = g.Green - cost.Green
	g.Blue = g.Blue - cost.Blue
	g.Pink = g.Pink - cost.Pink

	return nil
}

func (g *GemValues) strictlyGreater(other GemValues) bool {
	gArr := g.Array()
	otherArr := g.Array()

	for i := range gArr {
		if gArr[i] != 0 && otherArr[i] <= gArr[i] {
			return true
		}
	}

	return false
}

func (g *GemValues) Array() []int {
	return []int{g.Yellow, g.Green, g.Blue, g.Pink}
}

func (g GemCard) String() string {
	if g.Upgrades > 0 {
		return "\033[0m[ " + strings.Repeat(upgradeString, g.Upgrades) + "\033[0m]"
	}

	return "\033[0m[ " + g.Inputs.String() + yieldsString + g.Outputs.String() + "\033[0m ]"
}

func (g GemCard) Cost() GemValues {
	return g.Inputs
}



func parseGemCard(card string) *GemCard {
	// convert str to int
	strArr := strings.Split(card, ",")
	convertedArr := make([]int, len(strArr))

	for i, val := range strArr {
		convertedArr[i], _ = strconv.Atoi(val)
	}

	// Craft our card
	inputs := convertedArr[0:4]
	outputs := convertedArr[4:8]
	upgrades := convertedArr[8]
	return &GemCard{
		Inputs: GemValues{
			Yellow: inputs[0],
			Green:  inputs[1],
			Blue:   inputs[2],
			Pink:   inputs[3],
		},
		Outputs: GemValues{
			Yellow: outputs[0],
			Green:  outputs[1],
			Blue:   outputs[2],
			Pink:   outputs[3],
		},
		Upgrades: upgrades,
	}
}

func parseGemInput(gems string) (GemValues, error) {
	// Check it is a valid string
	gemValue := GemValues{}
	match, err := regexp.MatchString(gemRE, gems)
	if err != nil {
		return gemValue, err
	}

	if !match {
		return gemValue, errors.New("Invalid input string")
	}

	// iterate through
	for _, char := range gems {
		switch string(char) {
		case "y", "Y":
			gemValue.Yellow += 1
		case "g", "G":
			gemValue.Green += 1
		case "b", "B":
			gemValue.Blue += 1
		case "p", "P":
			gemValue.Pink += 1
		}
	}
	return gemValue, nil
}
