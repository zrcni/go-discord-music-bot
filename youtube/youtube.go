package youtube

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/rylio/ytdl"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
)

const youtubeURLMatcher = "^.*(?:(?:youtu\\.be\\/|v\\/|vi\\/|u\\/\\w\\/|embed\\/)|(?:(?:watch)?\\?v(?:i)?=|\\&v(?:i)?=))([^#\\&\\?]*).*"

// Audio stores data and info
type Audio struct {
	Data []byte
	Info *ytdl.VideoInfo
}

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
func GetMetadata(url string) (*ytdl.VideoInfo, error) {
	videoInfo, err := ytdl.GetVideoInfo(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return videoInfo, nil
}

// Download a youtube video
func downloadVideo(writer io.Writer, videoInfo *ytdl.VideoInfo) error {
	if len(videoInfo.Formats) == 0 {
		return errors.New("No available video formats")
	}

	format := videoInfo.Formats[0]

	log.Printf("Format: %v", format)

	if err := videoInfo.Download(format, writer); err != nil {
		return err
	}
	return nil
}

// Get downloads video from youtube and transcodes it to audio
func Get(url string) (*Audio, error) {
	videoInfo, err := GetMetadata(url)
	if err != nil {
		return nil, err
	}

	videoFilenameWithExtension := fmt.Sprintf("%s.mp4", videoInfo.ID)

	videoFile, err := os.Create(videoFilenameWithExtension)
	if err != nil {
		return nil, err
	}
	defer videoFile.Close()
	defer os.Remove(fmt.Sprintf("%s.mp4", videoInfo.ID))

	err = downloadVideo(videoFile, videoInfo)
	if err != nil {
		return nil, err
	}

	err = videoaudio.TranscodeVideoToAudio(videoInfo.ID)
	if err != nil {
		return nil, err
	}
	defer os.Remove(fmt.Sprintf("%s.mp3", videoInfo.ID))

	audioFilenameWithExtension := fmt.Sprintf("%s.mp3", videoInfo.ID)
	data, err := ioutil.ReadFile(audioFilenameWithExtension)
	if err != nil {
		return nil, err
	}

	return &Audio{
		Data: data,
		Info: videoInfo,
	}, nil
}
