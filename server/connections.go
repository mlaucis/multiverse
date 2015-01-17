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
	"github.com/tapglue/backend/entity"
	"github.com/tapglue/backend/mysql"
)

// getUserConnections handles requests to list a users connections
// Request: GET /app/:AppID/user/:Token/connections
// Test with: curl -i localhost/app/:AppID/user/:Token/connections
func getUserConnections(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user      *entity.User
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

	if user, err = mysql.GetApplicationUserWithConnections(appID, userToken); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	response := &struct {
		appID uint64 `json: "appId"`
		entity.User
	}{
		appID: appID,
		User:  *user,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createUserConnection handles requests to create a user connection
// Request: POST /app/:AppID/connection
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_id1":"123456","user_id2":"654321"}' localhost/app/:AppID/connection
func createUserConnection(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		connection = &entity.Connection{}
		appID      uint64
		err        error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(connection); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if err = mysql.AddApplicationUserConnection(appID, connection); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(nil, http.StatusCreated, 0, w, r)
}
