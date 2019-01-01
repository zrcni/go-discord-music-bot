package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zrcni/go-discord-music-bot/bot"
)

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("setupEnv: %v", err)
	}
}

func init() {
	setupEnv()
}

func main() {
	// spotify.Authenticate()
	bot.Start()
}
