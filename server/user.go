/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	var (
		appId, userId uint64
		err           error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	if userId, err = strconv.ParseUint(vars["userId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		appID uint64 `json: "appId"`
		entity.User
	}{
		appID: appId,
		User: entity.User{
			ID:          userId,
			Token:       "demoToken",
			DisplayName: "Demo User",
			URL:         "app://users/" + strconv.FormatUint(userId, 10),
		},
	}

	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getUserEvents(w http.ResponseWriter, r *http.Request) {
	var (
		appId, userId uint64
		err           error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	if userId, err = strconv.ParseUint(vars["userId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.User
		Events []entity.Event `json:"events"`
	}{
		User: entity.User{
			ID:          userId,
			AppID:       appId,
			Token:       "demoToken",
			DisplayName: "Demo User",
			URL:         "app://users/" + strconv.FormatUint(userId, 10),
		},
		Events: []entity.Event{
			entity.Event{
				ID:        1,
				EventType: "like",
				ItemID:    "1",
				ItemURL:   "app://item/1",
				CreatedAt: 1418839005,
			},
			entity.Event{
				ID:        2,
				EventType: "like",
				ItemID:    "2",
				ItemURL:   "app://item/2",
				CreatedAt: 1418839015,
			},
		},
	}

	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
