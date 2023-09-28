package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/yogenyslav/kokoc-hack/internal/config"
)

var Log *logrus.Logger

func NewLogger() error {
	Log = logrus.New()

	Log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006/01/02 - 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})

	logFile, err := os.OpenFile("backend.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	Log.SetOutput(mw)

	level, err := logrus.ParseLevel(config.Cfg.LoggingLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	Log.SetLevel(level)
	return nil
}
