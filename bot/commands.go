package bot

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/spotify"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
	"github.com/zrcni/go-discord-music-bot/youtube"
)

func repeatCommand(message *discordgo.MessageCreate) {
	str := strings.Split(message.Content, "repeat ")
	if len(str) != 2 {
		return
	}

	msg, err := state.session.ChannelMessageSend(message.ChannelID, str[1])
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", message.ChannelID, err)
		return
	}
	log.Printf("Message %v sent to channel %v", msg.ID, message.ChannelID)
}

func startCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	state.UpdateListeningStatus("")

	guild := session.State.Guilds[0]

	channels, err := session.GuildChannels(guild.ID)
	if err != nil {
		log.Printf("Could not get guild channels: %v", err)
		return
	}
	voiceChannels := filterChannels(channels, discordgo.ChannelTypeGuildVoice)

	voiceChannel := voiceChannels[0].ID

	if voiceChannel == state.audio.GetChannelID() {
		return
	}

	vc, err := joinChannel(session, guild.ID, voiceChannel)
	if err != nil {
		return
	}

	state.audio.SetConnection(vc)
}

func stopCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	state.UpdateListeningStatus("")

	if !state.audio.IsConnected() {
		return
	}

	channelID := state.audio.GetChannelID()

	if state.audio.IsStreaming() {
		state.audio.stream.SetPaused(true)
	}

	err := state.audio.connection.Disconnect()
	if err != nil {
		log.Printf("Could not disconnect from audio channel %v: %v", channelID, err)
		state.audio.SetConnection(nil)
		return
	}

	log.Printf("Disconnected from audio channel %v", channelID)
}

func playlistsCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	msg := strings.Split(message.Content, "playlist ")
	if len(msg) != 2 {
		return
	}

	searchTerm := msg[1]

	spotifyClient, _ := spotify.NewClient()
	playlists := spotifyClient.GetPlaylists(searchTerm)

	if len(playlists) > 0 {
		session.ChannelMessageSend(message.ChannelID, strings.Join(playlists, "\n"))
	}
}

func playCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	msg := strings.Split(message.Content, "play ")
	if len(msg) != 2 || msg[1] == "" {
		return
	}

	searchTerm := msg[1]

	if !state.audio.IsConnected() {
		log.Printf("Audio connection doesn't exist")
		return
	}

	state.UpdateListeningStatus("Preparing song")

	track, err := youtube.Get(searchTerm)
	if err != nil {
		log.Printf("error while downloading youtube audio: %v", err)
		return
	}

	encodeSession, err := videoaudio.EncodeAudioToDCA(track.Data)
	if err != nil {
		log.Printf("EncodeAudioToDCA: %v", err)
		return
	}
	defer encodeSession.Cleanup()

	state.UpdateListeningStatus(track.Info.Title)
	state.SetNowPlaying(track.Info.Title)
	state.audio.connection.Speaking(true)

	err = state.audio.Stream(encodeSession)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		log.Print("end of file, continuing")
		err = nil
	}
	if err != nil {
		log.Printf("CreateStream: %v", err)
		return
	}

	state.UpdateListeningStatus("")
	state.audio.connection.Speaking(false)
}

func pauseCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	if !state.audio.IsStreaming() {
		return
	}
	if state.audio.stream.Paused() {
		return
	}
	state.audio.stream.SetPaused(true)

	nowPlaying := state.GetNowPlaying()
	state.UpdateListeningStatus(fmt.Sprintf("%s %s", pausedPrefix, nowPlaying))
}

func continueCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	if !state.audio.IsStreaming() {
		return
	}
	if !state.audio.stream.Paused() {
		return
	}

	state.audio.stream.SetPaused(false)
	nowPlaying := state.GetNowPlaying()
	state.UpdateListeningStatus(nowPlaying)
}
