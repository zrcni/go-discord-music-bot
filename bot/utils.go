package bot

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rylio/ytdl"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/audiorepository"
	"github.com/zrcni/go-discord-music-bot/player"
	"github.com/zrcni/go-discord-music-bot/utils"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
	"github.com/zrcni/go-discord-music-bot/youtube"
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

func createFile(videoID string) (*os.File, error) {
	basePath, err := utils.GetBasePath()
	if err != nil {
		return nil, errors.Wrap(err, "utils.GetBasePath")
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", basePath, videoID))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func downloadYoutube(url string, cp commandParams) (*player.Track, error) {
	videoInfo, err := youtube.GetMetadata(url)
	if err != nil {
		log.Errorf("Couldn't get metadata for youtube video (%s)", url)
		return nil, err
	}

	thumbnailURL := videoInfo.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String()
	var track *player.Track

	// TODO: figure out optimal format
	format := videoInfo.Formats[0]

	// Download audio to buffer
	audioBuffer := &bytes.Buffer{}
	err = youtube.Download(videoInfo, format, audioBuffer)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("downloaded \"%s\"", videoInfo.Title)

	timestampMs := time.Now().UnixNano() / 1000000
	uniqueID := fmt.Sprintf("%s-%v", videoInfo.ID, timestampMs)
	// if queue is not empty: save audio data as buffer to file
	// else: assign audio data as dca.EncodeSession pointer to the track
	if bot.player.QueueLength() > 0 {
		err := videoaudio.SaveAudioToFile(uniqueID, audioBuffer.Bytes())
		if err != nil {
			log.Error(err)
			return nil, err
		}

		filename := fmt.Sprintf("%s.mp3", uniqueID)

		audiorepository.Add(&audiorepository.File{
			ID:   uniqueID,
			Name: filename,
		})

		track = &player.Track{
			Title:        videoInfo.Title,
			ID:           uniqueID,
			Duration:     videoInfo.Duration,
			ThumbnailURL: thumbnailURL,
			URL:          url,
			Filename:     filename,
			ChannelID:    cp.message.ChannelID,
		}

	} else {
		audiorepository.Add(&audiorepository.Raw{
			ID:  uniqueID,
			Raw: audioBuffer.Bytes(),
		})

		// assign whole track, because encodeSession can't be reassigned
		track = &player.Track{
			Title:        videoInfo.Title,
			ID:           uniqueID,
			Duration:     videoInfo.Duration,
			ThumbnailURL: thumbnailURL,
			URL:          url,
			ChannelID:    cp.message.ChannelID,
		}
	}

	return track, nil
}
