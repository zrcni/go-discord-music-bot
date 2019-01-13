package videoaudio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/jonas747/dca"
	"github.com/pkg/errors"
	"github.com/zrcni/go-discord-music-bot/utils"
)

// TranscodeVideoToAudio transcodes video to audio (mp4 to mp3 currently)
func TranscodeVideoToAudio(input io.Reader, output io.Writer) error {
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
	// defer log.Println("Stderr output:", stderrError)

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

	log.Print("encoding from buffer")
	encodeSession, err := dca.EncodeMem(input, options)
	if err != nil {
		log.Printf("dca.EncodeMem: %v", err)
		return nil, err
	}

	return encodeSession, nil
}

// EncodeAudioFileToDCA encodes audio (in memory) to dca format
func EncodeAudioFileToDCA(filePath string) (*dca.EncodeSession, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.BufferedFrames = 500
	options.Application = "lowdelay"

	log.Print("encoding file")
	encodeSession, err := dca.EncodeFile(filePath, options)
	if err != nil {
		log.Printf("dca.EncodeMem: %v", err)
		return nil, err
	}

	return encodeSession, nil
}

// SaveAudioToFile saves audio data to a file
func SaveAudioToFile(name string, b []byte) error {
	filename := fmt.Sprintf("%s.%s", name, "mp3")
	filePath, err := GetTrackFilepath(filename)
	if err != nil {
		return err
	}

	log.Printf("creating file to path: %s", filePath)
	err = ioutil.WriteFile(filePath, b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ReadAudioFile save file data to a writer
func ReadAudioFile(name string) (*dca.EncodeSession, error) {
	filename := fmt.Sprintf("%s.%s", name, "mp3")
	filePath, err := GetTrackFilepath(filename)
	if err != nil {
		return nil, err
	}

	encodeSession, err := EncodeAudioFileToDCA(filePath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Printf("couldn't remove file from path: %s.\n%v", filePath, err)
	}

	return encodeSession, nil
}

// GetTrackFilepath gets the filepath of a track
func GetTrackFilepath(filenameWithExt string) (string, error) {
	basePath, err := utils.GetBasePath()
	if err != nil {
		return "", err
	}

	trackDirectory := "temp"
	// /**/tracks/<nameWithExtension>
	filePath := fmt.Sprintf("%s/%s/%s", basePath, trackDirectory, filenameWithExt)

	return filePath, nil
}
