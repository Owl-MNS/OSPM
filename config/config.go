package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type OSPMConfig struct {
	API            *APISetting
	Logrus         *LogrusConfig
	RDBMS          *CockRoachDBConfig
	ClientPolicies *ClientPolicy
}

var OSPM *OSPMConfig

func LoadLocalEnvironments() {

	configFile := GetConfigFilePath()
	if err := godotenv.Load(configFile); err != nil {
		log.Printf("failed to load the config file under %s, Default Values will be used. error: %s", configFile, err)
	}

}

// LoadOSPMConfigs loads all of the configurations defined
// in the config file so then can be accessible fomrconfig.OSPMConfigs
func LoadOSPMConfigs() {
	LoadLocalEnvironments()

	OSPM = &OSPMConfig{
		API:            LoadAPISettings(),
		Logrus:         LoadLogrusConfigs(),
		RDBMS:          LoadCockroachDBConfigs(),
		ClientPolicies: LoadClientPolicies(),
	}
}

// GetConfigFilePath checks two paths for the config.env file
// It should be either at the root of the server (for contaierized deployment)
// or at the root of the project (for development or manual deployments)
func GetConfigFilePath() string {
	configFilePath := "/config.env"

	_, err := os.Stat(configFilePath)
	if err != nil || os.IsNotExist(err) {
		log.Printf("failed to load config file unsed %s, considering another path: config.env (config file at the root of execution directory), error: %s", configFilePath, err)
		return "config.env"
	}

	log.Printf("checking config file under %s (config file at the root of server/container/pod)", configFilePath)
	return configFilePath
}
