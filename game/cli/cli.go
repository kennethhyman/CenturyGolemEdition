package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli/actions"
  "context"
  "time"
	game_server "github.com/kennethhyman/CenturyGolemEdition/grpc"
	pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
)

type model struct {
  game Game
  client *pb.GameClient
}

var actionsChoices = []tea.Model{turns.PlayGemCard, turns.GetGemCard, turns.Rest, turns.GetGolemCard}

func InitialModel() model {
  client := CreateGameClient()

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  msg := &game_server.CreateGameMessage{
    PlayerCount: 3,
  }

  game, _ := client.NewGame(ctx, msg)

	return model {
    game: UnmarshallGame(game),
    client: &client,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func CardToString(card GemCard) string {
  return fmt.Sprintf("%vy%vg%vb%vp", card.In.Yellow, card.In.Green, card.In.Blue, card.In.Pink)
}

func (m model) View() string {
  return fmt.Sprintf("%v\n", m.game)
}
