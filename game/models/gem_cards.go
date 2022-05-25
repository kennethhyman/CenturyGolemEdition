package models

import "strconv"
import "strings"

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

func (gems GemValues) HasInputs(inputs GemValues) bool {
	return (gems.Yellow >= inputs.Yellow && gems.Green >= inputs.Green && gems.Blue >= inputs.Blue && gems.Pink >= inputs.Pink)
}

type GemCard struct {
	Inputs   GemValues
	Outputs  GemValues
	Upgrades int
}

func (g GemValues) String() string {
	output := ""
	output += strings.Repeat(yellowString, g.Yellow)
	output += strings.Repeat(greenString, g.Green)
	output += strings.Repeat(blueString, g.Blue)
	output += strings.Repeat(pinkString, g.Pink)

	return output
}

func (g GemCard) String() string {
	if g.Upgrades > 0 {
		return "\033[0m[ " + strings.Repeat(upgradeString, g.Upgrades) + "\033[0m]"
	}

	return "\033[0m[ " + g.Inputs.String() + yieldsString + g.Outputs.String() + "\033[0m ]"
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
