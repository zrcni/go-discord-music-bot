package spotify

import (
	"log"

	"github.com/zmb3/spotify"
)

var playlistAPIURL = "https://api.spotify.com/v1/playlists"

type Client struct {
	client *spotify.Client
}

// DoAction SEPARATE THIS WHEN YOU KNOW WHAT YOU'RE DOING
func (c *Client) DoAction(action string) {
	var playerState *spotify.PlayerState
	var err error

	switch action {
	case "play":
		err = c.client.Play()
	case "pause":
		err = c.client.Pause()
	case "next":
		err = c.client.Next()
	case "previous":
		err = c.client.Previous()
	case "shuffle":
		playerState.ShuffleState = !playerState.ShuffleState
		err = c.client.Shuffle(playerState.ShuffleState)
	}

	if err != nil {
		log.Print(err)
	}
}

func (c *Client) GetPlaylists(searchTerm string) []string {
	// isPlaylist := isPlaylistURI(searchTerm)
	// if isPlaylist {
	// }
	res, err := c.client.Search(searchTerm, spotify.SearchTypePlaylist)
	if err != nil {
		log.Printf("Search playlist: %v", err)
		return []string{}
	}

	var playlistNames []string

	for _, playlist := range res.Playlists.Playlists {
		log.Printf("%+v\n", playlist.Name)
		playlistNames = append(playlistNames, playlist.Name)
	}

	return playlistNames
}
