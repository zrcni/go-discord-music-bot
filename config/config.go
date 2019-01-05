package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Setup Sets environment variables from .env file
func Setup() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
