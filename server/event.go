/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

// getEvent handles requests to retrieve a single event
// Request: GET /application/:AppID/event/:ID
// Test with: curl -i localhost/application/:AppID/event/:ID
func getEvent(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		event   = &entity.Event{}
		appID   int64
		userID  int64
		eventID int64
		err     error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read eventID
	if eventID, err = strconv.ParseInt(vars["eventId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("eventId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read event
	if event, err = core.ReadEvent(appID, userID, eventID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(event, http.StatusOK, 10, w, r)
}

// getEventList handles requests to retrieve a users events
// Request: GET /application/:AppID/user/:UserID/events
// Test with: curl -i localhost/application/:AppID/user/:UserID/events
func getEventList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		events []*entity.Event
		user   *entity.User
		appID  int64
		userID int64
		err    error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read user
	if user, err = core.ReadUser(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Read events
	if events, err = core.ReadEventList(appID, userID); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Create response
	response := &struct {
		appID  int64 `json: "appId"`
		User   *entity.User
		Events []*entity.Event
	}{
		appID:  appID,
		User:   user,
		Events: events,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// getConnectionEventList handles requests to retrieve a users connections events
// Request: GET /application/:AppID/user/:UserID/connections/events
// Test with: curl -i localhost/application/:AppID/user/:UserID/connections/events
func getConnectionEventList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		events = []*entity.Event{}
		appID  int64
		userID int64
		err    error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read events
	if events, err = core.ReadConnectionEventList(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create response
	response := struct {
		AppID  int64           `json:"appId"`
		UserID int64           `json:"userId"`
		Events []*entity.Event `json:"events"`
	}{
		AppID:  appID,
		UserID: userID,
		Events: events,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createEvent handles requests to create an event
// Request: POST /application/:AppID/user/:UserID/events
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' localhost/application/:AppID/user/:UserID/events
func createEvent(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		event  = &entity.Event{}
		appID  int64
		userID int64
		err    error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Parse JSON
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	event.ApplicationID = appID
	event.UserID = userID
	event.ReceivedAt = time.Now().UTC().UnixNano()

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	// Write resource
	if event, err = core.WriteEvent(event, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(event, http.StatusCreated, 0, w, r)
}
