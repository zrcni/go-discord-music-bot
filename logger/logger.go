package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/config"
	"github.com/zrcni/go-discord-music-bot/utils"
)

// Setup sets up logger for writing to stdout and a file
func Setup() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(getLogLevel())

	logWriter := logger.Writer()
	// timestampMs := time.Now().UnixNano() / 1000000
	// path := fmt.Sprintf("log-%v", timestampMs)
	// sessionFile := openFile(path)

	fullLogFile := getFile("logs")

	mw := io.MultiWriter(logWriter, fullLogFile)

	logrus.SetOutput(mw)
}

func getLogLevel() logrus.Level {
	if config.Config.Debug {
		return logrus.DebugLevel
	}

	return logrus.InfoLevel
}

func getFile(filename string) *os.File {
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
