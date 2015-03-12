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
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if event, err = core.ReadEvent(accountID, applicationID, userID, eventID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, event, http.StatusOK, 10)
}

// updateEvent handles requests to update an event
// Request: PUT account/:AccountID/application/:ApplicationID/event/:EventID
func updateEvent(ctx *context.Context) {
	var (
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if event.ID == 0 {
		event.ID = eventID
	}
	if event.AccountID == 0 {
		event.AccountID = accountID
	}
	if event.ApplicationID == 0 {
		event.ApplicationID = applicationID
	}
	if event.UserID == 0 {
		event.UserID = userID
	}

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
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.Vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteEvent(accountID, applicationID, userID, eventID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// getEventList handles requests to retrieve a users events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/events
func getEventList(ctx *context.Context) {
	var (
		events        []*entity.Event
		user          *entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadApplicationUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	if events, err = core.ReadEventList(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		User          *entity.User
		Events        []*entity.Event
	}{
		ApplicationID: applicationID,
		User:          user,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections/events
func getConnectionEventList(ctx *context.Context) {
	var (
		events        = []*entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if events, err = core.ReadConnectionEventList(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		UserID        int64           `json:"userId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: applicationID,
		UserID:        userID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// createEvent handles requests to create an event
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/events
func createEvent(ctx *context.Context) {
	var (
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	event.AccountID = accountID
	event.ApplicationID = applicationID
	event.UserID = userID

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
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radis
func getGeoEventList(ctx *context.Context) {
	var (
		events                      = []*entity.Event{}
		accountID, applicationID    int64
		latitude, longitude, radius float64
		err                         error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

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

	if events, err = core.SearchGeoEvents(accountID, applicationID, latitude, longitude, radius); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: applicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radis
func getObjectEventList(ctx *context.Context) {
	var (
		events                   = []*entity.Event{}
		accountID, applicationID int64
		objectKey                string
		err                      error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	objectKey = ctx.Vars["objectKey"]

	if events, err = core.SearchObjectEvents(accountID, applicationID, objectKey); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: applicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET account/:accountID/application/:applicationID/events/geo/:latitude/:longitude/:radis
func getLocationEventList(ctx *context.Context) {
	var (
		events                   = []*entity.Event{}
		accountID, applicationID int64
		location                 string
		err                      error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	location = ctx.Vars["location"]

	if events, err = core.SearchLocationEvents(accountID, applicationID, location); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := struct {
		ApplicationID int64           `json:"applicationId"`
		Events        []*entity.Event `json:"events"`
	}{
		ApplicationID: applicationID,
		Events:        events,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}
