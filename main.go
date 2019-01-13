package main

import (
	"github.com/zrcni/go-discord-music-bot/bot"
	"github.com/zrcni/go-discord-music-bot/config"
	"github.com/zrcni/go-discord-music-bot/logger"
)

func init() {
	config.Setup()
	logger.Setup()
}

func main() {
	bot.Start()
}
