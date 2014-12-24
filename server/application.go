/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
 	"encoding/json"
 	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)


func getApplication(w http.ResponseWriter, r *http.Request) {
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
		*entity.Application
	}{
		Application: &entity.Application{
			ID:        appID,
			Key:       "abc123def",
			AccountID: 123456,
			Name:      "Demo App",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

/*
createApplication handles the POST request to create a new app
Test with following command:
curl -H "Content-Type: application/json" -d '{"name":"New App"}' localhost:8082/app
*/
func createApplication(w http.ResponseWriter, r *http.Request) {
	var (
		app entity.Application
		err error
	)

	// Check that data isn't to big
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//TBD: Write app in database here

	response := &struct {
		*entity.Application
	}{
		Application: &entity.Application{
			ID:        1,
			Key:       "abc123def",
			AccountID: 123456,
			Name:      "Demo App",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusCreated, 10, w, r)
}

