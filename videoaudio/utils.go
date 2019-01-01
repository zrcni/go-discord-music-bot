package videoaudio

import (
	"bytes"
	"fmt"
	"log"

	"github.com/jonas747/dca"
	"github.com/xfrr/goffmpeg/transcoder"
)

// TranscodeVideoToAudio transcodes video file to audio
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

// EncodeAudioToDCA encodes audio (in memory) to dca format
func EncodeAudioToDCA(data []byte) (*dca.EncodeSession, error) {
	reader := bytes.NewReader(data)

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.BufferedFrames = 500
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeMem(reader, options)
	if err != nil {
		log.Printf("dca.EncodeMem: %v", err)
		return nil, err
	}

	return encodeSession, nil
}
