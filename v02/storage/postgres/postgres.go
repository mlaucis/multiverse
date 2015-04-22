/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package postgres supports various functionality for using PostgreSQL as a storage backend
package postgres

import (
	"database/sql"

	"github.com/tapglue/backend/tgerrors"

	_ "github.com/lib/pq"
)

type (
	//Client interface to define PostgreSQL methods
	Client interface {
		// Returns the raw database connection
		Datastore() *sql.DB
	}

	cli struct {
		postgres *sql.DB
	}
)

func (c *cli) Datastore() {
	return c.postgres
}

// New constructs a new PostgreSQL client and returns it
//
// connectionURL can be the following:
// "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
//
// For information on what connection parameters can be set check
// http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func New(connectionURL string) Client {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		tgerrors.Fatal(err)
	}

	return &cli{
		postgres: db,
	}
}
