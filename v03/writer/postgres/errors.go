package postgres

import "github.com/tapglue/multiverse/errors"

var (
	errBadInputJSON = errors.NewInternalError(5800, "unable to decode the message", "")
)
