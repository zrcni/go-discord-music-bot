package videoaudio

import (
	"fmt"
	"log"

	"github.com/xfrr/goffmpeg/transcoder"
)

func TranscodeVideoToAudio(filename string) error {
	videoExtension := "mp4"
	audioExtension := "mp3"
	inputVideoPath := fmt.Sprintf("%s.%s", filename, videoExtension)
	outputAudioPath := fmt.Sprintf("%s.%s", filename, audioExtension)

	trans := &transcoder.Transcoder{}

	err := trans.Initialize(inputVideoPath, outputAudioPath)
	if err != nil {
		return err
	}

	done := trans.Run(false)

	err = <-done
	if err != nil {
		return err
	}

	log.Printf(fmt.Sprintf("Transcoded %s successfully from %s to %s", filename, videoExtension, audioExtension))
	return nil
}
