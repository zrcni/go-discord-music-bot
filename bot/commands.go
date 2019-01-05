package bot

import (
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/spotify"
	"github.com/zrcni/go-discord-music-bot/youtube"
)

func repeatCommand(cp commandParams) {
	str := strings.Split(cp.message.Content, "repeat ")
	if len(str) != 2 {
		return
	}

	msg, err := bot.session.ChannelMessageSend(cp.message.ChannelID, str[1])
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", cp.message.ChannelID, err)
		return
	}
	log.Printf("Message %v sent to channel %v", msg.ID, cp.message.ChannelID)
}

func startCommand(cp commandParams) {
	bot.UpdateListeningStatus("")

	guild := cp.session.State.Guilds[0]

	channels, err := cp.session.GuildChannels(guild.ID)
	if err != nil {
		log.Printf("Could not get guild channels: %v", err)
		return
	}
	voiceChannels := filterChannels(channels, discordgo.ChannelTypeGuildVoice)

	voiceChannelID := voiceChannels[0].ID

	err = bot.joinChannel(cp.session, guild.ID, voiceChannelID)
	if err != nil {
		return
	}
}

func stopCommand(cp commandParams) {
	bot.UpdateListeningStatus("")

	if bot.voiceConnection == nil {
		log.Print(errors.New("voice connection doesn't exist"))
		return
	}

	channelID := bot.voiceConnection.ChannelID

	bot.player.ClearQueue()

	err := bot.voiceConnection.Disconnect()
	if err != nil {
		log.Printf("Could not disconnect from audio channel %v: %v", channelID, err)
		return
	}

	log.Printf("Disconnected from audio channel %v", channelID)
}

func playlistsCommand(cp commandParams) {
	msg := strings.Split(cp.message.Content, "playlist ")
	if len(msg) != 2 {
		return
	}

	searchTerm := msg[1]

	spotifyClient, _ := spotify.NewClient()
	playlists := spotifyClient.GetPlaylists(searchTerm)

	if len(playlists) > 0 {
		cp.session.ChannelMessageSend(cp.message.ChannelID, strings.Join(playlists, "\n"))
	}
}

func playCommand(cp commandParams) {
	msg := strings.Split(cp.message.Content, "play ")
	if len(msg) != 2 || msg[1] == "" {
		return
	}

	searchTerm := msg[1]

	if !bot.isVoiceConnected() {
		log.Print(errors.New("voice connection isn't active"))
		return
	}

	track, err := youtube.Get(searchTerm)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("downloaded \"%s\"", track.Info.Title)

	track.ChannelID = cp.message.ChannelID

	go bot.player.Queue(track, bot.voiceConnection)
}

func pauseCommand(cp commandParams) {
	if !bot.player.IsPlaying() {
		log.Print("playback is already paused")
		return
	}
	bot.player.SetPaused(true)
}

func continueCommand(cp commandParams) {
	if bot.player.IsPlaying() {
		log.Print("playback is active")
		return
	}

	bot.player.SetPaused(false)
}
