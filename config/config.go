/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package config provides application configuration structure and loading logic
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

// Config structure for the application configuration
type Config struct {
	Env        string `json: "env"`
	ListenHost string `json: "listenHost"`
}

var cfg *Config

// GetConfig loads the configuration for the application.
// After the first call, it caches the values internally.
//
// The name of the config file must be "config.json"
//
// It first tries to load the config file from the environment variable that you pass as the argument.
// If the environment variable doesn't exist or it's empty it then tries to use the directory where the binary file is.
//
// If the file is not present or it's not a valid json file the the call fails as well.
//
func GetConfig(configPath string) *Config {
	if cfg != nil {
		return cfg
	}

	configDir := os.Getenv(configPath)

	if configDir == "" {
		_, currentFilename, _, _ := runtime.Caller(1)
		configDir = path.Dir(currentFilename)
	}

	file, err := ioutil.ReadFile(path.Join(configDir, "config.json"))
	if err != nil {
		panic(err)
	}

	cfg = &Config{}
	if err := json.Unmarshal(file, cfg); err != nil {
		panic(err)
	}

	return cfg
}
