/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getEvent handles requests to retrieve a single event
// Request: GET /application/:applicationId/event/:ID
// Test with: curl -i localhost/0.1/application/:applicationId/event/:ID
func getEvent(ctx *context) {
	var (
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if event, err = core.ReadEvent(accountID, applicationID, userID, eventID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, event, http.StatusOK, 10)
}

// updateEvent handles requests to update an event
// Request: PUT /application/:applicationId/user/:UserID/event/:EventID
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' -X PUT localhost/0.1/application/:applicationId/user/:UserID/event/:EventID
func updateEvent(ctx *context) {
	var (
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
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
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if event, err = core.UpdateEvent(event, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, event, http.StatusCreated, 0)
}

// deleteEvent handles requests to delete a single event
// Request: DELETE /application/:applicationId/user/:UserID/event/:EventID
// Test with: curl -i -X DELETE localhost/0.1/application/:applicationId/user/:UserID/event/:EventID
func deleteEvent(ctx *context) {
	var (
		accountID     int64
		applicationID int64
		userID        int64
		eventID       int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if eventID, err = strconv.ParseInt(ctx.vars["eventId"], 10, 64); err != nil {
		errorHappened(ctx, "eventId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if err = core.DeleteEvent(accountID, applicationID, userID, eventID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// getEventList handles requests to retrieve a users events
// Request: GET /application/:applicationId/user/:UserID/events
// Test with: curl -i localhost/0.1/application/:applicationId/user/:UserID/events
func getEventList(ctx *context) {
	var (
		events        []*entity.Event
		user          *entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if events, err = core.ReadEventList(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
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
// Request: GET /application/:applicationId/user/:UserID/connections/events
// Test with: curl -i localhost/0.1/application/:applicationId/user/:UserID/connections/events
func getConnectionEventList(ctx *context) {
	var (
		events        = []*entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if events, err = core.ReadConnectionEventList(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
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
// Request: POST /application/:applicationId/user/:UserID/events
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' localhost/0.1/application/:applicationId/user/:UserID/events
func createEvent(ctx *context) {
	var (
		event         = &entity.Event{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	event.AccountID = accountID
	event.ApplicationID = applicationID
	event.UserID = userID

	if err = validator.CreateEvent(event); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, event, http.StatusCreated, 0)
}
