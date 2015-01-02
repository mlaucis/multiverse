/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//"github.com/tapglue/backend/entity"
	"github.com/gorilla/mux"
	"github.com/tapglue/backend/db"
	"github.com/tapglue/backend/entity"
)

// getUserSession handles requests to retrieve a single session
// Request: GET /app/:AppID/user/:Token/session/:SessionID
// Test with: curl -i localhost/app/:AppID/user/:Token/session/:SessionID
func getUserSession(w http.ResponseWriter, r *http.Request) {
	var (
		sessionID uint64
		appID     uint64
		userToken string
		session   = &entity.Session{}
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if sessionID, err = strconv.ParseUint(vars["sessionId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userToken
	userToken = vars["userToken"]

	if session, err = db.GetSessionByID(sessionID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if session.AppID != appID || session.UserToken != userToken {
		errorHappened(fmt.Errorf("session doesn't match expected values"), http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(session, http.StatusOK, 10, w, r)
}

// getUserSessionList handles requests to retrieve all sessions of a user
// Request: GET /app/:AppID/user/:userToken/sessions
// Test with: curl -i localhost/app/:AppID/user/:userToken/sessions
func getUserSessionList(w http.ResponseWriter, r *http.Request) {
	var (
		user      *entity.User
		userToken string
		sessions  []*entity.Session
		appID     uint64
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

	if user, err = db.GetApplicationUserByToken(appID, userToken); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if sessions, err = db.GetAllUserSessions(appID, userToken); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create mock response
	response := &struct {
		*entity.User
		Sessions []*entity.Session
	}{
		User:     user,
		Sessions: sessions,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createUserSession handles requests create a user session
// Request: POST /app/:AppID/user/:userToken/session
// Test with: curl -i -H "Content-Type: application/json" -d '{"nth":1,"custom":"{}","gid":"tapglue_uid","model":"galaxy siv","manufacturer":"samsung","uuid":"uuid","idfa":"iddd","android_id":"1","platfrom":"android","os_version":"lollipop","browser":"","app_version":"1.0.1","sdk_version":"0.1","timezone":"+0100","language":"en","country":"de","city":"berlin","ip":"300.400.500.600","carrier":"vodasucks","network":"wifi"}' localhost/app/1/user/token1/session
func createUserSession(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		session   = &entity.Session{}
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

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(session); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	session.AppID = appID
	session.UserToken = userToken

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if session, err = db.AddUserSession(session); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(session, http.StatusCreated, 0, w, r)
}
