package server

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/context"
	v01_validator "github.com/tapglue/backend/v01/validator"
	v01_keys "github.com/tapglue/backend/v01/validator/keys"
	v01_tokens "github.com/tapglue/backend/v01/validator/tokens"
)

// ValidateAccountRequestToken validates that the request contains a valid request token
func ValidateAccountRequestToken(ctx *context.Context) (err tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	return v01_keys.VerifyRequest(ctx, 1)
}

// ValidateApplicationRequestToken validates that the request contains a valid request token
func ValidateApplicationRequestToken(ctx *context.Context) (err tgerrors.TGError) {
	if ctx.SkipSecurity {
		return nil
	}

	return v01_tokens.VerifyRequest(ctx, 3)
}

// CheckAccountSession checks if the session token is valid or not
func CheckAccountSession(ctx *context.Context) (err tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	ctx.SessionToken, err = v01_validator.CheckAccountSession(ctx.R)
	return
}

// CheckApplicationSession checks if the session token is valid or not
func CheckApplicationSession(ctx *context.Context) (err tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	ctx.SessionToken, err = v01_validator.CheckApplicationSimpleSession(ctx)
	return
}
