package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tapglue/multiverse/service/user"

	"github.com/tapglue/multiverse/service/session"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/tapglue/multiverse/platform/service"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/device"
)

const (
	pgDropIndex  = `DROP INDEX %s.%s`
	pgIndexNames = `
		SELECT
			indexrelname
		FROM
			pg_statio_user_indexes
		WHERE
			relname = '%s'
			AND schemaname = '%s'`
)

var namespaces = []string{
	"app_791_1139",
	"app_791_1115",
	"app_787_1105",
	"app_776_1081",
	"app_775_1078",
	"app_684_948",
	"app_628_863",
	"app_57_661",
	"app_57_452",
	"app_55_1047",
	"app_529_720",
	"app_515_922",
	"app_409_652",
	"app_374_501",
	"app_309_443",
	"app_309_440",
	"app_309_428",
	"app_309_425",
	"app_309_1020",
	"app_26_187",
	"app_261_626",
	"app_1_610",
	"app_1_1147",
	"app_1_1067",
}

func main() {
	var (
		pgURL = flag.String("pg.url", "", "")
	)
	flag.Parse()

	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatal(err)
	}

	services := map[string]service.Lifecycle{
		"connections": connection.NewPostgresService(db),
		"devices":     device.PostgresService(db),
		"sessions":    session.NewPostgresService(db),
		"users":       user.NewPostgresService(db),
	}

	for _, ns := range namespaces {
		for name, s := range services {
			is, err := getIndexes(db, ns, name)
			if err != nil {
				log.Fatal(err)
			}

			for _, i := range is {
				err := dropIndex(db, ns, i)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("DROP INDEX %s.%s\n", ns, i)
			}

			err = s.Setup(ns)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Setup for %s.%s\n", ns, name)
		}
	}
}

func dropIndex(db *sqlx.DB, ns, name string) error {
	_, err := db.Exec(fmt.Sprintf(pgDropIndex, ns, name))
	return err
}

func getIndexes(db *sqlx.DB, ns, name string) ([]string, error) {
	q := fmt.Sprintf(pgIndexNames, name, ns)

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	is := []string{}

	for rows.Next() {
		var index string

		err := rows.Scan(&index)
		if err != nil {
			return nil, err
		}

		is = append(is, index)
	}

	return is, nil
}
