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

// getUser handles requests to retrieve a single user
// Request: GET /application/:AppID/user/:ID
// Test with: curl -i localhost/0.1/application/:AppID/user/:ID
func getUser(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user   *entity.User
		appID  int64
		userID int64
		err    error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	response := &struct {
		appID int64 `json: "appId"`
		*entity.User
	}{
		appID: appID,
		User:  user,
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

// updateUser handles requests to update a user
// Request: PUT /application/:AppID/user/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}'  -X PUT localhost/0.1/application/:AppID/user/:ID
func updateUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user   = &entity.User{}
		appID  int64
		userID int64
		err    error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if user.ID == 0 {
		user.ID = userID
	}
	if user.ApplicationID == 0 {
		user.ApplicationID = appID
	}

	if err = validator.UpdateUser(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.UpdateUser(user, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(user, http.StatusCreated, 0, w, r)
}

// deleteUser handles requests to delete a single user
// Request: DELETE /application/:AppID/user/:ID
// Test with: curl -i -X DELETE localhost/0.1/application/:AppId/user/:ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := validateDeleteCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		appID  int64
		userID int64
		err    error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteUser(appID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getUserList handles requests to retrieve all users of an app
// THIS ROUTE IS NOT YET ACTIVATED
// Request: GET /application/:AppID/users
// Test with: curl -i localhost/0.1/application/:AppID/users
func getUserList(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		appID int64
		users []*entity.User
		err   error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if users, err = core.ReadUserList(appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
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

// createUser handles requests to create a user
// Request: POST /application/:AppID/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost/0.1/application/:AppID/users
func createUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user  = &entity.User{}
		appID int64
		err   error
	)

	vars := mux.Vars(r)

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	user.ApplicationID = appID
	user.Enabled = true

	if err = validator.CreateUser(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(user, http.StatusCreated, 0, w, r)
}
