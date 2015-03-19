/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/v1/core"
	"github.com/tapglue/backend/v1/entity"
	"github.com/tapglue/backend/v1/validator"
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
		utils.ErrorHappened(ctx, "failed to retrieve the event (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if event, err = core.ReadEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		utils.ErrorHappened(ctx, "failed to retrieve the event (2)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, event, http.StatusOK, 10)
}

// updateEvent handles requests to update an event
// Request: PUT account/:AccountID/application/:ApplicationID/event/:EventID
func updateEvent(ctx *context.Context) {
	var (
		eventID int64
		err     error
	)

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		utils.ErrorHappened(ctx, "failed to update the event (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	existingEvent, err := core.ReadEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the event (2)", http.StatusInternalServerError, err)
		return
	}

	event := *existingEvent
	if err = json.NewDecoder(ctx.Body).Decode(&event); err != nil {
		utils.ErrorHappened(ctx, "failed to update the event (3)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	event.ID = eventID
	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

	if err = validator.UpdateEvent(existingEvent, &event); err != nil {
		utils.ErrorHappened(ctx, "failed to update the event (4)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedEvent, err := core.UpdateEvent(*existingEvent, event, true)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the event (5)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, updatedEvent, http.StatusCreated, 0)
}

// deleteEvent handles requests to delete a single event
// Request: DELETE account/:AccountID/application/:ApplicationID/event/:EventID
func deleteEvent(ctx *context.Context) {
	var (
		eventID int64
		err     error
	)

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the event (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteEvent(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, eventID); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the event (2)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "", http.StatusNoContent, 10)
}

// getEventList handles requests to retrieve a users events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/events
func getEventList(ctx *context.Context) {
	var (
		events []*entity.Event
		err    error
	)

	if events, err = core.ReadEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		utils.ErrorHappened(ctx, "failed to retrieve the event list (1)", http.StatusBadRequest, err)
		return
	}

	utils.WriteResponse(ctx, events, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections/events
func getConnectionEventList(ctx *context.Context) {
	var (
		events = []*entity.Event{}
		err    error
	)

	if events, err = core.ReadConnectionEventList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		utils.ErrorHappened(ctx, "failed to retrieve the connections event list (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, events, http.StatusOK, 10)
}

// createEvent handles requests to create an event
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/events
func createEvent(ctx *context.Context) {
	var (
		event = &entity.Event{}
		err   error
	)

	if err = json.NewDecoder(ctx.Body).Decode(event); err != nil {
		utils.ErrorHappened(ctx, "failed to create the event (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	event.AccountID = ctx.AccountID
	event.ApplicationID = ctx.ApplicationID
	event.UserID = ctx.ApplicationUserID

	if err = validator.CreateEvent(event); err != nil {
		utils.ErrorHappened(ctx, "failed to create the event (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		utils.ErrorHappened(ctx, "failed to create the event (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, event, http.StatusCreated, 0)
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
		utils.ErrorHappened(ctx, "failed to get the geo event list (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if longitude, err = strconv.ParseFloat(ctx.Vars["longitude"], 64); err != nil {
		utils.ErrorHappened(ctx, "failed to get the geo event list (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if radius, err = strconv.ParseFloat(ctx.Vars["radius"], 64); err != nil {
		utils.ErrorHappened(ctx, "failed to get the geo event list (3)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if events, err = core.SearchGeoEvents(ctx.AccountID, ctx.ApplicationID, latitude, longitude, radius); err != nil {
		utils.ErrorHappened(ctx, "failed to get the geo event list (4)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, events, http.StatusOK, 10)
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
		utils.ErrorHappened(ctx, "failed to get the object event list (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, events, http.StatusOK, 10)
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
		utils.ErrorHappened(ctx, "failed to get the location event list (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, events, http.StatusOK, 10)
}
