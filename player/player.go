package player

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/queue"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
)

// Track stores audio data and info
type Track struct {
	Audio        *dca.EncodeSession
	Title        string
	Duration     time.Duration
	ID           string
	ThumbnailURL string
	URL          string
	Filename     string
	// Discord text channel ID where the track was queued from
	ChannelID string
}

// New returns a new player
func New() Player {
	return Player{
		queue:  queue.New(20),
		Events: make(chan Event),
	}
}

// Player handler audio playback
type Player struct {
	currentTrack    Track
	stream          *dca.StreamingSession
	VoiceConnection *discordgo.VoiceConnection
	queue           queue.Queue
	Events          chan Event
}

// SetNowPlaying sets currently playing track
func (p *Player) SetNowPlaying(track Track) {
	e := Event{
		Type:      PLAY,
		Track:     track,
		Message:   track.Title,
		ChannelID: track.ChannelID,
	}
	p.sendEvent(e)
	p.currentTrack = track
}

// GetNowPlaying gets currenly playing track
func (p *Player) GetNowPlaying() Track {
	return p.currentTrack
}

// IsPlaying returns stream status as boolean
func (p *Player) IsPlaying() bool {
	if p.stream == nil {
		return false
	}
	return !p.stream.Paused()
}

// Stop stops streaming and clears the queue
func (p *Player) Stop() {
	e := Event{
		Type:      STOP,
		ChannelID: p.currentTrack.ChannelID,
		Message:   "Stopped playing",
	}

	p.sendEvent(e)

	if p.stream != nil {
		p.stream.SetPaused(true)
	}
	p.ClearQueue()
	p.currentTrack = Track{}
}

// SetPaused sets stream's the pause state
func (p *Player) SetPaused(paused bool) {
	e := Event{
		Track:     p.currentTrack,
		Message:   p.currentTrack.Title,
		ChannelID: p.currentTrack.ChannelID,
	}

	if paused {
		e.Type = PAUSE
		p.sendEvent(e)
	} else {
		e.Type = UNPAUSE
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

// Queue adds a track to the queue, returns ok to channel if track starts playing
func (p *Player) Queue(track Track) {
	err := p.queue.Add(track)
	if err != nil {
		e := Event{
			Type:      ERROR,
			Track:     track,
			Message:   err.Error(),
			ChannelID: track.ChannelID,
		}
		p.sendEvent(e)
		log.Error(err)
		return
	}

	log.Infof("\"%s\" added to queue", track.Title)

	e := Event{
		Type:      QUEUE,
		Track:     track,
		Message:   track.Title,
		ChannelID: track.ChannelID,
	}
	p.sendEvent(e)

	p.processQueue()
	// if p.isStreaming() {
	// err := p.queue.Add(track)
	// if err != nil {
	// 	e := Event{
	// 		Type:      ERROR,
	// 		Track:     track,
	// 		Message:   err.Error(),
	// 		ChannelID: track.ChannelID,
	// 	}
	// 	p.sendEvent(e)
	// 	log.Error(err)
	// 	return
	// }

	// log.Infof("\"%s\" added to queue", track.Title)

	// e := Event{
	// 	Type:      QUEUE,
	// 	Track:     track,
	// 	Message:   track.Title,
	// 	ChannelID: track.ChannelID,
	// }
	// 	// p.sendEvent(e)
	// 	return
	// }

	// go p.play(track, vc)
}

// processQueue removes the first item from the queue and returns it
func (p *Player) processQueue() {
	log.Debug("queueing next track")
	if p.queue.Length() == 0 {
		e := Event{
			Type:      STOP,
			ChannelID: p.currentTrack.ChannelID,
			Message:   "Stopped playing - the queue is empty",
		}
		p.sendEvent(e)
		p.currentTrack = Track{}
		log.Debug("the queue is empty")
		return
	}

	if p.isStreaming() {
		log.Debug("already streaming")
		return
	}

	track := p.queue.Shift()

	t, ok := track.(Track)
	if !ok {
		panic(fmt.Sprintf("track is not of type Track: %+v", t))
	}

	if p.queue.Length() > 1 {
		go func(p *Player) {
			err := p.prepareNextTrack()
			if err != nil {
				log.Error(err)
			}
		}(p)
	}

	go p.play(t)
}

// play starts the process that streams the track
func (p *Player) play(track Track) {
	if !p.VoiceConnection.Ready {
		log.Debug("voice connection is not ready")
		return
	}
	// if p.IsPlaying() {
	// 	log.Debug("already playing")
	// 	return
	// }

	if track.Audio != nil {
		p.SetNowPlaying(track)

		err := p.startStream(track.Audio)
		if err != nil {
			log.Error(err)
		}
	} else {
		// TODO: fix audio sometimes being nil.
		// It happens when replacing track.
		log.Debugf("track.Audio is nil in: %v", track.Title)

		ok := p.queue.DeleteAt(0)
		if !ok {
			log.Debugf("could not delete: %v", track.Title)
			return
		}

		log.Debugf("deleted: %v", track.Title)
	}

	p.processQueue()
}

// startStream actually streams the audio to Discord
func (p *Player) startStream(source dca.OpusReader) error {
	done := make(chan error)
	p.stream = dca.NewStream(source, p.VoiceConnection, done)

	log.Infof("Started streaming \"%s\"", p.currentTrack.Title)

	err := <-done

	log.Infof("Stopped streaming \"%s\"", p.currentTrack.Title)

	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		// p.stream = nil
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
		log.Errorf("player.IsStreaming:", err)
		return false
	}

	return !finished
}

func (p *Player) prepareNextTrack() error {
	log.Debug("preparing next track")

	t, err := p.queue.GetAt(0)
	if err != nil {
		return errors.New("queue is empty, can't prepare next track")
	}

	track, ok := t.(Track)
	if !ok {
		panic(fmt.Sprintf("t is not of type Track: %+v", t))
	}

	// If track doesn't have audio (it's taken from queue)
	// probably dont need to check this because it's from the queue anyway?
	if track.Audio != nil {
		return errors.New("track already has audio data")
	}

	es, err := videoaudio.ReadAudioFile(track.Filename)
	if err != nil {
		log.Errorf("couldn't read audio file, removing track from queue: \"%s\"", track.Title)

		ok := p.queue.DeleteAt(0)
		if !ok {
			log.Errorf("Can't delete item. Queue length is %v", p.queue.Length())
		}
		return err
	}

	// Copy track, because audio/encodeSession can't be reassigned
	tr := Track{
		Title:        track.Title,
		Duration:     track.Duration,
		ID:           track.ID,
		ThumbnailURL: track.ThumbnailURL,
		URL:          track.URL,
		ChannelID:    track.ChannelID,
		Audio:        es,
	}

	ok = p.queue.ReplaceAt(0, tr)
	if !ok {
		log.Debugf("could not replace item at index %v, deleteting...", 0)
		ok = p.queue.DeleteAt(0)
		if !ok {
			log.Errorf("Can't delete item. Queue length is %v", p.queue.Length())
		}
		return errors.New("could not prepare track")
	}

	log.Debug("next track prepared")
	return nil
}

// QueueLength return queue length
func (p *Player) QueueLength() int {
	return p.queue.Length()
}
