/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"net/http"

	"github.com/motain/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

func humans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`/* TEAM */
Founder: Onur Akpolat
http://gluemobile.com
Twitter: @gluemobile
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
	r.HandleFunc("/app/{appId}/user/{userId}", getUser)
	r.HandleFunc("/app/{appId}/user/{userId}/events", getUserEvents)
	r.HandleFunc("/humans.txt", humans)
	r.HandleFunc("/robots.txt", robots)
	r.HandleFunc("/", home)

	return r
}
