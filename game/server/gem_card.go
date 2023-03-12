package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"

)

func MarshallGemCard(card *models.GemCard) *pb.GemCard {
  return &pb.GemCard{
    In: &pb.GemValues {
      Yellow: int32(card.Inputs.Yellow),
      Green: int32(card.Inputs.Green),
      Blue: int32(card.Inputs.Blue),
      Pink: int32(card.Inputs.Pink),
    },
    Out: &pb.GemValues {
      Yellow: int32(card.Outputs.Yellow),
      Green: int32(card.Outputs.Green),
      Blue: int32(card.Outputs.Blue),
      Pink: int32(card.Outputs.Pink),
    },
    Upgrades: int32(card.Upgrades),
  }
}

func UnmarshallGemCard(card *pb.GemCard) *models.GemCard {
  return &models.GemCard {
    Inputs: *UnmarshallGems(card.In),
    Outputs: *UnmarshallGems(card.Out), 
    Upgrades: int(card.Upgrades),
  }
}
