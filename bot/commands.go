package bot

import (
	"bytes"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/zrcni/go-discord-music-bot/spotify"
	"github.com/zrcni/go-discord-music-bot/youtube"
)

func repeatCommand(message *discordgo.MessageCreate) {
	str := strings.Split(message.Content, "repeat ")
	if len(str) != 2 {
		return
	}

	msg, err := session.ChannelMessageSend(message.ChannelID, str[1])
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", message.ChannelID, err)
		return
	}
	log.Printf("Message %v sent to channel %v", msg.ID, message.ChannelID)
}

func startCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	updateListeningStatus(session, "Waiting")

	guild := session.State.Guilds[0]

	channels, err := session.GuildChannels(guild.ID)
	if err != nil {
		log.Printf("Could not get guild channels: %v", err)
		return
	}
	voiceChannels := filterVoiceChannels(channels)

	vc, err := joinChannel(session, guild.ID, voiceChannels[0].ID)
	if err != nil {
		return
	}

	voice.connection = vc
}

func stopCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	updateListeningStatus(session, "")

	if voice.connection == nil {
		return
	}

	channelID := voice.connection.ChannelID

	err := voice.connection.Disconnect()
	if err != nil {
		log.Printf("Could not disconnect from voice channel %v: %v", channelID, err)
		voice.connection = nil
		return
	}
	log.Printf("Disconnected from voice channel %v", channelID)
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

	if voice.connection == nil {
		log.Printf("Voice connection doesn't exist")
		return
	}

	data, videoInfo, err := youtube.DownloadAudio(searchTerm)
	if err != nil {
		log.Printf("error while downloading youtube audio: %v", err)
		return
	}

	reader := bytes.NewReader(data)

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.BufferedFrames = 500
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeMem(reader, options)
	if err != nil {
		log.Printf("dca.EncodeMem: %v", err)
		return
	}
	defer encodeSession.Cleanup()

	updateListeningStatus(session, videoInfo.Title)
	voice.connection.Speaking(true)

	done := make(chan error)
	dca.NewStream(encodeSession, voice.connection, done)

	err = <-done
	if err != nil {
		log.Printf("NewStream: %v", err)
		return
	}

	updateListeningStatus(session, "Waiting")
	voice.connection.Speaking(false)
}

func joinChannel(session *discordgo.Session, guildID string, channelID string) (*discordgo.VoiceConnection, error) {
	voiceConnection, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Printf("Join voice channel: %v", err)
		return nil, err
	}

	log.Printf("Joined channel: %v", channelID)

	return voiceConnection, nil
}
