package member

import (
	"errors"
	"fmt"
)

const errFmt = "%s: %s"

// Common errors for User service implementations and validations.
var (
	ErrInvalidMember = errors.New("invalid member")
	ErrNotFound      = errors.New("member not found")
)

// Error wraps common User errors.
type Error struct {
	err error
	msg string
}

func (e Error) Error() string {
	return e.msg
}

// IsInvalidMember indicates if err is ErrInvalidMember.
func IsInvalidMember(err error) bool {
	return unwrapError(err) == ErrInvalidMember
}

// IsNotFound indicates if err is ErrNotFound.
func IsNotFound(err error) bool {
	return unwrapError(err) == ErrNotFound
}

func unwrapError(err error) error {
	switch e := err.(type) {
	case *Error:
		return e.err
	}

	return err
}

func wrapError(err error, format string, args ...interface{}) error {
	return &Error{
		err: err,
		msg: fmt.Sprintf(
			errFmt,
			err,
			fmt.Sprintf(format, args...),
		),
	}
}
