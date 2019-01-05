package player

import (
	"errors"
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"github.com/zrcni/go-discord-music-bot/queue"
)

// Track stores audio data and info
type Track struct {
	Audio *dca.EncodeSession
	Info  ytdl.VideoInfo
	// Discord channel ID where the track was queued from
	ChannelID string
}

// New returns a new player
func New() *Player {
	return &Player{
		queue:        queue.Queue{},
		EventChannel: make(chan Event),
	}
}

// Player handler audio playback
type Player struct {
	currentTrack Track
	stream       *dca.StreamingSession
	queue        queue.Queue
	EventChannel chan Event
}

// SetNowPlaying sets currently playing track
func (p *Player) SetNowPlaying(track Track) {
	e := Event{
		Type:      PLAY,
		TrackInfo: track.Info,
		Message:   track.Info.Title,
		ChannelID: track.ChannelID,
	}
	p.sendEvent(e)
	p.currentTrack = track
}

// GetNowPlaying gets currenly playing track
func (p *Player) GetNowPlaying() Track {
	return p.currentTrack
}

// Queue adds a track to the queue, returns ok to channel if track starts playing
func (p *Player) Queue(track Track, vc *discordgo.VoiceConnection) {
	if p.isStreaming() {
		p.queue.Add(track)
		log.Printf("\"%s\" added to queue", track.Info.Title)

		e := Event{
			Type:      QUEUE,
			TrackInfo: track.Info,
			Message:   track.Info.Title,
			ChannelID: track.ChannelID,
		}
		p.sendEvent(e)
		return
	}

	go p.play(track, vc)
}

// IsPlaying returns stream status as boolean
func (p *Player) IsPlaying() bool {
	if p.stream == nil {
		return false
	}
	return !p.stream.Paused()
}

// SetPaused sets stream's the pause state
func (p *Player) SetPaused(paused bool) {
	e := Event{
		TrackInfo: p.currentTrack.Info,
		Message:   p.currentTrack.Info.Title,
		ChannelID: p.currentTrack.ChannelID,
	}

	if paused {
		e.Type = PAUSE
		p.sendEvent(e)
	} else {
		e.Type = PLAY
		p.sendEvent(e)
	}

	if p.stream != nil {
		p.stream.SetPaused(paused)
	}
}

// ClearQueue clears the queue
func (p *Player) ClearQueue() {
	p.queue.Clear()
}

// processQueue removes the first item from the queue and returns it
func (p *Player) processQueue() (Track, error) {
	log.Printf("queueing next track")
	if p.queue.Length() == 0 {

		e := Event{
			Type:      STOP,
			ChannelID: p.currentTrack.ChannelID,
			Message:   "the queue is empty",
		}
		p.sendEvent(e)
		return Track{}, errors.New("the queue is empty")
	}

	track := p.queue.Shift()

	t, ok := track.(Track)
	if !ok {
		return Track{}, errors.New("track is not of type Track")
	}

	return t, nil
}

// play starts the process that streams the track
func (p *Player) play(track Track, vc *discordgo.VoiceConnection) {
	if !vc.Ready && p.IsPlaying() {
		return
	}

	p.SetNowPlaying(track)

	if err := p.startStream(track.Audio, vc); err != nil {
		log.Print(err)
	}

	track, err := p.processQueue()
	if err != nil {
		log.Print(err)
		return
	}

	go p.play(track, vc)
}

// startStream actually streams the audio to Discord
func (p *Player) startStream(source dca.OpusReader, vc *discordgo.VoiceConnection) error {
	done := make(chan error)
	p.stream = dca.NewStream(source, vc, done)

	log.Printf("Started streaming \"%s\"", p.currentTrack.Info.Title)

	err := <-done

	log.Printf("Stopped streaming \"%s\"", p.currentTrack.Info.Title)

	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		p.stream = nil
		return err
	}

	return nil
}

// IsStreaming returns stream status as boolean
func (p *Player) isStreaming() bool {
	if p.stream == nil {
		return false
	}

	finished, err := p.stream.Finished()
	if err != nil {
		log.Println("player.IsStreaming:", err)
		return false
	}

	return !finished
}
