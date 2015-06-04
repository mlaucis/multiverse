/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import "github.com/tapglue/backend/errors"

var (
	invalidAppIDError   = errors.NewBadRequestError("application id is not valid", "application id is not valid")
	invalidUserIDError  = errors.NewBadRequestError("user id is not valid", "user id is not valid")
	invalidEventIDError = errors.NewBadRequestError("event id is not valid", "event id is not valid")
)
