/**
 * @author Florin Patan <florinpatan@gmai.com>
 */

package db

import (
	"fmt"
	"sync"

	"github.com/gluee/backend/config"

	_ "github.com/go-sql-driver/mysql" // Get the MySQL driver
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Get a GORM dependency
)

type (
	// Keep the database connection and the number of times it was used
	dbSlave struct {
		Usage uint64
		DB    *gorm.DB
	}

	// Keep the database slave connections and a mutex to be able to safely keep track of the least used connection
	dbSlaves struct {
		sync.Mutex
		MinSlave int
		Slaves   []*dbSlave
	}
)

var (
	masterConnection = &gorm.DB{}
	slaveConnections = &dbSlaves{}
)

// Open a connection to master server
func openMasterConnection(cfg *config.Config) {
	var err error

	masterDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8&collation=utf8_general_ci",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Master.Host,
		cfg.DB.Master.Port,
		cfg.DB.Database,
	)

	if *masterConnection, err = gorm.Open("mysql", masterDSN); err != nil {
		panic(err)
	}

	masterConnection.DB().Ping()
	masterConnection.DB().SetMaxIdleConns(10)
	masterConnection.DB().SetMaxOpenConns(100)

	masterConnection.LogMode(cfg.DB.Master.Debug)
}

// Open the connections to the slave servers
func openSlaveConnections(cfg *config.Config) {
	var err error

	for _, slave := range cfg.DB.Slaves {
		slaveConnection := &gorm.DB{}
		slaveDSN := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8&collation=utf8_general_ci",
			cfg.DB.Username,
			cfg.DB.Password,
			slave.Host,
			slave.Port,
			cfg.DB.Database,
		)

		if *slaveConnection, err = gorm.Open("mysql", slaveDSN); err != nil {
			panic(err)
		}

		slaveConnection.DB().Ping()
		slaveConnection.DB().SetMaxIdleConns(10)
		slaveConnection.DB().SetMaxOpenConns(100)

		slaveConnection.LogMode(cfg.DB.Master.Debug)

		slaveConnections.Slaves = append(slaveConnections.Slaves, &dbSlave{DB: slaveConnection})
	}
}

// InitDatabases initializes the connections to the servers
func InitDatabases(cfg *config.Config) {
	openMasterConnection(cfg)
	openSlaveConnections(cfg)
}

// GetMaster returns the master database connection.
//
// You should use this when you want to write to the database
func GetMaster() *gorm.DB {
	return masterConnection
}

// GetSlave returns a slave database connection from the connection pool.
//
// You should use this when you want to only read from the database
//
// If there's no slave configured, it returns master
func GetSlave() *gorm.DB {
	if len(slaveConnections.Slaves) == 0 {
		return masterConnection
	}

	// Make sure we don't select the wrong master after we finish counting the current least used one
	slaveConnections.Lock()

	min := slaveConnections.MinSlave
	slaveConnections.Slaves[min].Usage = slaveConnections.Slaves[min].Usage + 1
	newMin := min
	minVal := slaveConnections.Slaves[min].Usage

	// Find the least used slave
	for key, slave := range slaveConnections.Slaves {
		if slave.Usage < minVal {
			newMin = key
			minVal = slave.Usage
		}
	}
	slaveConnections.MinSlave = newMin

	slaveConnections.Unlock()

	return slaveConnections.Slaves[min].DB
}
