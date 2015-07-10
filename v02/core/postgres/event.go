/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	event struct {
		pg     postgres.Client
		mainPg *sqlx.DB
		c      core.Connection
	}
)

const (
	createEventQuery = `INSERT INTO app_%d_%d.events(json_data, geo)
		VALUES($1, ST_SetSRID(ST_MakePoint($2, $3), 4326))`

	selectEventByIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> json_build_object('id', $1::text, 'user_id', $2::text, 'enabled', true)::jsonb  LIMIT 1`

	updateEventByIDQuery = `UPDATE app_%d_%d.events
		SET json_data = $1, geo = ST_GeomFromText('POINT(' || $2 || ' ' || $3 || ')', 4326)
		WHERE json_data @> json_build_object('id', $4::text, 'user_id', $5::text)::jsonb`

	listPublicEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> json_build_object('user_id', $1::text, 'enabled', true)::jsonb
			AND (json_data @> '{"visibility": 30}' OR json_data @> '{"visibility": 40}')
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listConnectionEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> json_build_object('user_id', $1::text, 'enabled', true)::jsonb
			AND (json_data @> '{"visibility": 20}' OR json_data @> '{"visibility": 30}' OR json_data @> '{"visibility": 40}')
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listAllEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> json_build_object('user_id', $1::text, 'enabled', true)::jsonb
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByUserFollowerIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (((%s) AND (json_data @> '{"visibility": 20}' OR json_data @> '{"visibility": 30}'))
			OR (json_data @> '{"visibility": 40}' AND json_data->>'user_id' != $1))
			AND json_data @> '{"enabled": true}'
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listUnreadEventsByUserFollowerIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (((%s) AND (json_data @> '{"visibility": 20}' OR json_data @> '{"visibility": 30}'))
			OR (json_data @> '{"visibility": 40}' AND json_data->>'user_id' != $1))
			AND json_data->>'created_at' > $2
			AND json_data @> '{"enabled": true}'
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	countUnreadEventsByUserFollowerIDQuery = `SELECT count(*)
		FROM (
			SELECT json_data
			FROM app_%d_%d.events
			WHERE (((%s) AND (json_data @> '{"visibility": 20}' OR json_data @> '{"visibility": 30}'))
				OR (json_data @> '{"visibility": 40}' AND json_data->>'user_id' != $1))
				AND json_data->>'created_at' > $2
				AND json_data @> '{"enabled": true}'
			ORDER BY json_data->>'created_at' DESC LIMIT 200) AS events`

	updateApplicationUserLastReadQuery = `UPDATE app_%d_%d.users
		SET last_read = now()
		WHERE json_data @> json_build_object('id', $1::text, 'enabled', true)::jsonb`

	listEventsByLocationQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> json_build_object('location', $1::text, 'enabled', true)::jsonb
			AND (
				json_data @> '{"visibility": 30}' OR json_data @> '{"visibility": 40}' OR
				%s
				json_data @> json_build_object('user_id', $2::text)::jsonb
			)
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByLatLonRadQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE ST_DWithin(geo, ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, true)
			AND json_data @> '{"enabled": true}'
			AND (
				json_data @> '{"visibility": 30}' OR json_data @> '{"visibility": 40}' OR
				%s
				json_data @> json_build_object('user_id', $4::text)::jsonb
			)
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByLatLonNearQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data @> '{"enabled": true}'
			AND (
				json_data @> '{"visibility": 30}' OR json_data @> '{"visibility": 40}' OR
				%s
				json_data @> json_build_object('user_id', $1::text)::jsonb
			)
		ORDER BY ST_Distance_Sphere(geo, ST_SetSRID(ST_MakePoint($2, $3), 4326)), json_data->>'created_at' DESC LIMIT $4`
)

func (e *event) Create(accountID, applicationID int64, currentUserID string, event *entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	if event.ID == "" {
		return nil, []errors.Error{errmsg.ErrInternalEventMissingID}
	}
	event.Enabled = true
	timeNow := time.Now()
	event.CreatedAt, event.UpdatedAt = &timeNow, &timeNow

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventCreation.UpdateInternalMessage(err.Error())}
	}

	_, err = e.mainPg.
		Exec(appSchema(createEventQuery, accountID, applicationID), string(eventJSON), event.Latitude, event.Longitude)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}
	return e.Read(accountID, applicationID, event.UserID, currentUserID, event.ID)
}

func (e *event) Read(accountID, applicationID int64, userID, currentUserID, eventID string) (*entity.Event, []errors.Error) {
	var JSONData string
	err := e.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectEventByIDQuery, accountID, applicationID), eventID, userID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrEventNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalEventRead.UpdateInternalMessage(err.Error())}
	}

	event := &entity.Event{}
	err = json.Unmarshal([]byte(JSONData), event)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventRead.UpdateInternalMessage(err.Error())}
	}
	event.ID = eventID

	return event, nil
}

