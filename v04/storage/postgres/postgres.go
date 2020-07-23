// Package postgres supports various functionality for using PostgreSQL as a storage backend
package postgres

import (
	"fmt"
	"math/rand"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/errors"

	"github.com/jmoiron/sqlx"
	// Well, we want to have PostgreSQL as database so we kinda need this..
	_ "github.com/lib/pq"
)

type (
	//Client interface to define PostgreSQL methods
	Client interface {
		// Returns the raw main database connection
		MainDatastore() *sqlx.DB

		// SubordinateDatastore returns a random subordinate database connection
		//
		// If the paramater is -1 then the returned connection is chosen out of the connection pool
		//
		// If there's no subordinate connection available, then the main connection is returned
		SubordinateDatastore(id int) *sqlx.DB

		SubordinateCount() int
	}

	cli struct {
		main *config.PostgresDB
		subordinates []config.PostgresDB

		mainPg  *sqlx.DB
		subordinatePg []*sqlx.DB
	}
)

func (c *cli) MainDatastore() *sqlx.DB {
	return c.mainPg
}

func (c *cli) SubordinateDatastore(id int) *sqlx.DB {
	if len(c.subordinatePg) == 0 {
		return c.mainPg
	}

	if id == -1 || id >= len(c.subordinatePg) {
		return c.subordinatePg[rand.Intn(len(c.subordinatePg))]
	}

	return c.subordinatePg[id]
}

func (c *cli) SubordinateCount() int {
	return len(c.subordinatePg)
}

func formatConnectionURL(database string, config *config.PostgresDB) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", config.Username, config.Password, config.Host, database, config.Options)
}

func composeConnection(database string, config *config.PostgresDB) *sqlx.DB {
	db, err := sqlx.Open("postgres", formatConnectionURL(database, config))
	if err != nil {
		errors.Fatal(err)
	}

	if err := db.Ping(); err != nil {
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
		mainPg: composeConnection(config.Database, &config.Main),

		main: &config.Main,
		subordinates: config.Subordinates,
	}

	for idx := range config.Subordinates {
		result.subordinatePg = append(result.subordinatePg, composeConnection(config.Database, &config.Subordinates[idx]))
	}

	return result
}
