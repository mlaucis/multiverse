package redis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	event struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (e *event) Create(event *entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	data, er := json.Marshal(event)
	if er != nil {
		return nil, errors.NewInternalError("error while creating the event (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", event.AccountID, event.ApplicationID)
	_, err = e.storage.PackAndPutRecord(kinesis.StreamEventCreate, partitionKey, data)

	return nil, err
}

func (e *event) Read(accountID, applicationID, userID, eventID int64) (event *entity.Event, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) Update(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	data, er := json.Marshal(updatedEvent)
	if er != nil {
		return nil, errors.NewInternalError("error while updating the event (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", updatedEvent.AccountID, updatedEvent.ApplicationID)
	_, err = e.storage.PackAndPutRecord(kinesis.StreamEventUpdate, partitionKey, data)

	return nil, err
}

func (e *event) Delete(event *entity.Event) (err errors.Error) {
	data, er := json.Marshal(event)
	if er != nil {
		return errors.NewInternalError("error while deleting the event (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", event.AccountID, event.ApplicationID)
	_, err = e.storage.PackAndPutRecord(kinesis.StreamEventDelete, partitionKey, data)

	return err
}

func (e *event) List(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) WriteToConnectionsLists(event *entity.Event, key string) (err errors.Error) {
	return errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID, userID int64, key string) (err errors.Error) {
	return errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (e *event) LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

// NewEvent creates a new Event
func NewEvent(storageClient kinesis.Client) core.Event {
	return &event{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
