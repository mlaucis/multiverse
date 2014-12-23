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
			ID:   eventId,
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
				ThumbnailUrl: "https://avatars2.githubusercontent.com/u/1712926?v=3&s=460",
				Custom:       `{"sound": "boo"}`,
				LastLogin:    "2014-12-20T12:10:10Z",
			},
			Custom: `{"key1": "value1"}`,
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}
