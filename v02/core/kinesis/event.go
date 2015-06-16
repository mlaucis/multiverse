package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	event struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (e *event) Create(accountID, applicationID int64, currentUserID string, event *entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	data, er := json.Marshal(event)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventCreate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (e *event) Read(accountID, applicationID int64, userID, currentUserID, eventID string) (event *entity.Event, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) Update(accountID, applicationID int64, currentUserID string, existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	data, er := json.Marshal(updatedEvent)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventUpdate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (e *event) Delete(accountID, applicationID int64, currentUserID string, event *entity.Event) []errors.Error {
	data, er := json.Marshal(event)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventDelete, partitionKey, data)

	return []errors.Error{err}
}

func (e *event) List(accountID, applicationID int64, userID, currentUserID string) (events []*entity.Event, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) UserFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, err []errors.Error) {
	return 0, nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) UnreadFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, err []errors.Error) {
	return 0, nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) UnreadFeedCount(accountID, applicationID int64, user *entity.ApplicationUser) (count int, err []errors.Error) {
	return 0, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) GeoSearch(accountID, applicationID int64, currentUserID string, latitude, longitude, radius float64, nearest int64) (events []*entity.Event, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) ObjectSearch(accountID, applicationID int64, currentUserID, objectKey string) ([]*entity.Event, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) LocationSearch(accountID, applicationID int64, currentUserID, locationKey string) ([]*entity.Event, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

// NewEvent creates a new Event
func NewEvent(storageClient kinesis.Client) core.Event {
	return &event{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
