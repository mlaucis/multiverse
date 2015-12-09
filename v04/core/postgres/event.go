package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type event struct {
	pg     postgres.Client
	mainPg *sqlx.DB
	conn   core.Connection
}

const (
	createEventQuery = `INSERT INTO app_%d_%d.events(json_data, geo)
		VALUES($1, ST_SetSRID(ST_MakePoint($2, $3), 4326))`

	selectEventByIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (json_data->>'id')::BIGINT = $1::BIGINT
			AND (json_data->>'user_id')::BIGINT = $2::BIGINT
			AND (json_data->>'enabled')::BOOL = true
		LIMIT 1`

	updateEventByIDQuery = `UPDATE app_%d_%d.events
		SET json_data = $1, geo = ST_GeomFromText('POINT(' || $2 || ' ' || $3 || ')', 4326)
		WHERE (json_data->>'id')::BIGINT = $4::BIGINT
			AND (json_data->>'user_id')::BIGINT = $5::BIGINT`

	listPublicEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (json_data->>'user_id')::BIGINT = $1::BIGINT
			AND ((json_data->>'visibility')::INT == 30 OR (json_data->>'visibility')::INT == 40)
			AND (json_data->'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listConnectionEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (json_data->>'user_id')::BIGINT = $1::BIGINT
			AND ((json_data->>'visibility')::INT = 20 OR (json_data->>'visibility')::INT = 30 OR (json_data->>'visibility')::INT = 40)
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listAllEventsByUserIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (json_data->>'user_id')::BIGINT = $1::BIGINT
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listAllEventsQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listAllEventsWithUserEventsAndWithUserFollowersEventsQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE ((%s AND (json_data->>'visibility')::INT = 20)
			OR (json_data->>'user_id')::BIGINT = $1::BIGINT
			OR (json_data->>'visibility')::INT = 30
			OR (json_data->>'visibility')::INT = 40
			OR (json_data->'target'->>'id') = $2::TEXT)
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listAllEventsWithUserEventsAndNoUserFollowersEventsQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE ((json_data->>'user_id')::BIGINT = $1::BIGINT
				OR (json_data->>'visibility')::INT = 30
				OR (json_data->>'visibility')::INT = 40
				OR (json_data->'target'->>'id') = $2::TEXT)
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByUserFollowerIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE ((%s AND ((json_data->>'visibility')::INT = 20 OR (json_data->>'visibility')::INT = 30))
			OR ((json_data->>'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
			OR (json_data->'target'->>'id') = $2::TEXT)
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsNoUserFollowersQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (((json_data->>'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
			OR (json_data->'target'->>'id') = $2::TEXT)
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listUnreadEventsByUserFollowerIDQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (((%s) AND ((json_data->>'visibility')::INT = 20 OR (json_data->>'visibility')::INT = 30))
			OR ((json_data->> 'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
			OR (json_data->'target'->>'id') = $2::TEXT)
			AND json_data->>'created_at' > $3
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listUnreadEventsNoUserFollowersQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (((json_data->>'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
			OR (json_data->'target'->>'id') = $2::TEXT)
			AND json_data->>'created_at' > $3
			AND (json_data->>'enabled')::BOOL = true
			%s
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	countUnreadEventsByUserFollowerIDQuery = `SELECT count(*) AS "events"
			FROM app_%d_%d.events
			WHERE (((%s) AND ((json_data->>'visibility')::INT = 20 OR (json_data->>'visibility')::INT = 30))
				OR ((json_data->>'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
				OR (json_data->'target'->>'id') = $2::TEXT)
				AND json_data->>'created_at' > $3
				%s
				AND (json_data->>'enabled')::BOOL = true`

	countUnreadEventsNoUserFollowersQuery = `SELECT count(*) AS "events"
			FROM app_%d_%d.events
			WHERE (((json_data->>'visibility')::INT = 40 AND (json_data->>'user_id')::BIGINT != $1::BIGINT)
				OR (json_data->'target'->>'id') = $2::TEXT)
				AND json_data->>'created_at' > $3
				%s
				AND (json_data->>'enabled')::BOOL = true`

	updateApplicationUserLastReadQuery = `UPDATE app_%d_%d.users
		SET last_read = now()
		WHERE (json_data->>'id')::BIGINT = $1::BIGINT
			AND (json_data->>'enabled')::BOOL = true`

	listEventsByLocationQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE json_data->>'location' = $1::text
			AND (
				(json_data->>'visibility')::INT = 30 OR (json_data->>'visibility')::INT = 40 OR
				%s
				(json_data->>'user_id')::BIGINT = $2::bigint
			)
			AND (json_data->>'enabled')::BOOL = true
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByLatLonRadQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE ST_DWithin(geo, ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, true)
			AND (
				(json_data->>'visibility')::INT = 30 OR (json_data->>'visibility')::INT = 40 OR
				%s
				(json_data->>'user_id')::BIGINT = $4::bigint
			)
			AND (json_data->>'enabled')::BOOL = true
		ORDER BY json_data->>'created_at' DESC LIMIT 200`

	listEventsByLatLonNearQuery = `SELECT json_data
		FROM app_%d_%d.events
		WHERE (
				(json_data->>'visibility')::INT = 30 OR (json_data->>'visibility')::INT = 40 OR
				%s
				(json_data->>'user_id')::BIGINT = $1::bigint
			)
			AND (json_data->>'enabled')::BOOL = true
		ORDER BY ST_Distance_Sphere(geo, ST_SetSRID(ST_MakePoint($2, $3), 4326)), json_data->>'created_at' DESC LIMIT $4`
)

func (e *event) Create(accountID, applicationID int64, currentUserID uint64, event *entity.Event) []errors.Error {
	if event.ID == 0 {
		var err error
		event.ID, err = tgflake.FlakeNextID(applicationID, "events")
		if err != nil {
			return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

	}
	event.Enabled = true
	timeNow := time.Now()
	event.CreatedAt, event.UpdatedAt = &timeNow, &timeNow

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalEventCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = e.mainPg.
		Exec(appSchema(createEventQuery, accountID, applicationID), string(eventJSON), event.Latitude, event.Longitude)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalEventCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return nil
}

func (e *event) Read(accountID, applicationID int64, userID, eventID uint64) (*entity.Event, []errors.Error) {
	var JSONData string
	err := e.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectEventByIDQuery, accountID, applicationID), eventID, userID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrEventNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalEventRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	event := &entity.Event{}
	err = json.Unmarshal([]byte(JSONData), event)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	event.ID = eventID

	return event, nil
}

func (e *event) Update(accountID, applicationID int64, currentUserID uint64, existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	timeNow := time.Now()
	updatedEvent.UpdatedAt = &timeNow
	eventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = e.mainPg.Exec(
		appSchema(updateEventByIDQuery, accountID, applicationID),
		string(eventJSON), updatedEvent.Latitude, updatedEvent.Longitude, existingEvent.ID, existingEvent.UserID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	if !retrieve {
		return nil, nil
	}

	return &updatedEvent, nil
}

func (e *event) Delete(accountID, applicationID int64, userID, eventID uint64) []errors.Error {
	event, err := e.Read(accountID, applicationID, userID, eventID)

	if err != nil {
		return err
	}
	event.Enabled = false
	_, err = e.Update(accountID, applicationID, userID, *event, *event, false)

	return err
}

func (e *event) ListAll(orgID, appID int64, condition core.EventCondition) ([]*entity.Event, []errors.Error) {
	requestCondition, requestParams, er := condition.Process(1)
	if er != nil {
		return nil, er
	}

	query := fmt.Sprintf(listAllEventsQuery, orgID, appID, requestCondition)

	rows, err := e.pg.SlaveDatastore(-1).Query(query, requestParams...)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return e.rowsToSlice(rows)
}

func (e *event) ListUser(accountID, applicationID int64, userID, currentUserID uint64, condition *core.EventCondition) (events []*entity.Event, er []errors.Error) {
	events = []*entity.Event{}

	requestCondition, requestParams, er := condition.Process(2)
	if er != nil {
		return nil, er
	}

	var query string
	if userID == currentUserID {
		query = listAllEventsByUserIDQuery
	} else {
		relation, errs := e.conn.Relation(accountID, applicationID, currentUserID, userID)
		if errs != nil {
			return nil, errs
		}

		if *relation.IsFriend || *relation.IsFollowed {
			query = listConnectionEventsByUserIDQuery
		} else {
			query = listPublicEventsByUserIDQuery
		}
	}

	query = fmt.Sprintf(query, accountID, applicationID, requestCondition)

	var (
		rows *sql.Rows
		err  error
	)

	if len(requestParams) == 0 {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, userID)
	} else {
		requestParams = append([]interface{}{userID}, requestParams...)
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, requestParams...)
	}

	if err != nil {
		return events, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	return e.rowsToSlice(rows)
}

func (e *event) List(accountID, applicationID int64, userID uint64, condition *core.EventCondition) (events []*entity.Event, er []errors.Error) {
	connectionCondition, er := e.composeConnectionCondition(accountID, applicationID, userID, " OR ")
	if er != nil {
		return nil, er
	}

	requestCondition, requestParams, er := condition.Process(3)
	if er != nil {
		return nil, er
	}

	query := ""
	if connectionCondition == "" {
		query = fmt.Sprintf(listAllEventsWithUserEventsAndNoUserFollowersEventsQuery, accountID, applicationID, requestCondition)
	} else {
		query = fmt.Sprintf(listAllEventsWithUserEventsAndWithUserFollowersEventsQuery, accountID, applicationID, connectionCondition, requestCondition)
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(requestParams) == 0 {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, userID, strconv.FormatUint(userID, 10))
	} else {
		requestParams = append([]interface{}{userID, strconv.FormatUint(userID, 10)}, requestParams...)
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, requestParams...)
	}

	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	return e.rowsToSlice(rows)
}

func (e *event) UserFeed(accountID, applicationID int64, user *entity.ApplicationUser, condition *core.EventCondition) (count int, events []*entity.Event, er []errors.Error) {
	connectionCondition, er := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if er != nil {
		return 0, nil, er
	}

	requestCondition, requestParams, er := condition.Process(3)
	if er != nil {
		return 0, nil, er
	}

	query := ""
	if connectionCondition == "" {
		query = fmt.Sprintf(listEventsNoUserFollowersQuery, accountID, applicationID, requestCondition)
	} else {
		query = fmt.Sprintf(listEventsByUserFollowerIDQuery, accountID, applicationID, connectionCondition, requestCondition)
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(requestParams) == 0 {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, user.ID, strconv.FormatUint(user.ID, 10))
	} else {
		requestParams = append([]interface{}{user.ID, strconv.FormatUint(user.ID, 10)}, requestParams...)
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, requestParams...)
	}

	if err != nil {
		return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()

	unread := 0

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		if event.CreatedAt.Sub(*user.LastRead) > 0 {
			unread++
		}

		events = append(events, event)
	}

	if len(requestParams) == 0 {
		if err := e.updateApplicationUserLastRead(accountID, applicationID, user); err != nil {
			return 0, nil, err
		}
	}

	return unread, events, nil
}

