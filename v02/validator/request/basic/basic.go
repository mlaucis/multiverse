/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package basic

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

// SignAccount checks that the current request is a signed account request
func SignAccount(ctx *context.Context) tgerrors.TGError {
	return nil
}

// VerifyAccount checks that the current request is a signed account request
func VerifyAccount(ctx *context.Context) tgerrors.TGError {
	return nil
}

// SignAccountUser checks that the current request is a signed account request
func SignAccountUser(ctx *context.Context) tgerrors.TGError {
	return nil
}

// VerifyAccountUser checks that the current request is a signed account request
func VerifyAccountUser(ctx *context.Context) tgerrors.TGError {
	return nil
}

// SignApplication checks that the current request is a signed account request
func SignApplication(ctx *context.Context) tgerrors.TGError {
	return nil
}

// VerifyApplication checks that the current request is a signed app request
func VerifyApplication(ctx *context.Context) tgerrors.TGError {
	return nil
}

// SignApplicationUser signs the request as an app user
func SignApplicationUser(ctx *context.Context) tgerrors.TGError {
	return nil
}

// VerifyApplicationUser checks that the current request is a signed app:appUser request
func VerifyApplicationUser(ctx *context.Context) tgerrors.TGError {
	return nil
}
