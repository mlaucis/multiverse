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

// updateConnection handles requests to update a user connection
// Request: PUT /application/:AppID/user/:UserFromID/connection/:UserToID
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":2, "enabled":false}' -X PUT localhost/0.1/application/:AppID/user/:UserFromID/connection/:UserToID
func updateConnection(w http.ResponseWriter, r *http.Request) {
	var (
		connection = &entity.Connection{}
		appID      int64
		userFromID int64
		userToID   int64
		err        error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userFromID, err = strconv.ParseInt(vars["userFromId"], 10, 64); err != nil {
		errorHappened("userFromId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userToID, err = strconv.ParseInt(vars["userToId"], 10, 64); err != nil {
		errorHappened("userToId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(connection); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if connection.ApplicationID == 0 {
		connection.ApplicationID = appID
	}
	if connection.UserFromID == 0 {
		connection.UserFromID = userFromID
	}
	if connection.UserToID == 0 {
		connection.UserToID = userToID
	}

	if err = validator.UpdateConnection(connection); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if connection, err = core.UpdateConnection(connection, false); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(connection, http.StatusCreated, 0, w, r)
}

// deleteConnection handles requests to delete a single connection
// Request: DELETE /application/:AppID/user/:UserFromID/connection/:UserToID
// Test with: curl -i -X DELETE localhost/0.1/application/:AppID/user/:UserFromID/connection/:UserToID
func deleteConnection(w http.ResponseWriter, r *http.Request) {
	var (
		appID      int64
		userFromID int64
		userToID   int64
		err        error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userFromID, err = strconv.ParseInt(vars["userFromId"], 10, 64); err != nil {
		errorHappened("userFromId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userToID, err = strconv.ParseInt(vars["userToId"], 10, 64); err != nil {
		errorHappened("userToId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteConnection(appID, userFromID, userToID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getConnectionList handles requests to list a users connections
// Request: GET /application/:AppID/user/:UserID/connections
// Test with: curl -i localhost/0.1/application/:AppID/user/:UserID/connections
func getConnectionList(w http.ResponseWriter, r *http.Request) {
	var (
		users  []*entity.User
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

	if users, err = core.ReadConnectionList(appID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	response := &struct {
		appID int64 `json: "appId"`
		Users []*entity.User
	}{
		appID: appID,
		Users: users,
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

// TODO: GET FOLLOWER LIST (followedbyid users)
// TODO: GET FOLLOWING USERS LIST

// createConnection handles requests to create a user connection
// Request: POST /application/:AppID/user/:UserID/connections
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":2}' localhost/0.1/application/:AppID/user/:UserID/connections
func createConnection(w http.ResponseWriter, r *http.Request) {
	var (
		connection = &entity.Connection{}
		appID      int64
		userFromID int64
		err        error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userFromID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(connection); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if connection.UserFromID == 0 {
		connection.UserFromID = userFromID
	}
	connection.ApplicationID = appID

	if err = validator.CreateConnection(connection); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if connection, err = core.WriteConnection(connection, false); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(connection, http.StatusCreated, 0, w, r)
}
