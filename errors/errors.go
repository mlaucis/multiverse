/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package errors holds the custom error
package errors

import (
	"fmt"
	"net/http"
	"runtime"
)

type (
	// errorType defines our custom type
	errorType uint16

	// Error holds our custom error
	Error interface {
		Type() errorType
		Error() string
		Raw() error
		ErrorWithLocation() string
		InternalErrorWithLocation() string
	}

	tgError struct {
		internalMessage string
		message         string
		location        string
		errType         errorType
	}
)

// TGInternalError represents an error that is caused by us
// TGClienterror represents an error that is caused by the client
// TGAuthenticationError represents an error that is because the client is not authenticated
const (
	TGInternalError     errorType = http.StatusInternalServerError
	TGBadRequestError   errorType = http.StatusBadRequest
	TGUnauthorizedError errorType = http.StatusUnauthorized
	TGNotFoundError     errorType = http.StatusNotFound
)

var dbgMode = false

// New generates a new error message
func New(errorType errorType, message, internalMessage string, withLocation bool) Error {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, message, internalMessage, stackDepth)
}

// NewFromError generates a new error message from an existing error
func NewFromError(errorType errorType, err error, withLocation bool) Error {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, err.Error(), err.Error(), stackDepth)
}

// NewBadRequestError generatates a new client error
func NewBadRequestError(message, internalMessage string) Error {
	return newError(TGBadRequestError, message, internalMessage, -1)
}

// NewInternalError generatates a new client error
func NewInternalError(message, internalMessage string) Error {
	return newError(TGInternalError, message, internalMessage, -1)
}

// NewUnauthorizedError generatates a new client error
func NewUnauthorizedError(message, internalMessage string) Error {
	return newError(TGUnauthorizedError, message, internalMessage, -1)
}

// NewNotFoundError generatates a new client error
func NewNotFoundError(message, internalMessage string) Error {
	return newError(TGNotFoundError, message, internalMessage, -1)
}

// Fatal will cause the message to be printed and then panic
func Fatal(message error) {
	panic(message)
}

func newError(errorType errorType, message, internalMessage string, stackDepth int) Error {
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
func (err *tgError) Type() errorType {
	return err.errType
}

// RawError generates a go error out of the existing error
func (err *tgError) Raw() error {
	return fmt.Errorf(err.Error())
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
