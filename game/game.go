package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kennethhyman/CenturyGolemEdition/cli"
	"os"
)

func main() {
	// Create the Bubble tea application
	if err := tea.NewProgram(cli.InitialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
