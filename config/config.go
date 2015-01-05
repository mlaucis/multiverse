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
	// DB interface
	DB interface {
		IsMasterDebug() bool
		IsSlaveDebug(slaveID uint) bool
		MasterDSN() string
		SlavesDSN() []string
		MaxIdleConnections() int
		MaxOpenConnections() int
	}

	// Config interface
	Config interface {
		Load(configEnvPath string)
		Valiate()
		Env() string
		ListenHost() string
		NewRelic() (string, string)
		DB() *DB
	}

	// Db structure
	Db struct {
		Username string `json:"username"`
		Password string `json:"password`
		Database string `json:"database"`
		MaxIdle  int    `json:"max_idle"`
		MaxOpen  int    `json:"max_open"`
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
	}

	// Cfg structure for the application configuration
	Cfg struct {
		Environment    string `json:"env"`
		ListenHostPort string `json:"listenHost"`
		Newrelic       struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"newrelic"`
		Database *Db `json:"db"`
	}
)

var cfg *Cfg

// getDefaultConfig returns the default configuration. It will be overwritten by the config from the user
func defaultConfig() *Cfg {
	cfg := &Cfg{}
	cfg.Environment = "dev"
	cfg.ListenHostPort = ":8082"

	cfg.Newrelic.Key = "demo"
	cfg.Newrelic.Name = "tapglue - stub"

	cfg.Database = &Db{}
	cfg.Database.Username = ""
	cfg.Database.Password = ""
	cfg.Database.Database = ""

	cfg.Database.MaxIdle = 10
	cfg.Database.MaxOpen = 300

	cfg.Database.Master.Debug = true
	cfg.Database.Master.Host = ""
	cfg.Database.Master.Port = 0

	cfg.Database.Slaves = append(cfg.Database.Slaves, struct {
		Debug bool   `json:"debug"`
		Host  string `json:"host"`
		Port  uint   `json:"port"`
	}{
		Debug: true,
		Host:  "",
		Port:  0,
	},
	)

	return cfg
}

// Env returns the environment of the application
func (config *Cfg) Env() string {
	return config.Environment
}

// ListenHost returns the host:port combination for the main server
func (config *Cfg) ListenHost() string {
	if os.Getenv("IS_HEROKU_ENV") != "" {
		return fmt.Sprintf(":%s", os.Getenv("PORT"))
	}
	return config.ListenHostPort
}

func (config *Cfg) NewRelic() (string, string) {
	return config.Newrelic.Key, config.Newrelic.Name
}

// DB returns a database interface
func (config *Cfg) DB() DB {
	return config.Database
}

// IsMasterDebug returns if the master database is set to debug
func (database *Db) IsMasterDebug() bool {
	return database.Master.Debug
}

// IsSlaveDebug returns if the specified slave is set to debug
func (database *Db) IsSlaveDebug(slaveID uint) bool {
	if uint(len(database.Slaves)-1) > slaveID {
		return false
	}

	return database.Slaves[slaveID].Debug
}

// MasterDSN returns the master DSN connection string
func (database *Db) MasterDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8&collation=utf8_general_ci",
		database.Username,
		database.Password,
		database.Master.Host,
		database.Master.Port,
		database.Database,
	)
}

// SlavesDSN returns the DSN for all the slaves
func (database *Db) SlavesDSN() []string {
	result := []string{}
	for _, slave := range database.Slaves {

		slaveDSN := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8&collation=utf8_general_ci",
			database.Username,
			database.Password,
			slave.Host,
			slave.Port,
			database.Database,
		)

		result = append(result, slaveDSN)
	}

	return result
}

// MaxIdleConnections returns the number of maximum idle connections for the database
func (database *Db) MaxIdleConnections() int {
	return database.MaxIdle
}

// MaxOpenConnections returns the number of maximum open connections to the database
func (database *Db) MaxOpenConnections() int {
	return database.MaxOpen
}

// Validate should be implemented to add config validation or panic if needed
func (config *Cfg) Validate() {

}

// Load loads the configuration for the application.
//
// The name of the config file must be "config.json"
//
// It first tries to load the config file from the environment variable that you pass as the argument.
// If the environment variable doesn't exist or it's empty it then tries to use the directory where the binary file is.
//
// If the file is not present or it's not a valid json file the the call fails as well.
func (config *Cfg) Load(configEnvPath string) {
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
	config.Validate()
}

// NewConf will load and return the config
func NewConf(configEnvPath string) *Cfg {
	cfg.Load(configEnvPath)

	return cfg
}

// Conf will return the config
//
// If the config is not loaded already, it will attempt to load it from the directory of the binary
func Conf() *Cfg {
	// Last mile defence, try to load the config from the current binary directory if it's not loaded yet
	if cfg.Env() == "" {
		cfg.Load("")
	}

	return cfg
}
