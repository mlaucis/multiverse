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

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/tgflake"
	v02_server "github.com/tapglue/backend/v02/server"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	// Common holds common used fields
	Common struct {
		Metadata  interface{}       `json:"metadata,omitempty"`
		Images    map[string]*Image `json:"images,omitempty"`
		CreatedAt *time.Time        `json:"created_at,omitempty"`
		UpdatedAt *time.Time        `json:"updated_at,omitempty"`
		Enabled   bool              `json:"enabled,omitempty"`
	}

	// UserCommon holds common used fields for users
	UserCommon struct {
		Username         string     `json:"user_name"`
		OriginalPassword string     `json:"-"`
		Password         string     `json:"password,omitempty"`
		FirstName        string     `json:"first_name,omitempty"`
		LastName         string     `json:"last_name,omitempty"`
		Email            string     `json:"email,omitempty"`
		URL              string     `json:"url,omitempty"`
		LastLogin        *time.Time `json:"last_login,omitempty"`
		Activated        bool       `json:"activated,omitempty"`
	}

	// Image structure
	Image struct {
		URL    string `json:"url"`
		Type   string `json:"type,omitempty"` // image/jpeg image/png
		Width  int    `json:"width,omitempty"`
		Heigth int    `json:"height,omitempty"`
	}

	// Object structure
	Object struct {
		ID           string            `json:"id"`
		Type         string            `json:"type"`
		URL          string            `json:"url,omitempty"`
		DisplayNames map[string]string `json:"display_names"` // ["en"=>"article", "de"=>"artikel"]
	}

	// Participant structure
	Participant struct {
		ID     string            `json:"id"`
		URL    string            `json:"url,omitempty"`
		Images map[string]*Image `json:"images,omitempty"`
	}

	// ApplicationUser structure
	ApplicationUser struct {
		ID                   interface{}         `json:"id"`
		CustomID             string              `json:"custom_id,omitempty"`
		SessionToken         string              `json:"-"`
		SocialIDs            map[string]string   `json:"social_ids,omitempty"`
		SocialConnectionsIDs map[string][]string `json:"social_connections_ids,omitempty"`
		SocialConnectionType string              `json:"connection_type,omitempty"`
		DeviceIDs            []string            `json:"device_ids,omitempty"`
		Events               []*Event            `json:"events,omitempty"`
		Connections          []*ApplicationUser  `json:"connections,omitempty"`
		LastRead             *time.Time          `json:"-"`
		UserCommon
		Common
	}

	// Connection structure holds the connections of the users
	Connection struct {
		UserFromID  interface{} `json:"user_from_id"`
		UserToID    interface{} `json:"user_to_id"`
		Type        string      `json:"type"`
		ConfirmedAt *time.Time  `json:"confirmed_at,omitempty"`
		Common
	}

	// Event structure
	Event struct {
		ID                 interface{}    `json:"id"`
		UserID             interface{}    `json:"user_id"`
		Type               string         `json:"type"`
		Language           string         `json:"language,omitempty"`
		Priority           string         `json:"priority,omitempty"`
		Location           string         `json:"location,omitempty"`
		Latitude           float64        `json:"latitude,omitempty"`
		Longitude          float64        `json:"longitude,omitempty"`
		DistanceFromTarget float64        `json:"-"`
		Visibility         uint8          `json:"visibility,omitempty"`
		Object             *Object        `json:"object"`
		Target             *Object        `json:"target,omitempty"`
		Instrument         *Object        `json:"instrument,omitempty"`
		Participant        []*Participant `json:"participant,omitempty"`
		Common
	}
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_INTAKER_CONFIG_PATH"

	listApplicationUsersByApplicationIDQuery = `SELECT json_data FROM app_%d_%d.users`
	listConnectionsByApplicationIDQuery      = `SELECT json_data FROM app_%d_%d.Connections`
	listEventsByApplicationIDQuery           = `SELECT json_data FROM app_%d_%d.Events`

	updateApplicationUserByIDQuery = `UPDATE app_%d_%d.users
		SET json_data = $1
		WHERE json_data @> json_build_object('id', $2::text)::jsonb`

	updateConnectionQuery = `UPDATE app_%d_%d.connections
		SET json_data = $1
		WHERE json_data @> json_build_object('user_from_id', $2::text, 'user_to_id', $3::text)::jsonb`

	updateEventByIDQuery = `UPDATE app_%d_%d.events
		SET json_data = $1, geo = ST_GeomFromText('POINT(' || $2 || ' ' || $3 || ')', 4326)
		WHERE json_data @> json_build_object('id', $4::text, 'user_id', $5::text)::jsonb`

	destroyAllApplicationSessionsQuery = `DELETE FROM app_%d_%d.sessions`
)

var (
	startTime         time.Time
	db                *sqlx.DB
	conf              *config.Config
	v02PostgresClient v02_postgres.Client

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

	v02PostgresClient = v02_postgres.New(conf.Postgres)
	db = v02PostgresClient.MainDatastore()

	v02_server.SetupFlakes(v02PostgresClient)
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

		accID, err := strconv.ParseInt(details[1], 10, 64)
		if err != nil {
			panic(err)
		}

		appID, err := strconv.ParseInt(details[2], 10, 64)
		if err != nil {
			panic(err)
		}

		apps[accID] = append(apps[accID], appID)
	}

	fmt.Printf("Migrating %# v\n", apps)

	for accID, aps := range apps {
		for idx := 0; idx < len(aps); idx++ {
			appID := aps[idx]
			existingUsers := fetchAllUsers(accID, appID)
			if len(existingUsers) == 0 {
				continue
			}

			existingConnections := fetchAllConnections(accID, appID)
			existingEvents := fetchAllEvents(accID, appID)
			fmt.Printf("Got %d user(s), %d conection(s), %d event(s)\n", len(existingUsers), len(existingConnections), len(existingEvents))

			flakedUserPairs := generateUserFlakes(appID, existingUsers)
			flakedEventPairs := generateEventFlakes(appID, existingEvents)

			remapUsers(accID, appID, existingUsers, flakedUserPairs)
			remapConnections(accID, appID, existingConnections, flakedUserPairs)
			remapEvents(accID, appID, existingEvents, flakedUserPairs, flakedEventPairs)
		}
	}

	duration := time.Now().Sub(startTime)
	if *doRun {
		time.Now().Add(-5 * time.Second).Sub(startTime)
	}
	fmt.Printf("\n\nFinished in %s\n", duration)
}

