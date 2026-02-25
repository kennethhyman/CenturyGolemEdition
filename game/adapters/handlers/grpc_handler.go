package handlers

import (
	"context"
	"sync"

	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcHandler implements pb.GameServer and delegates to the GameService.
// It maintains a mapping from the proto's int32 game_id to the service's UUID strings.
type GrpcHandler struct {
	pb.UnimplementedGameServer
	service ports.GameService
	mu      sync.Mutex
	nextID  int32
	idMap   map[int32]string
}

func New(service ports.GameService) *GrpcHandler {
	return &GrpcHandler{
		service: service,
		idMap:   make(map[int32]string),
	}
}

func (h *GrpcHandler) GetCard(ctx context.Context, msg *pb.GetCardMessage) (*pb.GemCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCard not implemented")
}

func (h *GrpcHandler) NewGame(ctx context.Context, msg *pb.CreateGameMessage) (*pb.CreateGameResponse, error) {
	game, err := h.service.Create(int(msg.PlayerCount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create game: %v", err)
	}

	h.mu.Lock()
	h.nextID++
	h.idMap[h.nextID] = game.ID
	h.mu.Unlock()

	return domain.MarshallGame(&game), nil
}

func (h *GrpcHandler) PlayGemCard(ctx context.Context, msg *pb.PlayGemCardMessage) (*pb.PlayGemCardResponse, error) {
	h.mu.Lock()
	gameUUID, ok := h.idMap[msg.GameId]
	h.mu.Unlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "game %d not found", msg.GameId)
	}

	if msg.Card == nil {
		return nil, status.Errorf(codes.InvalidArgument, "card is required")
	}

	card := domain.UnmarshallGemCard(msg.Card)

	var discards domain.GemValues
	if msg.Discarded != nil {
		discards = *domain.UnmarshallGems(msg.Discarded)
	}

	game, err := h.service.PlayGemCard(gameUUID, 0, *card, discards)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to play gem card: %v", err)
	}

	gs := domain.MarshallGame(&game)
	return &pb.PlayGemCardResponse{GameState: gs.GameState}, nil
}
