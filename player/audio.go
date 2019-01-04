package player

import (
	"log"

	"github.com/jonas747/dca"
)

// Audio struct stores discordgo voice connection
type Audio struct {
	stream *dca.StreamingSession
}

// IsStreaming returns stream status as boolean
func (a *Audio) IsStreaming() bool {
	finished, err := a.stream.Finished()
	if err != nil {
		log.Println("player.IsStreaming:", err)
		return false
	}

	return finished
}
