/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"log"
	"net/http"
	"time"
)

/**
 * Logger logs all server requests and prints to console
 * @param inner, http.Handler that is beeing used
 */
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}