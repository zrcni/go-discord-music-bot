package player

import (
	"log"

	"github.com/rylio/ytdl"
)

// Event has the event info
type Event struct {
	Type    event
	Track   ytdl.VideoInfo
	Message string
}

type event int

const (
	// PLAY - track started playing
	PLAY event = iota
	// PAUSE - track was paused
	PAUSE
	// STOP - streaming stopped
	STOP
	// ERROR means an error occurred
	ERROR
)

// sendEvent sends and event to the event channel
func (p *Player) sendEvent(e Event) {
	p.EventChannel <- e
	logEvent(e)
}

func logEvent(e Event) {
	var eventName string

	switch e.Type {
	case PLAY:
		eventName = "PLAY"
	case PAUSE:
		eventName = "PAUSE"
	case STOP:
		eventName = "STOP"
	case ERROR:
		eventName = "ERROR"
	default:
		return
	}

	log.Printf("%s event: [Message: \"%v\"]", eventName, e.Message)
}
