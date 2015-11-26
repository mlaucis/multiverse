package response

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
)

const (
	// APIVersion holds which API Version does this module holds
	APIVersion = "0.4"

	appRateLimitProduction int64 = 20000
	appRateLimitStaging    int64 = 100
	appRateLimitSeconds    int64 = 60
)

var currentRevision, currentHostname string

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
	ctx.StatusCode = code

	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx)
	CORSHandler(ctx)

	// TODO here it would be nice if we would consider the requested format when the stuff happens and deliver
	// either JSON or XML or FlatBuffers or whatever
	ctx.W.Header().Add("Content-Type", "application/json; charset=UTF-8")

	output := new(bytes.Buffer)
	err := json.NewEncoder(output).Encode(response)
	if err != nil {
		ctx.LogError(err)
		ErrorHappened(ctx, []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(err.Error()).SetCurrentLocation()})
		return
	}

	// We should only check these for valid responses, I think. Future me blame past me for this decision
	if (ctx.R.Method == "GET" || ctx.R.Method == "HEAD") && ctx.StatusCode < 300 {
		// We implememt the etag check first, aka the hard check, because we don't know if something else was changed
		// in the response, not just what we calculate the etag for.
		// Example situation: we compute the etag for getFeed as being the highest LastUpdated date but maybe
		// a user was updated meanwhile which would mean that the feed might be the same, event wise, but the user wise
		// it will be different so the app should retrieve the feed and process it as maybe the display name of the user
		// was changed or something else (thumbnail or whatever)

		h := md5.New()
		h.Write(output.Bytes())
		etag := h.Sum(nil)
		etagString := fmt.Sprintf("%x", etag)
		ctx.W.Header().Set("ETag", etagString)

		if requestEtag := ctx.R.Header.Get("If-None-Match"); requestEtag != "" {
			if requestEtag == etagString {
				ctx.StatusCode = http.StatusNotModified
				ctx.W.WriteHeader(ctx.StatusCode)
				return
			}
		}

		if ctx.R.Method == "HEAD" {
			ctx.W.WriteHeader(code)
			return
		}
	}

	// Write response

	// No gzip support
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		ctx.W.WriteHeader(code)
		io.Copy(ctx.W, output)
		return
	}

	// gzip support
	ctx.W.Header().Set("Content-Encoding", "gzip")
	ctx.W.WriteHeader(code)
	gz := gzip.NewWriter(ctx.W)
	io.Copy(gz, output)
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

	appRateLimit := appRateLimitStaging
	if ctx.Application.InProduction {
		appRateLimit = appRateLimitProduction
	}

	ctx.W.Header().Set("X-RateLimit-Limit", strconv.FormatInt(appRateLimit, 10))
	ctx.W.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Bag["rateLimit.limit"].(int64), 10))
	ctx.W.Header().Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Bag["rateLimit.refreshTime"].(time.Time).Unix(), 10))
}

// ErrorHappened handles the error message
func ErrorHappened(ctx *context.Context, errs []errors.Error) {
	errorMessage := entity.ErrorsResponse{Errors: []entity.ErrorResponse{}}
	for idx := range errs {
		errorMessage.Errors = append(errorMessage.Errors, entity.ErrorResponse{Code: errs[idx].Code(), Message: errs[idx].Error()})
	}
	WriteResponse(ctx, errorMessage, int(errs[0].Type()), 0)
	ctx.LogError(errs)
}

// Setup initializes the response
func Setup(revision, hostname string) {
	currentRevision, currentHostname = revision, hostname
}
