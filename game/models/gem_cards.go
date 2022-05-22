package models

import "strings"

type GemValues struct {
	Yellow int
	Green  int
	Blue   int
	Pink   int
}

const yellowString = "\033[33m○"
const greenString = "\033[32m○"
const blueString = "\033[34m○"
const pinkString = "\033[35m○"
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
