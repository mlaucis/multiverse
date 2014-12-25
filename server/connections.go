/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

/**
 * getUserConnections handles requests to list a users connections
 * Request: GET /app/:AppID/user/:Token/connections
 * Test with: curl -i localhost/app/:AppID/user/:Token/connections
 * @param w, response writer
 * @param r, http request
 */
func getUserConnections(w http.ResponseWriter, r *http.Request) {
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
		entity.User
	}{
		appID: appID,
		User: entity.User{
			Token:    userToken,
			Username: "Demo User",
			URL:      "app://users/2",
			Connections: []*entity.User{
				&entity.User{
					Token:        "DemoToken1",
					Username:     "Onur",
					URL:          "app://user/1",
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				&entity.User{
					Token:        "DemoToken2",
					Username:     "Florin",
					URL:          "app://user/2",
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				&entity.User{
					Token:        "DemoToken3",
					Username:     "Norman",
					URL:          "app://user/3",
					Custom:       `{"sound": "boo"}`,
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
			},
		},
	}

	// Read users connections from database
	// TBD Query draft

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * createUserConnection handles requests to create a user connection
 * Request: POST /app/:AppID/connection
 * Test with: curl -H "Content-Type: application/json" -d '{"user_id1":"123456","user_id2":"654321"}' localhost/app/:AppID/connection
 * @param w, response writer
 * @param r, http request
 */
func createUserConnection(w http.ResponseWriter, r *http.Request) {

}