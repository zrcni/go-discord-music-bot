package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// State manages the state of the bot
type State struct {
	audio   Audio
	session *discordgo.Session
	botID   string
}

// SetBotID sets botID
func (st *State) SetBotID(userID string) {
	st.botID = userID
}

// SetSession sets session
func (st *State) SetSession(sess *discordgo.Session) {
	st.session = sess
}

// Audio struct stores discordgo voice connection
type Audio struct {
	connection *discordgo.VoiceConnection
}

// IsConnected sets discord voice connection
func (a *Audio) IsConnected() bool {
	return a.connection != nil
}

// SetConnection sets discord voice connection
func (a *Audio) SetConnection(vc *discordgo.VoiceConnection) {
	a.connection = vc
}

// Implements io.Writer to be able to write to the connection
func (a Audio) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return len(data), errors.New("Voice.Write data length: 0")
	}

	a.connection.OpusSend <- data

	return len(data), nil
}
