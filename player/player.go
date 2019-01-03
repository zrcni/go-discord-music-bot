package player

import (
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// Track stores audio data and info
type Track struct {
	Audio *dca.EncodeSession
	Info  *ytdl.VideoInfo
}
