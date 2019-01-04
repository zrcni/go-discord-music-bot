package bot

import (
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/player"
	"github.com/zrcni/go-discord-music-bot/spotify"
	"github.com/zrcni/go-discord-music-bot/youtube"
)

func repeatCommand(message *discordgo.MessageCreate) {
	str := strings.Split(message.Content, "repeat ")
	if len(str) != 2 {
		return
	}

	msg, err := bot.session.ChannelMessageSend(message.ChannelID, str[1])
	if err != nil {
		log.Printf("Could not send a message to channel %v: %v", message.ChannelID, err)
		return
	}
	log.Printf("Message %v sent to channel %v", msg.ID, message.ChannelID)
}

func startCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	bot.UpdateListeningStatus("")

	guild := session.State.Guilds[0]

	channels, err := session.GuildChannels(guild.ID)
	if err != nil {
		log.Printf("Could not get guild channels: %v", err)
		return
	}
	voiceChannels := filterChannels(channels, discordgo.ChannelTypeGuildVoice)

	voiceChannelID := voiceChannels[0].ID

	vc, err := bot.joinChannel(session, guild.ID, voiceChannelID)
	if err != nil {
		return
	}

	bot.setConnection(vc)
}

func stopCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	bot.UpdateListeningStatus("")

	if bot.voiceConnection == nil {
		log.Print(errors.New("voice connection doesn't exist"))
		return
	}

	channelID := bot.voiceConnection.ChannelID

	err := bot.voiceConnection.Disconnect()
	if err != nil {
		log.Printf("Could not disconnect from audio channel %v: %v", channelID, err)
		bot.voiceConnection = nil
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

	if !bot.isVoiceConnected() {
		log.Print(errors.New("voice connection isn't active"))
		return
	}

	track, err := youtube.Get(searchTerm)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("\"%s\" downloaded", track.Info.Title)

	go func(bot *Bot, track player.Track) {
		ok := make(chan bool, 1)

		bot.player.Queue(track, bot.voiceConnection, ok)

		success := <-ok
		if success == true {
			bot.UpdateListeningStatus(track.Info.Title)
		}
	}(bot, track)
}

func pauseCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	if !bot.player.IsPlaying() {
		log.Print("playback is already paused")
		return
	}
	bot.player.SetPaused(true)
}

func continueCommand(message *discordgo.MessageCreate, session *discordgo.Session) {
	if bot.player.IsPlaying() {
		log.Print("playback is active")
		return
	}

	bot.player.SetPaused(false)
}
