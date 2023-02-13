package server

import (
	"context"
	//"flag"
	"fmt"
	"log"
  "net"

	"google.golang.org/grpc"
	//"github.com/golang/protobuf/proto"
  models "github.com/kennethhyman/CenturyGolemEdition/models"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
)

type gameServer struct {
	pb.UnimplementedGameServer

  gemcards models.Deck[models.GemCard]
}


func marshallCard(card models.GemCard) (*pb.GemCard) {
  return &pb.GemCard{
    In: &pb.GemValue {
      Yellow: int32(card.Inputs.Yellow),
      Green: int32(card.Inputs.Green),
      Blue: int32(card.Inputs.Blue),
      Pink: int32(card.Inputs.Pink),
    },
    Out: &pb.GemValue {
      Yellow: int32(card.Outputs.Yellow),
      Green: int32(card.Outputs.Green),
      Blue: int32(card.Outputs.Blue),
      Pink: int32(card.Outputs.Pink),
    },
  }
}

func (s *gameServer) GetCard(ctx context.Context, _ *pb.GetCardMessage) (*pb.GemCard, error) {
  card, _ := s.gemcards.DrawCard()

  return marshallCard(card), nil
}


func newServer() *gameServer {
	s := &gameServer{
    gemcards: *models.NewGemDeck(),
  }
	return s
}

func StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGameServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
