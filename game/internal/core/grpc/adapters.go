package grpc

// The GameServer implementation (translating grpc requests to
// game_service.Service calls and back) does not live in this file/package.
//
// Why: this package holds the generated pb/grpc code, and
// internal/core/domain already imports it (domain/game.go et al. import
// this package as `pb` to build response messages via MarshallGame /
// MarshallGemCard / MarshallGolemCard). A GameServer implementation needs
// domain types too — at minimum to call domain.MarshallGame on whatever
// game_service.Service.Create returns. Importing internal/core/domain from
// this package would therefore close an import cycle
// (grpc -> domain -> grpc), which Go rejects outright; this isn't a style
// choice to revisit later, it's a hard compiler constraint.
//
// The actual implementation is internal/adapters/grpc (package
// grpcadapter), which sits above both this package and domain/game_service
// and depends on each one-way. cmd/server wires
// grpcadapter.NewGameServer(...) into this package's RegisterGameServer.
