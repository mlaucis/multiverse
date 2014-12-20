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

func getEvent(w http.ResponseWriter, r *http.Request) {
	var (
		appId, eventId uint64
		err            error
	)
	vars := mux.Vars(r)

	if appId, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	if eventId, err = strconv.ParseUint(vars["eventId"], 10, 64); err != nil {
		errorHappened("eventId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		appID uint64 `json: "appId"`
		*entity.Event
	}{
		appID: appId,
		Event: &entity.Event{
			ID:        eventId,
			EventType: "read news",
			ItemID:    "1",
			ItemName:  "Demo news",
			ItemURL:   "app://news/1",
			CreatedAt: "2014-12-20T10:20:30Z",
			User: &entity.User{
				DisplayName:  "Onur",
				URL:          "app://user/1",
				ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
				Custom:       `{"sound": "boo"}`,
				LastLogin:    "2014-12-20T12:10:10Z",
			},
			Custom: `{"key1": "value1"}`,
		},
	}

	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
