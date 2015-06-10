/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import "github.com/tapglue/backend/errors"

var (
	notImplementedYet = []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
)
