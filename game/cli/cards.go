package cli

import (
	//"github.com/kennethhyman/CenturyGolemEdition/server"
	"fmt"
	"strings"

	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

const yellowString = "\033[1;33m○\033[0m"
const greenString = "\033[2;32m○"
const blueString = "\033[2;34m○"
const pinkString = "\033[38;5;206m○"
const upgradeString = "🌟"
const yieldsString = "\033[0m->"

type GemValue struct {
  Yellow int
  Green int
  Blue int
  Pink int
}

type GemCard struct {
  In GemValue
  Out GemValue
  Upgrades int
}

type GolemCard struct {
  Cost GemValue
  Points int
}

func (g GemValue) String() string {
	output := ""
	output += strings.Repeat(yellowString, g.Yellow)
	output += strings.Repeat(greenString, g.Green)
	output += strings.Repeat(blueString, g.Blue)
	output += strings.Repeat(pinkString, g.Pink)

	return output
}

func (g GemValue) Count() int {
  return g.Yellow + g.Green + g.Blue + g.Pink
}

func (g GemCard) String() string {
  if g.Upgrades > 0 {
		return "\033[0m[ " + strings.Repeat(upgradeString, g.Upgrades) + "\033[0m]"
	}

  return "\033[0m[" + g.In.String() + yieldsString + g.Out.String() + "\033[0m]"
}

func (g GolemCard) String() string {
  return fmt.Sprintf("\033[0m[%v|%v\033[0m]", g.Points, g.Cost)
}

func (g GolemCard) StringLength() int {
  // Number of gems
  length := g.Cost.Count()

  // Length of point string
  if g.Points >= 10 {
    length += 2
  } else {
    length += 1
  }
  
  // Extra for the |
  return length + 1
}

func (g GemCard) StringLength() int {
  if g.Upgrades > 0 {
    return g.Upgrades
  } else {
    // Number of gems plus the -> string
    return g.In.Count() + g.Out.Count() + 2
  }
}


func UnmarshallGemCard(card *pb.GemCard) GemCard {
  return GemCard{
    In: GemValue {
      Yellow: int(card.In.Yellow),
      Green: int(card.In.Green),
      Blue: int(card.In.Blue),
      Pink: int(card.In.Pink),
    },
    Out: GemValue {
      Yellow: int(card.Out.Yellow),
      Green: int(card.Out.Green),
      Blue: int(card.Out.Blue),
      Pink: int(card.Out.Pink),
    },
    Upgrades: int(card.Upgrades),
  }
}

func UnmarshallGolemCard(card *pb.GolemCard) GolemCard {
  return GolemCard{
    Cost: GemValue {
      Yellow: int(card.Cost.Yellow),
      Green: int(card.Cost.Green),
      Blue: int(card.Cost.Blue),
      Pink: int(card.Cost.Pink),
    },
    Points: int(card.Points),
  }
}
