// Command cli is the interactive terminal client: it connects to a running
// Game grpc server (local or deployed — see -addr/GAME_SERVER_ADDR, which
// mirrors cmd/server's own -addr/GAME_SERVER_ADDR convention) and drives a
// game over the network via the bubbletea program in package cli.
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kennethhyman/CenturyGolemEdition/cli"
)

const defaultAddr = "localhost:50051"

func main() {
	addr := flag.String("addr", envOr("GAME_SERVER_ADDR", defaultAddr), "address of the game server to connect to (env GAME_SERVER_ADDR)")
	players := flag.Int("players", 2, "number of players in the game to create")
	flag.Parse()

	if err := tea.NewProgram(cli.InitialModel(*addr, *players)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
