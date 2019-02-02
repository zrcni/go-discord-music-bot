package player

import (
	log "github.com/sirupsen/logrus"
)

// Event has the event info
type Event struct {
	Type      event
	Track     *Track
	Message   string
	ChannelID string
}

type event int

const (
	// PLAY - track started playing
	PLAY event = iota
	// QUEUE - track was queued
	QUEUE
	// UNPAUSE - track was unpaused
	UNPAUSE
	// PAUSE - track was paused
	PAUSE
	// STOP - streaming stopped
	STOP
	// ERROR means an error occurred
	ERROR
)

// sendEvent sends and event to the event channel
func (p *Player) sendEvent(e Event) {
	p.Events <- e
	logEvent(e)
}

func logEvent(e Event) {
	var eventName string

	switch e.Type {
	case PLAY:
		eventName = "PLAY"
	case QUEUE:
		eventName = "QUEUE"
	case UNPAUSE:
		eventName = "UNPAUSE"
	case PAUSE:
		eventName = "PAUSE"
	case STOP:
		eventName = "STOP"
	case ERROR:
		eventName = "ERROR"
	default:
		return
	}

	log.Debugf("%s event: [Message: \"%v\"]", eventName, e.Message)
}
