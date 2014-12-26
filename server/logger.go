/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

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

		// we should use this to sanitize any other headers that should not be exposed to the logs
		headers := r.Header

		inner.ServeHTTP(w, r)

		// TODO sanitize headers that shouldn't not appear in the logs

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
