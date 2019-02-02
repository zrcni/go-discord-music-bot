package videoaudio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/jonas747/dca"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/utils"
)

var encodingOptions = &dca.EncodeOptions{
	Volume:           256,
	Channels:         2,
	FrameRate:        48000,
	FrameDuration:    20,
	CompressionLevel: 10,
	PacketLoss:       1,
	VBR:              true,
	Application:      dca.AudioApplicationLowDelay,
	RawOutput:        true,
	Bitrate:          128,
	BufferedFrames:   500,
}

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
	// defer log.Debugf("Stderr output:", stderrError)

	err := ffmpeg.Start()
	if err != nil {
		return errors.Wrap(err, "ffmpeg.Start")
	}

	err = ffmpeg.Wait()
	if err != nil {
		log.Errorf("ffmpeg wait %v", err)
		return errors.Wrap(err, "ffmpeg.Wait")
	}
	return nil
}

// EncodeAudioToDCA encodes audio (in memory) to dca format
func EncodeAudioToDCA(input io.Reader) (*dca.EncodeSession, error) {
	log.Debug("encoding from buffer")
	encodeSession, err := dca.EncodeMem(input, encodingOptions)
	if err != nil {
		log.Errorf("dca.EncodeMem: %v", err)
		return nil, err
	}

	return encodeSession, nil
}

// EncodeAudioFileToDCA encodes audio (in memory) to dca format
func EncodeAudioFileToDCA(filePath string) (*dca.EncodeSession, error) {
	log.Debug("encoding file")
	encodeSession, err := dca.EncodeFile(filePath, encodingOptions)
	if err != nil {
		log.Errorf("dca.EncodeMem: %v", err)
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

	log.Debugf("creating file to path: %s", filePath)
	err = ioutil.WriteFile(filePath, b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ReadAudioFile save file data to a writer
func ReadAudioFile(filenameWithExt string) (*dca.EncodeSession, error) {
	filePath, err := GetTrackFilepath(filenameWithExt)
	if err != nil {
		return nil, err
	}

	encodeSession, err := EncodeAudioFileToDCA(filePath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Errorf("couldn't remove file from path: %s.\n%v", filePath, err)
	}

	return encodeSession, nil
}

// ReadAudioFilePath save file data from path to a writer
func ReadAudioFilePath(path string) (*dca.EncodeSession, error) {
	encodeSession, err := EncodeAudioFileToDCA(path)
	if err != nil {
		return nil, err
	}

	err = os.Remove(path)
	if err != nil {
		log.Errorf("couldn't remove file from path: %s.\n%v", path, err)
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
