package server

import (
  "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
)

func MarshallGems(gems *models.GemValues) *pb.GemValues {
  return &pb.GemValues {
    Yellow: int32(gems.Yellow),
    Green: int32(gems.Green),
    Blue: int32(gems.Blue),
    Pink: int32(gems.Pink),
  }
}

func UnmarshallGems(gems *pb.GemValues) *models.GemValues {
  return &models.GemValues {
    Yellow: int(gems.Yellow),
    Green: int(gems.Green),
    Blue: int(gems.Blue),
    Pink: int(gems.Pink),
  }
}
