package main

import (
	"kyri56xcaesar/discord_bots_app/internal/server"
)

func main() {
	server := server.NewServer("server.env")
	server.Start()
}
