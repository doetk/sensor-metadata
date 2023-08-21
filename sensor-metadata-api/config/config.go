package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ServerConfig *ServerConfig `json:"server_config"`
	DBConfig     *DBConfig     `json:"db_config"`
}

type ServerConfig struct {
	Addr            string `json:"addr"`
	ReadTimeoutSec  int    `json:"read_timeout_sec"`
	WriteTimeoutSec int    `json:"write_timeout_sec"`
	IdleTimeoutSec  int    `json:"idle_timeout_sec"`
}

type DBConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DBName     string `json:"db_name"`
	SchemaName string `json:"schema_name"`
	SSLMode    string `json:"ssl_mode"`
}

// InitConfig initializes and returns the configuration data based on defaults and config.json file.
func InitConfig(workingDir string) *Configuration {
	var err error
	cfg, err := loadConfig(workingDir + "/config.json")
	if err != nil {
		log.Printf("failed to load config file, only defaults are active: %v", err)
	}
	if cfg.ServerConfig == nil {
		log.Fatal("missing required configurations")
	}
	return cfg
}

// loadConfig loads default config and config from file, if found.
func loadConfig(baseConfig string) (*Configuration, error) {
	configuration := defaultConfig()

	confFile, err := os.Open(baseConfig)
	if err != nil {
		return configuration, err
	}

	decoder := json.NewDecoder(confFile)
	err = decoder.Decode(configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

// defaultConfig returns default configuration values.
func defaultConfig() *Configuration {
	return &Configuration{
		ServerConfig: &ServerConfig{
			Addr:            ":8080",
			ReadTimeoutSec:  90,
			WriteTimeoutSec: 90,
			IdleTimeoutSec:  0,
		},
	}
}
