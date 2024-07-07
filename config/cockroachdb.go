package config

import (
	"fmt"
	"net/url"
	"os"
)

type CockRoachDBConfig struct {
	DBName         string
	Username       string
	Password       string
	Address        string
	Port           string
	SSLMode        string
	ClientKeyPath  string
	ClientCertPath string
	CACertPath     string
}

func LoadCockroachDBConfigs() *CockRoachDBConfig {
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
