/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package context provides the request context of the server request
package context

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/utils"
)

type (
	// Context struct holds the request, response and additional informations about the context of the request
	Context struct {
		SessionToken string
		StatusCode   int
		Vars         map[string]string
		Body         []byte
		MainLog      chan *logger.LogMsg
		ErrorLog     chan *logger.LogMsg
		W            http.ResponseWriter
		R            *http.Request
		StartTime    time.Time
		RouteName    string
		Scope        string
		Version      string
		Environment  string
		DebugMode    bool
		SkipSecurity bool
		Bag          map[string]interface{}
	}

	// Filter is a callback that helps updating the context with extra information
	Filter func(*Context) errors.Error
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
func (ctx *Context) LogError(err interface{}) {
	var msg *logger.LogMsg
	if tgError, ok := err.(errors.Error); ok {
		msg = ctx.newLogMessage(-1)
		msg.RawError = tgError.InternalErrorWithLocation()
		msg.Message = tgError.Error()
	} else if er, ok := err.(error); ok {
		msg = ctx.newLogMessage(-1)
		msg.RawError = er.Error()
		msg.Message = er.Error()
	}

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
		Payload:    string(ctx.Body),
		Name:       ctx.RouteName,
		Start:      ctx.StartTime,
		End:        time.Now(),
		Location:   location,
	}
}

// New creates a new context from the current request
func New(
	w http.ResponseWriter,
	r *http.Request,
	mainLog, errorLog chan *logger.LogMsg,
	routeName, scope, version string,
	contextFilters []Filter,
	environment string,
	debugMode bool) (ctx *Context, err errors.Error) {

	ctx = new(Context)
	ctx.StartTime = time.Now()
	ctx.R = r
	ctx.W = w
	ctx.MainLog = mainLog
	ctx.ErrorLog = errorLog
	if r.Method == "POST" || r.Method == "PUT" {
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
