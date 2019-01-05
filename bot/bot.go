package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/zrcni/go-discord-music-bot/player"
)

const (
	commandPrefix = "!"
	pausedPrefix  = "[Paused]"
)

var bot = &Bot{}

type commandParams struct {
	message *discordgo.MessageCreate
	session *discordgo.Session
}

// Bot manages the state of the bot
type Bot struct {
	ID              string
	session         *discordgo.Session
	voiceConnection *discordgo.VoiceConnection
	player          player.Player
}

func (b *Bot) joinChannel(session *discordgo.Session, guildID string, channelID string) error {
	voiceConnection, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return errors.Wrap(err, "could not join voice channel")
	}

	log.Printf("Joined channel: %v", channelID)

	b.voiceConnection = voiceConnection

	return nil
}

func (b *Bot) leaveChannel(session *discordgo.Session, channelID string) error {
	if bot.voiceConnection == nil {
		return errors.New("voice connection doesn't exist")
	}

	err := bot.voiceConnection.Disconnect()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not leave voice channel %v", channelID))
	}

	log.Printf("Disconnected from audio channel %v", channelID)

	return nil
}

// SetConnection sets discord voice connection
func (b *Bot) setConnection(vc *discordgo.VoiceConnection) {
	b.voiceConnection = vc
}

func (b *Bot) isVoiceConnected() bool {
	if bot.voiceConnection == nil {
		return false
	}
	return bot.voiceConnection.Ready
}

// SetSession sets session
func (b *Bot) SetSession(sess *discordgo.Session) {
	b.session = sess
}

// UpdateListeningStatus sets discord listening status and stores it locally
func (b *Bot) UpdateListeningStatus(status string) {
	if err := b.session.UpdateListeningStatus(status); err != nil {
		fmt.Printf("Could not set listening status: %v", err)
	}
}

// Start discord bot
func Start() {
	bot.player = *player.New()
	go bot.listenForPlayerEvents()

	sess, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))
	if err != nil {
		log.Printf("Create session: %v", err)
		return
	}

	bot.SetSession(sess)

	user, err := bot.session.User("@me")
	if err != nil {
		log.Printf("Create user: %v", err)
		return
	}

	bot.ID = user.ID

	bot.session.AddHandler(readyHandler)
	bot.session.AddHandler(commandHandler)

	if err := bot.session.Open(); err != nil {
		log.Printf("Error opening connection to Discord: %v", err)
	}

	defer bot.session.Close()

	// Keep process running indefinitely, because channel
	// keeps waiting for message that will never be sent
	<-make(chan struct{})
}

func readyHandler(session *discordgo.Session, ready *discordgo.Ready) {
	user, err := session.User("@me")
	if err != nil {
		log.Printf("readyHandler user: %v", err)
	}
	guilds := session.State.Guilds

	fmt.Printf("%s has started on %d server(s)\n", user.Username, len(guilds))

	bot.UpdateListeningStatus("")
}

func commandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == bot.ID || user.Bot {
		return
	}

	if !strings.HasPrefix(message.Content, commandPrefix) {
		return
	}

	params := commandParams{message, session}
	var command func(commandParams)

	switch {
	case messageHasCommand(message.Content, "join "):
		command = joinCommand

	case messageHasCommand(message.Content, "repeat "):
		command = repeatCommand

	case messageHasCommand(message.Content, "start"):
		command = startCommand

	case messageHasCommand(message.Content, "stop"):
		command = stopCommand

	case messageHasCommand(message.Content, "playlist "):
		command = playlistsCommand

	case messageHasCommand(message.Content, "play "):
		command = playCommand

	case messageHasCommand(message.Content, "play"):
		command = continueCommand

	case messageHasCommand(message.Content, "pause"):
		command = pauseCommand

	default:
		return
	}

	callCommand(command, params)
}

func callCommand(fn func(commandParams), cp commandParams) {
	log.Printf("Received command \"%s\"", cp.message.Content)
	fn(cp)
}
