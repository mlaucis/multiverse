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

	"github.com/gorilla/mux"
)

// getEvent handles requests to retrieve a single event
// Request: GET /application/:AppID/event/:ID
// Test with: curl -i localhost/0.1/application/:AppID/event/:ID
func getEvent(w http.ResponseWriter, r *http.Request) {
	var (
		event   = &entity.Event{}
		appID   int64
		userID  int64
		eventID int64
		err     error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if eventID, err = strconv.ParseInt(vars["eventId"], 10, 64); err != nil {
		errorHappened("eventId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if event, err = core.ReadEvent(appID, userID, eventID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(event, http.StatusOK, 10, w, r)
}

// updateEvent handles requests to update an event
// Request: PUT /application/:AppID/user/:UserID/event/:EventID
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' -X PUT localhost/0.1/application/:AppID/user/:UserID/event/:EventID
func updateEvent(w http.ResponseWriter, r *http.Request) {
	var (
		event   = &entity.Event{}
		appID   int64
		userID  int64
		eventID int64
		err     error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if eventID, err = strconv.ParseInt(vars["eventId"], 10, 64); err != nil {
		errorHappened("eventId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if event.ID == 0 {
		event.ID = eventID
	}
	if event.ApplicationID == 0 {
		event.ApplicationID = appID
	}
	if event.UserID == 0 {
		event.UserID = userID
	}

	if err = validator.UpdateEvent(event); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if event, err = core.UpdateEvent(event, true); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(event, http.StatusCreated, 0, w, r)
}

// deleteEvent handles requests to delete a single event
// Request: DELETE /application/:AppID/user/:UserID/event/:EventID
// Test with: curl -i -X DELETE localhost/0.1/application/:AppID/user/:UserID/event/:EventID
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	var (
		appID   int64
		userID  int64
		eventID int64
		err     error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if eventID, err = strconv.ParseInt(vars["eventId"], 10, 64); err != nil {
		errorHappened("eventId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteEvent(appID, userID, eventID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getEventList handles requests to retrieve a users events
// Request: GET /application/:AppID/user/:UserID/events
// Test with: curl -i localhost/0.1/application/:AppID/user/:UserID/events
func getEventList(w http.ResponseWriter, r *http.Request) {
	var (
		events []*entity.Event
		user   *entity.User
		appID  int64
		userID int64
		err    error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(appID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	if events, err = core.ReadEventList(appID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	response := &struct {
		appID  int64 `json: "appId"`
		User   *entity.User
		Events []*entity.Event
	}{
		appID:  appID,
		User:   user,
		Events: events,
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET /application/:AppID/user/:UserID/connections/events
// Test with: curl -i localhost/0.1/application/:AppID/user/:UserID/connections/events
func getConnectionEventList(w http.ResponseWriter, r *http.Request) {
	var (
		events = []*entity.Event{}
		appID  int64
		userID int64
		err    error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if events, err = core.ReadConnectionEventList(appID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	response := struct {
		AppID  int64           `json:"appId"`
		UserID int64           `json:"userId"`
		Events []*entity.Event `json:"events"`
	}{
		AppID:  appID,
		UserID: userID,
		Events: events,
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

// createEvent handles requests to create an event
// Request: POST /application/:AppID/user/:UserID/events
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' localhost/0.1/application/:AppID/user/:UserID/events
func createEvent(w http.ResponseWriter, r *http.Request) {
	var (
		event  = &entity.Event{}
		appID  int64
		userID int64
		err    error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	event.ApplicationID = appID
	event.UserID = userID

	if err = validator.CreateEvent(event); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if event, err = core.WriteEvent(event, true); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(event, http.StatusCreated, 0, w, r)
}
