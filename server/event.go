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
 * getApplicationEvent handles requests to retrieve a single event
 * Request: GET /app/:AppID/event/:EventID
 * Test with: curl -i localhost/app/:AppID/event/:EventID
 * @param w, response writer
 * @param r, http request
 */
func getApplicationEvent(w http.ResponseWriter, r *http.Request) {
	var (
		appID, eventID uint64
		err            error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	// Read eventID
	if eventID, err = strconv.ParseUint(vars["eventId"], 10, 64); err != nil {
		errorHappened("eventId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	// Create mock response
	response := &struct {
		appID uint64 `json: "appId"`
		*entity.Event
	}{
		appID: appID,
		Event: &entity.Event{
			ID:   eventID,
			Type: "read news",
			Item: &entity.Item{
				ID:   "1",
				Name: "Demo news",
				URL:  "app://news/1",
			},
			User: &entity.User{
				Username:     "Onur",
				URL:          "app://user/1",
				ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
				Custom:       `{"sound": "boo"}`,
				LastLogin:    api_demo_time,
			},
			Custom:    `{"key1": "value1"}`,
			CreatedAt: api_demo_time,
		},
	}

	// Read event and user from database

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getApplicationUserEvents handles requests to retrieve a users events
 * Request: GET /app/:AppID/user/:Token/events
 * Test with: curl -i localhost/app/:AppID/user/:Token/events
 * @param w, response writer
 * @param r, http request
 */
func getApplicationUserEvents(w http.ResponseWriter, r *http.Request) {
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
		entity.User
		Events []*entity.Event `json:"events"`
	}{
		User: entity.User{
			AppID:        appID,
			Token:        "demoToken",
			Username:     "Demo User",
			URL:          "app://users/2",
			ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
			Custom:       `{"sound": "boo"}`,
			CreatedAt:    api_demo_time,
			UpdatedAt:    api_demo_time,
			LastLogin:    api_demo_time,
		},
		Events: []*entity.Event{
			&entity.Event{
				ID:   1,
				Type: "read news",
				Item: &entity.Item{
					ID:   "1",
					Name: "Demo news",
					URL:  "app://news/1",
				},
				CreatedAt: api_demo_time,
				Custom:    `{"key1": "value1"}`,
			},
			&entity.Event{
				ID:   2,
				Type: "like",
				Item: &entity.Item{
					ID:   "2",
					Name: "Demo news",
					URL:  "app://item/2",
				},
				CreatedAt: api_demo_time,
			},
			&entity.Event{
				ID:   0,
				Type: "ad",
				Item: &entity.Item{
					ID:   "0",
					Name: "Get more Gluee",
					URL:  "http://gluee.co",
				},
				CreatedAt: api_demo_time,
			},
		},
	}

	// Read events and user from database

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getSessionEvents handles requests to retrieve a sessions events
 * Request: GET /app/:AppID/user/:userToken/session/:SessionID/events
 * Test with: curl -i localhost/app/:AppID/user/:userToken/session/:SessionID/events
 * @param w, response writer
 * @param r, http request
 */
func getSessionEvents(w http.ResponseWriter, r *http.Request) {
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
		appID     uint64          `json: "appId"`
		userToken string          `json: "userToken"`
		sessionID string          `json: "sessionId"`
		Events    []*entity.Event `json:"events"`
	}{
		appID:     appID,
		userToken: userToken,
		sessionID: sessionID,
		Events: []*entity.Event{
			&entity.Event{
				ID:   1,
				Type: "read news",
				Item: &entity.Item{
					ID:   "1",
					Name: "Demo news",
					URL:  "app://news/1",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Onur",
					URL:          "app://user/1",
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				ID:   2,
				Type: "like",
				Item: &entity.Item{
					ID:   "2",
					Name: "Demo news",
					URL:  "app://item/2",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Florin",
					URL:          "app://user/2",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
			},
			&entity.Event{
				ID:   0,
				Type: "ad",
				Item: &entity.Item{
					ID:   "0",
					Name: "Get more Gluee",
					URL:  "http://gluee.co",
				},
				CreatedAt: api_demo_time,
			},
			&entity.Event{
				ID:   3,
				Type: "shared",
				Item: &entity.Item{
					ID:   "3",
					Name: "Gluee works",
					URL:  "app://item/3",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				ID:           4,
				Type:         "picture",
				ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				Item: &entity.Item{
					ID:   "4",
					Name: "Summer in Berlin",
					URL:  "app://item/4",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
			&entity.Event{
				ID:           5,
				Type:         "pictures",
				ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				Item: &entity.Item{
					ID:   "5",
					Name: "Winter in London",
					URL:  "app://item/5",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
		},
	}

	// Read events and users from database

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getUserConnectionsEvents handles requests to retrieve a users connections events
 * Request: GET /app/:AppID/user/:Token/connections/events
 * Test with: curl -i localhost/app/:AppID/user/:Token/connections/events
 * @param w, response writer
 * @param r, http request
 */
func getUserConnectionsEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appID     uint64
		userToken string
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

	// Create mock response
	response := &struct {
		appID     uint64          `json: "appId"`
		userToken string          `json: "userToken"`
		Events    []*entity.Event `json:"events"`
	}{
		appID:     appID,
		userToken: userToken,
		Events: []*entity.Event{
			&entity.Event{
				ID:   1,
				Type: "read news",
				Item: &entity.Item{
					ID:   "1",
					Name: "Demo news",
					URL:  "app://news/1",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Onur",
					URL:          "app://user/1",
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				ID:   2,
				Type: "like",
				Item: &entity.Item{
					ID:   "2",
					Name: "Demo news",
					URL:  "app://item/2",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Florin",
					URL:          "app://user/2",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
			},
			&entity.Event{
				ID:   0,
				Type: "ad",
				Item: &entity.Item{
					ID:   "0",
					Name: "Get more Gluee",
					URL:  "http://gluee.co",
				},
				CreatedAt: api_demo_time,
			},
			&entity.Event{
				ID:   3,
				Type: "shared",
				Item: &entity.Item{
					ID:   "3",
					Name: "Gluee works",
					URL:  "app://item/3",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				ID:           4,
				Type:         "picture",
				ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				Item: &entity.Item{
					ID:   "4",
					Name: "Summer in Berlin",
					URL:  "app://item/4",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
			&entity.Event{
				ID:           5,
				Type:         "pictures",
				ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				Item: &entity.Item{
					ID:   "5",
					Name: "Winter in London",
					URL:  "app://item/5",
				},
				CreatedAt: api_demo_time,
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    api_demo_time,
				},
				Custom: `{
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
		},
	}

	// Read events and users from database

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * createApplicationEvent handles requests to create an event
 * Request: POST /app/:AppID/user/:userToken/session/:SessionID/event/:EventID
 * Test with: curl -H "Content-Type: application/json" -d '{"TBD"}' localhost/app/:AppID/user/:userToken/session/:SessionID/event/:EventID
 * @param w, response writer
 * @param r, http request
 */
func createApplicationEvent(w http.ResponseWriter, r *http.Request) {

}
