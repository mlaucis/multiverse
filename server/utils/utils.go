package utils

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/v1/validator"
	"github.com/tapglue/backend/v1/validator/keys"
	"github.com/tapglue/backend/v1/validator/tokens"
)

type (
	// RouteFunc defines the pattern for a route handling function
	RouteFunc func(*context.Context)

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
func ErrorHappened(ctx *context.Context, message string, code int, internalError error) {
	WriteCommonHeaders(0, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(code)
		fmt.Fprintf(ctx.W, "%d %s", code, message)
	} else {
		ctx.W.Header().Set("Content-Encoding", "gzip")
		ctx.W.WriteHeader(code)
		gz := gzip.NewWriter(ctx.W)
		fmt.Fprintf(gz, "%d %s", code, message)
		gz.Close()
	}

	ctx.StatusCode = code
	ctx.LogErrorWithMessage(internalError, message, 1)
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
func CorsHandler(ctx *context.Context) {
	if ctx.R.Method != "OPTIONS" {
		return
	}

	WriteCommonHeaders(100, ctx)
	WriteCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

// ValidateAccountRequestToken validates that the request contains a valid request token
func ValidateAccountRequestToken(ctx *context.Context) {
	if ctx.SkipSecurity {
		return
	}

	if errMsg, err := keys.VerifyRequest(ctx, 1); err != nil {
		ErrorHappened(ctx, errMsg, http.StatusUnauthorized, err)
	}
}

// ValidateApplicationRequestToken validates that the request contains a valid request token
func ValidateApplicationRequestToken(ctx *context.Context) {
	if ctx.SkipSecurity {
		return
	}

	var (
		errMsg string
		err    error
	)
	if ctx.Version == "0.1" {
		errMsg, err = tokens.VerifyRequest(ctx, 3)
	} else {
		errMsg, err = keys.VerifyRequest(ctx, 2)
	}

	if err != nil {
		ErrorHappened(ctx, errMsg, http.StatusUnauthorized, err)
	}
}

// CheckAccountSession checks if the session token is valid or not
func CheckAccountSession(ctx *context.Context) {
	if ctx.SkipSecurity {
		return
	}

	sessionToken, errMsg, err := validator.CheckAccountSession(ctx.R)
	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	ErrorHappened(ctx, errMsg, http.StatusUnauthorized, err)
}

// CheckApplicationSession checks if the session token is valid or not
func CheckApplicationSession(ctx *context.Context) {
	if ctx.SkipSecurity {
		return
	}

	var (
		errMsg, sessionToken string
		err                  error
	)

	if ctx.Version == "0.1" {
		sessionToken, errMsg, err = validator.CheckApplicationSimpleSession(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, ctx.R)
	} else {
		sessionToken, errMsg, err = validator.CheckApplicationSession(ctx.R)
	}

	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	ErrorHappened(ctx, errMsg, http.StatusUnauthorized, err)
}
