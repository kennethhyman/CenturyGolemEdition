package domain

import (
  "strconv"
  "strings"
  "fmt"
  pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

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

func MarshallGolemCard(card *GolemCard) *pb.GolemCard {
  return &pb.GolemCard{
    Cost: MarshallGems(&card.Inputs),
    Points: int32(card.Points),
  }
}

func UnmarshallGolemCard(card *pb.GolemCard) *GolemCard {
  return &GolemCard {
    Inputs: *UnmarshallGems(card.Cost),
    Points:  int(card.Points),
  }
}
