/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sort"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

// GetEvent handles requests to retrieve a single event
// Request: GET account/:AccountID/application/:ApplicationID/event/:EventID
func GetEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		event   = &entity.Event{}
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("read event failed (1)\n"+er.Error(), er.Error())
	}

	if event, err = core.ReadEvent(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), eventID); err != nil {
		return
	}

	WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

// UpdateEvent handles requests to update an event
// Request: PUT account/:AccountID/application/:ApplicationID/event/:EventID
func UpdateEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to update the event (1)\n"+er.Error(), er.Error())
	}

	existingEvent, err := core.ReadEvent(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return tgerrors.NewBadRequestError("failed to update the event (2)\n"+er.Error(), er.Error())
	}

	event.ID = eventID
	event.AccountID = ctx.Bag["accountID"].(int64)
	event.ApplicationID = ctx.Bag["applicationID"].(int64)
	event.UserID = ctx.Bag["applicationUserID"].(int64)

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		return
	}

	updatedEvent, err := core.UpdateEvent(*existingEvent, event, true)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
	return
}

// DeleteEvent handles requests to delete a single event
// Request: DELETE account/:AccountID/application/:ApplicationID/event/:EventID
func DeleteEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to delete the event (1)\n"+er.Error(), er.Error())
	}

	if err = core.DeleteEvent(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), eventID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// GetEventList handles requests to retrieve a users events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/events
func GetEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var events []*entity.Event

	if events, err = core.ReadEventList(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// GetConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections/events
func GetConnectionEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var events = []*entity.Event{}

	if events, err = core.ReadConnectionEventList(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// CreateEvent handles requests to create an event
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/events
func CreateEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return tgerrors.NewBadRequestError("failed to create the event (1)\n"+er.Error(), er.Error())
	}

	event.AccountID = ctx.Bag["accountID"].(int64)
	event.ApplicationID = ctx.Bag["applicationID"].(int64)
	event.UserID = ctx.Bag["applicationUserID"].(int64)

	if err = validator.CreateEvent(event); err != nil {
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		return
	}

	WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

// GetGeoEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func GetGeoEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		er                          error
	)

	if latitude, er = strconv.ParseFloat(ctx.Vars["latitude"], 64); er != nil {
		return tgerrors.NewBadRequestError("failed to read the event by geo (1)\n"+er.Error(), er.Error())
	}

	if longitude, er = strconv.ParseFloat(ctx.Vars["longitude"], 64); er != nil {
		return tgerrors.NewBadRequestError("failed to read the event by geo (2)\n"+er.Error(), er.Error())
	}

	if radius, er = strconv.ParseFloat(ctx.Vars["radius"], 64); er != nil {
		return tgerrors.NewBadRequestError("failed to read the event by geo (3)\n"+er.Error(), er.Error())
	}

	if radius < 1 {
		return tgerrors.NewBadRequestError("failed to read the event by geo (4)\nLocation radius can't be smaller than 2 meters", "radius smaller than 2")
	}

	if events, err = core.SearchGeoEvents(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), latitude, longitude, radius); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// GetObjectEventList handles requests to retrieve events in a certain location / radius
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func GetObjectEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		events    = []*entity.Event{}
		objectKey string
	)

	objectKey = ctx.Vars["objectKey"]

	if events, err = core.SearchObjectEvents(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), objectKey); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// GetLocationEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func GetLocationEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		events   = []*entity.Event{}
		location string
	)

	location = ctx.Vars["location"]

	if events, err = core.SearchLocationEvents(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), location); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

type byIDDesc []*entity.Event

func (a byIDDesc) Len() int           { return len(a) }
func (a byIDDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIDDesc) Less(i, j int) bool { return a[i].ID > a[j].ID }
