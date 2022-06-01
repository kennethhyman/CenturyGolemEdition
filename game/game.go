package main

import "fmt"
import . "github.com/kennethhyman/CenturyGolemEdition/models"

func main() {
	//gem_stack := GemValues{
	//	Yellow: 0,
	//	Green:  0,
	//	Blue:   0,
	//	Pink:   0,
	//}
	game := NewGame(2)
	fmt.Printf("%v\n", game.String())
}
