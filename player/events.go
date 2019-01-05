package player

import "github.com/rylio/ytdl"

// Event has the event info
type Event struct {
	Type  event
	Track ytdl.VideoInfo
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
