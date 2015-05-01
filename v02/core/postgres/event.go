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
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
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
	createEventQuery                = `INSERT INTO app_%d_%d.events(json_data) VALUES($1)`
	selectEventByIDQuery            = `SELECT json_data FROM app_%d_%d.events WHERE json_data->>'id' = $1 AND json_data->>'user_id' = $2`
	updateEventByIDQuery            = `UPDATE app_%d_%d.events SET json_data = $1 WHERE json_data->>'id' = $2 AND json_data->>'user_id' = $3`
	listEventsByUserIDQuery         = `SELECT json_data FROM app_%d_%d.events WHERE json_data->>'user_id' = $1`
	listEventsByUserFollowerIDQuery = `SELECT json_data FROM app_%d_%d.events WHERE %s ORDER BY json_data->>'created_at' DESC LIMIT 200`
)

func (e *event) Create(accountID, applicationID int64, event *entity.Event, retrieve bool) (*entity.Event, errors.Error) {
	event.ID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, errors.NewInternalError("error whiel saving the event", err.Error())
	}

	_, err = e.mainPg.
		Exec(appSchema(createEventQuery, accountID, applicationID), string(eventJSON))
	if err != nil {
		return nil, errors.NewInternalError("error while saving the event", err.Error())
	}

	if !retrieve {
		return nil, nil
	}
	return e.Read(accountID, applicationID, event.UserID, event.ID)
}

func (e *event) Read(accountID, applicationID int64, userID, eventID string) (*entity.Event, errors.Error) {
	var JSONData string
	err := e.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectEventByIDQuery, accountID, applicationID), eventID, userID).
		Scan(&JSONData)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the event", err.Error())
	}

	event := &entity.Event{}
	err = json.Unmarshal([]byte(JSONData), event)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account user", err.Error())
	}
	event.ID = eventID

	return event, nil
}

func (e *event) Update(accountID, applicationID int64, existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, errors.Error) {
	eventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		return nil, errors.NewInternalError("failed to update the event", err.Error())
	}

	_, err = e.mainPg.Exec(
		appSchema(updateEventByIDQuery, accountID, applicationID),
		string(eventJSON), existingEvent.ID, existingEvent.UserID)
	if err != nil {
		return nil, errors.NewInternalError("failed to update the event", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return e.Read(accountID, applicationID, existingEvent.UserID, existingEvent.ID)
}

func (e *event) Delete(accountID, applicationID int64, event *entity.Event) errors.Error {
	event.Enabled = false
	_, err := e.Update(accountID, applicationID, *event, *event, false)

	return err
}

func (e *event) List(accountID, applicationID int64, userID string) (events []*entity.Event, er errors.Error) {
	events = []*entity.Event{}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(appSchema(listEventsByUserIDQuery, accountID, applicationID), userID)
	if err != nil {
		return events, errors.NewInternalError("failed to read the events", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}

		events = append(events, event)
	}

	return events, nil
}

func (e *event) ConnectionList(accountID, applicationID int64, userID string) (events []*entity.Event, er errors.Error) {
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
		Query(fmt.Sprintf(listEventsByUserFollowerIDQuery, accountID, applicationID, strings.Join(condition, " AND ")))
	if err != nil {
		return events, errors.NewInternalError("failed to read the events", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the events", err.Error())
		}

		events = append(events, event)
	}

	return events, nil
}

func (e *event) WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID, key string) (err errors.Error) {
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
