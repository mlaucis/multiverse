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
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

// getConnectionList handles requests to list a users connections
// Request: GET /application/:AppID/user/:UserID/connections
// Test with: curl -i localhost/application/:AppID/user/:UserID/connections
func getConnectionList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		users  []*entity.User
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

	// Read users
	if users, err = core.ReadConnectionList(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create response
	response := &struct {
		appID int64 `json: "appId"`
		Users []*entity.User
	}{
		appID: appID,
		Users: users,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createConnection handles requests to create a user connection
// Request: POST /application/:AppID/connections
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":2}' localhost/application/:AppID/connections
func createConnection(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		connection = &entity.Connection{}
		appID      int64
		err        error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Parse JSON
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(connection); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	connection.ApplicationID = appID
	connection.Enabled = true

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	// Write resource
	if connection, err = core.WriteConnection(connection, false); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(nil, http.StatusCreated, 0, w, r)
}
