package config

import (
	"fmt"
	"os"
)

type APISetting struct {
	Port          string
	ListenAddress string
	AllowOrigins  string
	AllowMethods  string
	AllowHeaders  string
}

func (a *APISetting) GetListenAddress() string {
	return fmt.Sprintf("%s:%s", a.ListenAddress, a.Port)
}

func LoadAPISettings() *APISetting {
	loadedConfigs := &APISetting{}

	loadedConfigs.ListenAddress = os.Getenv("OSPM_API_LISTEN_ADDRESS")
	if loadedConfigs.ListenAddress == "" {
		loadedConfigs.ListenAddress = "127.0.0.1"
	}

	loadedConfigs.Port = os.Getenv("OSPM_API_LISTEN_PORT")
	if loadedConfigs.Port == "" {
		loadedConfigs.Port = "9898"
	}

	loadedConfigs.AllowOrigins = os.Getenv("OSPM_API_ALLOW_ORIGIN")
	if loadedConfigs.AllowOrigins == "" {
		loadedConfigs.AllowOrigins = "*"
	}

	loadedConfigs.AllowMethods = os.Getenv("OSPM_API_ALLOW_METHOS")
	if loadedConfigs.AllowMethods == "" {
		loadedConfigs.AllowMethods = "*"
	}

	loadedConfigs.AllowHeaders = os.Getenv("OSPM_API_ALLOW_HEADERS")
	if loadedConfigs.AllowHeaders == "" {
		loadedConfigs.AllowHeaders = "*"
	}

	return loadedConfigs
}
