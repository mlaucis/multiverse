/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package context provides the request context of the server request
package context

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
)

// NewContext creates a new context from the current request
func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	mainLog, errorLog chan *logger.LogMsg,
	routeName, scope, version string,
	contextFilters []context.Filter,
	environment string,
	debugMode bool) (ctx *context.Context, err *tgerrors.TGError) {

	ctx = new(context.Context)
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
	ctx.Bag = map[string]interface{}{}

	for _, extraContext := range contextFilters {
		err = extraContext(ctx)
		if err != nil {
			break
		}
	}

	return ctx, err
}
