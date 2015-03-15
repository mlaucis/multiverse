/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getEvent handles requests to retrieve a single event
// Request: GET account/:AccountID/application/:ApplicationID/event/:EventID
func getEvent(ctx *context.Context) {
	var (
		event   = &entity.Event{}
		eventID int64
		err     error
	)

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if event, err = core.ReadEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, event, http.StatusOK, 10)
}

// updateEvent handles requests to update an event
// Request: PUT account/:AccountID/application/:ApplicationID/event/:EventID
func updateEvent(ctx *context.Context) {
	var (
		event   = &entity.Event{}
		eventID int64
		err     error
	)

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	event.ID = eventID
	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

	if err = validator.UpdateEvent(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if event, err = core.UpdateEvent(event, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, event, http.StatusCreated, 0)
}

// deleteEvent handles requests to delete a single event
// Request: DELETE account/:AccountID/application/:ApplicationID/event/:EventID
func deleteEvent(ctx *context.Context) {
	var (
		eventID int64
		err     error
	)

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// getEventList handles requests to retrieve a users events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/events
func getEventList(ctx *context.Context) {
	var (
		events []*entity.Event
		err    error
	)

	if events, err = core.ReadEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		UserID        int64 `json:"userId"`
		Events        []*entity.Event
	}{
		ApplicationID: ctx.ApplicationID,
		UserID:        ctx.ApplicationUserID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections/events
func getConnectionEventList(ctx *context.Context) {
	var (
		events = []*entity.Event{}
		err    error
	)

	if events, err = core.ReadConnectionEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		UserID        int64           `json:"userId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: ctx.ApplicationID,
		UserID:        ctx.ApplicationUserID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// createEvent handles requests to create an event
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/events
func createEvent(ctx *context.Context) {
	var (
		event = &entity.Event{}
		err   error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

	if err = validator.CreateEvent(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, event, http.StatusCreated, 0)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getGeoEventList(ctx *context.Context) {
	var (
		events                      = []*entity.Event{}
		latitude, longitude, radius float64
		err                         error
	)

	if latitude, err = strconv.ParseFloat(ctx.Vars["latitude"], 64); err != nil {
		errorHappened(ctx, "latitude is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if longitude, err = strconv.ParseFloat(ctx.Vars["longitude"], 64); err != nil {
		errorHappened(ctx, "longitude is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if radius, err = strconv.ParseFloat(ctx.Vars["radius"], 64); err != nil {
		errorHappened(ctx, "radius is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if events, err = core.SearchGeoEvents(ctx.AccountID, ctx.ApplicationID, latitude, longitude, radius); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: ctx.ApplicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getObjectEventList(ctx *context.Context) {
	var (
		events    = []*entity.Event{}
		objectKey string
		err       error
	)

	objectKey = ctx.Vars["objectKey"]

	if events, err = core.SearchObjectEvents(ctx.AccountID, ctx.ApplicationID, objectKey); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: ctx.ApplicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radius
func getLocationEventList(ctx *context.Context) {
	var (
		events   = []*entity.Event{}
		location string
		err      error
	)

	location = ctx.Vars["location"]

	if events, err = core.SearchLocationEvents(ctx.AccountID, ctx.ApplicationID, location); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: ctx.ApplicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}
