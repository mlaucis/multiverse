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

	myError struct {
		internalMessage string
		message         string
		location        string
		extraData       map[string]string
		errType         errorType
	}
)

// This represents a list of common errors and their status codes
const (
	InternalError     errorType = http.StatusInternalServerError
	BadRequestError   errorType = http.StatusBadRequest
	UnauthorizedError errorType = http.StatusUnauthorized
	NotFoundError     errorType = http.StatusNotFound
	ConflictError     errorType = http.StatusConflict
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
	return newError(BadRequestError, message, internalMessage, -1)
}

// NewInternalError generatates a new client error
func NewInternalError(message, internalMessage string) Error {
	return newError(InternalError, message, internalMessage, -1)
}

// NewUnauthorizedError generatates a new client error
func NewUnauthorizedError(message, internalMessage string) Error {
	return newError(UnauthorizedError, message, internalMessage, -1)
}

// NewNotFoundError generatates a new client error
func NewNotFoundError(message, internalMessage string) Error {
	return newError(NotFoundError, message, internalMessage, -1)
}

// Fatal will cause the message to be printed and then panic
func Fatal(message error) {
	panic(message)
}

func newError(errorType errorType, message, internalMessage string, stackDepth int) Error {
	err := &myError{message: message, internalMessage: internalMessage, errType: errorType}
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
	err.extraData = map[string]string{}
	return err
}

// Type returns the type of the error
func (err *myError) Type() errorType {
	return err.errType
}

// RawError generates a go error out of the existing error
func (err *myError) Raw() error {
	return fmt.Errorf(err.Error())
}

// Error returns the error message
func (err *myError) Error() string {
	return err.message
}

// ErrorWithLocation returns the error and the location where it happened if that information is present
func (err *myError) ErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%q in %s", err.message, err.location)
	}
	return err.message
}

// InternalErrorWithLocation returns the internal error message and the location where it happened, if that exists
func (err *myError) InternalErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%q in %s", err.internalMessage, err.location)
	}
	return err.internalMessage
}

// ExtraData is usefull if you want to attach more data to the error while producing the error.
// This MUST be considered debug data and not be displayed to the user!
func (err *myError) ExtraData(key, value string) Error {
	err.extraData[key] = value
	return err
}

// Init initializes the logging module
func Init(debugMode bool) {
	dbgMode = debugMode
}
