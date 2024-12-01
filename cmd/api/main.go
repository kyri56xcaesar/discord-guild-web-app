package main

import (
	"kyri56xcaesar/discord-guild-web-app/internal/server"
)

// This will be reconfigured to be able to run multiple servers
// Argument given
// Goroutines
func main() {
	server, err := server.NewServer("configs/server.env")
	if err == nil {
		server.Start()
	}
}
