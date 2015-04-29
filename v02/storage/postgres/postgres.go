/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package postgres supports various functionality for using PostgreSQL as a storage backend
package postgres

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"

	// Well, we want to have PostgreSQL as database so we kinda need this..
	_ "github.com/lib/pq"
)

type (
	//Client interface to define PostgreSQL methods
	Client interface {
		// Returns the raw main database connection
		MainDatastore() *sql.DB

		// SlaveDatastore returns a random slave database connection
		//
		// If the paramater is -1 then the returned connection is chosen out of the connection pool
		//
		// If there's no slave connection available, then the main connection is returned
		SlaveDatastore(id int) *sql.DB

		// Database returns a specific database connection. It might be nil so check for it
		Database(database string) Client
	}

	cli struct {
		master *config.PostgresDB
		slaves []config.PostgresDB

		connections map[string]*cli

		mainPg  *sql.DB
		slavePg []*sql.DB
	}
)

func (c *cli) MainDatastore() *sql.DB {
	return c.mainPg
}

func (c *cli) SlaveDatastore(id int) *sql.DB {
	if len(c.slavePg) == 0 {
		return c.mainPg
	}

	if id == -1 || id >= len(c.slavePg) {
		return c.slavePg[rand.Intn(len(c.slavePg))]
	}

	return c.slavePg[id]
}

func (c *cli) Database(database string) Client {
	if db, ok := c.connections[database]; ok {
		return db
	}

	db := &cli{
		mainPg: composeConnection(database, c.master),
	}

	for idx := range c.slaves {
		db.slavePg = append(db.slavePg, composeConnection(database, &c.slaves[idx]))
	}

	c.connections[database] = db
	return db
}

func formatConnectionURL(database string, config *config.PostgresDB) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", config.Username, config.Password, config.Host, database, config.Options)
}

func composeConnection(database string, config *config.PostgresDB) *sql.DB {
	db, err := sql.Open("postgres", formatConnectionURL(database, config))
	if err != nil {
		errors.Fatal(err)
	}

	return db
}

// New constructs a new PostgreSQL client and returns it
//
// A ConnectionURL can be the following:
// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
//
// For information on what connection parameters can be set check
// http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func New(config *config.Postgres) Client {
	result := &cli{
		mainPg: composeConnection(config.Database, &config.Master),

		master: &config.Master,
		slaves: config.Slaves,
	}

	for idx := range config.Slaves {
		result.slavePg = append(result.slavePg, composeConnection(config.Database, &config.Slaves[idx]))
	}

	return result
}
