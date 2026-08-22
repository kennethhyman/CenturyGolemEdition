// Command server is the composition root: it picks the ports.GameRepository
// implementation (sqlite today), builds game_service.Service around it,
// wraps that in the grpc adapter, and serves the Game grpc service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	grpcadapter "github.com/kennethhyman/CenturyGolemEdition/internal/adapters/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/adapters/sqlite"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
	"github.com/kennethhyman/CenturyGolemEdition/internal/core/services/game_service"
)

const defaultAddr = ":50051"
const defaultDBPath = "games.db"

func main() {
	addr := flag.String("addr", envOr("GAME_SERVER_ADDR", defaultAddr), "address to listen on (env GAME_SERVER_ADDR)")
	dbPath := flag.String("db", envOr("GAME_SERVER_DB", defaultDBPath), "path to the sqlite database file (env GAME_SERVER_DB)")
	flag.Parse()

	lis, srv, repo, err := build(*addr, *dbPath)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	defer repo.Close()

	fmt.Printf("server mode: listening on %s (sqlite db: %s)\n", *addr, *dbPath)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("server: serve: %v", err)
	}
}

// build wires the composition root's dependency graph and binds addr,
// without starting to serve — split out from main so tests can drive a real
// grpc.Server over a real net.Listener without forking a process.
//
// Post: on success, the returned listener is already bound to addr, and the
// returned *sql-backed repository has been opened (caller owns Close), and
// srv is registered and ready for Serve(lis).
func build(addr, dbPath string) (net.Listener, *grpc.Server, *sqlite.Repository, error) {
	repo, err := sqlite.New(dbPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("open sqlite repository %q: %w", dbPath, err)
	}

	service := game_service.New(repo)
	adapter := grpcadapter.NewGameServer(service)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		repo.Close()
		return nil, nil, nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGameServer(grpcServer, adapter)

	return lis, grpcServer, repo, nil
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
