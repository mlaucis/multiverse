/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/tapglue/backend/logger"
	v01_server "github.com/tapglue/backend/v01/server"
	//v02_server "github.com/tapglue/backend/v02/server"

	"fmt"
	"time"

	"github.com/gorilla/mux"
)

/*const (
	errUserAgentNotSet           = "User-Agent header must be set (1)"
	errContentLengthNotSet       = "Content-Length header must be set (1)"
	errContentTypeNotSet         = "Content-Type header must be set (1)"
	errContentLengthNotDecodable = "Content-Length header value could not be decoded (2)"
	errContentLengthSizeNotMatch = "Content-Length header value is different from the received payload size (3)"
	errRequestBodyCannotBeEmpty  = "Request body cannot be empty (1)"
	errWrongContentType          = "Wrong Content-Type header value (1)"
)*/

var (
	mainLogChan  = make(chan *logger.LogMsg, 100000)
	errorLogChan = make(chan *logger.LogMsg, 100000)
)

// writeCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func writeCommonHeaders(cacheTime uint, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Strict-Transport-Security", "max-age=63072000")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	if cacheTime > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		w.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
}

// GetRouter creates the router
func GetRouter(environment string, debugMode, skipSecurityChecks bool) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	router := mux.NewRouter().StrictSlash(true)

	v01_server.Init(router, mainLogChan, errorLogChan, environment, debugMode, skipSecurityChecks)
	//v02_server.Init(router, mainLogChan, errorLogChan, environment, debugMode, skipSecurityChecks)

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	for routeName := range generalRoutes {
		router.
			Methods(generalRoutes[routeName].method, "OPTIONS").
			Path(generalRoutes[routeName].path).
			Name(routeName).
			HandlerFunc(generalRoutes[routeName].handler)
	}

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	return router, mainLogChan, errorLogChan, nil
}
