package object

import (
	"errors"
	"fmt"
)

// Common errors for Object validation and Service.
var (
	ErrInvalidAttachment = errors.New("invalid attachment")
	ErrInvalidObject     = errors.New("invalid object")
	ErrMissingReference  = errors.New("referenced object missing")
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrNotFound          = errors.New("object not found")
)

func wrapError(err error, format string, args ...interface{}) error {
	return fmt.Errorf(
		"%s: %s",
		err.Error(),
		fmt.Sprintf(format, args...),
	)
}
