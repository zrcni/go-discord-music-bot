package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageHasCommand(msgContent string, command string) bool {
	commandWithPrefix := fmt.Sprintf("%s%s", commandPrefix, command)
	return strings.HasPrefix(msgContent, commandWithPrefix)
}

func filterChannels(channels []*discordgo.Channel, chanType discordgo.ChannelType) []discordgo.Channel {
	var voiceChannels []discordgo.Channel

	for _, channel := range channels {
		if chanType == channel.Type {
			voiceChannels = append(voiceChannels, *channel)
		}
	}

	return voiceChannels
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
