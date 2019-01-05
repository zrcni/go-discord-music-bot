package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageHasCommand(msgContent string, command string) bool {
	commandWithPrefix := fmt.Sprintf("%s%s", commandPrefix, command)

	if msgContent == commandWithPrefix {
		return true
	}

	return strings.HasPrefix(msgContent, commandWithPrefix)
}

func filterChannels(channels []*discordgo.Channel, chanType discordgo.ChannelType) []*discordgo.Channel {
	var voiceChannels []*discordgo.Channel

	for _, channel := range channels {
		if chanType == channel.Type {
			voiceChannels = append(voiceChannels, channel)
		}
	}

	return voiceChannels
}

func findChannelByName(channels []*discordgo.Channel, name string) *discordgo.Channel {
	for _, c := range channels {
		if c.Name == name {
			return c
		}
	}
	return nil
}
