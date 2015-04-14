/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package request

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"

	httpBasic "github.com/tapglue/intaker/v02/validator/request/basic"
	"strings"
)

var errAuthMethodNotSupported = tgerrors.NewBadRequestError("authorization method not supported", "auth method not supported")

// SignAccount checks that the current request is a signed account request
func SignAccount(ctx *context.Context) tgerrors.TGError {
	return httpBasic.SignAccount(ctx)
}

// VerifyAccount checks that the current request is a signed account request
func VerifyAccount(ctx *context.Context) tgerrors.TGError {
	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if (strings.HasPrefix(httpAuthorization, "Basic ")) {
			return httpBasic.VerifyAccount(ctx)
		}
	}

	return errAuthMethodNotSupported
}

// SignAccountUser checks that the current request is a signed account request
func SignAccountUser(ctx *context.Context) tgerrors.TGError {
	return httpBasic.SignAccountUser(ctx)
}

// VerifyAccountUser checks that the current request is a signed account request
func VerifyAccountUser(ctx *context.Context) tgerrors.TGError {
	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if (strings.HasPrefix(httpAuthorization, "Basic ")) {
			return httpBasic.VerifyAccountUser(ctx)
		}
	}

	return errAuthMethodNotSupported
}

// SignApplication checks that the current request is a signed account request
func SignApplication(ctx *context.Context) tgerrors.TGError {
	return httpBasic.SignApplication(ctx)
}

// VerifyApplication checks that the current request is a signed app request
func VerifyApplication(ctx *context.Context) tgerrors.TGError {
	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if (strings.HasPrefix(httpAuthorization, "Basic ")) {
			return httpBasic.VerifyApplication(ctx)
		}
	}

	return errAuthMethodNotSupported
}

// SignApplicationUser signs the request as an app user
func SignApplicationUser(ctx *context.Context) tgerrors.TGError {
	return httpBasic.SignApplicationUser(ctx)
}

// VerifyApplicationUser checks that the current request is a signed app:appUser request
func VerifyApplicationUser(ctx *context.Context) tgerrors.TGError {
	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if (strings.HasPrefix(httpAuthorization, "Basic ")) {
			return httpBasic.VerifyApplicationUser(ctx)
		}
	}

	return errAuthMethodNotSupported
}
