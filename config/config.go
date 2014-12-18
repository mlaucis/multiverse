/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type Config struct {
	Env        string `json: "env"`
	ListenHost string `json: "listenHost"`
}

var cfg *Config

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
