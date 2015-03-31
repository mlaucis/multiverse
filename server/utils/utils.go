package utils

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/validator"
	"github.com/tapglue/backend/v01/validator/keys"
	"github.com/tapglue/backend/v01/validator/tokens"
)

type (
	// RouteFunc defines the pattern for a route handling function
	RouteFunc func(*context.Context) *tgerrors.TGError

	// Route holds the route pattern
	Route struct {
		Method   string
		Pattern  string
		CPattern string
		Scope    string
		Handlers []RouteFunc
		Filters  []context.Filter
	}
)

// RoutePattern returns the route pattern for a certain version
func (r *Route) RoutePattern(version string) string {
	if version == "" {
		return r.Pattern
	}
	return "/" + version + r.Pattern
}

// ComposePattern returns the composed pattern for a route
func (r *Route) ComposePattern(version string) string {
	return "/" + version + r.CPattern
}

// WriteResponse handles the http responses and returns the data
func WriteResponse(ctx *context.Context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx)
	WriteCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//Check if we have a session enable and if so write it back
	if ctx.SessionToken != "" {
		ctx.W.Header().Set("x-tapglue-session", ctx.SessionToken)
	}

	ctx.StatusCode = code

	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(code)
		json.NewEncoder(ctx.W).Encode(response)
		return
	}

	ctx.W.Header().Set("Content-Encoding", "gzip")
	ctx.W.WriteHeader(code)
	gz := gzip.NewWriter(ctx.W)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

// ErrorHappened handles the error message
func ErrorHappened(ctx *context.Context, err *tgerrors.TGError) {
	WriteCommonHeaders(0, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(int(err.Type))
		fmt.Fprintf(ctx.W, "%d %s", err.Type, err.Error())
	} else {
		ctx.W.Header().Set("Content-Encoding", "gzip")
		ctx.W.WriteHeader(int(err.Type))
		gz := gzip.NewWriter(ctx.W)
		fmt.Fprintf(gz, "%d %s", int(err.Type), err.Error())
		gz.Close()
	}

	ctx.StatusCode = int(err.Type)
	ctx.LogError(err)
}

// WriteCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func WriteCommonHeaders(cacheTime uint, ctx *context.Context) {
	ctx.W.Header().Set("Strict-Transport-Security", "max-age=63072000")
	ctx.W.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.W.Header().Set("X-Frame-Options", "DENY")

	if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.W.Header().Set("Pragma", "no-cache")
		ctx.W.Header().Set("Expires", "0")
	}
}

// WriteCorsHeaders will write the needed CORS headers
func WriteCorsHeaders(ctx *context.Context) {
	ctx.W.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding, x-tapglue-id, x-tapglue-date, x-tapglue-session, x-tapglue-payload-hash, x-tapglue-signature")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
}

// CorsHandler will handle the CORS requests
func CorsHandler(ctx *context.Context) (err *tgerrors.TGError) {
	WriteCommonHeaders(100, ctx)
	WriteCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return
}

// ValidateAccountRequestToken validates that the request contains a valid request token
func ValidateAccountRequestToken(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	return keys.VerifyRequest(ctx, 1)
}

// ValidateApplicationRequestToken validates that the request contains a valid request token
func ValidateApplicationRequestToken(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return nil
	}

	if ctx.Version == "0.1" {
		return tokens.VerifyRequest(ctx, 3)
	}
	return keys.VerifyRequest(ctx, 2)
}

// CheckAccountSession checks if the session token is valid or not
func CheckAccountSession(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	ctx.SessionToken, err = validator.CheckAccountSession(ctx.R)
	return
}

// CheckApplicationSession checks if the session token is valid or not
func CheckApplicationSession(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.Version == "0.1" {
		ctx.SessionToken, err = validator.CheckApplicationSimpleSession(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, ctx.R)
	} else {
		ctx.SessionToken, err = validator.CheckApplicationSession(ctx.R)
	}
	return
}
