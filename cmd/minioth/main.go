package main

import (
	"log"

	"kyri56xcaesar/discord_bots_app/internal/minioth"
)

func main() {
	haveDb := false
	mauth := minioth.NewMinioth("root", haveDb, "minioth.db")
	err := mauth.Init()
	if err != nil {
		log.Print(err)
	}
	mauth.Useradd("koularos", "", "", "", "")
	mauth.Useradd("patatas", "", "", "", "")
	mauth.Useradd("koularos", "", "", "", "")
	mauth.Useradd("ntomatas", "", "", "", "")
	mauth.Useradd("rengar", "j4ngl3r", "JGKing", "/home/rengar", "/bin/gshell")
	mauth.Useradd("diego", "", "", "", "")
	mauth.Useradd("viego", "", "", "", "")
	mauth.Useradd("Jarvan IV", "", "", "", "")

	mauth.Userdel("viego")
	mauth.Userdel("ntomatas")
	mauth.Useradd("Katsaplias", "", "OKatsapleas", "", "")

	m := minioth.NewMSerivce(&mauth)
	m.ServeHTTP()
}
