package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/tapglue/multiverse/service/event"
)

const (
	namespace      = "app_409_652"
	pgSelectEvents = `SELECT json_data FROM app_409_652.events`
	pgUpdateEvent  = `UPDATE app_409_652.events SET json_data = $1
		WHERE (json_data->>'id')::BIGINT = $2::BIGINT`
)

type inputEvent struct {
	Enabled    bool             `json:"enabled"`
	ID         uint64           `json:"id"`
	Language   string           `json:"language,omitempty"`
	Object     *inputObject     `json:"object,omitempty"`
	ObjectID   uint64           `json:"object_id"`
	Owned      bool             `json:"owned"`
	Target     *event.Target    `json:"target"`
	Type       string           `json:"type"`
	UserID     uint64           `json:"user_id"`
	Visibility event.Visibility `json:"visibility"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type inputObject struct {
	ID   string
	Type string
	URL  string
}

func (p *inputObject) UnmarshalJSON(raw []byte) error {
	f := struct {
		ID   interface{} `json:"id"`
		Type string      `json:"type"`
		URL  string      `json:"url"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.Type = f.Type
	p.URL = f.URL

	id, err := parseID(f.ID)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

func main() {
	var (
		pgURL = flag.String("postgres.url", "postgres://xla@127.0.0.1:5432/tapglue_dev?sslmode=disable&connect_timeout=5", "postgres db conection")
	)
	flag.Parse()

	log.SetFlags(log.Llongfile)

	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(pgSelectEvents)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	is := []*inputEvent{}

	for rows.Next() {
		var (
			input = &inputEvent{}

			raw []byte
		)

		err := rows.Scan(&raw)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(raw, input)
		if err != nil {
			log.Fatal(err)
		}

		is = append(is, input)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, updated := range is {
		ev := &event.Event{
			Enabled:    updated.Enabled,
			ID:         updated.ID,
			Language:   updated.Language,
			ObjectID:   updated.ObjectID,
			Owned:      updated.Owned,
			Target:     updated.Target,
			Type:       updated.Type,
			UserID:     updated.UserID,
			Visibility: updated.Visibility,
			CreatedAt:  updated.CreatedAt,
			UpdatedAt:  updated.UpdatedAt,
		}

		if updated.Object != nil {
			ev.Object = &event.Object{
				ID:   updated.Object.ID,
				Type: updated.Object.Type,
				URL:  updated.Object.URL,
			}
		}

		data, err := json.Marshal(ev)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(pgUpdateEvent, data, ev.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseID(input interface{}) (string, error) {
	var id string

	switch t := input.(type) {
	case float64:
		id = fmt.Sprintf("%d", int64(t))
	case int:
		id = strconv.Itoa(t)
	case int64:
		id = strconv.FormatInt(t, 10)
	case uint64:
		id = strconv.FormatUint(t, 10)
	case string:
		id = t
	default:
		return "", fmt.Errorf("unexpected value for id")
	}

	return id, nil
}
