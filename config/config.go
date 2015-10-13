// Package config provides application configuration structure and loading logic
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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

	// Kinesis structure
	Kinesis struct {
		AuthKey    string `json:"auth_key"`
		SecretKey  string `json:"secret_key"`
		Region     string `json:"region"`
		Endpoint   string `json:"endpoint"`
		StreamName string `json:"stream_name"`
	}

	// PostgresDB structure
	PostgresDB struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Options  string `json:"options"`
	}

	// Postgres structure
	Postgres struct {
		Database string       `json:"database"`
		Master   PostgresDB   `json:"master"`
		Slaves   []PostgresDB `json:"slaves"`
	}

	// Config structure for the application configuration
	Config struct {
		Environment    string    `json:"env"`
		UseArtwork     bool      `json:"use_artwork"`
		UseSysLog      bool      `json:"use_syslog"`
		UseSSL         bool      `json:"use_ssl"`
		SkipSecurity   bool      `json:"skip_security"`
		JSONLogs       bool      `json:"json_logs"`
		ListenHostPort string    `json:"listenHost"`
		TelemetryAddr  string    `json:"telemetry_addr"`
		CacheApp       *Redis    `json:"cache_app"`
		RateLimiter    *Redis    `json:"rate_limiter"`
		Kinesis        *Kinesis  `json:"kinesis"`
		Postgres       *Postgres `json:"postgres"`
	}
)

var cfg *Config

// defaultConfig returns the default configuration. It will be overwritten by the config from the user
func defaultConfig() *Config {
	cfg := &Config{}
	cfg.Environment = "dev"
	cfg.UseArtwork = true
	cfg.UseSSL = false
	cfg.UseSysLog = true
	cfg.SkipSecurity = false
	cfg.JSONLogs = false
	cfg.ListenHostPort = ":8082"
	cfg.TelemetryAddr = ":9000"

	cfg.RateLimiter = &Redis{}
	cfg.RateLimiter.Hosts = append(cfg.RateLimiter.Hosts, "127.0.0.1:6379")

	cfg.CacheApp = &Redis{}
	cfg.CacheApp.Hosts = append(cfg.CacheApp.Hosts, "127.0.0.1:6379")

	cfg.Kinesis = &Kinesis{}
	cfg.Kinesis.AuthKey = ""
	cfg.Kinesis.SecretKey = ""
	cfg.Kinesis.Region = "eu-central-1"
	cfg.Kinesis.Endpoint = ""

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
		var err error
		configDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the default configuration
	cfg = defaultConfig()

	// Read config.json
	configContents, err := ioutil.ReadFile(path.Join(configDir, "config.json"))
	if err != nil {
		_, currentFilename, _, ok := runtime.Caller(2)
		if !ok {
			panic("Could not retrieve the caller for loading config")
		}

		configDir = path.Dir(currentFilename)
		configContents, err = ioutil.ReadFile(path.Join(configDir, "config.json"))
		if err != nil {
			configContents = []byte(os.Getenv("TAPGLUE_CONFIG_VARS"))

			if len(configContents) < 1 {
				panic(fmt.Errorf("no suitable config file was found"))
			}
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
