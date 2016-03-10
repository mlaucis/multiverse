package event

import (
	"errors"
	"fmt"
)

// Common erorrs for Event service implementations.
var (
	ErrNamespaceNotFound = errors.New("namespace not found")
	ErrNotFound          = errors.New("event not found")
)

func wrapError(err error, format string, args ...interface{}) error {
	return fmt.Errorf(
		"%s: %s",
		err.Error(),
		fmt.Sprintf(format, args...),
	)
}
