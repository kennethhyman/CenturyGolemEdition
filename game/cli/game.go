package cli

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Game struct {
  gem_deck_size int
  gem_lineup []GemCard
  golem_deck_size int
  golem_lineup []GolemCard
  player Player
  silver_coins int
  gold_coins int
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
  gold := fmt.Sprintf("%vGold: %v", strings.Repeat(" ", 17 - len(silver)), g.gold_coins)
  coins := fmt.Sprintf("\t%v%v%v", padding, silver, gold)

  // Golems
  golem_cards := fmt.Sprintf("[%v]\t", g.golem_deck_size)
  for _, card := range(g.golem_lineup) {
    card_string := card.String()
    padding := strings.Repeat(" ", 15 - card.StringLength())
    golem_cards += card_string + padding
  }

  // Gem Cards
  gem_cards := fmt.Sprintf("[%v]\t", g.gem_deck_size)
  for _, card := range(g.gem_lineup) {
    card_string := card.String()
    padding := strings.Repeat(" ", 15 - card.StringLength())
    gem_cards += card_string + padding
  }

  return fmt.Sprintf("%v\n%v\n%v\n\n\n", coins, golem_cards, gem_cards)
}

func UnmarshallGame(game *pb.CreateGameResponse) Game {
  var gem_lineup []GemCard
  var golem_lineup []GolemCard

  for _, card := range(game.GameState.GemLineup) {
    gem_lineup = append(gem_lineup, UnmarshallGemCard(card))
  }

  for _, card := range(game.GameState.GolemLineup) {
    golem_lineup = append(golem_lineup, UnmarshallGolemCard(card))
  }

  return Game {
    gem_lineup: gem_lineup,
    gem_deck_size: int(game.GameState.GemDeckSize),
    golem_lineup: golem_lineup,
    golem_deck_size: int(game.GameState.GolemDeckSize),
    gold_coins: int(game.GameState.GoldCoins),
    silver_coins: int(game.GameState.SilverCoins),
  }
}
