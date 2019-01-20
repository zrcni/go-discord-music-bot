package spotify

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"github.com/zrcni/go-discord-music-bot/config"
	"golang.org/x/oauth2/clientcredentials"
)

func NewClient() (*Client, error) {
	config := &clientcredentials.Config{
		ClientID:     config.Config.Spotify.ID,
		ClientSecret: config.Config.Spotify.Secret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Errorf("couldn't get token: %v", err)
		return nil, err
	}

	client := spotify.Authenticator{}.NewClient(token)

	return &Client{
		client: &client,
	}, nil
}
