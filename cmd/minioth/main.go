package main

import (
	"log"

	"kyri56xcaesar/discord_bots_app/internal/minioth"
)

func main() {
	mauth := minioth.NewMinioth("root", true, "minioth.db")
	err := mauth.Init()
	if err != nil {
		log.Print(err)
	}
}
