package kinesis

import (
	"net/http"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
)

type organization struct {
	writeStorage core.Organization
	readStorage  core.Organization
}

func (org *organization) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (org *organization) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (org *organization) Delete(ctx *context.Context) (err []errors.Error) {
	if ctx.Organization.PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch.SetCurrentLocation()}
	}

	if err = org.writeStorage.Delete(ctx.Organization); err != nil {
		return err
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (org *organization) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (org *organization) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// NewOrganization returns a new account handler tweaked specifically for Kinesis
func NewOrganization(writeStorage, readStorage core.Organization) handlers.Organization {
	return &organization{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
