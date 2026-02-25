package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kennethhyman/CenturyGolemEdition/adapters/handlers"
	"github.com/kennethhyman/CenturyGolemEdition/adapters/repository"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/services/game_service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewInMemoryRepository()
	svc := game_service.New(repo)
	handler := handlers.New(svc)

	s := grpc.NewServer()
	pb.RegisterGameServer(s, handler)

	fmt.Printf("server mode: Listening on port 50051\n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
