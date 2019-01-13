package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zrcni/go-discord-music-bot/player"
)

const (
	// PURPLE color code
	PURPLE int = 0x631b68
)

func createPlayingMessage(e player.Event) *discordgo.MessageSend {
	trackDuration := fmt.Sprintf("Duration: %s", e.Track.Duration.String())
	image := &discordgo.MessageEmbedImage{
		URL: e.Track.ThumbnailURL,
	}
	footer := &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("🕜 %s", time.Now().Format(time.ANSIC)),
	}

	return &discordgo.MessageSend{
		Content: "Now playing:",
		Embed: &discordgo.MessageEmbed{
			Color:       PURPLE,
			Title:       e.Track.Title,
			Description: trackDuration,
			URL:         e.Track.URL,
			Image:       image,
			Footer:      footer,
		},
	}
}
