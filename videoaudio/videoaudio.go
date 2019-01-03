package videoaudio

import (
	"bytes"
	"io"
	"log"
	"os/exec"

	"github.com/jonas747/dca"
	"github.com/pkg/errors"
)

// TranscodeVideoToAudio transcodes video to audio (mp4 to mp3 currently)
func TranscodeVideoToAudio(input io.Reader, filename string, output io.Writer) error {
	args := []string{
		"-i", "pipe:0",
		"-f", "mp3",
		"-ab", "128000",
		"-vn",
		"pipe:1",
	}

	ffmpeg := exec.Command("ffmpeg", args...)

	ffmpeg.Stdin = input

	ffmpeg.Stdout = output

	stderrError := &bytes.Buffer{}
	ffmpeg.Stderr = stderrError
	defer log.Println("Stderr output:", stderrError)

	err := ffmpeg.Start()
	if err != nil {
		return errors.Wrap(err, "ffmpeg.Start")
	}

	err = ffmpeg.Wait()
	if err != nil {
		log.Println("ffmpeg wait", err)
		return errors.Wrap(err, "ffmpeg.Wait")
	}
	return nil
}

// EncodeAudioToDCA encodes audio (in memory) to dca format
func EncodeAudioToDCA(input io.Reader) (*dca.EncodeSession, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.BufferedFrames = 500
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeMem(input, options)
	if err != nil {
		log.Printf("dca.EncodeMem: %v", err)
		return nil, err
	}

	return encodeSession, nil
}
