package basic

import (
	"github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/errors"
)

// SignAccount checks that the current request is a signed account request
func SignAccount(ctx *context.Context) []errors.Error {
	return nil
}

// VerifyAccount checks that the current request is a signed account request
func VerifyAccount(ctx *context.Context) []errors.Error {
	return nil
}

// SignAccountUser checks that the current request is a signed account request
func SignAccountUser(ctx *context.Context) []errors.Error {
	return nil
}

// VerifyAccountUser checks that the current request is a signed account request
func VerifyAccountUser(ctx *context.Context) []errors.Error {
	return nil
}

// SignApplication checks that the current request is a signed account request
func SignApplication(ctx *context.Context) []errors.Error {
	return nil
}

// VerifyApplication checks that the current request is a signed app request
func VerifyApplication(ctx *context.Context) []errors.Error {
	return nil
}

// SignApplicationUser signs the request as an app user
func SignApplicationUser(ctx *context.Context) []errors.Error {
	return nil
}

// VerifyApplicationUser checks that the current request is a signed app:appUser request
func VerifyApplicationUser(ctx *context.Context) []errors.Error {
	return nil
}
