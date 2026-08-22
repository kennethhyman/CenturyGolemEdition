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

// GameService is the subset of game_service.Service this adapter calls.
// Declared here, at the point of use, rather than depending on
// ports.GameService directly: ports.GameService's fuller contract
// (PlayGemCard/GetGemCard/GetGolemCard/Rest) isn't implemented by
// game_service.Service yet, so requiring it here would make this adapter
// fail to compile against the real Service for no benefit to what this
// adapter actually calls today.
type GameService interface {
	Create(players int) (domain.Game, error)
	Get(id string) (domain.Game, error)
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

// GetCard and PlayGemCard are intentionally left Unimplemented (via the
// embedded pb.UnimplementedGameServer default) rather than given a
// best-effort implementation here:
//
//   - GetGardMessage carries only an int32 game_id, but domain game IDs are
//     UUID strings (see game_service.GetUID) — there's no meaningful way to
//     look up a real game from it, and GetCard returns a single GemCard
//     rather than game state, which no GameService method produces.
//   - PlayGemCardMessage carries no player index, and GameService (as of
//     Stage 2) exposes no PlayGemCard method for this adapter to call —
//     game_service.Service only implements Create/Get so far.
//
// Wiring these up is follow-on work once game_service.Service grows the
// matching methods (and, for GetCard, once the .proto itself carries enough
// information to serve it) — tracked in PLAN.md rather than guessed at here.
