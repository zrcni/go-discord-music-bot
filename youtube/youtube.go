package youtube

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	ggyoutube "github.com/knadh/go-get-youtube/youtube"
	"github.com/rylio/ytdl"
)

// MetaData of a youtube video
type MetaData struct {
	Title           string `json:"title"`
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	Version         string `json:"version"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	Type            string `json:"type"`
	HTML            string `json:"html"`
}

const PlaceholderVideoURL = "https://www.youtube.com/watch?v=JrIhZPAAcjI"

const youtubeURLMatcher = "^.*(?:(?:youtu\\.be\\/|v\\/|vi\\/|u\\/\\w\\/|embed\\/)|(?:(?:watch)?\\?v(?:i)?=|\\&v(?:i)?=))([^#\\&\\?]*).*"

// func makeMetaDataURL(url string) string {
// 	 return fmt.Sprintf("https://www.youtube.com/oembed?url=%s&format=json", url)
// }

func getVideoID(url string) string {
	re := regexp.MustCompile(youtubeURLMatcher)
	match := re.FindStringSubmatch(url)

	if len(match) == 2 {
		videoID := match[1]
		return videoID
	}

	return ""
}

func getMetaData(url string) (*ggyoutube.Video, error) {
	videoID := getVideoID(url)
	video, err := ggyoutube.Get(videoID)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func GetVideoByURL(url string) *ggyoutube.Video {
	video, err := getMetaData(url)
	if err != nil {
		log.Printf("Could not get metadata: %v", err)
		return &ggyoutube.Video{}
	}

	options := &ggyoutube.Option{
		Rename: true,
		Resume: true,
		Mp3:    true,
	}

	filename := fmt.Sprintf("%s.mp3", url)

	if err := video.Download(1, filename, options); err != nil {
		os.Remove(filename)
		log.Printf("youtube video download: %v", err)
		return nil
	}

	return video
}

func Download(writer io.Writer, url string) {
	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		log.Println(err)
		return
	}

	if len(vid.Formats) == 0 {
		return
	}

	format := vid.Formats[0]

	log.Printf("Format: %v", format)
	if err := vid.Download(format, writer); err != nil {
		log.Println(err)
	}
}
