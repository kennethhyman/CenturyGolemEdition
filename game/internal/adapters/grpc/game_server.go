// Package grpcadapter implements the generated pb.GameServer interface,
// translating grpc requests into game_service.Service calls and marshalling
// the domain.Game results back into grpc responses.
//
// It lives outside internal/core/grpc (the generated-code package) because
// internal/core/domain already imports that package for its pb message
// types; a GameServer implementation needs domain types too, and importing
// domain from internal/core/grpc would close an import cycle
// (grpc -> domain -> grpc). This package depends one-way on both
// internal/core/grpc and internal/core/domain, breaking that cycle.
package grpcadapter

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kennethhyman/CenturyGolemEdition/internal/core/domain"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

// GameService is the subset of game_service.Service this adapter calls —
// the full ports.GameService contract, now that game_service.Service
// implements all of it (Stage 3).
type GameService interface {
	Create(players int) (domain.Game, error)
	Get(id string) (domain.Game, error)
	PlayGemCard(id string, player int, card domain.GemCard) (domain.Game, error)
	GetGemCard(id string, player int, card domain.GemCard, payment string) (domain.Game, error)
	GetGolemCard(id string, player int, golem domain.GolemCard) (domain.Game, error)
	Rest(id string, player int, discards string) (domain.Game, error)
}

// GameServer adapts a GameService to the grpc Game service (pb.GameServer).
type GameServer struct {
	pb.UnimplementedGameServer
	service GameService
}

var _ pb.GameServer = (*GameServer)(nil)

// NewGameServer builds a GameServer that serves grpc requests off of
// service.
//
// Pre: service is non-nil.
func NewGameServer(service GameService) *GameServer {
	return &GameServer{service: service}
}

// NewGame implements pb.GameServer.
//
// Pre: req.PlayerCount > 0.
// Post: on success, returns the freshly created game's state, already saved
// to the backing repository via GameService.Create.
func (a *GameServer) NewGame(ctx context.Context, req *pb.CreateGameMessage) (*pb.CreateGameResponse, error) {
	if req.GetPlayerCount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "player_count must be positive")
	}

	game, err := a.service.Create(int(req.GetPlayerCount()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create game: %v", err)
	}

	return domain.MarshallGame(&game), nil
}

// GetCard implements pb.GameServer: picks up a gem card from the lineup.
//
// Pre: req.GameId names an existing game; req.Player is that game's current
// player; req.Card matches (by value) an entry in the gem lineup; req.Payment
// carries one gem letter per lineup card bypassed to reach it.
// Post: on success, returns the updated game state.
func (a *GameServer) GetCard(ctx context.Context, req *pb.GetCardMessage) (*pb.GetCardResponse, error) {
	if req.GetGameId() == "" {
		return nil, status.Error(codes.InvalidArgument, "game_id is required")
	}

	card := domain.UnmarshallGemCard(req.GetCard())
	game, err := a.service.GetGemCard(req.GetGameId(), int(req.GetPlayer()), *card, req.GetPayment())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "get card: %v", err)
	}

	return &pb.GetCardResponse{GameState: domain.MarshallGameState(&game)}, nil
}

// PlayGemCard implements pb.GameServer: plays a gem card from the current
// player's hand.
//
// Pre: req.GameId names an existing game; req.Player is that game's current
// player; req.Card matches (by value) a card in that player's hand.
// Post: on success, returns the updated game state.
func (a *GameServer) PlayGemCard(ctx context.Context, req *pb.PlayGemCardMessage) (*pb.PlayGemCardResponse, error) {
	if req.GetGameId() == "" {
		return nil, status.Error(codes.InvalidArgument, "game_id is required")
	}

	card := domain.UnmarshallGemCard(req.GetCard())
	game, err := a.service.PlayGemCard(req.GetGameId(), int(req.GetPlayer()), *card)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "play gem card: %v", err)
	}

	return &pb.PlayGemCardResponse{GameState: domain.MarshallGameState(&game)}, nil
}

// BuyGolem implements pb.GameServer: buys a golem card from the lineup.
//
// Pre: req.GameId names an existing game; req.Player is that game's current
// player; req.Golem matches (by value) an entry in the golem lineup and the
// player can afford it — an unaffordable (or non-matching) golem is
// correctly rejected with FailedPrecondition, not silently ignored.
// Post: on success, returns the updated game state.
func (a *GameServer) BuyGolem(ctx context.Context, req *pb.BuyGolemMessage) (*pb.BuyGolemResponse, error) {
	if req.GetGameId() == "" {
		return nil, status.Error(codes.InvalidArgument, "game_id is required")
	}

	golem := domain.UnmarshallGolemCard(req.GetGolem())
	game, err := a.service.GetGolemCard(req.GetGameId(), int(req.GetPlayer()), *golem)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "buy golem: %v", err)
	}

	return &pb.BuyGolemResponse{GameState: domain.MarshallGameState(&game)}, nil
}

// Rest implements pb.GameServer: ends the current player's turn.
//
// Pre: req.GameId names an existing game; req.Player is that game's current
// player; if that player holds more than 10 gems, req.Discards must bring
// them to 10 or fewer.
// Post: on success, returns the updated game state (CurrentPlayer advanced).
func (a *GameServer) Rest(ctx context.Context, req *pb.RestMessage) (*pb.RestResponse, error) {
	if req.GetGameId() == "" {
		return nil, status.Error(codes.InvalidArgument, "game_id is required")
	}

	game, err := a.service.Rest(req.GetGameId(), int(req.GetPlayer()), req.GetDiscards())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "rest: %v", err)
	}

	return &pb.RestResponse{GameState: domain.MarshallGameState(&game)}, nil
}
