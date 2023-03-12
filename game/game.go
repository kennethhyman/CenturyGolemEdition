package main

import (
	"fmt"
	//tea "github.com/charmbracelet/bubbletea"
  "github.com/kennethhyman/CenturyGolemEdition/cli"
  "flag"
  "github.com/kennethhyman/CenturyGolemEdition/server"
	//"os"
)

//func main() {
//	// Create the Bubble tea application
//	if err := tea.NewProgram(cli.InitialModel()).Start(); err != nil {
//		fmt.Printf("could not start program: %s\n", err)
//		os.Exit(1)
//	}
//}


var (
  	serverMode        = flag.Bool("s", false, "Connection uses TLS if true, else plain TCP")
)

func main() {

  flag.Parse()

  if (*serverMode) {
    fmt.Printf("server mode: Listening on port 50051\n")
    server.StartServer()
  } else {
    game := cli.NewGame()
    fmt.Printf("%v\n", game)
    //if err := tea.NewProgram(cli.InitialModel()).Start(); err != nil {
    //  fmt.Printf("could not start program: %s\n", err)
    //  os.Exit(1)
    //}
  }
}
