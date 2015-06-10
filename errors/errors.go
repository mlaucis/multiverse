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
		Code() int
		SetCode(errorCode int) Error
		Error() string
		Raw() error
		ErrorWithLocation() string
		InternalErrorWithLocation() string
		SetCurrentLocation() Error
		UpdateMessage(message string) Error
		UpdateInternalMessage(message string) Error
	}

	myError struct {
		internalMessage string
		message         string
		location        string
		errType         errorType
		errCode         int
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
func New(errorType errorType, errorCode int, message, internalMessage string, withLocation bool) Error {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, errorCode, message, internalMessage, stackDepth)
}

// NewFromError generates a new error message from an existing error
func NewFromError(errorType errorType, errorCode int, err error, withLocation bool) Error {
	stackDepth := -1
	if withLocation {
		stackDepth = 2
	}
	return newError(errorType, errorCode, err.Error(), err.Error(), stackDepth)
}

// NewBadRequestError generatates a new client error
func NewBadRequestError(errorCode int, message, internalMessage string) Error {
	return newError(BadRequestError, errorCode, message, internalMessage, -1)
}

// NewInternalError generatates a new client error
func NewInternalError(errorCode int, message, internalMessage string) Error {
	return newError(InternalError, errorCode, message, internalMessage, -1)
}

// NewUnauthorizedError generatates a new client error
func NewUnauthorizedError(errorCode int, message, internalMessage string) Error {
	return newError(UnauthorizedError, errorCode, message, internalMessage, -1)
}

// NewNotFoundError generatates a new client error
func NewNotFoundError(errorCode int, message, internalMessage string) Error {
	return newError(NotFoundError, errorCode, message, internalMessage, -1)
}

// Fatal will cause the message to be printed and then panic
func Fatal(message error) {
	panic(message)
}

func newError(errorType errorType, errorCode int, message, internalMessage string, stackDepth int) Error {
	if internalMessage == "" {
		internalMessage = message
	}

	err := myError{message: message, internalMessage: internalMessage, errType: errorType, errCode: errorCode}
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
func (err myError) Type() errorType {
	return err.errType
}

// RawError generates a go error out of the existing error
func (err myError) Raw() error {
	return fmt.Errorf(err.Error())
}

// Error returns the error message
func (err myError) Error() string {
	return err.message
}

// ErrorWithLocation returns the error and the location where it happened if that information is present
func (err myError) ErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%q in %s", err.errCode, err.message, err.location)
	}
	return err.message
}

// InternalErrorWithLocation returns the internal error message and the location where it happened, if that exists
func (err myError) InternalErrorWithLocation() string {
	if err.location != "" {
		return fmt.Sprintf("%d %q in %s", err.errCode, err.internalMessage, err.location)
	}
	return err.internalMessage
}

// SetCurrentLocation will update the error location to point to the invokation line rather that the creation line
func (err myError) SetCurrentLocation() Error {
	_, filename, line, _ := runtime.Caller(1)
	err.location = fmt.Sprintf("%s:%d", filename, line)
	return err
}

// UpdateMessage updates the error message
func (err myError) UpdateMessage(message string) Error {
	err.message = message
	return err
}

// UpdateInternalMessage updates the internal message of the error
func (err myError) UpdateInternalMessage(message string) Error {
	err.internalMessage = message
	return err
}

func (err myError) Code() int {
	return err.errCode
}

func (err myError) SetCode(errorCode int) Error {
	err.errCode = errorCode
	return err
}

// Init initializes the logging module
func Init(debugMode bool) {
	dbgMode = debugMode
}
