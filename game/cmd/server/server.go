package main

import (
	"fmt"
  "github.com/kennethhyman/CenturyGolemEdition/server"
)

func main() {
  fmt.Printf("server mode: Listening on port 50051\n")
  server.StartServer()
}
