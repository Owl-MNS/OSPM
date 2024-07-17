package logger

import (
	"os"
	"ospm/config"

	"github.com/sirupsen/logrus"
)

var OSPMLogger *logrus.Logger

func InitLogger() {
	OSPMLogger = logrus.New()

	// Set the log level based on an environment variable
	level, err := logrus.ParseLevel(config.OSPM.Logrus.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	OSPMLogger.SetLevel(level)

	// Set output to stdout
	OSPMLogger.SetOutput(os.Stdout)

	// Set a custom format (JSON is also an option)
	OSPMLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
