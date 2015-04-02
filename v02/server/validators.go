package server

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/context"
	v02_validator "github.com/tapglue/backend/v02/validator"
	v02_keys "github.com/tapglue/backend/v02/validator/keys"
	v02_tokens "github.com/tapglue/backend/v02/validator/tokens"
)

// ValidateAccountRequestToken validates that the request contains a valid request token
func ValidateAccountRequestToken(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	return v02_keys.VerifyRequest(ctx, 1)
}

// ValidateApplicationRequestToken validates that the request contains a valid request token
func ValidateApplicationRequestToken(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return nil
	}

	return v02_tokens.VerifyRequest(ctx, 3)
}

// CheckAccountSession checks if the session token is valid or not
func CheckAccountSession(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	ctx.SessionToken, err = v02_validator.CheckAccountSession(ctx.R)
	return
}

// CheckApplicationSession checks if the session token is valid or not
func CheckApplicationSession(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	ctx.SessionToken, err = v02_validator.CheckApplicationSimpleSession(ctx)
	return
}
