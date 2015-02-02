/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package config provides application configuration structure and loading logic
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type (
	// Redis struture
	Redis struct {
		Hosts    []string `json:"hosts"`
		Password string   `json:"password"`
		DB       int64    `json:"db"`
		PoolSize int      `json:"pool_size"`
	}

	// Config structure for the application configuration
	Config struct {
		Environment    string `json:"env"`
		UseArtwork     bool   `json:"use_artwork"`
		ListenHostPort string `json:"listenHost"`
		Newrelic       struct {
			Key     string `json:"key"`
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
		} `json:"newrelic"`
		Redis *Redis `json:"redis"`
	}
)

var cfg *Config

// defaultConfig returns the default configuration. It will be overwritten by the config from the user
func defaultConfig() *Config {
	cfg := &Config{}
	cfg.Environment = "dev"
	cfg.UseArtwork = true
	cfg.ListenHostPort = ":8082"

	cfg.Newrelic.Key = "demo"
	cfg.Newrelic.Name = "tapglue - stub"

	cfg.Redis = &Redis{}
	cfg.Redis.Hosts = append(cfg.Redis.Hosts, "127.0.0.1:6379")

	return cfg
}

// validate the config or panic if needed
func (config *Config) validate() {
	// TODO Implement this
}

// Load loads the configuration for the application.
//
// The name of the config file must be "config.json"
//
// It first tries to load the config file from the environment variable that you pass as the argument.
// If the environment variable doesn't exist or it's empty it then tries to use the directory where the binary file is.
//
// If the file is not present or it's not a valid json file the the call fails as well.
func (config *Config) Load(configEnvPath string) {
	// Read config path from environment variable
	configDir := ""
	if configEnvPath != "" {
		configDir = os.Getenv(configEnvPath)
	}

	// If empty set path to path of current file
	if configDir == "" {
		_, currentFilename, _, ok := runtime.Caller(2)
		if !ok {
			panic("Could not retrieve the caller for loading config")
		}

		configDir = path.Dir(currentFilename)
	}

	// Get the default configuration
	cfg = defaultConfig()

	// Read config.json
	configContents, err := ioutil.ReadFile(path.Join(configDir, "config.json"))
	if err != nil {
		configContents = []byte(os.Getenv("TAPGLUE_CONFIG_VARS"))

		if len(configContents) < 1 {
			panic(fmt.Errorf("no suitable config file was found"))
		}
	}

	// Overwrite with user configuration from file
	if err := json.Unmarshal(configContents, cfg); err != nil {
		panic(err)
	}

	// Validate configuration
	config.validate()
}

// NewConf will load and return the config
func NewConf(configEnvPath string) *Config {
	cfg.Load(configEnvPath)

	return cfg
}

// Conf will return the config
//
// If the config is not loaded already, it will attempt to load it from the directory of the binary
func Conf() *Config {
	// Last mile defence, try to load the config from the current binary directory if it's not loaded yet
	if cfg.Environment == "" {
		cfg.Load("")
	}

	return cfg
}