func (e *event) UnreadFeed(accountID, applicationID int64, user *entity.ApplicationUser, condition *core.EventCondition) (count int, events []*entity.Event, er []errors.Error) {
	connectionCondition, er := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if er != nil {
		return 0, nil, er
	}

	requestCondition, requestParams, er := condition.Process(3)
	if er != nil {
		return 0, nil, er
	}

	query := ""
	if connectionCondition == "" {
		query = fmt.Sprintf(listUnreadEventsNoUserFollowersQuery, accountID, applicationID, requestCondition)
	} else {
		query = fmt.Sprintf(listUnreadEventsByUserFollowerIDQuery, accountID, applicationID, connectionCondition, requestCondition)
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(requestParams) == 0 {
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, user.ID, strconv.FormatUint(user.ID, 10), user.LastRead)
	} else {
		requestParams = append([]interface{}{user.ID, strconv.FormatUint(user.ID, 10), user.LastRead}, requestParams...)
		rows, err = e.pg.SlaveDatastore(-1).
			Query(query, requestParams...)
	}

	if err != nil {
		return 0, nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	events, er = e.rowsToSlice(rows)
	if er != nil {
		return
	}

	if len(requestParams) == 0 {
		if err := e.updateApplicationUserLastRead(accountID, applicationID, user); err != nil {
			return 0, nil, err
		}
	}

	return len(events), events, nil
}

func (e *event) UnreadFeedCount(accountID, applicationID int64, user *entity.ApplicationUser, condition *core.EventCondition) (count int, err []errors.Error) {
	friendCondition, err := e.composeConnectionCondition(accountID, applicationID, user.ID, " OR ")
	if err != nil {
		return 0, err
	}

	requestCondition, requestParams, err := condition.Process(3)
	if err != nil {
		return 0, err
	}

	query := ""
	if friendCondition == "" {
		query = fmt.Sprintf(countUnreadEventsNoUserFollowersQuery, accountID, applicationID, requestCondition)
	} else {
		query = fmt.Sprintf(countUnreadEventsByUserFollowerIDQuery, accountID, applicationID, friendCondition, requestCondition)
	}

	unread := 0
	var er error

	if len(requestParams) == 0 {
		er = e.pg.SlaveDatastore(-1).
			QueryRow(query, user.ID, strconv.FormatUint(user.ID, 10), user.LastRead).
			Scan(&unread)
	} else {
		requestParams = append([]interface{}{user.ID, strconv.FormatUint(user.ID, 10), user.LastRead}, requestParams...)
		er = e.pg.SlaveDatastore(-1).
			QueryRow(query, requestParams...).
			Scan(&unread)
	}

	if er != nil {
		return 0, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(er.Error()).SetCurrentLocation()}
	}
	// We cap the number of unread events to the number of maxium events we can retrieve, for now
	if unread > 200 {
		unread = 200
	}

	return unread, nil
}

