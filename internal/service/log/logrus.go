package log

import (
	"os"
	"ospm/config"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set the log level based on an environment variable
	level, err := logrus.ParseLevel(config.OSPM.Logrus.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Set output to stdout
	Log.SetOutput(os.Stdout)

	// Set a custom format (JSON is also an option)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
