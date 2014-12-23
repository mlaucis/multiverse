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

// Get the default configuration. It will be overwritten by the config from the user
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

// This function should be implemented to add config validation
// or panic if needed
func validateConfig() {

}

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

	cfg = getDefaultConfig()

	if err := json.Unmarshal(file, cfg); err != nil {
		panic(err)
	}

	validateConfig()

	return cfg
}
