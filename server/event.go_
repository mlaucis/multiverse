/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/mysql"
)

// getApplicationEvent handles requests to retrieve a single event
// Request: GET /app/:AppID/event/:EventID
// Test with: curl -i localhost/app/:AppID/event/:EventID
func getApplicationEvent(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		event   = &entity.Event{}
		appID   uint64
		eventID uint64
		err     error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read eventID
	if eventID, err = strconv.ParseUint(vars["eventId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("eventId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if event, err = mysql.GetEventByID(eventID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if event.ApplicationID != appID {
		errorHappened(fmt.Errorf("event doesn't match expected values"), http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(event, http.StatusOK, 10, w, r)
}

// getApplicationUserEvents handles requests to retrieve a users events
// Request: GET /app/:AppID/user/:Token/events
// Test with: curl -i localhost/app/:AppID/user/:Token/events
func getApplicationUserEvents(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user      = &entity.User{}
		appID     uint64
		userToken string
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	userToken = vars["userToken"]

	if user, err = mysql.GetAllUserAppEvents(appID, userToken); err != nil {
		if config.Conf().Env() != "dev" {
			err = fmt.Errorf("could not retrieve the user events")
		}
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Write response
	writeResponse(user, http.StatusOK, 10, w, r)
}

// getUserConnectionsEvents handles requests to retrieve a users connections events
// Request: GET /app/:AppID/user/:Token/connections/events
// Test with: curl -i localhost/app/:AppID/user/:Token/connections/events
func getUserConnectionsEvents(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		events    = []*entity.Event{}
		appID     uint64
		userToken string
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userToken
	userToken = vars["userToken"]

	if events, err = mysql.GetUserConnectionsEvents(appID, userToken); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create mock response
	response := struct {
		AppID     uint64          `json:"appId"`
		UserToken string          `json:"userToken"`
		Events    []*entity.Event `json:"events"`
	}{
		AppID:     appID,
		UserToken: userToken,
		Events:    events,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createApplicationEvent handles requests to create an event
// Request: POST /app/:AppID/user/:userId/event
// Test with: curl -i -H "Content-Type: application/json" -d '{"type": "like", "item_id": "item1", "item_name": "item-name", "item_url": "app://url", "thumbnail_url": "gravatar", "custom": "{}", "nth": 1}' localhost/app/:appId/user/:userId/event
func createApplicationEvent(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		event  = &entity.Event{}
		appID  uint64
		userID uint64
		err    error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userToken
	if userID, err = strconv.ParseUint(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(event); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	event.ApplicationID = appID
	event.UserID = userID

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if event, err = mysql.AddSessionEvent(event); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(event, http.StatusCreated, 0, w, r)
}
