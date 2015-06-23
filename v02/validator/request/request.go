/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package request

import (
	"net/http"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/errmsg"
	httpBasic "github.com/tapglue/backend/v02/validator/request/basic"
)

var errAuthMethodNotSupported = []errors.Error{errmsg.ErrAuthMethodNotSupported}

// SignAccount checks that the current request is a signed account request
func SignAccount(ctx *context.Context) []errors.Error {
	return httpBasic.SignAccount(ctx)
}

// VerifyAccount checks that the current request is a signed account request
func VerifyAccount(ctx *context.Context) []errors.Error {
	if ctx.SkipSecurity {
		return nil
	}

	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if strings.HasPrefix(httpAuthorization, "Basic ") {
			return nil
		}
	}

	return errAuthMethodNotSupported
}

// SignAccountUser checks that the current request is a signed account request
func SignAccountUser(ctx *context.Context) []errors.Error {
	return httpBasic.SignAccountUser(ctx)
}

// VerifyAccountUser checks that the current request is a signed account request
func VerifyAccountUser(ctx *context.Context) []errors.Error {
	if ctx.SkipSecurity {
		return nil
	}

	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if strings.HasPrefix(httpAuthorization, "Basic ") {
			return nil
		}
	}

	return errAuthMethodNotSupported
}

// SignApplication checks that the current request is a signed account request
func SignApplication(ctx *context.Context) []errors.Error {
	return httpBasic.SignApplication(ctx)
}

// VerifyApplication checks that the current request is a signed app request
func VerifyApplication(ctx *context.Context) []errors.Error {
	if ctx.SkipSecurity {
		return nil
	}

	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if strings.HasPrefix(httpAuthorization, "Basic ") {
			return nil
		}
	}

	return errAuthMethodNotSupported
}

// SignApplicationUser signs the request as an app user
func SignApplicationUser(ctx *context.Context) []errors.Error {
	return httpBasic.SignApplicationUser(ctx)
}

// VerifyApplicationUser checks that the current request is a signed app:appUser request
func VerifyApplicationUser(ctx *context.Context) []errors.Error {
	if ctx.SkipSecurity {
		return nil
	}

	if ctx.R.Header.Get("Authorization") != "" {
		httpAuthorization := ctx.R.Header.Get("Authorization")
		if strings.HasPrefix(httpAuthorization, "Basic ") {
			return nil
		}
	}

	return errAuthMethodNotSupported
}

// SignRequest is just a dummy placeholder
func SignRequest(r *http.Request, key string) error {
	return nil
}
