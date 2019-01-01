package main

import (
	"github.com/joho/godotenv" // "github.com/zrcni/go-discord-music-bot/youtube"
	"github.com/zrcni/go-discord-music-bot/bot"
	"log"
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
	// youtube.GetVideoByURL(youtube.PlaceholderVideoURL)
	// spotify.Authenticate()
	bot.Start()
}
