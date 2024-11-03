package main

import (
	"kyri56xcaesar/discord_bots_app/internal/server"
)

// This will be reconfigured to be able to run multiple servers
// Argument given
// Goroutines
// 
func main() {
	server, err := server.NewServer("server.env")
  if err == nil {
    server.Start()   
  }

}
