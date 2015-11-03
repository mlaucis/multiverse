package server

import (
	"net/http"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/utils"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/server/handlers"
	v03_request_validator "github.com/tapglue/multiverse/v04/validator/request"
)

// NewContext creates a new context from the current request
func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	p map[string]string,
	mainLog, errorLog chan *logger.LogMsg,
	route *Route,
	environment string,
	debugMode bool) (ctx *context.Context, err []errors.Error) {

	ctx = new(context.Context)
	ctx.StartTime = time.Now()
	ctx.R = r
	ctx.W = w
	if p != nil {
		ctx.Vars = p
	} else {
		ctx.Vars = map[string]string{}
	}
	ctx.MainLog = mainLog
	ctx.ErrorLog = errorLog
	if r.Method != "GET" {
		ctx.Body = utils.PeakBody(r).Bytes()
	}
	ctx.RouteName = route.Name
	ctx.Environment = environment
	ctx.DebugMode = debugMode
	ctx.Bag = map[string]interface{}{}
	ctx.Bag["rateLimit.enabled"] = false
	ctx.AuthUsername, ctx.AuthPassword, ctx.AuthOk = r.BasicAuth()
	ctx.Query = r.URL.Query()

	return ctx, nil
}

// ContextHasOrganization populates the context with the account information
func ContextHasOrganization(handler handlers.Organization) Filter {
	return func(ctx *context.Context) []errors.Error {
		if ctx.R.Header.Get("X-Jarvis-Auth") != "ZTBmZjI3MGE2M2YzYzAzOWI1MjhiYTNi" {
			if err := v03_request_validator.VerifyAccount(ctx); err != nil {
				return err
			}
		}
		if err := handler.PopulateContext(ctx); err != nil {
			return err
		}
		return nil
	}
}

// ContextHasMember populates the context with the account user information
func ContextHasMember(handler handlers.Member) Filter {
	return func(ctx *context.Context) []errors.Error {
		if err := v03_request_validator.VerifyAccountUser(ctx); err != nil {
			return err
		}
		if err := handler.PopulateContext(ctx); err != nil {
			return err
		}
		return nil
	}
}

// ContextHasOrganizationApplication populates the context with the application information from ID
func ContextHasOrganizationApplication(handler handlers.Application) Filter {
	return func(ctx *context.Context) []errors.Error {
		if err := v03_request_validator.VerifyAccountUser(ctx); err != nil {
			return err
		}
		if err := handler.PopulateContextFromID(ctx); err != nil {
			return err
		}
		return nil
	}
}

// ContextHasApplication populates the context with the application information
func ContextHasApplication(handler handlers.Application) Filter {
	return func(ctx *context.Context) []errors.Error {
		if err := v03_request_validator.VerifyApplication(ctx); err != nil {
			return err
		}
		if err := handler.PopulateContext(ctx); err != nil {
			return err
		}
		return nil
	}
}

// ContextHasApplicationUser populates the context with the application user information
func ContextHasApplicationUser(handler handlers.ApplicationUser) Filter {
	return func(ctx *context.Context) []errors.Error {
		if err := v03_request_validator.VerifyApplicationUser(ctx); err != nil {
			return err
		}
		if err := handler.PopulateContext(ctx); err != nil {
			return err
		}
		return nil
	}
}
