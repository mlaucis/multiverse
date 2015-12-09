package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// Handler is the gateway specific http.HandlerFunc expecting a context.Context.
type Handler func(context.Context, http.ResponseWriter, *http.Request)

// Middleware can be used to chain Handlers with different responsibilities.
type Middleware func(Handler) Handler

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

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func mapUsers(us user.Users) map[string]*v04_entity.ApplicationUser {
	m := map[string]*v04_entity.ApplicationUser{}

	for _, u := range us {
		m[strconv.FormatUint(u.ID, 10)] = u
	}

	return m
}

func respondError(w http.ResponseWriter, code int, err error) {
	statusCode := http.StatusInternalServerError

	e := unwrapError(err)

	switch {
	case e == ErrBadRequest:
		statusCode = http.StatusBadRequest
	case e == ErrLimitExceeded:
		statusCode = 429
	case e == ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case e == controller.ErrNotFound:
		statusCode = http.StatusNotFound
	}

	respondJSON(w, statusCode, struct {
		Errors []apiError `json:"errors"`
	}{
		Errors: []apiError{
			{Code: code, Message: err.Error()},
		},
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
