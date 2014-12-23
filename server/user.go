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

func getUser(w http.ResponseWriter, r *http.Request) {
	var (
		appID     uint64
		userToken string
		err       error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userToken = vars["userToken"]

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
			LastLogin:    "2014-12-20T12:10:10Z",
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getUserEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

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
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
			LastLogin:    "2014-12-20T12:10:10Z",
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
				CreatedAt: "2014-12-20T10:20:30Z",
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
				CreatedAt: "2014-12-20T10:23:30Z",
			},
			&entity.Event{
				ID:   0,
				Type: "ad",
				Item: &entity.Item{
					ID:   "0",
					Name: "Get more Gluee",
					URL:  "http://gluee.co",
				},
				CreatedAt: "2014-12-20T10:23:30Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getUserConnections(w http.ResponseWriter, r *http.Request) {
	var (
		appID     uint64
		userToken string
		err       error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userToken = vars["userToken"]

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

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getUserConnectionsEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	vars := mux.Vars(r)

	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		appID  uint64          `json: "appId"`
		Events []*entity.Event `json:"events"`
	}{
		appID: appID,
		Events: []*entity.Event{
			&entity.Event{
				ID:   1,
				Type: "read news",
				Item: &entity.Item{
					ID:   "1",
					Name: "Demo news",
					URL:  "app://news/1",
				},
				CreatedAt: "2014-12-20T10:20:30Z",
				User: &entity.User{
					Username:     "Onur",
					URL:          "app://user/1",
					ThumbnailURL: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
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
				CreatedAt: "2014-12-20T10:23:30Z",
				User: &entity.User{
					Username:     "Florin",
					URL:          "app://user/2",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
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
				CreatedAt: "2014-12-20T10:23:30Z",
			},
			&entity.Event{
				ID:   3,
				Type: "shared",
				Item: &entity.Item{
					ID:   "3",
					Name: "Gluee works",
					URL:  "app://item/3",
				},
				CreatedAt: "2014-12-20T10:30:30Z",
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
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
				CreatedAt: "2014-12-20T10:31:30Z",
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
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
				CreatedAt: "2014-12-20T10:35:30Z",
				User: &entity.User{
					Username:     "Norman",
					URL:          "app://user/3",
					ThumbnailURL: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				Custom: `{
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}
