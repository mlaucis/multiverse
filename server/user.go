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

// getUser handles requests to retrieve a single user
// Request: GET /application/:AppID/user/:ID
// Test with: curl -i localhost/application/:AppID/user/:ID
func getUser(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		user   *entity.User
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

	// Read user
	if user, err = core.ReadUser(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create response
	response := &struct {
		appID int64 `json: "appId"`
		*entity.User
	}{
		appID: appID,
		User:  user,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// getUserList handles requests to retrieve all users of an app
// THIS ROUTE IS NOT YET ACTIVATED
// Request: GET /application/:AppID/users
// Test with: curl -i localhost/application/:AppID/users
func getUserList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		appID int64
		users []*entity.User
		err   error
	)
	// Read vars
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read users
	if users, err = core.ReadUserList(appID); err != nil {
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

// createUser handles requests create an application user
// Request: POST /application/:AppId/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost/application/:AppID/users
func createUser(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		user  = &entity.User{}
		appID int64
		err   error
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
	if err = decoder.Decode(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	user.ApplicationID = appID
	user.Enabled = true

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	// Write resource
	if user, err = core.WriteUser(user, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(user, http.StatusCreated, 0, w, r)
}
