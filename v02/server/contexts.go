/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
)

// ContextHasAccountID adds the accountID to the context
func ContextHasAccountID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.Bag["accountID"], er = strconv.ParseInt(ctx.Vars["accountId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse account id\n"+er.Error(), er.Error())
}

// ContextHasAccount adds the account to the context
func ContextHasAccount(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Bag["account"], err = core.ReadAccount(ctx.Bag["accountID"].(int64))
	return
}

// ContextHasAccountUserID adds the accountUserID to the context
func ContextHasAccountUserID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.Bag["accountUserID"], er = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse account user id\n"+er.Error(), er.Error())
}

// ContextHasAccountUser adds the accountUser to the context
func ContextHasAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Bag["accountUser"], err = core.ReadAccountUser(ctx.Bag["accountID"].(int64), ctx.Bag["accountUserID"].(int64))
	return
}

// ContextHasApplicationID adds the applicationID to the context
func ContextHasApplicationID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.Bag["applicationID"], er = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse application id\n"+er.Error(), er.Error())
}

// ContextHasApplication adds the application to the context
func ContextHasApplication(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Bag["application"], err = core.ReadApplication(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64))
	return
}

// ContextHasApplicationUserID adds the applicationUserID to the context
func ContextHasApplicationUserID(ctx *context.Context) (err *tgerrors.TGError) {
	var er error
	ctx.Bag["applicationUserID"], er = strconv.ParseInt(ctx.Vars["userId"], 10, 64)
	if er == nil {
		return
	}
	return tgerrors.NewBadRequestError("failed to parse application user id\n"+er.Error(), er.Error())
}

// ContextHasApplicationUser adds the applicationUser to the context
func ContextHasApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	ctx.Bag["applicationUser"], err = core.ReadApplicationUser(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64))
	return
}
