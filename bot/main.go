package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/config"
)

var (
	voice   Voice
	session *discordgo.Session
	botID   string
)

type Voice struct {
	connection *discordgo.VoiceConnection
}

func (v Voice) Write(data []byte) (n int, err error) {
	v.connection.OpusSend <- data

	if len(data) > 0 {
		return len(data), nil
	}

	return len(data), errors.New("Voice.Write data length: 0")
}

func init() {
	config.SetupEnv()
}

// Start discord bot
func Start() {
	sess, err := discordgo.New(fmt.Sprintf("Bot %s", config.BotToken))
	if err != nil {
		log.Printf("Create session: %v", err)
		return
	}

	session = sess

	user, err := session.User("@me")
	if err != nil {
		log.Printf("Create user: %v", err)
		return
	}

	botID = user.ID

	session.AddHandler(commandHandler)

	session.AddHandler(func(session *discordgo.Session, ready *discordgo.Ready) {
		updateListeningStatus(session, "")
		guilds := session.State.Guilds
		fmt.Printf("%s has started on %d server(s)\n", user.Username, len(guilds))
	})

	if err := session.Open(); err != nil {
		log.Printf("Error opening connection to Discord: %v", err)
	}

	defer session.Close()

	// Keep process running indefinitely, because channel
	// keeps waiting for message that will never be sent
	<-make(chan struct{})
}

func commandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		// TODO: queue
		return
	}

	if !strings.HasPrefix(message.Content, commandPrefix) {
		return
	}

	// Repeat string after !repeat command
	if strings.Contains(message.Content, "repeat ") {
		repeatCommand(message)
		return
	}

	if strings.Contains(message.Content, "start") {
		startCommand(message, session)
		return
	}

	if strings.Contains(message.Content, "stop") {
		stopCommand(message, session)
		return
	}

	if strings.Contains(message.Content, "playlist ") {
		playlistsCommand(message, session)
		return
	}

	if strings.Contains(message.Content, "play ") {
		playCommand(message, session)
		return
	}
}

func updateListeningStatus(discord *discordgo.Session, status string) {
	if err := discord.UpdateListeningStatus(status); err != nil {
		fmt.Printf("Could not set listening status: %v", err)
	}
}

func filterVoiceChannels(channels []*discordgo.Channel) []discordgo.Channel {
	var voiceChannels []discordgo.Channel
	for _, channel := range channels {
		if discordgo.ChannelTypeGuildVoice == channel.Type {
			voiceChannels = append(voiceChannels, *channel)
		}
	}
	return voiceChannels
}
