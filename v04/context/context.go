package context

import (
	"fmt"
	"runtime"
	"time"

	ctx "github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/v04/entity"
)

type (
	tokenType uint8

	// Context holds the context for the current request
	Context struct {
		ctx.Context

		TokenType         tokenType
		OrganizationID    int64
		MemberID          int64
		ApplicationID     int64
		ApplicationUserID uint64
		Organization      *entity.Organization
		Member            *entity.Member
		Application       *entity.Application
		ApplicationUser   *entity.ApplicationUser
	}
)

// Here we define the supported token types
const (
	TokenTypeUnknown tokenType = iota
	TokenTypeOrganization
	TokenTypeOrganizationID
	TokenTypeApplication
	TokenTypeBackend
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

// LogError provides the ability to log an error
func (ctx *Context) LogError(err interface{}) {
	var msg *logger.LogMsg
	if tgError, ok := err.([]errors.Error); ok {
		for _, tgErr := range tgError {
			msg := ctx.newLogMessage(-1)
			msg.StatusCode = int(tgErr.Type())
			msg.RawError = tgErr.InternalErrorWithLocation()
			msg.Message = tgErr.Error()
			ctx.ErrorLog <- msg
		}
		return
	} else if tgError, ok := err.(errors.Error); ok {
		msg = ctx.newLogMessage(-1)
		msg.StatusCode = int(tgError.Type())
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

	queryString := ""
	if ctx.R.URL.RawQuery != "" {
		queryString = "?" + ctx.R.URL.RawQuery
	}

	return &logger.LogMsg{
		Host:       ctx.R.Host,
		OrgID:      ctx.OrganizationID,
		AppID:      ctx.ApplicationID,
		UserID:     ctx.ApplicationUserID,
		RemoteAddr: ctx.R.RemoteAddr,
		Method:     ctx.R.Method,
		RequestURI: ctx.R.URL.Path + queryString,
		Headers:    ctx.R.Header,
		Payload:    string(ctx.Body),
		Name:       ctx.RouteName,
		Start:      ctx.StartTime,
		End:        time.Now(),
		Location:   location,
	}
}
