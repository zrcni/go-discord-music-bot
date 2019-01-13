package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/utils"
)

// Setup sets up logger for writing to stdout and a file
func Setup() {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{}
	logger := l.Writer()

	// timestampMs := time.Now().UnixNano() / 1000000
	// path := fmt.Sprintf("log-%v", timestampMs)
	// sessionFile := openFile(path)

	fullLogFile := openFile("logs")

	mw := io.MultiWriter(logger, fullLogFile)

	log.SetOutput(mw)
}

func openFile(filename string) *os.File {
	basePath, err := utils.GetBasePath()
	if err != nil {
		panic(err)
	}

	logPath := fmt.Sprintf("%s/logs/%s.txt", basePath, filename)

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return file
}
