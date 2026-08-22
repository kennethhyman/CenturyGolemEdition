package cli

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kennethhyman/CenturyGolemEdition/cli/actions"
)

// stage is which screen model.View renders and which keys model.Update
// routes.
type stage int

const (
	// stageConnecting: waiting on the initial CreateGame call.
	stageConnecting stage = iota
	// stageMenu: game loaded: showing state + the action menu.
	stageMenu
	// stageInput: an action was picked; collecting its input(s) via a
	// TurnView.
	stageInput
	// stageWaiting: an action's rpc is in flight.
	stageWaiting
	// stageError: the last rpc (or local validation) failed; err is set.
	stageError
)

// actionsChoices is the fixed action menu, in the order the user selects
// them by number (1-4), and how each maps to the TurnView that collects
// its argument(s).
var actionsChoices = []actions.GameAction{
	actions.PlayGemCard,
	actions.GetGemCard,
	actions.Rest,
	actions.GetGolemCard,
}

func newTurnView(a actions.GameAction) actions.TurnView {
	switch a {
	case actions.PlayGemCard:
		return actions.NewPlayCardView()
	case actions.GetGemCard:
		return actions.NewGemCardView()
	case actions.GetGolemCard:
		return actions.NewGolemCardView()
	case actions.Rest:
		return actions.NewRestView()
	}
	return actions.TurnView{}
}

// model is the CLI's bubbletea root: it owns the connection address, the
// last known Game state (server is the single source of truth — model
// never mutates game locally, only replaces it with what the server
// returns), and which screen is showing.
type model struct {
	addr    string
	players int

	game Game
	err  error

	stage stage
	input actions.TurnView
}

// InitialModel builds a model that will connect to addr and start a
// playerCount-player game as soon as the bubbletea program starts.
//
// Pre: playerCount > 0 (server-validated on the first call; a bad count
// surfaces as stageError, not a crash).
func InitialModel(addr string, playerCount int) model {
	return model{
		addr:    addr,
		players: playerCount,
		stage:   stageConnecting,
	}
}

func (m model) Init() tea.Cmd {
	return createGame(m.addr, m.players)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameStateMsg:
		m.game = msg.game
		m.err = nil
		m.stage = stageMenu
		return m, nil

	case gameErrMsg:
		m.err = msg.err
		m.stage = stageError
		return m, nil

	case actions.SubmitMsg:
		cmd, err := m.dispatch(msg)
		if err != nil {
			m.err = err
			m.stage = stageError
			return m, nil
		}
		m.stage = stageWaiting
		return m, cmd

	case tea.KeyMsg:
		switch m.stage {
		case stageMenu:
			return m.updateMenu(msg)
		case stageInput:
			return m.updateInput(msg)
		case stageError:
			return m.updateError(msg)
		default: // stageConnecting, stageWaiting: only allow quitting
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	i, err := strconv.Atoi(msg.String())
	if err != nil || i < 1 || i > len(actionsChoices) {
		return m, nil
	}

	m.input = newTurnView(actionsChoices[i-1])
	m.stage = stageInput

	return m, nil
}

func (m model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.stage = stageMenu
		return m, nil
	}
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	newInput, cmd := m.input.Update(msg)
	m.input = newInput.(actions.TurnView)
	return m, cmd
}

func (m model) updateError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	// A game we already have keeps playing after acknowledging the error;
	// a connect failure (no game yet) has nothing to return to.
	if m.game.id == "" {
		return m, tea.Quit
	}
	m.stage = stageMenu
	m.err = nil
	return m, nil
}

// dispatch turns a validated SubmitMsg into the tea.Cmd that performs its
// rpc, resolving the typed index against the current game state (the
// bounds-check a TurnView's own local validation can't do, since it has no
// access to cli.Game).
//
// Pre: msg.Values were already accepted by the originating TurnView's
// local (format-only) validator.
// Post: on success, the returned Cmd performs exactly one rpc call. On
// error, the index/value was out of range for the current game state.
func (m model) dispatch(msg actions.SubmitMsg) (tea.Cmd, error) {
	fields := msg.Values

	switch msg.Action {
	case actions.PlayGemCard:
		idx, err := strconv.Atoi(fields[0])
		if err != nil || idx < 0 || idx >= len(m.game.player.hand) {
			return nil, fmt.Errorf("play gem card: %q is not a valid hand index (0-%d)", fields[0], len(m.game.player.hand)-1)
		}
		return playGemCard(m.addr, m.game, m.game.player.hand[idx]), nil

	case actions.GetGemCard:
		idx, err := strconv.Atoi(fields[0])
		if err != nil || idx < 0 || idx >= len(m.game.gem_lineup) {
			return nil, fmt.Errorf("get gem card: %q is not a valid lineup index (0-%d)", fields[0], len(m.game.gem_lineup)-1)
		}
		payment := ""
		if len(fields) > 1 {
			payment = fields[1]
		}
		return getGemCard(m.addr, m.game, m.game.gem_lineup[idx], payment), nil

	case actions.GetGolemCard:
		idx, err := strconv.Atoi(fields[0])
		if err != nil || idx < 0 || idx >= len(m.game.golem_lineup) {
			return nil, fmt.Errorf("buy golem: %q is not a valid lineup index (0-%d)", fields[0], len(m.game.golem_lineup)-1)
		}
		return buyGolem(m.addr, m.game, m.game.golem_lineup[idx]), nil

	case actions.Rest:
		discards := ""
		if len(fields) > 0 {
			discards = fields[0]
		}
		return rest(m.addr, m.game, discards), nil
	}

	return nil, fmt.Errorf("unknown action %q", msg.Action)
}

func (m model) View() string {
	switch m.stage {
	case stageConnecting:
		return fmt.Sprintf("connecting to %s...\n", m.addr)
	case stageWaiting:
		return "working...\n"
	case stageError:
		hint := "press any key to continue"
		if m.game.id == "" {
			hint = "press any key to quit"
		}
		return fmt.Sprintf("error: %v\n%s\n", m.err, hint)
	case stageInput:
		return fmt.Sprintf("%v\n%v\n(enter to submit, esc to cancel)\n", m.game, m.input.View())
	default: // stageMenu
		var choices []string
		for i, a := range actionsChoices {
			choices = append(choices, fmt.Sprintf("%d:%v", i+1, a.View()))
		}
		return fmt.Sprintf("%v\n%v\n", m.game, strings.Join(choices, "  "))
	}
}

