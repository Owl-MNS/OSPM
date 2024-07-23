package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type CockRoachDBConfig struct {
	DBName                    string
	Username                  string
	Password                  string
	Address                   string
	Port                      string
	SSLMode                   string
	ClientKeyPath             string
	ClientCertPath            string
	CACertPath                string
	MaxIdleConnection         int
	MaxIdleConnectionLifeTime int
	MaxOpenConnection         int
	MaxConnectionLifeTime     int
}

// OSPM_COCKROACHDB_MAX_IDLE_CONNECTIONS=20
// OSPM_COCKROACHDB_MAX_OPEN_CONNECTIONS=100
// OSPM_COCKROACHDB_MAX_CONNECTION_LIFE_TIME=1800

func LoadCockroachDBConfigs() *CockRoachDBConfig {
	var err error
	loadedConfig := &CockRoachDBConfig{}

	loadedConfig.DBName = os.Getenv("OSPM_COCKROACHDB_DB_NAME")
	if loadedConfig.DBName == "" {
		loadedConfig.DBName = "ospm"
	}

	loadedConfig.Username = os.Getenv("OSPM_COCKROACHDB_USERNAME")
	if loadedConfig.Username == "" {
		loadedConfig.Username = "root"
	}

	// password can be empty
	loadedConfig.Password = os.Getenv("OSPM_COCKROACHDB_PASSWORD")

	loadedConfig.Address = os.Getenv("OSPM_COCKROACHDB_ADDRESS")
	if loadedConfig.Address == "" {
		loadedConfig.Address = "127.0.0.1"
	}

	loadedConfig.Port = os.Getenv("OSPM_COCKROACHDB_PORT")
	if loadedConfig.Port == "" {
		loadedConfig.Port = "26257"
	}

	loadedConfig.SSLMode = os.Getenv("OSPM_COCKROACHDB_SSL_MODE")
	if !(loadedConfig.SSLMode != "" || loadedConfig.SSLMode == "disabled" || loadedConfig.SSLMode == "verify-full") {
		loadedConfig.SSLMode = "disabled"
	}

	loadedConfig.ClientKeyPath = os.Getenv("OSPM_COCKROACHDB_SSL_CLIENT_KEY_PATH")
	if loadedConfig.ClientKeyPath == "" {
		loadedConfig.ClientKeyPath = "/etc/roachCerts/client.key"
	}

	loadedConfig.ClientCertPath = os.Getenv("OSPM_COCKROACHDB_SSL_CLIENT_CERT_PATH")
	if loadedConfig.ClientCertPath == "" {
		loadedConfig.ClientCertPath = "/etc/roachCerts/client.crt"
	}

	loadedConfig.CACertPath = os.Getenv("OSPM_COCKROACHDB_SSL_CA_CERT_PATH")
	if loadedConfig.CACertPath == "" {
		loadedConfig.CACertPath = "/etc/roachCerts/ca.crt"
	}

	loadedConfig.MaxIdleConnection, err = strconv.Atoi(os.Getenv("OSPM_COCKROACHDB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		loadedConfig.MaxIdleConnection = 20
	}

	loadedConfig.MaxIdleConnectionLifeTime, err = strconv.Atoi(os.Getenv("OSPM_COCKROACHDB_MAX_IDLE_CONNECTIONS_LIFE_TIME"))
	if err != nil {
		loadedConfig.MaxIdleConnection = 1800
	}

	loadedConfig.MaxOpenConnection, err = strconv.Atoi(os.Getenv("OSPM_COCKROACHDB_MAX_OPEN_CONNECTIONS"))
	if err != nil {
		loadedConfig.MaxIdleConnection = 100
	}

	loadedConfig.MaxConnectionLifeTime, err = strconv.Atoi(os.Getenv("OSPM_COCKROACHDB_MAX_CONNECTION_LIFE_TIME"))
	if err != nil {
		loadedConfig.MaxIdleConnection = 1800
	}

	return loadedConfig
}

func (roach *CockRoachDBConfig) DSN() string {
	if roach.SSLMode == "disabled" {
		return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			roach.Username,
			url.QueryEscape(roach.Password),
			roach.Address,
			roach.Port,
			roach.DBName,
		)
	} else {
		return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s&sslcert=%s&sslkey=%s",
			roach.Username,
			url.QueryEscape(roach.Password),
			roach.Address,
			roach.Port,
			roach.DBName,
			roach.CACertPath,
			roach.ClientCertPath,
			roach.ClientKeyPath)
	}
}
