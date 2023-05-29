package main

import (
	"os"
	"os/signal"
	"syscall"

	"phoenixDiscordBot/api"
	"phoenixDiscordBot/discord"
)

var (
	Token string
)

func init() {
	Token = os.Getenv("DISCORD-TOKEN")
}

func main() {
	if err := discord.StartDiscord(Token); err != nil {
		return
	}

	api.StartApi()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Shutdown()
}
