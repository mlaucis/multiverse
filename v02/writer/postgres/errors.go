/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import "github.com/tapglue/backend/errors"

var (
	errBadInputJSON = errors.NewInternalError(5800, "unable to decode the message", "")
)
