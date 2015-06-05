/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import "github.com/tapglue/backend/errors"

var (
	badUserAgentError              = errors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
	contentTypeMissingError        = errors.NewBadRequestError("Content-Type header empty", "missing content-type header")
	contentTypeMismatchError       = errors.NewBadRequestError("Content-Type header mismatch", "content-type header mismatch")
	contentLengthMissingError      = errors.NewBadRequestError("Content-Length header missing", "missing content-length header")
	contentLengthInvalidError      = errors.NewBadRequestError("Content-Length header is invalid", "content-length header is not an int")
	contentLengthSizeMismatchError = errors.NewBadRequestError("Content-Length header size mismatch", "content-length header size mismatch")
	requestBodyEmpty               = errors.NewBadRequestError("Empty request body", "empty request body")
)
