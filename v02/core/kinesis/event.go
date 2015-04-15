package redis

import (
	"github.com/tapglue/backend/tgerrors"
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

func (e *event) Create(event *entity.Event, retrieve bool) (evn *entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) Read(accountID, applicationID, userID, eventID int64) (event *entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) Update(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) Delete(accountID, applicationID, userID, eventID int64) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (e *event) List(accountID, applicationID, userID int64) (events []*entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) WriteToConnectionsLists(event *entity.Event, key string) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID, userID int64, key string) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (e *event) GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (e *event) LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

// NewEvent creates a new Event
func NewEvent(storageClient kinesis.Client) core.Event {
	return &event{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
