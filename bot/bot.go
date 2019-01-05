package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/config"
	"github.com/zrcni/go-discord-music-bot/player"
)

const commandPrefix = "!"

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

func (b *Bot) joinChannel(session *discordgo.Session, guildID string, channelID string) (*discordgo.VoiceConnection, error) {
	voiceConnection, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Printf("Join voice channel: %v", err)
		return nil, err
	}

	log.Printf("Joined channel: %v", channelID)

	b.voiceConnection = voiceConnection

	return b.voiceConnection, nil
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

func init() {
	config.SetupEnv()
}

// Start discord bot
func Start() {
	bot.player = *player.New()
	bot.player.UpdateBotStatus = bot.UpdateListeningStatus

	sess, err := discordgo.New(fmt.Sprintf("Bot %s", config.BotToken))
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

	cp := commandParams{message, session}

	switch {
	case messageHasCommand(message.Content, "repeat "):
		callCommand(repeatCommand, cp)

	case messageHasCommand(message.Content, "start"):
		callCommand(startCommand, cp)

	case messageHasCommand(message.Content, "stop"):
		callCommand(stopCommand, cp)

	case messageHasCommand(message.Content, "playlist "):
		callCommand(playlistsCommand, cp)

	case messageHasCommand(message.Content, "play "):
		callCommand(playCommand, cp)

	case messageHasCommand(message.Content, "pause"):
		callCommand(pauseCommand, cp)

	case messageHasCommand(message.Content, "continue"):
		callCommand(continueCommand, cp)
	}
}

func callCommand(fn func(commandParams), cp commandParams) {
	log.Printf("Received command \"%s\"", cp.message.Content)
	fn(cp)
}
