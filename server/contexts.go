/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
)

func contextHasAccountID(ctx *context.Context) (err error) {
	ctx.AccountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64)
	return
}

func contextHasAccount(ctx *context.Context) (err error) {
	ctx.Account, err = core.ReadAccount(ctx.AccountID)
	return
}

func contextHasAccountUserID(ctx *context.Context) (err error) {
	ctx.AccountUserID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	return
}

func contextHasAccountUser(ctx *context.Context) (err error) {
	ctx.AccountUser, err = core.ReadAccountUser(ctx.AccountID, ctx.AccountUserID)
	return
}

func contextHasApplicationID(ctx *context.Context) (err error) {
	ctx.ApplicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64)
	return
}

func contextHasApplication(ctx *context.Context) (err error) {
	ctx.Application, err = core.ReadApplication(ctx.AccountID, ctx.ApplicationID)
	return
}

func contextHasApplicationUserID(ctx *context.Context) (err error) {
	ctx.ApplicationUserID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	return
}

func contextHasApplicationUser(ctx *context.Context) (err error) {
	ctx.ApplicationUser, err = core.ReadApplicationUser(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID)
	return
}
