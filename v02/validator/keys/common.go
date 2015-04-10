/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package keys

import (
	"github.com/tapglue/backend/v02/core"
)

var (
	acc core.Account
	app core.Application
)

// getScope returns the full-formed request scope and version
func getScope(date, scope, requestVersion string) string {
	return date + "/" + scope + "/" + requestVersion
}

// Init initializes the modules
func Init(account core.Account, application core.Application) {
	acc = account
	app = application
}
