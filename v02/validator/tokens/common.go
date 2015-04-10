package tokens

import (
	"github.com/tapglue/backend/v02/core"
)

var (
	acc core.Account
	app core.Application
)

// Init initializes the modules
func Init(account core.Account, application core.Application) {
	acc = account
	app = application
}
