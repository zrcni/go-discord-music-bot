package bot

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// State manages the state of the bot
type State struct {
	audio      Audio
	session    *discordgo.Session
	botID      string
	nowPlaying string
}

// SetBotID sets botID
func (st *State) SetBotID(userID string) {
	st.botID = userID
}

// SetSession sets session
func (st *State) SetSession(sess *discordgo.Session) {
	st.session = sess
}

// UpdateListeningStatus sets discord listening status and stores it locally
func (st *State) UpdateListeningStatus(status string) {
	if err := st.session.UpdateListeningStatus(status); err != nil {
		fmt.Printf("Could not set listening status: %v", err)
	}
}

// SetNowPlaying gets local listening status
func (st *State) SetNowPlaying(nowPlaying string) {
	st.nowPlaying = nowPlaying
}

// GetNowPlaying gets local listening status
func (st *State) GetNowPlaying() string {
	return st.nowPlaying
}

// Audio struct stores discordgo voice connection
type Audio struct {
	connection *discordgo.VoiceConnection
	stream     *dca.StreamingSession
}

// GetChannelID gets voice connection channel ID
func (a *Audio) GetChannelID() string {
	if a.connection == nil {
		return ""
	}
	return a.connection.ChannelID
}

// IsConnected gets discord voice connection status as boolean
func (a *Audio) IsConnected() bool {
	return a.connection != nil
}

// SetConnection sets discord voice connection
func (a *Audio) SetConnection(vc *discordgo.VoiceConnection) {
	a.connection = vc
}

// IsStreaming returns stream status as boolean
func (a *Audio) IsStreaming() bool {
	return a.stream != nil
}

// Stream audio to discord
func (a *Audio) Stream(source dca.OpusReader) error {
	done := make(chan error)

	stream := dca.NewStream(source, state.audio.connection, done)

	a.stream = stream

	err := <-done
	if err != nil {
		return err
	}

	return nil
}

// Implements io.Writer to be able to write to the connection
func (a Audio) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return len(data), errors.New("Voice.Write data length: 0")
	}

	a.connection.OpusSend <- data

	return len(data), nil
}
