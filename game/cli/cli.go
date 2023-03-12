package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli/turns"
  "context"
  "time"
	game_server "github.com/kennethhyman/CenturyGolemEdition/grpc"

	//. "github.com/kennethhyman/CenturyGolemEdition/models"
	pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type model struct {
  game *pb.CreateGameResponse
  client *pb.GameClient
}


var actionsChoices = []Viewable{turns.PlayGemCard, turns.GetGemCard, turns.Rest, turns.GetGolemCard}

func InitialModel() model {
  client := connectToClient("localhost:50051")

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  msg := &game_server.CreateGameMessage{
    PlayerCount: 3,
  }

  game, _ := (*client).NewGame(ctx, msg)
  fmt.Printf("game: %v\n", game)

	//views := []turns.TurnView{turns.NewPlayCardView(game), turns.NewGemCardView(game), turns.NewRestView(game), turns.NewGolemCardView(game)}
	return model {
    game: game,
    client: client,
	}
}

func (m model) Init() tea.Cmd {
  fmt.Printf("init success\n")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  fmt.Printf("Updated\n")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	// pass the message along (and cast to view selector)
	return m, nil
}

func CardToString(card *pb.GemCard) string {
  card_string := "CARD" //fmt.Sprintf("%vy%vg%vb%vp", card.In.Yellow, card.In.Green, card.In.Blue, card.In.Pink)
  fmt.Printf("%v\n", card_string)

  return card_string
}

func (m model) View() string {
  var view string
  fmt.Printf("%v\n", m.game)
  for _, card := range(m.game.GameState.GemLineup) {
    if card != nil {
      view = fmt.Sprintf("%v\t", CardToString(card))
    }
  }

	return fmt.Sprintf("%v\n", view)
}

func connectToClient(server string) *pb.GameClient {

  var opts []grpc.DialOption
  opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

  fmt.Printf("client mode\n")

  conn, err := grpc.Dial("localhost:50051", opts...)
  fmt.Printf("connecting to server\n")
  if err != nil {
    log.Fatalf("fail to dial: %v", err)
  }

  defer conn.Close()
  client := pb.NewGameClient(conn)
  fmt.Print("connected successfully\n");
  return &client
}
