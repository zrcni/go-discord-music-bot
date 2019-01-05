package main

import (
	"github.com/zrcni/go-discord-music-bot/bot"
	"github.com/zrcni/go-discord-music-bot/config"
)

func init() {
	config.Setup()
}

func main() {
	// spotify.Authenticate()
	bot.Start()
}
