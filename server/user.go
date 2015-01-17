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

// getApplicationUser handles requests to retrieve a single user
// Request: GET /app/:AppID/user/:Token
// Test with: curl -i localhost/app/:AppID/user/:Token
func getApplicationUser(w http.ResponseWriter, r *http.Request) {
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

	if user, err = mysql.GetApplicationUserByToken(appID, userToken); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	response := &struct {
		appID uint64 `json: "appId"`
		*entity.User
	}{
		appID: appID,
		User:  user,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// getApplicationUserList handles requests to retrieve all users of an app
// Request: GET /app/:AppID/users
// Test with: curl -i localhost/app/:AppID/users
func getApplicationUserList(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		appID uint64
		users []*entity.User
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if users, err = mysql.GetApplicationUsers(appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create mock response
	response := &struct {
		appID uint64 `json: "appId"`
		Users []*entity.User
	}{
		appID: appID,
		Users: users,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createApplicationUser handles requests create an application user
// Request: POST /app/:AppId/user
// Test with: curl -i -H "Content-Type: application/json" -d '{"token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "thumbnail_url": "gravatar", "custom": "{}"}' localhost/app/:AppID/user
func createApplicationUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		user  = &entity.User{}
		appID uint64
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// TODO validation should be added here, for example, name shouldn't be empty ;)
	user.ApplicationID = appID

	if user, err = mysql.AddApplicationUser(appID, user); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(user, http.StatusCreated, 0, w, r)
}
