package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// BotToken is the bot's authentication token
	BotToken string
	// ClientID is the client ID
	ClientID string
)

// SetupEnv Sets environment variables from .env file
func SetupEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BotToken = os.Getenv("BOT_TOKEN")
	ClientID = os.Getenv("CLIENT_ID")
}
