/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

func getApplication(w http.ResponseWriter, r *http.Request) {
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
		*entity.Application
	}{
		Application: &entity.Application{
			AppID:     appId,
			Key:       "abc123def",
			AccountID: 123456,
			Name:      "Demo App",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, w)
}
