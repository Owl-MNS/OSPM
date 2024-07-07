package config

import (
	"os"
	"strings"
)

type LogrusConfig struct {
	LogLevel string
}

func LoadLogrusConfigs() *LogrusConfig {
	loadedConfig := &LogrusConfig{}

	loadedConfig.LogLevel = strings.ToUpper(os.Getenv("OSPM_LOG_LEVEL"))
	if loadedConfig.LogLevel == "" || !IsValid(loadedConfig.LogLevel) {
		loadedConfig.LogLevel = "INFO"
	}

	return loadedConfig
}

// IsValid checks the loaded log level config agains the valid options
func IsValid(logLevel string) bool {
	validOptions := []string{"INFO", "WARNING", "ERROR", "DEBUG"}

	for _, opt := range validOptions {
		if logLevel == opt {
			return true
		}
	}

	return false
}
