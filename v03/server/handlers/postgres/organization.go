package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type organization struct {
	storage core.Organization
}

func (org *organization) Read(ctx *context.Context) (err []errors.Error) {
	if ctx.Bag["account"] == nil {
		return []errors.Error{errmsg.ErrAccountMissingInContext}
	}

	if ctx.Bag["account"].(*entity.Organization).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	response.ComputeOrganizationLastModified(ctx, ctx.Bag["account"].(*entity.Organization))

	response.WriteResponse(ctx, ctx.Bag["account"].(*entity.Organization), http.StatusOK, 10)
	return
}

func (org *organization) Update(ctx *context.Context) (err []errors.Error) {
	account := *(ctx.Bag["account"].(*entity.Organization))

	if account.PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	account.ID = ctx.Bag["accountID"].(int64)

	if err := validator.UpdateOrganization(ctx.Bag["account"].(*entity.Organization), &account); err != nil {
		return err
	}

	updatedAccount, err := org.storage.Update(*(ctx.Bag["account"].(*entity.Organization)), account, true)
	if err != nil {
		return err
	}

	response.WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
	return nil
}

func (org *organization) Delete(ctx *context.Context) (err []errors.Error) {
	if ctx.Bag["account"].(*entity.Organization).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if err = org.storage.Delete(ctx.Bag["account"].(*entity.Organization)); err != nil {
		return err
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (org *organization) Create(ctx *context.Context) (err []errors.Error) {
	var account = &entity.Organization{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if err = validator.CreateOrganization(account); err != nil {
		return
	}

	if account, err = org.storage.Create(account, true); err != nil {
		return
	}

	response.WriteResponse(ctx, account, http.StatusCreated, 0)
	return
}

func (org *organization) PopulateContext(ctx *context.Context) (err []errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidAccountCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", user, pass))}
	}
	account, err := org.storage.FindByKey(user)
	if account == nil {
		return []errors.Error{errmsg.ErrAccountNotFound}
	}
	if err == nil {
		ctx.Bag["account"] = account
		ctx.Bag["accountID"] = account.ID
	}
	return
}

// NewOrganization returns a new account handler tweaked specifically for Kinesis
func NewOrganization(datastore core.Organization) handlers.Organization {
	return &organization{
		storage: datastore,
	}
}
