package youtube

import (
	"bytes"
	"errors"
	"io"
	"regexp"

	"github.com/rylio/ytdl"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
)

const youtubeURLMatcher = "^.*(?:(?:youtu\\.be\\/|v\\/|vi\\/|u\\/\\w\\/|embed\\/)|(?:(?:watch)?\\?v(?:i)?=|\\&v(?:i)?=))([^#\\&\\?]*).*"

func getVideoID(url string) string {
	re := regexp.MustCompile(youtubeURLMatcher)
	match := re.FindStringSubmatch(url)

	if len(match) == 2 {
		videoID := match[1]
		return videoID
	}

	return ""
}

// GetMetadata fetches metadata for a video
func GetMetadata(url string) (ytdl.VideoInfo, error) {
	videoInfo, err := ytdl.GetVideoInfo(url)
	if err != nil {
		log.Error(err)
		return ytdl.VideoInfo{}, err
	}
	return *videoInfo, nil
}

// Download a youtube video
func downloadVideo(writer io.Writer, videoInfo ytdl.VideoInfo) error {
	if len(videoInfo.Formats) == 0 {
		return errors.New("No available video formats")
	}

	format := videoInfo.Formats[0]

	log.Debugf("Format: %v", format)

	if err := videoInfo.Download(format, writer); err != nil {
		return err
	}
	return nil
}

// // Get downloads video from youtube and transcodes it to audio
// func Get(url string) (player.Track, error) {
// 	videoInfo, err := GetMetadata(url)
// 	if err != nil {
// 		return player.Track{}, err
// 	}

// 	videoData := &bytes.Buffer{}

// 	err = downloadVideo(videoData, videoInfo)
// 	if err != nil {
// 		return player.Track{}, err
// 	}

// 	// Transcoding mp4 to mp3 before encoding the result to DCA,
// 	// I have no idea what I'm doing so the sound quality is shit with mp4->DCA
// 	audioData := &bytes.Buffer{}
// 	err = videoaudio.TranscodeVideoToAudio(videoData, audioData)
// 	if err != nil {
// 		return player.Track{}, err
// 	}

// 	encodeSession, err := videoaudio.EncodeAudioToDCA(audioData)
// 	if err != nil {
// 		return player.Track{}, err
// 	}

// 	thumbnailURL := videoInfo.GetThumbnailURL(ytdl.ThumbnailQualityDefault).String()

// 	return player.Track{
// 		Audio:        encodeSession,
// 		Title:        videoInfo.Title,
// 		ID:           videoInfo.ID,
// 		Duration:     videoInfo.Duration,
// 		ThumbnailURL: thumbnailURL,
// 		URL:          url,
// 	}, nil
// }

// Download youtube video, transcode the video into audio to the writer
func Download(videoInfo ytdl.VideoInfo, format ytdl.Format, writer io.Writer) error {
	videoBuffer := &bytes.Buffer{}
	if err := videoInfo.Download(format, videoBuffer); err != nil {
		return err
	}

	err := videoaudio.TranscodeVideoToAudio(videoBuffer, writer)
	if err != nil {
		return err
	}

	return nil
}