func (e *event) Update(accountID, applicationID int64, currentUserID string, existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	timeNow := time.Now()
	updatedEvent.UpdatedAt = &timeNow
	eventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = e.mainPg.Exec(
		appSchema(updateEventByIDQuery, accountID, applicationID),
		string(eventJSON), updatedEvent.Latitude, updatedEvent.Longitude, existingEvent.ID, existingEvent.UserID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return e.Read(accountID, applicationID, existingEvent.UserID, currentUserID, existingEvent.ID)
}

func (e *event) Delete(accountID, applicationID int64, currentUserID string, event *entity.Event) []errors.Error {
	event.Enabled = false
	_, err := e.Update(accountID, applicationID, currentUserID, *event, *event, false)

	return err
}

func (e *event) List(accountID, applicationID int64, userID, currentUserID string) (events []*entity.Event, er []errors.Error) {
	events = []*entity.Event{}

	var query string
	if userID == currentUserID {
		query = listAllEventsByUserIDQuery
	} else if _, err := e.c.Read(accountID, applicationID, currentUserID, userID); err != nil {
		if err[0] == errmsg.ErrConnectionNotFound {
			query = listPublicEventsByUserIDQuery
		} else {
			return []*entity.Event{}, err
		}
	} else {
		query = listConnectionEventsByUserIDQuery
	}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(appSchema(query, accountID, applicationID), userID)
	if err != nil {
		return events, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
	}
	return e.rowsToSlice(rows)
}

func (e *event) UserFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, er []errors.Error) {
	condition, er := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if er != nil {
		return 0, nil, er
	}

	if condition == "" {
		return 0, []*entity.Event{}, nil
	}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(fmt.Sprintf(listEventsByUserFollowerIDQuery, accountID, applicationID, condition), user.ID)
	if err != nil {
		return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()

	unread := 0

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
		}

		if event.CreatedAt.Sub(*user.LastRead) > 0 {
			unread++
		}

		events = append(events, event)
	}

	go e.updateApplicationUserLastRead(accountID, applicationID, user)

	return unread, events, nil
}

func (e *event) UnreadFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, er []errors.Error) {
	condition, er := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if er != nil {
		return 0, nil, er
	}

	if condition == "" {
		return 0, nil, nil
	}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(fmt.Sprintf(listUnreadEventsByUserFollowerIDQuery, accountID, applicationID, condition), user.ID, user.LastRead)
	if err != nil {
		return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
	}
	events, er = e.rowsToSlice(rows)
	if er != nil {
		return
	}

	go e.updateApplicationUserLastRead(accountID, applicationID, user)

	return len(events), events, nil
}

func (e *event) UnreadFeedCount(accountID, applicationID int64, user *entity.ApplicationUser) (count int, err []errors.Error) {
	condition, err := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if err != nil {
		return 0, err
	}

	if condition == "" {
		return 0, nil
	}

	unread := 0
	er := e.pg.SlaveDatastore(-1).
		QueryRow(fmt.Sprintf(countUnreadEventsByUserFollowerIDQuery, accountID, applicationID, condition), user.ID, user.LastRead).
		Scan(&unread)
	if er != nil {
		return 0, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(er.Error())}
	}

	return unread, nil
}

func (e *event) WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (e *event) GeoSearch(accountID, applicationID int64, currentUserID string, latitude, longitude, radius float64, nearest int64) ([]*entity.Event, []errors.Error) {
	var (
		rows *sql.Rows
		err  error
	)

	condition, er := e.composeConnectionCondition(accountID, applicationID, currentUserID, " OR ")
	if er != nil {
		return nil, er
	}

	if condition != "" {
		condition = `(json_data @> '{"visibility": 20}' AND (` + condition + `)) OR`
	}

	if nearest == 0 {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(appSchemaWithParams(listEventsByLatLonRadQuery, accountID, applicationID, condition), latitude, longitude, radius, currentUserID)
	} else {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(appSchemaWithParams(listEventsByLatLonNearQuery, accountID, applicationID, condition), currentUserID, latitude, longitude, nearest)
	}

	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
	}
	return e.rowsToSlice(rows)
}

func (e *event) ObjectSearch(accountID, applicationID int64, currentUserID string, objectKey string) ([]*entity.Event, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (e *event) LocationSearch(accountID, applicationID int64, currentUserID string, locationKey string) ([]*entity.Event, []errors.Error) {
	condition, er := e.composeConnectionCondition(accountID, applicationID, currentUserID, " OR ")
	if er != nil {
		return nil, er
	}

	if condition != "" {
		condition = `(json_data @> '{"visibility": 20}' AND (` + condition + `)) OR`
	}

	rows, err := e.pg.SlaveDatastore(-1).
		Query(appSchemaWithParams(listEventsByLocationQuery, accountID, applicationID, condition), locationKey, currentUserID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
	}
	return e.rowsToSlice(rows)
}

func (e *event) updateApplicationUserLastRead(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error {
	_, err := e.mainPg.Exec(appSchema(updateApplicationUserLastReadQuery, accountID, applicationID), user.ID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error())}
	}

	return nil
}

func (e *event) rowsToSlice(rows *sql.Rows) (events []*entity.Event, err []errors.Error) {
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error())}
		}

		events = append(events, event)
	}
	return
}

func (e *event) composeConnectionCondition(accountID, applicationID int64, userID, joinOperator string) (string, []errors.Error) {
	connections, er := e.c.FriendsAndFollowing(accountID, applicationID, userID)
	if er != nil {
		return "", er
	}

	if len(connections) == 0 {
		return "", nil
	}

	condition := []string{}
	for idx := range connections {
		condition = append(condition, fmt.Sprintf(`json_data @> '{"user_id": %q}'`, connections[idx].ID))
	}

	return strings.Join(condition, joinOperator), nil
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
