package cli

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/kennethhyman/CenturyGolemEdition/internal/core/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Game struct {
	gem_deck_size   int
	gem_lineup      []GemCard
	golem_deck_size int
	golem_lineup    []GolemCard
	player          Player
	silver_coins    int
	gold_coins      int
}

func UnmarshallGame(resp *pb.CreateGameResponse) Game {
	gs := resp.GameState
	var gemLineup []GemCard
	var golemLineup []GolemCard

	for _, c := range gs.GemLineup {
		gemLineup = append(gemLineup, UnmarshallGemCard(c))
	}

	for _, c := range gs.GolemLineup {
		golemLineup = append(golemLineup, UnmarshallGolemCard(c))
	}

	return Game{
		gem_deck_size:   int(gs.GemDeckSize),
		gem_lineup:      gemLineup,
		golem_deck_size: int(gs.GolemDeckSize),
		golem_lineup:    golemLineup,
		player:          UnmarshallPlayer(gs.Player),
		silver_coins:    int(gs.SilverCoins),
		gold_coins:      int(gs.GoldCoins),
	}
}

func createGame() tea.Msg {
	client := CreateGameClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &pb.CreateGameMessage{
		PlayerCount: 3,
	}

	game, _ := client.NewGame(ctx, msg)

	return UnmarshallGame(game)
}

func CreateGameClient() pb.GameClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	fmt.Printf("client mode\n")

	conn, err := grpc.Dial("localhost:50051", opts...)
	fmt.Printf("connecting to server\n")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return pb.NewGameClient(conn)
}

func NewGame() Game {
	client := CreateGameClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := pb.CreateGameMessage{
		PlayerCount: 3,
	}

	game, _ := client.NewGame(ctx, &msg)

	return UnmarshallGame(game)
}

func (g Game) String() string {
	// Coins
	padding := strings.Repeat(" ", 51)
	silver := fmt.Sprintf("Silver: %v", g.silver_coins)
	gold := fmt.Sprintf("%vGold: %v", strings.Repeat(" ", 17-len(silver)), g.gold_coins)
	coins := fmt.Sprintf("\t%v%v%v", padding, silver, gold)

	// Golems
	golem_cards := fmt.Sprintf("[%v]\t", g.golem_deck_size)
	for _, card := range g.golem_lineup {
		card_string := card.String()
		padding := strings.Repeat(" ", 15-card.StringLength())
		golem_cards += card_string + padding
	}

	// Gem Cards
	gem_cards := fmt.Sprintf("[%v]\t", g.gem_deck_size)
	for _, card := range g.gem_lineup {
		card_string := card.String()
		padding := strings.Repeat(" ", 15-card.StringLength())
		gem_cards += card_string + padding
	}

	return fmt.Sprintf("%v\n%v\n%v\n\n\n\n%v\n", coins, golem_cards, gem_cards, g.player)
}