func (e *event) WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID uint64, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (e *event) GeoSearch(accountID, applicationID int64, currentUserID uint64, latitude, longitude, radius float64, nearest int64) ([]*entity.Event, []errors.Error) {
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
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	return e.rowsToSlice(rows)
}

func (e *event) LocationSearch(accountID, applicationID int64, currentUserID uint64, locationKey string) ([]*entity.Event, []errors.Error) {
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
		return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	return e.rowsToSlice(rows)
}

func (e *event) updateApplicationUserLastRead(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error {
	_, err := e.mainPg.Exec(appSchema(updateApplicationUserLastReadQuery, accountID, applicationID), user.ID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return nil
}

func (e *event) rowsToSlice(rows *sql.Rows) (events []*entity.Event, err []errors.Error) {
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		event := &entity.Event{}
		err = json.Unmarshal([]byte(JSONData), event)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalEventsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		events = append(events, event)
	}
	return
}

func (e *event) composeConnectionCondition(accountID, applicationID int64, userID uint64, joinOperator string) (string, []errors.Error) {
	connections, er := e.conn.FriendsAndFollowingIDs(accountID, applicationID, userID)
	if er != nil {
		return "", er
	}

	if len(connections) == 0 {
		return "", nil
	}

	condition := "(json_data->>'user_id')::BIGINT IN ("
	con := connections[len(connections)-1]
	if len(connections) > 1 {
		connections = connections[:len(connections)-1]

		for idx := range connections {
			condition += fmt.Sprintf(`%d::BIGINT, `, connections[idx])
		}
	}
	condition += fmt.Sprintf(`%d::BIGINT)`, con)

	return condition, nil
}

// NewEvent returns a new event handler with PostgreSQL as storage driver
func NewEvent(pgsql postgres.Client) core.Event {
	return &event{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		conn:   NewConnection(pgsql),
	}
}
