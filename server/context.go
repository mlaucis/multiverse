package server

import (
	"bytes"
	"net/http"
	"time"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"

	"github.com/gorilla/mux"
)

type (
	context struct {
		accID        int64
		acc          *entity.Account
		appID        int64
		app          *entity.Application
		userID       int64
		user         *entity.User
		sessionToken string
		vars         map[string]string
		body         *bytes.Buffer
		mainLog      chan *LogMsg
		errorLog     chan *LogMsg
		w            http.ResponseWriter
		r            *http.Request
		startTime    time.Time
		routeName    string
		scope        string
		version      string
	}
)

// NewContext creates a new context from the current request
func NewContext(w http.ResponseWriter, r *http.Request, mainLog, errorLog chan *LogMsg, routeName, scope, version string) (*context, error) {
	ctx := new(context)

	ctx.startTime = time.Now()
	ctx.r = r
	ctx.w = w
	ctx.mainLog = mainLog
	ctx.errorLog = errorLog
	ctx.vars = mux.Vars(r)
	ctx.body = utils.PeakBody(r)
	ctx.routeName = routeName
	ctx.scope = scope
	ctx.version = version

	return ctx, nil
}
