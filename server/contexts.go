/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
)

func contextHasAccount(ctx *context.Context) (err error) {
	ctx.AccountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64)
	if err != nil {
		return
	}

	ctx.Account, err = core.ReadAccount(ctx.AccountID)
	return
}

func contextHasApplication(ctx *context.Context) (err error) {
	ctx.ApplicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64)
	if err != nil {
		return
	}

	ctx.App, err = core.ReadApplication(ctx.AccountID, ctx.ApplicationID)
	return
}
