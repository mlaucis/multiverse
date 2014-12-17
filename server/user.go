/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gluee/backend/entity"
	"github.com/motain/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appId := vars["appId"]
	userId := vars["userId"]

	response := &struct {
		appID string `json: "appId"`
		entity.User
	}{
		appID: appId,
		User: entity.User{
			ID:          userId,
			Token:       "demoToken",
			DisplayName: "Demo User",
			URL:         "app://users/" + userId,
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
	vars := mux.Vars(r)
	appId := vars["appId"]
	userId := vars["userId"]

	response := &struct {
		appID string `json: "appId"`
		entity.User
		Events []entity.Event
	}{
		appID: appId,
		User: entity.User{
			ID:          userId,
			Token:       "demoToken",
			DisplayName: "Demo User",
			URL:         "app://users/" + userId,
		},
		Events: []entity.Event{
			entity.Event{
				ID:        "1",
				EventType: "like",
				ItemID:    "1",
				ItemURL:   "app://item/1",
				CreatedAt: "1418839005",
			},
			entity.Event{
				ID:        "2",
				EventType: "like",
				ItemID:    "2",
				ItemURL:   "app://item/2",
				CreatedAt: "1418839015",
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
