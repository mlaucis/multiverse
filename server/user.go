/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

/**
 * getApplicationUser handles requests to retrieve a single user
 * Request: GET /app/:AppID/user/:Token
 * Test with: curl -i localhost/app/:AppID/user/:Token
 * @param w, response writer
 * @param r, http request
 */
func getApplicationUser(w http.ResponseWriter, r *http.Request) {
	var (
		appID     uint64
		userToken string
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Read userToken
	userToken = vars["userToken"]

	// Create mock response
	response := &struct {
		appID uint64 `json: "appId"`
		*entity.User
	}{
		appID: appID,
		User: &entity.User{
			Token:        userToken,
			Username:     "GlueUser123",
			Name:         "Demo User",
			Email:        "demouser@demo.com",
			URL:          "app://users/2",
			ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
			Custom:       `{"sound": "boo"}`,
			LastLogin:    api_demo_time,
			CreatedAt:    api_demo_time,
			UpdatedAt:    api_demo_time,
		},
	}

	// Read user from database

	// Query draft
	/**
	 * SELECT token, username, name, email, url, thumbnail_url, custom, last_login, created_at, updated_at
	 * FROM users
	 * WHERE app_id={appID} AND token={userToken};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getApplicationUserList handles requests to retrieve all users of an app
 * Request: GET /app/:AppID/users
 * Test with: curl -i localhost/app/:AppID/users
 * @param w, response writer
 * @param r, http request
 */
func getApplicationUserList(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Create mock response
	response := &struct {
		appID uint64 `json: "appId"`
		*entity.User
	}{
		appID: appID,
		User: &entity.User{
			Token:        "1",
			Username:     "GlueUser123",
			Name:         "Demo User",
			Email:        "demouser@demo.com",
			URL:          "app://users/2",
			ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
			Custom:       `{"sound": "boo"}`,
			LastLogin:    api_demo_time,
			CreatedAt:    api_demo_time,
			UpdatedAt:    api_demo_time,
		},
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

/**
 * createApplicationUser handles requests create an application user
 * Request: POST /app/:AppID/user/:userToken
 * Test with: curl -H "Content-Type: application/json" -d '{"name":"User name"}' localhost/app/:AppID/user/:userToken
 * @param w, response writer
 * @param r, http request
 */
func createApplicationUser(w http.ResponseWriter, r *http.Request) {

}
