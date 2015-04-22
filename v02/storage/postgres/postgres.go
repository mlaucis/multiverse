/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package postgres supports various functionality for using PostgreSQL as a storage backend
package postgres

import (
	"database/sql"
	"math/rand"

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
	}

	cli struct {
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

// New constructs a new PostgreSQL client and returns it
//
// A ConnectionURL can be the following:
// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
//
// For information on what connection parameters can be set check
// http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func New(mainConnectionURL string, slaveConnetionsURLs []string) Client {
	db, err := sql.Open("postgres", mainConnectionURL)
	if err != nil {
		errors.Fatal(err)
	}

	result := &cli{
		mainPg: db,
	}

	for idx := range slaveConnetionsURLs {
		db, err := sql.Open("postgres", slaveConnetionsURLs[idx])
		if err != nil {
			errors.Fatal(err)
		}
		result.slavePg = append(result.slavePg, db)
	}

	return result
}
