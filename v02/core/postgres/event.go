/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"

	"fmt"
	"strings"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	event struct {
		pg     postgres.Client
		mainPg *sql.DB
		c      core.Connection
	}
)

const (
	createEventQuery                = `INSERT INTO app_$1_$2.events(json_data, enabled) VALUES($3, $4) RETURNING id`
	selectEventByIDQuery            = `SELECT json_data, enabled FROM app_$1_$2.events WHERE id = $3 AND json_data @> '{"user_id": $4}'`
	updateEventByIDQuery            = `UPDATE app_$1_$2.events SET json_data = $3, enabled = $4 WHERE id = $5 AND json_data @> '{"user_id": $6}'`
	deleteEventByIDQuery            = `UPDATE app_$1_$2.events SET enabled = 0 WHERE id = $4 AND json_data @> '{"user_id": $4}'`
	listEventsByUserIDQuery         = `SELECT id, json_data, enabled FROM app_$1_$2.events WHERE json_data @> '{"user_id": $3}'`
	listEventsByUserFollowerIDQuery = `SELECT id, json_data, enabled FROM app_$1_$2.events WHERE %s ORDER BY json_data->>'created_at' DESC LIMIT 200`
)

func (e *event) Create(event *entity.Event, retrieve bool) (*entity.Event, errors.Error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, errors.NewInternalError("error whiel saving the event", err.Error())
	}

	var eventID int64
	err = e.mainPg.
		QueryRow(createEventQuery, event.AccountID, event.ApplicationID, string(eventJSON), event.Enabled).
		Scan(&eventID)
	if err != nil {
		return nil, errors.NewInternalError("error while saving the event", err.Error())
	}

	if !retrieve {
		return nil, nil
	}
	return e.Read(event.AccountID, event.ApplicationID, event.UserID, eventID)
}

func (e *event) Read(accountID, applicationID, userID, eventID int64) (*entity.Event, errors.Error) {
	eventJSON := &struct {
		ID       int64
		JSONData string
		Enabled  bool
	}{}
	err := e.pg.SlaveDatastore(-1).
		QueryRow(selectEventByIDQuery, accountID, applicationID, eventID, userID).
		Scan(eventJSON)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the event", err.Error())
	}

	event := &entity.Event{}
	err = json.Unmarshal([]byte(eventJSON.JSONData), event)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account user", err.Error())
	}
	event.ID = eventJSON.ID
	event.Enabled = eventJSON.Enabled

	return event, nil
}

func (e *event) Update(existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, errors.Error) {
	eventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		return nil, errors.NewInternalError("failed to update the event", err.Error())
	}

	_, err = e.mainPg.Exec(updateEventByIDQuery, existingEvent.AccountID, existingEvent.ApplicationID, string(eventJSON), updatedEvent.Enabled, existingEvent.ID, existingEvent.UserID)
	if err != nil {
		return nil, errors.NewInternalError("failed to update the event", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return e.Read(existingEvent.AccountID, existingEvent.ApplicationID, existingEvent.UserID, existingEvent.ID)
}

func (e *event) Delete(event *entity.Event) errors.Error {
	_, err := e.mainPg.Exec(deleteEventByIDQuery, event.AccountID, event.ApplicationID, event.ID, event.UserID)
	if err != nil {
		return errors.NewInternalError("error while deleting the event", err.Error())
	}
	return nil
}

func (e *event) List(accountID, applicationID, userID int64) (events []*entity.Event, er errors.Error) {
	events = []*entity.Event{}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(listEventsByUserIDQuery, accountID, applicationID, userID)
	if err != nil {
		return events, errors.NewInternalError("failed to read the events", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		eventJSON := &struct {
			ID       int64
			JSONData string
			Enabled  bool
		}{}
		err := rows.Scan(eventJSON)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(eventJSON.JSONData), event)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event.ID = eventJSON.ID
		event.Enabled = eventJSON.Enabled

		events = append(events, event)
	}

	return events, nil
}

func (e *event) ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, er errors.Error) {
	events = []*entity.Event{}

	connections, er := e.c.List(accountID, applicationID, userID)
	if er != nil {
		return events, er
	}

	if len(connections) == 0 {
		return events, nil
	}

	condition := []string{}
	for idx := range connections {
		condition = append(condition, fmt.Sprintf(`json_data @> '{"user_id": %d}'`, connections[idx].ID))
	}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(fmt.Sprintf(listEventsByUserFollowerIDQuery, strings.Join(condition, " AND ")), accountID, applicationID, userID)
	if err != nil {
		return events, errors.NewInternalError("failed to read the events", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		eventJSON := &struct {
			ID       int64
			JSONData string
			Enabled  bool
		}{}
		err := rows.Scan(eventJSON)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(eventJSON.JSONData), event)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event.ID = eventJSON.ID
		event.Enabled = eventJSON.Enabled

		events = append(events, event)
	}

	return events, nil
}

func (e *event) WriteToConnectionsLists(event *entity.Event, key string) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID, userID int64, key string) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err errors.Error) {
	return []*entity.Event{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, errors.Error) {
	return []*entity.Event{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, errors.Error) {
	return []*entity.Event{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewEventWithConnection returns a new event handler with PostgreSQL as storage driver
func NewEventWithConnection(pgsql postgres.Client, connection core.Connection) core.Event {
	return &event{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		c:      connection,
	}
}

// NewEvent returns a new event handler with PostgreSQL as storage driver
func NewEvent(pgsql postgres.Client) core.Event {
	return &event{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		c:      NewConnection(pgsql),
	}
}