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

  games []models.Game
}

func (s *gameServer) NewGame(ctx context.Context, newGameInfo *pb.CreateGameMessage) (*pb.CreateGameResponse, error) {
  newGame := models.NewGame(int(newGameInfo.PlayerCount))
  fmt.Print("CREATING NEW GAME\n")
  return MarshallGame(newGame), nil
}

func (s *gameServer) PlayGemCard(ctx context.Context, playGemCardInfo *pb.PlayGemCardMessage) (*pb.PlayGemCardResponse, error){
  // TODO Check for correct player making the action
  // find game info
  game := s.games[playGemCardInfo.GameId]
  card := *UnmarshallGemCard(playGemCardInfo.Card)
  discards := *UnmarshallGems(playGemCardInfo.Discarded)
   
  // update game with appropriate play gem card playGemCardInfo
  _ = game.PlayGemCard(card, discards)
  
  // update the game state stored on the server

  // return new gamestate
  return (*pb.PlayGemCardResponse)(MarshallGame(&game)), nil
}

func newServer() *gameServer {
	s := &gameServer{}
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
