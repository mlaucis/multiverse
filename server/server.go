/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func writeResponse(response interface{}, code int, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func errorHappened(message string, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%d %s", code, message)))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

func humans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
http://gluee.co
Location: Berlin, Germany.

/* THANKS */
Name: @dlsniper

/* SITE */
Last update: 2014/12/17
Standards: HTML5
Components: None
Software: Golang`))
}

func robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/app/{appId}", getApplication)
	r.HandleFunc("/app/{appId}/user/{userToken}", getUser)
	r.HandleFunc("/app/{appId}/user/{userToken}/events", getUserEvents)
	r.HandleFunc("/app/{appId}/user/{userToken}/connections", getUserConnections)
	r.HandleFunc("/app/{appId}/user/{userToken}/connections/events", getUserConnectionsEvents)
	r.HandleFunc("/app/{appId}/event/{eventId}", getEvent)
	r.HandleFunc("/account/{accountId}", getAccount)
	r.HandleFunc("/account/{accountId}/applications", getApplicationList)
	r.HandleFunc("/account/{accountId}/user/{userId}", getAccountUser)
	r.HandleFunc("/account/{accountId}/users", getAccountUserList)
	r.HandleFunc("/humans.txt", humans)
	r.HandleFunc("/robots.txt", robots)
	r.HandleFunc("/", home)
	r.StrictSlash(true)

	return r
}
