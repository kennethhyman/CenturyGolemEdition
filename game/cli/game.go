package cli

import (
	//"github.com/kennethhyman/CenturyGolemEdition/server"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const yellowString = "\033[1;33m○\033[0m"
const greenString = "\033[2;32m○"
const blueString = "\033[2;34m○"
const pinkString = "\033[38;5;206m○"
const upgradeString = "🌟"
const yieldsString = "\033[0m->"

type GemValue struct {
  Yellow int
  Green int
  Blue int
  Pink int
}

type GemCard struct {
  In GemValue
  Out GemValue
  Upgrades int
}

type Game struct {
  gem_lineup []GemCard
}

func (g GemValue) String() string {
	output := ""
	output += strings.Repeat(yellowString, g.Yellow)
	output += strings.Repeat(greenString, g.Green)
	output += strings.Repeat(blueString, g.Blue)
	output += strings.Repeat(pinkString, g.Pink)

	return output
}
func (g GemCard) String() string {
  if g.Upgrades > 0 {
		return "\033[0m[ " + strings.Repeat(upgradeString, g.Upgrades) + "\033[0m]"
	}

	return "\033[0m[ " + g.In.String() + yieldsString + g.Out.String() + "\033[0m ]"
}
func (g Game) String() string {
  output := ""

  for i, card := range(g.gem_lineup) {
    output += card.String()
    if i != len(g.gem_lineup) {
      output += "\t"
    }
  }

  return output
}



var client pb.GameClient

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

func UnmarshallGemCard(card *pb.GemCard) GemCard {
  return GemCard{
    In: GemValue {
      Yellow: int(card.In.Yellow),
      Green: int(card.In.Green),
      Blue: int(card.In.Blue),
      Pink: int(card.In.Pink),
    },
    Out: GemValue {
      Yellow: int(card.Out.Yellow),
      Green: int(card.Out.Green),
      Blue: int(card.Out.Blue),
      Pink: int(card.Out.Pink),
    },
    Upgrades: int(card.Upgrades),
  }
}

func UnmarshallGame(game *pb.CreateGameResponse) Game {
  var gem_lineup []GemCard

  for _, card := range(game.GameState.GemLineup) {
    gem_lineup = append(gem_lineup, UnmarshallGemCard(card))
  }

  return Game {
    gem_lineup: gem_lineup,
  }
}

func NewGame() Game {
  if client == nil {
    client = CreateGameClient()
  }
  
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  msg := pb.CreateGameMessage{
    PlayerCount: 3,
  }

  game, _ := client.NewGame(ctx, &msg)
  
  return UnmarshallGame(game)
}