func remapUsers(accountID, applicationID int64, existingUsers []ApplicationUser, flakedUserPairs map[string]uint64) {
	for idx := 0; idx < len(existingUsers); idx++ {
		if _, ok := existingUsers[idx].ID.(string); !ok {
			continue
		}

		oldID := existingUsers[idx].ID.(string)
		existingUsers[idx].ID = flakedUserPairs[oldID]

		if !*doRun {
			continue
		}

		userJSON, err := json.Marshal(existingUsers[idx])
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(appSchema(updateApplicationUserByIDQuery, accountID, applicationID), string(userJSON), oldID)
		if err != nil {
			panic(err)
		}
	}

	_, _ = db.Exec(appSchema(destroyAllApplicationSessionsQuery, accountID, applicationID))
}

func remapConnections(accountID, applicationID int64, existingConnections []Connection, flakedUserPairs map[string]uint64) {
	for idx := 0; idx < len(existingConnections); idx++ {
		if _, ok := existingConnections[idx].UserFromID.(string); !ok {
			continue
		}

		oldFromID, oldToID := existingConnections[idx].UserFromID.(string), existingConnections[idx].UserToID.(string)

		existingConnections[idx].UserFromID = flakedUserPairs[existingConnections[idx].UserFromID.(string)]
		existingConnections[idx].UserToID = flakedUserPairs[existingConnections[idx].UserToID.(string)]

		if !*doRun {
			continue
		}

		connectionJSON, err := json.Marshal(existingConnections[idx])
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(
			appSchema(updateConnectionQuery, accountID, applicationID),
			string(connectionJSON), oldFromID, oldToID)
		if err != nil {
			panic(err)
		}
	}
}

func remapEvents(accountID, applicationID int64, existingEvents []Event, flakedUserPairs, flakedEventPairs map[string]uint64) {
	for idx := 0; idx < len(existingEvents); idx++ {
		if _, ok := existingEvents[idx].ID.(string); !ok {
			continue
		}

		oldID := existingEvents[idx].ID.(string)
		oldUserID := existingEvents[idx].UserID.(string)
		existingEvents[idx].ID = flakedEventPairs[oldID]
		existingEvents[idx].UserID = flakedUserPairs[oldUserID]

		if !*doRun {
			continue
		}

		eventJSON, err := json.Marshal(existingEvents[idx])
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(
			appSchema(updateEventByIDQuery, accountID, applicationID),
			string(eventJSON), existingEvents[idx].Latitude, existingEvents[idx].Longitude, oldID, oldUserID)
		if err != nil {
			panic(err)
		}
	}
}

func generateUserFlakes(appID int64, users []ApplicationUser) map[string]uint64 {
	var err error
	flakedUsers := map[string]uint64{}

	for idx := 0; idx < len(users); idx++ {
		if _, ok := users[idx].ID.(string); !ok {
			continue
		}
		flakedUsers[users[idx].ID.(string)], err = tgflake.FlakeNextID(appID, "users")
		if err != nil {
			panic(err)
		}
	}

	return flakedUsers
}

func generateEventFlakes(appID int64, events []Event) map[string]uint64 {
	var err error
	flakedEvents := map[string]uint64{}

	for idx := 0; idx < len(events); idx++ {
		if _, ok := events[idx].ID.(string); !ok {
			continue
		}

		flakedEvents[events[idx].ID.(string)], err = tgflake.FlakeNextID(appID, "events")
		if err != nil {
			panic(err)
		}
	}

	return flakedEvents
}

func fetchAllUsers(accountID, applicationID int64) []ApplicationUser {
	users := []ApplicationUser{}

	fmt.Printf("running query: %s\n", appSchema(listApplicationUsersByApplicationIDQuery, accountID, applicationID))
	rows, err := db.Query(appSchema(listApplicationUsersByApplicationIDQuery, accountID, applicationID))
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
		user := ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), &user)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users
}

func fetchAllConnections(accountID, applicationID int64) []Connection {
	users := []Connection{}

	fmt.Printf("running query: %s\n", appSchema(listConnectionsByApplicationIDQuery, accountID, applicationID))
	rows, err := db.Query(appSchema(listConnectionsByApplicationIDQuery, accountID, applicationID))
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

		users = append(users, user)
	}

	return users
}

func fetchAllEvents(accountID, applicationID int64) []Event {
	users := []Event{}

	fmt.Printf("running query: %s\n", appSchema(listEventsByApplicationIDQuery, accountID, applicationID))
	rows, err := db.Query(appSchema(listEventsByApplicationIDQuery, accountID, applicationID))
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
		user := Event{}
		err = json.Unmarshal([]byte(JSONData), &user)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users
}

func appSchema(query string, accountID, applicationID int64) string {
	return fmt.Sprintf(query, accountID, applicationID)
}
