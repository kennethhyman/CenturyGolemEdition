package main

import (
	"fmt"
	//tea "github.com/charmbracelet/bubbletea"
  //"github.com/kennethhyman/CenturyGolemEdition/cli"
  "flag"
  "github.com/kennethhyman/CenturyGolemEdition/server"
  "google.golang.org/grpc"
  "log"
  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
  "time"
  "context"
	_ "os"
  "google.golang.org/grpc/credentials/insecure"
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
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    fmt.Printf("client mode\n")

    conn, err := grpc.Dial("localhost:50051", opts...)

    if err != nil {
      log.Fatalf("fail to dial: %v", err)
    }

    defer conn.Close()
    client := pb.NewGameClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    card, err := client.GetCard(ctx, &pb.GetCardMessage{})

    if err != nil {
      log.Fatalf("client.GetCard failed: %v", err)
    }

    fmt.Printf("card: %v\n", card)
  }
}
