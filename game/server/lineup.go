package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"

)

func UnmarshallGameState(card *pb.GameState) *models.Game {
  return &models.Game {
  }
}
