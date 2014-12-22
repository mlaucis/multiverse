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
		appId     uint64
		userToken string
		err       error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userToken = vars["userToken"]

	response := &struct {
		appID uint64 `json: "appId"`
		*entity.User
	}{
		appID: appId,
		User: &entity.User{
			Token:        userToken,
			Username:  "GlueUser123",
			Name: "Demo User",
			Email: "demouser@demo.com",
			URL:          "app://users/2",
			ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
			Custom:       `{"sound": "boo"}`,
			LastLogin:    "2014-12-20T12:10:10Z",
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, w)
}

func getUserEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appId uint64
		err   error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.User
		Events []*entity.Event `json:"events"`
	}{
		User: entity.User{
			AppID:        appId,
			Token:        "demoToken",
			Username:  "Demo User",
			URL:          "app://users/2",
			ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
			Custom:       `{"sound": "boo"}`,
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
			LastLogin:    "2014-12-20T12:10:10Z",
		},
		Events: []*entity.Event{
			&entity.Event{
				EventID:        1,
				EventType: "read news",
				ItemID:    "1",
				ItemName:  "Demo news",
				ItemURL:   "app://news/1",
				CreatedAt: "2014-12-20T10:20:30Z",
				Custom:    `{"key1": "value1"}`,
			},
			&entity.Event{
				EventID:        2,
				EventType: "like",
				ItemID:    "2",
				ItemName:  "Demo news",
				ItemURL:   "app://item/2",
				CreatedAt: "2014-12-20T10:23:30Z",
			},
			&entity.Event{
				EventID:        0,
				EventType: "ad",
				ItemID:    "0",
				ItemName:  "Get more Gluee",
				ItemURL:   "http://gluee.co",
				CreatedAt: "2014-12-20T10:23:30Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, w)
}

func getUserConnections(w http.ResponseWriter, r *http.Request) {
	var (
		appId     uint64
		userToken string
		err       error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userToken = vars["userToken"]

	response := &struct {
		appID uint64 `json: "appId"`
		entity.User
	}{
		appID: appId,
		User: entity.User{
			Token:       userToken,
			Username: "Demo User",
			URL:         "app://users/2",
			Connections: []*entity.User{
				&entity.User{
					Token:        "DemoToken1",
					Username:  "Onur",
					URL:          "app://user/1",
					ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				&entity.User{
					Token:        "DemoToken2",
					Username:  "Florin",
					URL:          "app://user/2",
					ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				&entity.User{
					Token:        "DemoToken3",
					Username:  "Norman",
					URL:          "app://user/3",
					Custom:       `{"sound": "boo"}`,
					ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					CreatedAt:    "2014-12-15T10:10:10Z",
					UpdatedAt:    "2014-12-20T12:10:10Z",
					LastLogin:    "2014-12-20T12:10:10Z",
				},
			},
		},
	}

	writeResponse(response, http.StatusOK, w)
}

func getUserConnectionsEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appId uint64
		err   error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		appID  uint64          `json: "appId"`
		Events []*entity.Event `json:"events"`
	}{
		appID: appId,
		Events: []*entity.Event{
			&entity.Event{
				EventID:        1,
				EventType: "read news",
				ItemID:    "1",
				ItemName:  "Demo news",
				ItemURL:   "app://news/1",
				CreatedAt: "2014-12-20T10:20:30Z",
				User: &entity.User{
					Username:  "Onur",
					URL:          "app://user/1",
					ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				EventID:        2,
				EventType: "like",
				ItemID:    "2",
				ItemName:  "Demo news",
				ItemURL:   "app://item/2",
				CreatedAt: "2014-12-20T10:23:30Z",
				User: &entity.User{
					Username:  "Florin",
					URL:          "app://user/2",
					ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
				},
			},
			&entity.Event{
				EventID:        0,
				EventType: "ad",
				ItemID:    "0",
				ItemName:  "Get more Gluee",
				ItemURL:   "http://gluee.co",
				CreatedAt: "2014-12-20T10:23:30Z",
			},
			&entity.Event{
				EventID:        3,
				EventType: "shared",
				ItemID:    "3",
				ItemName:  "Gluee works",
				ItemURL:   "app://item/3",
				CreatedAt: "2014-12-20T10:30:30Z",
				User: &entity.User{
					Username:  "Norman",
					URL:          "app://user/3",
					ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				Custom: `{"key1": "value1"}`,
			},
			&entity.Event{
				EventID:           4,
				EventType:    "picture",
				ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				ItemID:       "4",
				ItemName:     "Summer in Berlin",
				ItemURL:      "app://item/4",
				CreatedAt:    "2014-12-20T10:31:30Z",
				User: &entity.User{
					Username:  "Norman",
					URL:          "app://user/3",
					ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
					Custom:       `{"sound": "boo"}`,
					LastLogin:    "2014-12-20T12:10:10Z",
				},
				Custom: `{"largeUrl": "https://avatars0.githubusercontent.com/u/607868?v=3&s=460"}`,
			},
			&entity.Event{
				EventID:           5,
				EventType:    "pictures",
				ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
				ItemID:       "5",
				ItemName:     "Winter in London",
				ItemURL:      "app://item/5",
				CreatedAt:    "2014-12-20T10:35:30Z",
				User: &entity.User{
					Username:  "Norman",
					URL:          "app://user/3",
					ThumbnailUrl: "https://avatars0.githubusercontent.com/u/607868?v=3&s=460",
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

	writeResponse(response, http.StatusOK, w)
}
