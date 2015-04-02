/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package context provides the request context of the server request
package context

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	ct "github.com/tapglue/backend/context"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/entity"
)

type (
	// Context struct holds the request, response and additional informations about the context of the request
	Context struct {
		AccountID         int64
		Account           *entity.Account
		AccountUserID     int64
		AccountUser       *entity.AccountUser
		ApplicationID     int64
		Application       *entity.Application
		ApplicationUserID int64
		ApplicationUser   *entity.User
		ct.Context
	}

	// Filter is a callback that helps updating the context with extra information
	Filter func(*Context) *tgerrors.TGError
)

// NewContext creates a new context from the current request
func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	mainLog, errorLog chan *logger.LogMsg,
	routeName, scope, version string,
	contextFilters []Filter,
	environment string,
	debugMode bool) (ctx *Context, err *tgerrors.TGError) {

	ctx = new(Context)
	ctx.StartTime = time.Now()
	ctx.R = r
	ctx.W = w
	ctx.MainLog = mainLog
	ctx.ErrorLog = errorLog
	ctx.Vars = mux.Vars(r)
	if r.Method != "GET" {
		ctx.Body = utils.PeakBody(r).Bytes()
	}
	ctx.RouteName = routeName
	ctx.Scope = scope
	ctx.Version = version
	ctx.Environment = environment
	ctx.DebugMode = debugMode

	for _, extraContext := range contextFilters {
		err = extraContext(ctx)
		if err != nil {
			break
		}
	}

	return ctx, err
}
