/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package tgerrors holds the custom error
package tgerrors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type (
	// TGErrorType defines our custom type
	tgErrorType uint16

	// TGError holds our custom error
	TGError interface {
		Type() tgErrorType
		Error() string
		RawError() error
		ErrorWithLocation() string
		InternalErrorWithLocation() string
	}

	tgError struct {
		internalMessage string
		message         string
		location        string
		errType         tgErrorType
	}
)

// TGInternalError represents an error that is caused by us
// TGClienterror represents an error that is caused by the client
// TGAuthenticationError represents an error that is because the client is not authenticated
const (
	TGInternalError     tgErrorType = http.StatusInternalServerError
	TGBadRequestError   tgErrorType = http.StatusBadRequest
	TGUnauthorizedError tgErrorType = http.StatusUnauthorized
	TGNotFoundError     tgErrorType = http.StatusNotFound
)

var dbgMode = false

// New generates a new error message
func New(errorType tgErrorType, message, internalMessage string, withLocation bool) TGError {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, message, internalMessage, stackDepth)
}

// NewFromError generates a new error message from an existing error
func NewFromError(errorType tgErrorType, err error, withLocation bool) TGError {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, err.Error(), err.Error(), stackDepth)
}

// NewBadRequestError generatates a new client error
func NewBadRequestError(message, internalMessage string) TGError {
	return newError(TGBadRequestError, message, internalMessage, -1)
}

// NewInternalError generatates a new client error
func NewInternalError(message, internalMessage string) TGError {
	return newError(TGInternalError, message, internalMessage, -1)
}

// NewUnauthorizedError generatates a new client error
func NewUnauthorizedError(message, internalMessage string) TGError {
	return newError(TGUnauthorizedError, message, internalMessage, -1)
}

// NewNotFoundError generatates a new client error
func NewNotFoundError(message, internalMessage string) TGError {
	return newError(TGNotFoundError, message, internalMessage, -1)
}

func newError(errorType tgErrorType, message, internalMessage string, stackDepth int) TGError {
	err := &tgError{message: message, internalMessage: internalMessage, errType: errorType}
	if stackDepth == -1 && !dbgMode {
		return err
	}
	if dbgMode && stackDepth == -1 {
		stackDepth = 2
	}

	_, filename, line, _ := runtime.Caller(stackDepth)
	err.location = fmt.Sprintf("%s:%d", filename, line)
	if dbgMode {
		_, filename, line, _ := runtime.Caller(stackDepth + 1)
		err.location = fmt.Sprintf("%s from %s:%d", err.location, filename, line)

	}
	return err
}

// Type returns the type of the error
func (err *tgError) Type() tgErrorType {
	return err.errType
}

// RawError generates a go error out of the existing error
func (err *tgError) RawError() error {
	return errors.New(err.Error())
}

// Error returns the error message
func (err *tgError) Error() string {
	return err.message
}

// ErrorWithLocation returns the error and the location where it happened if that information is present
func (err *tgError) ErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%q in %s", err.message, err.location)
	}
	return err.message
}

// InternalErrorWithLocation returns the internal error message and the location where it happened, if that exists
func (err *tgError) InternalErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%q in %s", err.internalMessage, err.location)
	}
	return err.internalMessage
}

// Init initializes the logging module
func Init(debugMode bool) {
	dbgMode = debugMode
}
