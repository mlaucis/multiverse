/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"strconv"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/context"
	"github.com/tapglue/backend/v02/core"
)

func ContextHasAccountID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.AccountID, er = strconv.ParseInt(ctx.Vars["accountId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse account id\n"+er.Error(), er.Error())
}

func ContextHasAccount(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Account, err = core.ReadAccount(ctx.AccountID)
	return
}

func ContextHasAccountUserID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.AccountUserID, er = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse account user id\n"+er.Error(), er.Error())
}

func ContextHasAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.AccountUser, err = core.ReadAccountUser(ctx.AccountID, ctx.AccountUserID)
	return
}

func ContextHasApplicationID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.ApplicationID, er = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse application id\n"+er.Error(), er.Error())
}

func ContextHasApplication(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Application, err = core.ReadApplication(ctx.AccountID, ctx.ApplicationID)
	return
}

func ContextHasApplicationUserID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.ApplicationUserID, er = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse application user id\n"+er.Error(), er.Error())
}

func ContextHasApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.ApplicationUser, err = core.ReadApplicationUser(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID)
	return
}
