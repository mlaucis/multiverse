/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"log"
	"net/http"
	"time"
)

// Logger logs all server requests and prints to console
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// we should use this to sanitize any other headers that should not be exposed to the logs
		headers := getSanitizedHeaders(r)

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%+v\t%s\t%s\n",
			r.Method,
			r.RequestURI,
			headers,
			name,
			time.Since(start),
		)
	})
}
