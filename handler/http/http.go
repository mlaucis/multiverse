package http

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is the gateway specific http.HandlerFunc expecting a context.Context.
type Handler func(context.Context, http.ResponseWriter, *http.Request)
