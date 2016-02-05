package main

// ratios
// user - connection: 1 - 1
// user - object: 1 - 5
// user - event: 1 - 10

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/user"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	pgInsertEvent      = `INSERT INTO app_1_1.events(json_data) VALUES($1)`
	pgInsertConnection = `INSERT INTO app_1_1.connections(json_data) VALUES ($1)`
	pgInsertObject     = `INSERT INTO app_1_1.objects(json_data) VALUES($1)`
	pgInsertUser       = `INSERT INTO app_1_1.users(json_data) VALUES($1)`

	daysPast = 61
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(
		"postgres://%s@127.0.0.1:5432/tapglue_dev?sslmode=disable&connect_timeout=5",
		user.Username,
	)

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	var (
		amountUsers       = rand.Intn(5000) + 10000
		amountConnections = amountUsers * (rand.Intn(2) + 1)
		amountEvents      = amountUsers * (rand.Intn(10) + 5)
		amountObjects     = amountUsers * (rand.Intn(5) + 5)
	)

	fmt.Printf(
		"USERS: %d CONNECTIONS: %d EVENTS: %d OBJECTS: %d\n",
		amountUsers,
		amountConnections,
		amountEvents,
		amountObjects,
	)

	// CONNECTIONS
	for i := 0; i < amountConnections; i++ {
		from, err := flake.NextID("app_1_1_connections")
		if err != nil {
			log.Fatal(err)
		}

		to, err := flake.NextID("app_1_1_events")
		if err != nil {
			log.Fatal(err)
		}

		d, err := time.ParseDuration(fmt.Sprintf("%dh", rand.Intn(daysPast)*24))
		if err != nil {
			log.Fatal(err)
		}

		var (
			stamp = time.Now().Add(-d).UTC()
			con   = v04_entity.Connection{
				UserFromID: from,
				UserToID:   to,
				CreatedAt:  &stamp,
				UpdatedAt:  &stamp,
			}
		)

		data, err := json.Marshal(con)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(pgInsertConnection, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	//  EVENTS
	for i := 0; i < amountEvents; i++ {
		id, err := flake.NextID("app_1_1_events")
		if err != nil {
			log.Fatal(err)
		}

		d, err := time.ParseDuration(fmt.Sprintf("%dh", rand.Intn(daysPast)*24))
		if err != nil {
			log.Fatal(err)
		}

		var (
			stamp = time.Now().Add(-d).UTC()
			event = v04_entity.Event{
				ID: id,
				Common: v04_entity.Common{
					CreatedAt: &stamp,
					UpdatedAt: &stamp,
				},
			}
		)

		data, err := json.Marshal(event)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(pgInsertEvent, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	// OBJECTS
	for i := 0; i < amountEvents; i++ {
		id, err := flake.NextID("app_1_1_objects")
		if err != nil {
			log.Fatal(err)
		}

		d, err := time.ParseDuration(fmt.Sprintf("%dh", rand.Intn(daysPast)*24))
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.Marshal(&object.Object{
			ID:        id,
			CreatedAt: time.Now().Add(-d).UTC(),
			UpdatedAt: time.Now().Add(-d).UTC(),
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(pgInsertObject, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	// USERS
	for i := 0; i < amountUsers; i++ {
		id, err := flake.NextID("app_1_1_users")
		if err != nil {
			log.Fatal(err)
		}

		d, err := time.ParseDuration(fmt.Sprintf("%dh", rand.Intn(daysPast)*24))
		if err != nil {
			log.Fatal(err)
		}

		stamp := time.Now().Add(-d).UTC()

		data, err := json.Marshal(&v04_entity.ApplicationUser{
			ID: id,
			Common: v04_entity.Common{
				CreatedAt: &stamp,
				UpdatedAt: &stamp,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(pgInsertUser, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
