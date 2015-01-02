/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package db provides configuration for database connection
package db

import (
	"sync"

	"github.com/tapglue/backend/config"

	_ "github.com/go-sql-driver/mysql" // Get the MySQL driver
	"github.com/jmoiron/sqlx"
)

type (
	// Keep the database connection and the number of times it was used
	dbSlave struct {
		Usage uint64
		DB    *sqlx.DB
	}

	// Keep the database slave connections and a mutex to be able to safely keep track of the least used connection
	dbSlaves struct {
		sync.Mutex
		MinSlave int
		Slaves   []*dbSlave
	}
)

var (
	masterConnection = &sqlx.DB{}
	slaveConnections = &dbSlaves{}
)

// openMasterConnection opens a connection to the master database
func openMasterConnection(cfg config.DB) {
	// Read settings from configuration
	masterDSN := cfg.MasterDSN()

	// Establish connection to master
	masterConnection = sqlx.MustConnect("mysql", masterDSN)
	masterConnection.DB.Ping()
	masterConnection.DB.SetMaxIdleConns(cfg.MaxIdleConnections())
	masterConnection.DB.SetMaxOpenConns(cfg.MaxOpenConnections())
}

// openSlaveConnections opens a connection to each of the slave databases
func openSlaveConnections(cfg config.DB) {
	slaves := cfg.SlavesDSN()

	for _, slaveDSN := range slaves {
		slaveConnection := &sqlx.DB{}

		// Establish connection to slaves
		slaveConnection = sqlx.MustConnect("mysql", slaveDSN)
		slaveConnection.DB.Ping()
		slaveConnection.DB.SetMaxIdleConns(cfg.MaxIdleConnections())
		slaveConnection.DB.SetMaxOpenConns(cfg.MaxOpenConnections())

		slaveConnections.Slaves = append(slaveConnections.Slaves, &dbSlave{DB: slaveConnection})
	}
}

// InitDatabases initializes the connections to the servers
func InitDatabases(cfg config.DB) {
	openMasterConnection(cfg)
	openSlaveConnections(cfg)
}

// GetMaster returns the connection to the master database.
// It is used to write to the database.
func GetMaster() *sqlx.DB {
	return masterConnection
}

// GetSlave is used to read from database.
// If there's no slave configured, it returns master.
func GetSlave() *sqlx.DB {
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
