package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zrcni/go-discord-music-bot/utils"
	"gopkg.in/yaml.v2"
)

var configFileName = "config"

// Config has all the config
var Config config

type config struct {
	Debug    bool `yaml:"debug"`
	BotToken string
	Spotify  spotify
}

type spotify struct {
	ID     string
	Secret string
}

// Setup Sets environment variables from .env file
func Setup() {
	parseDotEnv()
	populateConfigWithEnv()
	parseYAML()
}

func parseYAML() {
	basePath, err := utils.GetBasePath()
	if err != nil {
		panic(err)
	}

	configFile := fmt.Sprintf("%s.yaml", configFileName)
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", basePath, configFile))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
}

func parseDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Sets env variables to Config struct
func populateConfigWithEnv() {
	Config.BotToken = os.Getenv("BOT_TOKEN")
	Config.Spotify = spotify{
		ID:     os.Getenv("SPOTIFY_ID"),
		Secret: os.Getenv("SPOTIFY_SECRET"),
	}
}
