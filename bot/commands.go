package bot

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/bot/commands"
	"github.com/zrcni/go-discord-music-bot/downloader"
	"github.com/zrcni/go-discord-music-bot/spotify"
	"github.com/zrcni/go-discord-music-bot/youtube"
)

var GOPATH = os.Getenv("GOPATH")
var soundboardBinaryPath = fmt.Sprintf("%s/src/github.com/zrcni/go-discord-music-bot/bin/soundboard", GOPATH)

func repeatCommand(cp commandParams) {
	str := strings.Split(cp.message.Content, commands.REPEAT_TEXT)
	if len(str) != 2 {
		return
	}

	msg, err := bot.session.ChannelMessageSend(cp.message.ChannelID, str[1])
	if err != nil {
		log.Errorf("Could not send a message to channel %v: %v", cp.message.ChannelID, err)
		return
	}
	log.Debugf("Message %v sent to channel %v", msg.ID, cp.message.ChannelID)
}

func joinCommand(cp commandParams) {
	msg := strings.Split(cp.message.Content, commands.JOIN_CHANNEL)
	if len(msg) != 2 {
		return
	}
	channelName := msg[1]

	guild := cp.session.State.Guilds[0]

	channels, err := cp.session.GuildChannels(guild.ID)
	if err != nil {
		log.Errorf("Could not get guild channels: %v", err)
		return
	}

	voiceChannels := filterChannels(channels, discordgo.ChannelTypeGuildVoice)

	channel := findChannelByName(voiceChannels, channelName)

	if channel == nil {
		message := fmt.Sprintf("Could not find voice channel by name \"%s\"", channelName)
		_, err := bot.session.ChannelMessageSend(cp.message.ChannelID, message)
		if err != nil {
			log.Errorf("Could not send a message to channel %v: %v", cp.message.ChannelID, err)
		}
		return
	}

	err = bot.joinChannel(cp.session, guild.ID, channel.ID)
	if err != nil {
		log.Error(err)
		return
	}

	message := fmt.Sprintf("Joined \"%s\"", channelName)
	_, err = bot.session.ChannelMessageSend(cp.message.ChannelID, message)
	if err != nil {
		log.Errorf("Could not send a message to channel %v: %v", cp.message.ChannelID, err)
	}
}

func startCommand(cp commandParams) {
	bot.UpdateListeningStatus("")

	guild := cp.session.State.Guilds[0]

	channels, err := cp.session.GuildChannels(guild.ID)
	if err != nil {
		log.Errorf("Could not get guild channels: %v", err)
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
		log.Error(errors.New("voice connection doesn't exist"))
		return
	}

	bot.player.Stop()

	channelID := bot.voiceConnection.ChannelID

	err := bot.leaveChannel(cp.session, channelID)
	if err != nil {
		log.Error(err)
		return
	}
}

func playlistsCommand(cp commandParams) {
	msg := strings.Split(cp.message.Content, commands.FIND_PLAYLIST)
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
	msg := strings.Split(cp.message.Content, commands.PLAY_TRACK)
	if len(msg) != 2 || msg[1] == "" {
		return
	}
	url := msg[1]

	if !youtube.IsYoutubeURL(url) {
		message := fmt.Sprintf("Not a valid Youtube URL: %s", url)
		_, err := bot.session.ChannelMessageSend(cp.message.ChannelID, message)
		if err != nil {
			log.Errorf("Could not send a message to channel %v: %v", cp.message.ChannelID, err)
		}
		return
	}

	if !bot.isVoiceConnected() {
		log.Debug("voice connection isn't active")
		return
	}

	timestampMs := time.Now().UnixNano() / 1000000
	ID := fmt.Sprintf("%s-%v", url, timestampMs)

	downloadable := downloader.Downloadable{ID: ID}
	downloadable.Get = func() {
		track, err := downloadYoutube(url, cp)
		if err != nil {
			log.Error(err)
		}
		go bot.player.Queue(track)
	}

	go bot.downloader.Queue(downloadable)
}

func pauseCommand(cp commandParams) {
	if !bot.player.IsPlaying() {
		log.Debug("playback is already paused")
		return
	}
	bot.player.SetPaused(true)
}

func continueCommand(cp commandParams) {
	if bot.player.IsPlaying() {
		log.Debug("playback is active")
		return
	}

	bot.player.SetPaused(false)
}

func soundCommand(cp commandParams) {
	if bot.player.IsPlaying() {
		log.Debug("player is active")
		return
	}

	msg := strings.Split(cp.message.Content, commands.SOUND)
	if len(msg) != 2 || msg[1] == "" {
		return
	}
	soundName := msg[1]

	cmd := exec.Command(soundboardBinaryPath, "write", soundName)

	buf, err := cmd.Output()
	if err != nil {
		log.Error(err)
	}

	log.Printf("sound clip buffer length: %v", buf)
	// TODO stream sound to discord
}
