package bot

import (
	"fmt"
	"log"

	"github.com/zrcni/go-discord-music-bot/player"
)

func (b *Bot) handlePlayerEvent(e player.Event) {
	switch e.Type {
	case player.PLAY:
		b.handlePlayEvent(e)
		return

	case player.QUEUE:
		b.handleQueueEvent(e)
		return

	case player.PAUSE:
		b.handlePauseEvent(e)
		return

	case player.STOP:
		b.handleStopEvent(e)
		return

	default:
		log.Printf("invalid player event: %+v", e)
		return
	}
}

func (b *Bot) listenForPlayerEvents() {
	for event := range b.player.EventChannel {
		b.handlePlayerEvent(event)
	}
}

func (b *Bot) handlePlayEvent(e player.Event) {
	bot.UpdateListeningStatus(e.TrackInfo.Title)

	_, err := b.session.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Now playing: \"%s\"", e.TrackInfo.Title))
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", e.ChannelID, err)
	}
}

func (b *Bot) handleQueueEvent(e player.Event) {
	_, err := b.session.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Queued: \"%s\"", e.TrackInfo.Title))
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", e.ChannelID, err)
	}
}

func (b *Bot) handlePauseEvent(e player.Event) {
	status := fmt.Sprintf("%s %s", pausedPrefix, e.TrackInfo.Title)
	b.UpdateListeningStatus(status)
}

func (b *Bot) handleStopEvent(e player.Event) {
	b.UpdateListeningStatus("")

	var message string

	if e.Message != "" {
		message = fmt.Sprintf("Stopped playing - %s.", e.Message)
	} else {
		message = "Stopped plaing"
	}

	_, err := b.session.ChannelMessageSend(e.ChannelID, message)
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", e.ChannelID, err)
	}
}
