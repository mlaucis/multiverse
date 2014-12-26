/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"net/http"
	"strconv"

	//"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

// getUserSession handles requests to retrieve a single session
// Request: GET /app/:AppID/user/:Token/session/:SessionID
// Test with: curl -i localhost/app/:AppID/user/:Token/session/:SessionID
func getUserSession(w http.ResponseWriter, r *http.Request) {
	var (
		appID     uint64
		userToken string
		sessionID string
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	// Read userToken
	userToken = vars["userToken"]

	// Read sessionID
	sessionID = vars["sessionId"]

	// Create mock response
	response := &struct {
		appID     uint64 `json: "appId"`
		userToken string `json: "userToken"`
		sessionID string `json: "sessionId"`
	}{
		appID:     appID,
		userToken: userToken,
		sessionID: sessionID,
	}

	// Read session from database

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// getUserSessionList handles requests to retrieve all sessions of a user
// Request: GET /app/:AppID/user/:userToken/sessions
// Test with: curl -i localhost/app/:AppID/user/:userToken/sessions
func getUserSessionList(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	// Create mock response
	response := &struct {
		appID uint64 `json: "appId"`
	}{
		appID: appID,
	}

	// Read user from database

	// Query draft
	/**
	 * SELECT token, username, name, email, url, thumbnail_url, custom, last_login, created_at, updated_at
	 * FROM users
	 * WHERE app_id={appID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createUserSession handles requests create a user session
// Request: POST /app/:AppID/user/:userToken/session/:SessionID
// Test with: curl -H "Content-Type: application/json" -d '{"TBD"}' localhost/app/:AppID/user/:userToken/session/:SessionID
func createUserSession(w http.ResponseWriter, r *http.Request) {

}
