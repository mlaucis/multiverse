/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import "github.com/tapglue/backend/errors"

var (
	invalidHandlerError = []errors.Error{errors.NewNotFoundError("not found", "invalid handler specified")}
)
