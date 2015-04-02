/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sort"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/context"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

// getEvent handles requests to retrieve a single event
// Request: GET account/:AccountID/application/:ApplicationID/event/:EventID
func getEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		event   = &entity.Event{}
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("read event failed (1)\n"+er.Error(), er.Error())
	}

	if event, err = core.ReadEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		return
	}

	WriteResponse(ctx, event, http.StatusOK, 10)
	return
}

// updateEvent handles requests to update an event
// Request: PUT account/:AccountID/application/:ApplicationID/event/:EventID
func updateEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to update the event (1)\n"+er.Error(), er.Error())
	}

	existingEvent, err := core.ReadEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID)
	if err != nil {
		return
	}

	event := *existingEvent
	if er = json.Unmarshal(ctx.Body, &event); er != nil {
		return tgerrors.NewBadRequestError("failed to update the event (2)\n"+er.Error(), er.Error())
	}

	event.ID = eventID
	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

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

// deleteEvent handles requests to delete a single event
// Request: DELETE account/:AccountID/application/:ApplicationID/event/:EventID
func deleteEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		eventID int64
		er      error
	)

	if eventID, er = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to delete the event (1)\n"+er.Error(), er.Error())
	}

	if err = core.DeleteEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// getEventList handles requests to retrieve a users events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/events
func getEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var events []*entity.Event

	if events, err = core.ReadEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections/events
func getConnectionEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var events = []*entity.Event{}

	if events, err = core.ReadConnectionEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// createEvent handles requests to create an event
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/events
func createEvent(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		event = &entity.Event{}
		er    error
	)

	if er = json.Unmarshal(ctx.Body, event); er != nil {
		return tgerrors.NewBadRequestError("failed to create the event (1)\n"+er.Error(), er.Error())
	}

	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

	if err = validator.CreateEvent(event); err != nil {
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		return
	}

	WriteResponse(ctx, event, http.StatusCreated, 0)
	return
}

// getGeoEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getGeoEventList(ctx *context.Context) (err *tgerrors.TGError) {
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

	if events, err = core.SearchGeoEvents(ctx.AccountID, ctx.ApplicationID, latitude, longitude, radius); err != nil {
		return
	}

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// getObjectEventList handles requests to retrieve events in a certain location / radius
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getObjectEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		events    = []*entity.Event{}
		objectKey string
	)

	objectKey = ctx.Vars["objectKey"]

	if events, err = core.SearchObjectEvents(ctx.AccountID, ctx.ApplicationID, objectKey); err != nil {
		return
	}

	sort.Sort(byIDDesc(events))

	WriteResponse(ctx, events, http.StatusOK, 10)
	return
}

// getLocationEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getLocationEventList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		events   = []*entity.Event{}
		location string
	)

	location = ctx.Vars["location"]

	if events, err = core.SearchLocationEvents(ctx.AccountID, ctx.ApplicationID, location); err != nil {
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
