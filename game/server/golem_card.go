package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"

)

func MarshallGolemCard(card *models.GolemCard) *pb.GolemCard {
  return &pb.GolemCard{
    Inputs: MarshallGems(&card.Inputs),
    Points: int32(card.Points),
  }
}

func UnmarshallGolemCard(card *pb.GolemCard) *models.GolemCard {
  return &models.GolemCard {
    Inputs: *UnmarshallGems(card.Inputs),
    Points:  int(card.Points),
  }
}
