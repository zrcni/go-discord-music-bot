package player

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/audiorepository"
	"github.com/zrcni/go-discord-music-bot/listqueue"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
)

// Track stores audio data and info
type Track struct {
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
		queue:  listqueue.New(20),
		Events: make(chan Event),
	}
}

// Player handler audio playback
type Player struct {
	currentTrack    *Track
	stream          *dca.StreamingSession
	VoiceConnection *discordgo.VoiceConnection
	queue           listqueue.Queue
	Events          chan Event
}

// SetNowPlaying sets currently playing track
func (p *Player) SetNowPlaying(track *Track) {
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
func (p *Player) GetNowPlaying() *Track {
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

	if p.currentTrack != nil {
		e := Event{
			Type:      STOP,
			ChannelID: p.currentTrack.ChannelID,
			Message:   "Stopped playing",
		}

		p.sendEvent(e)
	}

	if p.stream != nil {
		p.stream.SetPaused(true)
	}
	p.ClearQueue()
	p.currentTrack = nil
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
func (p *Player) Queue(track *Track) {
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
		p.currentTrack = nil
		log.Debug("the queue is empty")
		return
	}

	if p.isStreaming() {
		log.Debug("already streaming")
		return
	}

	track, err := p.queue.Process()
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("TRACKKKKKK BOYYYY: %T - %+v", track, track)
	t, ok := track.(*Track)
	if !ok {
		panic(fmt.Sprintf("track is not of type Track: %+v", t))
	}

	// if p.queue.Length() > 1 {
	// 	go func(p *Player) {
	// 		err := p.prepareNextTrack()
	// 		if err != nil {
	// 			log.Error(err)
	// 		}
	// 	}(p)
	// }

	go p.play(t)
}

// play starts the process that streams the track
func (p *Player) play(track *Track) {
	log.Infof("QUEUEUE LISTTTTTTT: %+v", p.queue.GetList())
	if !p.VoiceConnection.Ready {
		log.Debug("voice connection is not ready")
		return
	}
	// if p.IsPlaying() {
	// 	log.Debug("already playing")
	// 	return
	// }

	audio, err := audiorepository.Get(track.ID)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("AUDIO: %v", len(audio))
	buffer := bytes.NewReader(audio)

	es, err := videoaudio.EncodeAudioToDCA(buffer)
	if err != nil {
		log.Error(err)
		return
	}

	p.SetNowPlaying(track)

	err = p.startStream(es)
	if err != nil {
		log.Error(err)
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
		log.Errorf("player.IsStreaming: %v", err)
		return false
	}

	return !finished
}

// QueueLength return queue length
func (p *Player) QueueLength() int {
	return p.queue.Length()
}
