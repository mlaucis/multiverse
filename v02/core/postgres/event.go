/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	event struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

func (e *event) Create(event *entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) Read(accountID, applicationID, userID, eventID int64) (event *entity.Event, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) Update(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) Delete(*entity.Event) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) List(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error) {
	return []*entity.Event{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (e *event) ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error) {
	return []*entity.Event{}, errors.NewInternalError("not implemented yet", "not implemented yet")
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

// NewEvent returns a new event handler with PostgreSQL as storage driver
func NewEvent(pgsql postgres.Client) core.Event {
	return &event{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
