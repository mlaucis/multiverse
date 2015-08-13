package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	event struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (e *event) Create(accountID, applicationID int64, currentUserID uint64, event *entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	if event.ID == 0 {
		return nil, []errors.Error{errmsg.ErrInternalEventMissingID}
	}
	evt := entity.EventWithIDs{}
	evt.AccountID = accountID
	evt.ApplicationID = applicationID
	event.Enabled = true
	evt.Event = *event
	data, er := json.Marshal(evt)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventCreate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return event, nil
	}

	return nil, nil
}

func (e *event) Read(accountID, applicationID int64, userID, eventID uint64) (event *entity.Event, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) Update(accountID, applicationID int64, currentUserID uint64, existingEvent, updatedEvent entity.Event, retrieve bool) (*entity.Event, []errors.Error) {
	evt := entity.EventWithIDs{}
	evt.AccountID = accountID
	evt.ApplicationID = applicationID
	evt.Event = updatedEvent
	data, er := json.Marshal(evt)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventUpdate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return &updatedEvent, nil
	}

	return nil, nil
}

func (e *event) Delete(accountID, applicationID int64, userID, eventID uint64) []errors.Error {
	evt := entity.EventWithIDs{}
	evt.AccountID = accountID
	evt.ApplicationID = applicationID
	evt.ID = eventID
	evt.UserID = userID
	data, er := json.Marshal(evt)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the event (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := e.storage.PackAndPutRecord(kinesis.StreamEventDelete, partitionKey, data)
	if err != nil {
		return []errors.Error{err}
	}

	return nil
}

func (e *event) List(accountID, applicationID int64, userID, currentUserID uint64) (events []*entity.Event, err []errors.Error) {
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

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID uint64, key string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) GeoSearch(accountID, applicationID int64, currentUserID uint64, latitude, longitude, radius float64, nearest int64) (events []*entity.Event, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) ObjectSearch(accountID, applicationID int64, currentUserID uint64, objectKey string) ([]*entity.Event, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (e *event) LocationSearch(accountID, applicationID int64, currentUserID uint64, locationKey string) ([]*entity.Event, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

// NewEvent creates a new Event
func NewEvent(storageClient kinesis.Client) core.Event {
	return &event{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
