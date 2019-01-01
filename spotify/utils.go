package spotify

import (
	"strings"
)

func isPlaylistURI(spotifyURI string) bool {
	uriParts := strings.Split(spotifyURI, ":")

	if len(uriParts) != 5 {
		return false
	}

	playlistID := uriParts[4]

	isNotEmptyString := len(playlistID) > 0
	return isNotEmptyString
}

func parseSpotifyPlaylistURI(spotifyURI string) (string, string) {
	uriParts := strings.Split(spotifyURI, ":")

	if len(uriParts) != 5 {
		return "", ""
	}

	user := uriParts[2]
	playlistID := uriParts[4]

	return user, playlistID
}
