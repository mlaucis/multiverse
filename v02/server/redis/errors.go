/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import "github.com/tapglue/backend/errors"

var (
	deprecatedStorageError = []errors.Error{errors.NewInternalError("deprecated storage used", "redis storage is deprecated")}
)
