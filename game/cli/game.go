package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
)

// rpcTimeout bounds every individual grpc call the CLI makes — long enough
// for a local or lightly-loaded remote server, short enough that a dead
// server doesn't hang the TUI forever.
const rpcTimeout = 10 * time.Second

// Game is the CLI's display/view model for a game, unmarshalled from the
// last GameState the server sent. id and currentPlayer are threaded back
// out on every subsequent request (game_id/player), which is the only way
// this client ever learns them — the server is the single source of truth
// for game state; the CLI never mutates Game locally.
type Game struct {
	id              string
	currentPlayer   int
	gem_deck_size   int
	gem_lineup      []GemCard
	golem_deck_size int
	golem_lineup    []GolemCard
	player          Player
	silver_coins    int
	gold_coins      int
}

// UnmarshallGame converts a wire GameState into the CLI's display type. A
// nil state yields the zero Game rather than panicking.
func UnmarshallGame(state *pb.GameState) Game {
	if state == nil {
		return Game{}
	}

	var gemLineup []GemCard
	for _, c := range state.GetGemLineup() {
		gemLineup = append(gemLineup, UnmarshallGemCard(c))
	}

	var golemLineup []GolemCard
	for _, c := range state.GetGolemLineup() {
		golemLineup = append(golemLineup, UnmarshallGolemCard(c))
	}

	return Game{
		id:              state.GetGameId(),
		currentPlayer:   int(state.GetCurrentPlayer()),
		gem_deck_size:   int(state.GetGemDeckSize()),
		gem_lineup:      gemLineup,
		golem_deck_size: int(state.GetGolemDeckSize()),
		golem_lineup:    golemLineup,
		player:          UnmarshallPlayer(state.GetPlayer()),
		silver_coins:    int(state.GetSilverCoins()),
		gold_coins:      int(state.GetGoldCoins()),
	}
}

func (g Game) String() string {
	// Coins
	silver := fmt.Sprintf("Silver: %v", g.silver_coins)
	gold := fmt.Sprintf("Gold: %v", g.gold_coins)
	coins := fmt.Sprintf("%v\t%v", silver, gold)

	// Golems
	golem_cards := fmt.Sprintf("[%v]\t", g.golem_deck_size)
	var golemStrs []string
	for i, card := range g.golem_lineup {
		golemStrs = append(golemStrs, fmt.Sprintf("%v:%v", i, card.String()))
	}
	golem_cards += strings.Join(golemStrs, " ")

	// Gem Cards
	gem_cards := fmt.Sprintf("[%v]\t", g.gem_deck_size)
	var gemStrs []string
	for i, card := range g.gem_lineup {
		gemStrs = append(gemStrs, fmt.Sprintf("%v:%v", i, card.String()))
	}
	gem_cards += strings.Join(gemStrs, " ")

	return fmt.Sprintf("%v\n%v\n%v\n\nplayer %v's turn\n%v\n",
		coins, golem_cards, gem_cards, g.currentPlayer, g.player)
}

// gameStateMsg carries a freshly fetched/updated Game back into the
// bubbletea Update loop after a network call completes.
type gameStateMsg struct {
	game Game
}

// gameErrMsg carries a network/rpc error back into the bubbletea Update
// loop — kept distinct from gameStateMsg so Update can render it without
// needing a Go error inside a struct value comparison.
type gameErrMsg struct {
	err error
}

// dial connects to the grpc server at addr and returns a Game client. addr
// comes from ServerAddr() (flag/env), which is how the CLI supports
// pointing at either a local or a deployed server.
func dial(addr string) (pb.GameClient, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial %s: %w", addr, err)
	}

	return pb.NewGameClient(conn), conn, nil
}

// createGame dials addr and starts a new game with playerCount players,
// returning a tea.Cmd that performs the call and reports its result.
//
// Pre: playerCount > 0 (server-validated; a non-positive count comes back
// as a gameErrMsg, not a crash).
func createGame(addr string, playerCount int) tea.Cmd {
	return func() tea.Msg {
		client, conn, err := dial(addr)
		if err != nil {
			return gameErrMsg{err}
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		resp, err := client.NewGame(ctx, &pb.CreateGameMessage{PlayerCount: int32(playerCount)})
		if err != nil {
			return gameErrMsg{fmt.Errorf("create game: %w", err)}
		}

		return gameStateMsg{UnmarshallGame(resp.GetGameState())}
	}
}

// playGemCard sends the play-gem-card RPC for game g's current player and
// card, returning a tea.Cmd that performs the call and reports its result.
//
// Pre: card is a card in g.player's hand (by value); g.currentPlayer is
// still up by the time the server processes this — a stale/duplicate
// submission comes back as a gameErrMsg, not a panic or silent no-op.
func playGemCard(addr string, g Game, card GemCard) tea.Cmd {
	return func() tea.Msg {
		client, conn, err := dial(addr)
		if err != nil {
			return gameErrMsg{err}
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		resp, err := client.PlayGemCard(ctx, &pb.PlayGemCardMessage{
			GameId: g.id,
			Player: int32(g.currentPlayer),
			Card:   MarshallGemCard(card),
		})
		if err != nil {
			return gameErrMsg{fmt.Errorf("play gem card: %w", err)}
		}

		return gameStateMsg{UnmarshallGame(resp.GetGameState())}
	}
}

// getGemCard sends the get-gem-card RPC for game g's current player,
// picking up card and paying payment (one gem letter y/g/b/p per lineup
// card bypassed) for cards bypassed to reach it.
func getGemCard(addr string, g Game, card GemCard, payment string) tea.Cmd {
	return func() tea.Msg {
		client, conn, err := dial(addr)
		if err != nil {
			return gameErrMsg{err}
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		resp, err := client.GetCard(ctx, &pb.GetCardMessage{
			GameId:  g.id,
			Player:  int32(g.currentPlayer),
			Card:    MarshallGemCard(card),
			Payment: payment,
		})
		if err != nil {
			return gameErrMsg{fmt.Errorf("get gem card: %w", err)}
		}

		return gameStateMsg{UnmarshallGame(resp.GetGameState())}
	}
}

// buyGolem sends the buy-golem RPC for game g's current player and golem.
func buyGolem(addr string, g Game, golem GolemCard) tea.Cmd {
	return func() tea.Msg {
		client, conn, err := dial(addr)
		if err != nil {
			return gameErrMsg{err}
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		resp, err := client.BuyGolem(ctx, &pb.BuyGolemMessage{
			GameId: g.id,
			Player: int32(g.currentPlayer),
			Golem:  MarshallGolemCard(golem),
		})
		if err != nil {
			return gameErrMsg{fmt.Errorf("buy golem: %w", err)}
		}

		return gameStateMsg{UnmarshallGame(resp.GetGameState())}
	}
}

// rest sends the rest RPC for game g's current player, ending their turn.
// discards (gem letters y/g/b/p) is only required — and validated — when
// the current player holds more than 10 gems; pass "" otherwise.
func rest(addr string, g Game, discards string) tea.Cmd {
	return func() tea.Msg {
		client, conn, err := dial(addr)
		if err != nil {
			return gameErrMsg{err}
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		resp, err := client.Rest(ctx, &pb.RestMessage{
			GameId:   g.id,
			Player:   int32(g.currentPlayer),
			Discards: discards,
		})
		if err != nil {
			return gameErrMsg{fmt.Errorf("rest: %w", err)}
		}

		return gameStateMsg{UnmarshallGame(resp.GetGameState())}
	}
}
