/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package utils handles utils related things
package response

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/limiter"
)

type (
	errorResponse struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url,omitempty"`
	}
)

const (
	// Which API Version does this module holds
	APIVersion = "0.2"

	appRateLimit        int64 = 1000
	appRateLimitSeconds int64 = 60
)

var (
	appRateLimiter limiter.Limiter

	currentRevision, currentHostname string
)

// CORSHandler handles the OPTIONS requests to all defined paths
func CORSHandler(ctx *context.Context) []errors.Error {
	ctx.W.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding, Authorization")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
	return nil
}

// WriteResponse handles the http responses and returns the data
func WriteResponse(ctx *context.Context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx)
	CORSHandler(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

// WriteCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func WriteCommonHeaders(cacheTime uint, ctx *context.Context) {
	ctx.W.Header().Set("Strict-Transport-Security", "max-age=63072000")
	ctx.W.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.W.Header().Set("X-Frame-Options", "DENY")

	ctx.W.Header().Set("X-Tapglue-Hash", currentRevision)
	ctx.W.Header().Set("X-Tapglue-Server", currentHostname)

	if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`max-age=%d, public`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.W.Header().Set("Pragma", "no-cache")
		ctx.W.Header().Set("Expires", "0")
	}

	if !ctx.Bag["rateLimit.enabled"].(bool) {
		return
	}
	ctx.W.Header().Set("X-RateLimit-Limit", strconv.FormatInt(appRateLimit, 10))
	ctx.W.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Bag["rateLimit.limit"].(int64), 10))
	ctx.W.Header().Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Bag["rateLimit.refreshTime"].(time.Time).Unix(), 10))
}

// ErrorHappened handles the error message
func ErrorHappened(ctx *context.Context, errs []errors.Error) {
	errorMessage := map[string][]errorResponse{
		"errors": []errorResponse{},
	}
	for idx := range errs {
		errorMessage["errors"] = append(errorMessage["errors"], errorResponse{Code: errs[idx].Code(), Message: errs[idx].Error()})
	}
	WriteResponse(ctx, errorMessage, int(errs[0].Type()), 0)
	go ctx.LogError(errs)
}
