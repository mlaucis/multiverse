/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package context provides the request context of the server request
package context

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/utils"

	"github.com/gorilla/mux"
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
		SessionToken      string
		StatusCode        int
		Vars              map[string]string
		Body              *bytes.Buffer
		Body2             *bytes.Buffer
		MainLog           chan *logger.LogMsg
		ErrorLog          chan *logger.LogMsg
		W                 http.ResponseWriter
		R                 *http.Request
		StartTime         time.Time
		RouteName         string
		Scope             string
		Version           string
		Environment       string
		DebugMode         bool
	}

	// Filter is a callback that helps updating the context with extra information
	Filter func(*Context) error
)

// LogRequest will generate a log message with the request status
// It is usable to log the request itself
func (ctx *Context) LogRequest(statusCode, stackDepth int) {
	if stackDepth != -1 {
		stackDepth++
	}
	msg := ctx.newLogMessage(stackDepth)
	msg.StatusCode = statusCode

	ctx.MainLog <- msg
}

// LogMessage logs a message from the application and includes all the required information
func (ctx *Context) LogMessage(message string, stackDepth int) {
	if stackDepth != -1 {
		stackDepth++
	}
	msg := ctx.newLogMessage(stackDepth)
	msg.Message = message

	ctx.MainLog <- msg
}

// LogError provides the ability to log and error
func (ctx *Context) LogError(err error, stackDepth int) {
	ctx.LogErrorWithMessage(err, "", stackDepth+1)
}

// LogErrorWithMessage will log an internal error from the app along with the custom message for it
func (ctx *Context) LogErrorWithMessage(err error, message string, stackDepth int) {
	if stackDepth != -1 {
		stackDepth++
	}
	msg := ctx.newLogMessage(stackDepth)
	msg.RawError = err
	msg.Message = message

	ctx.ErrorLog <- msg
}

// returns a new log message with the standard fields already populated
func (ctx *Context) newLogMessage(stackDepth int) *logger.LogMsg {
	location := ""
	if stackDepth != -1 {
		_, filename, line, _ := runtime.Caller(stackDepth + 1)
		location = fmt.Sprintf("%s:%d", filename, line)
	}

	requestPath := ctx.R.RequestURI
	if requestPath == "" {
		requestPath = ctx.R.URL.Path
	}

	return &logger.LogMsg{
		RemoteAddr: ctx.R.RemoteAddr,
		Method:     ctx.R.Method,
		RequestURI: requestPath,
		Headers:    ctx.R.Header,
		Payload:    ctx.Body2.String(),
		Name:       ctx.RouteName,
		Start:      ctx.StartTime,
		End:        time.Now(),
		Location:   location,
	}
}

// NewContext creates a new context from the current request
func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	mainLog, errorLog chan *logger.LogMsg,
	routeName, scope, version string,
	contextFilters []Filter,
	environment string,
	debugMode bool) (ctx *Context, err error) {

	ctx = new(Context)
	ctx.StartTime = time.Now()
	ctx.R = r
	ctx.W = w
	ctx.MainLog = mainLog
	ctx.ErrorLog = errorLog
	ctx.Vars = mux.Vars(r)
	if r.Method != "GET" {
		ctx.Body = utils.PeakBody(r)
		ctx.Body2 = utils.PeakBody(r)
	}
	ctx.RouteName = routeName
	ctx.Scope = scope
	ctx.Version = version
	ctx.Environment = environment
	ctx.DebugMode = debugMode

	for contextIndex, extraContext := range contextFilters {
		err = extraContext(ctx)
		if err != nil {
			err = fmt.Errorf("%s in context filter %d", err.Error(), contextIndex+1)
			break
		}
	}

	return ctx, err
}
