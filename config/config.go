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
	Env        string `json:"env"`
	ListenHost string `json:"listenHost"`
	DB         struct {
		Username string `json:"username"`
		Password string `json:"password`
		Database string `json:"database"`
		MaxIdle  uint   `json:"max_idle"`
		MaxOpen  uint   `json:"max_open"`
		Master   struct {
			Debug bool   `json:"debug"`
			Host  string `json:"host"`
			Port  uint   `json:"port"`
		} `json:"master"`
		Slaves []struct {
			Debug bool   `json:"debug"`
			Host  string `json:"host"`
			Port  uint   `json:"port"`
		} `json:"slaves"`
	} `json:"db"`
}

var cfg *Config

// getDefaultConfig returns the default configuration. It will be overwritten by the config from the user
func getDefaultConfig() *Config {
	cfg := &Config{}
	cfg.Env = "dev"
	cfg.ListenHost = ":8082"

	cfg.DB.Username = "gluee"
	cfg.DB.Password = "x"
	cfg.DB.Database = "gluee"

	cfg.DB.MaxIdle = 10
	cfg.DB.MaxOpen = 300

	cfg.DB.Master.Debug = true
	cfg.DB.Master.Host = "127.0.0.1"
	cfg.DB.Master.Port = 3306

	cfg.DB.Slaves = append(cfg.DB.Slaves, struct {
		Debug bool   `json:"debug"`
		Host  string `json:"host"`
		Port  uint   `json:"port"`
	}{
		Debug: true,
		Host:  "127.0.0.1",
		Port:  3306,
	},
	)

	return cfg
}

// validateConfig should be implemented to add config validation or panic if needed
func validateConfig() {

}

//GetConfig loads the configuration for the application.
//After the first call, it caches the values internally.
//
//The name of the config file must be "config.json"
//
//It first tries to load the config file from the environment variable that you pass as the argument.
//If the environment variable doesn't exist or it's empty it then tries to use the directory where the binary file is.
//
//If the file is not present or it's not a valid json file the the call fails as well.
func GetConfig(configPath string) *Config {
	// Return config if it's not empty
	if cfg != nil {
		return cfg
	}

	// Read config path from environment variable
	configDir := os.Getenv(configPath)

	// If empty set path to path of current file
	if configDir == "" {
		_, currentFilename, _, ok := runtime.Caller(1)
		if !ok {
			panic("Could not retrieve the caller for loading config")
		}

		configDir = path.Dir(currentFilename)
	}

	// Read config.json
	file, err := ioutil.ReadFile(path.Join(configDir, "config.json"))
	if err != nil {
		panic(err)
	}

	// Get the default configuration
	cfg = getDefaultConfig()

	// Overwrite with user configuration from file
	if err := json.Unmarshal(file, cfg); err != nil {
		panic(err)
	}

	// Validate configuration
	validateConfig()

	return cfg
}
