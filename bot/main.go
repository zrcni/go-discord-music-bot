package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/config"
)

var state State

func init() {
	config.SetupEnv()
}

// Start discord bot
func Start() {
	state = State{}

	sess, err := discordgo.New(fmt.Sprintf("Bot %s", config.BotToken))
	if err != nil {
		log.Printf("Create session: %v", err)
		return
	}

	state.SetSession(sess)

	user, err := state.session.User("@me")
	if err != nil {
		log.Printf("Create user: %v", err)
		return
	}

	state.SetBotID(user.ID)

	state.session.AddHandler(readyHandler)
	state.session.AddHandler(commandHandler)

	if err := state.session.Open(); err != nil {
		log.Printf("Error opening connection to Discord: %v", err)
	}

	defer state.session.Close()

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

	state.UpdateListeningStatus("")
}

func commandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == state.botID || user.Bot {
		return
	}

	// TODO queue messages

	if !strings.HasPrefix(message.Content, commandPrefix) {
		return
	}

	switch {
	case messageHasCommand(message.Content, "repeat "):
		repeatCommand(message)

	case messageHasCommand(message.Content, "start"):
		startCommand(message, session)

	case messageHasCommand(message.Content, "stop"):
		stopCommand(message, session)

	case messageHasCommand(message.Content, "playlist "):
		playlistsCommand(message, session)

	case messageHasCommand(message.Content, "play "):
		playCommand(message, session)

	case messageHasCommand(message.Content, "pause"):
		pauseCommand(message, session)

	case messageHasCommand(message.Content, "continue"):
		continueCommand(message, session)
	}
}
