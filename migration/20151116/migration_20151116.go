package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	mr "math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/entity"
	v04_postgres "github.com/tapglue/multiverse/v04/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	// Connection structure holds the connections of the users
	Connection struct {
		UserFromID uint64                     `json:"user_from_id"`
		UserToID   uint64                     `json:"user_to_id"`
		Type       entity.ConnectionTypeType  `json:"type"`
		State      entity.ConnectionStateType `json:"state"`
		Enabled    *bool                      `json:"enabled,omitempty"`
		CreatedAt  *time.Time                 `json:"created_at,omitempty"`
		UpdatedAt  *time.Time                 `json:"updated_at,omitempty"`
	}
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_INTAKER_CONFIG_PATH"

	listConnectionsByApplicationIDQuery = `SELECT json_data FROM app_%d_%d.connections`

	updateConnectionQuery = `UPDATE app_%d_%d.connections
		SET json_data = $1
		WHERE (json_data->>'user_from_id')::BIGINT = $2::BIGINT
			AND (json_data->>'user_to_id')::BIGINT = $3::BIGINT
			AND json_data->>'type' = $4::TEXT`

	deleteConnectionQuery = `DELETE FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
		 	AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT
		 	AND json_data->>'type' = $3::TEXT`
)

var (
	startTime         time.Time
	db                *sqlx.DB
	conf              *config.Config
	v04PostgresClient v04_postgres.Client

	doRun = flag.Bool("run", false, "Set to true to run, else do a dry-run")
)

func init() {
	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

	log.SetFlags(0)

	if conf.UseSysLog {
		syslogWriter, err := syslog.New(syslog.LOG_INFO, "intaker")
		if err == nil {
			log.Printf("logging to syslog is enabled. Please tail your syslog for intaker app for further logs\n")
			log.SetOutput(syslogWriter)
		} else {
			log.Printf("%v\n", err)
			log.Printf("logging to syslog failed reverting to stdout logging\n")
		}
		conf.UseArtwork = false
	}

	errors.Init(true)

	v04PostgresClient = v04_postgres.New(conf.Postgres)
	db = v04PostgresClient.MainDatastore()
}

func main() {
	apps := map[int64][]int64{}

	if *doRun {
		fmt.Printf("LAUNCHING WITH ACTIVE CHANGES MODE!!!\n")
		fmt.Printf("LAUNCHING WITH ACTIVE CHANGES MODE!!!\n")
		fmt.Printf("LAUNCHING WITH ACTIVE CHANGES MODE!!!\n\n")
		fmt.Printf("If you don't want to continue you have 5 seconds to abort\n\n")

		time.Sleep(5 * time.Second)
	}

	existingSchemas, err := db.Query(`SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname ILIKE 'app_%_%'`)
	if err != nil {
		panic(err)
	}
	defer existingSchemas.Close()
	for existingSchemas.Next() {
		schemaName := ""
		err := existingSchemas.Scan(&schemaName)
		if err != nil {
			panic(err)
		}
		details := strings.Split(schemaName, "_")
		if len(details) != 3 || details[0] != "app" {
			panic(fmt.Sprintf("%# v", details))
		}

		orgID, err := strconv.ParseInt(details[1], 10, 64)
		if err != nil {
			panic(err)
		}

		appID, err := strconv.ParseInt(details[2], 10, 64)
		if err != nil {
			panic(err)
		}

		apps[orgID] = append(apps[orgID], appID)
	}

	fmt.Printf("Migrating %d applications\n\n", len(apps))

	apps = map[int64][]int64{374: []int64{501}}

	for orgID, aps := range apps {
		for idx := 0; idx < len(aps); idx++ {
			appID := aps[idx]
			existingConnections := fetchAllConnections(orgID, appID)
			fmt.Printf("Got %d conection(s) for app_%d_%d\n", len(existingConnections), orgID, appID)
			remapConnections(orgID, appID, existingConnections)
		}
	}

	duration := time.Now().Sub(startTime)
	if *doRun {
		time.Now().Add(-5 * time.Second).Sub(startTime)
	}
	fmt.Printf("\n\nFinished in %s\n", duration)
}

func remapConnections(orgID, appID int64, connections []Connection) {
	query := appSchema(updateConnectionQuery, orgID, appID)
	for idx := 0; idx < len(connections); idx++ {
		connection := connections[idx]

		if connection.State != "" {
			continue
		}

		connection.State = entity.ConnectionStateConfirmed

		connectionJSON, err := json.Marshal(connection)
		if err != nil {
			panic(err)
		}

		fmt.Printf(
			"running query: %s with args: %s, %d, %d, %s\n",
			query, string(connectionJSON), connection.UserFromID, connection.UserToID, string(connection.Type),
		)
		if !*doRun {
			continue
		}
		_, err = db.Exec(query, string(connectionJSON), connection.UserFromID, connection.UserToID, string(connection.Type))
		if err != nil {
			panic(err)
		}
	}
}

func fetchAllConnections(orgID, appID int64) []Connection {
	var connections []Connection

	query := appSchema(listConnectionsByApplicationIDQuery, orgID, appID)
	fmt.Printf("running query: %s\n", query)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			panic(err)
		}
		user := Connection{}
		err = json.Unmarshal([]byte(JSONData), &user)
		if err != nil {
			panic(err)
		}

		connections = append(connections, user)
	}

	return connections
}

func appSchema(query string, accountID, applicationID int64) string {
	return fmt.Sprintf(query, accountID, applicationID)
}
