package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
)

// Chain takes a varidatic number of Middlewares and returns a combined
// Middleware.
func Chain(ms ...Middleware) Middleware {
	return func(handler Handler) Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			handler = ms[i](handler)
		}

		return handler
	}
}

// Wrap takes a Middleware and Handler and returns an http.HandlerFunc.
func Wrap(
	middleware Middleware,
	handler Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware(handler)(context.Background(), w, r)
	}
}

func respondError(w http.ResponseWriter, statusCode, code int, err error) {
	respondJSON(w, statusCode, struct {
		Errors map[string]string `json:"errors"`
	}{
		Errors: map[string]string{
			strconv.Itoa(code): err.Error(),
		},
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
