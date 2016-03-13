package server

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tapglue/multiverse/context"
)

type generalRoute struct {
	name    string
	path    string
	method  string
	handler func(*context.Context)
}

var generalRoutes = map[string]generalRoute{
	"humans": {
		name:    "humans",
		path:    "/humans.txt",
		method:  "GET",
		handler: humansHandler,
	},
	"robots": {
		name:    "robots",
		path:    "/robots.txt",
		method:  "GET",
		handler: robotsHandler,
	},
	"versions": {
		name:    "versions",
		path:    "/versions",
		method:  "GET",
		handler: versionsHandler,
	},
	"analytics": {
		name:    "analytics",
		path:    "/deprecated/analytics",
		method:  "POST",
		handler: analyticsHandler,
	},
	"healthcheck": {
		name:    "healthcheck",
		path:    "/deprecated/health-45016490610398192",
		method:  "GET",
		handler: healthCheckHandler,
	},
	"index": {
		name:    "home",
		path:    "/",
		method:  "GET",
		handler: homeHandler,
	},
}

var (
	notFoundResponse = []byte("{\"errors\":[{\"code\":0,\"message\":\"requested resource was not found\"}]}")
	robotsResponse   = []byte(`User-agent: *
Disallow: /`)
	humansResponse = []byte(`/* TEAM */
Founders: Norman Wiese, Onur Akpolat
Team: Florin Patan, Alexander Simmerl, Rafael Aviles
https://www.tapglue.com
Location: Berlin, Germany

/* SITE */
Last update: 2015/07/15
Software: Go, AWS, PostgreSQL, Redis, node.js`)
	analyticsOKResponse = []byte("ok")
)

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func homeHandler(ctx *context.Context) {
	WriteCommonHeaders(10*24*3600, ctx.W, ctx.R)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Header().Set("Refresh", "3; url=https://tapglue.com")
	ctx.W.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humansHandler(ctx *context.Context) {
	WriteCommonHeaders(10*24*3600, ctx.W, ctx.R)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write(humansResponse)
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robotsHandler(ctx *context.Context) {
	WriteCommonHeaders(10*24*3600, ctx.W, ctx.R)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write(robotsResponse)
}

// versionsHandler endpoint handles the api status for each api version
// Request: GET /versions
// Test with: curl -i localhost/versions
func versionsHandler(ctx *context.Context) {
	response := struct {
		Version map[string]struct {
			Version string `json:"version"`
			Status  string `json:"status"`
		} `json:"versions"`
		Revision string `json:"revision"`
	}{
		Version: map[string]struct {
			Version string `json:"version"`
			Status  string `json:"status"`
		}{
			"0.1": {"0.1", "disabled"},
			"0.2": {"0.2", "deprecated"},
			"0.3": {"0.3", "current"},
			"0.4": {"0.4", "alpha"},
		},
		Revision: currentRevision,
	}

	WriteCommonHeaders(7200, ctx.W, ctx.R)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(200)
		json.NewEncoder(ctx.W).Encode(response)
		return
	}

	ctx.W.Header().Set("Content-Encoding", "gzip")
	ctx.W.WriteHeader(200)
	gz := gzip.NewWriter(ctx.W)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

func analyticsHandler(ctx *context.Context) {
	WriteCommonHeaders(0, ctx.W, ctx.R)
	ctx.W.WriteHeader(200)
	ctx.W.Write(analyticsOKResponse)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(5, w, r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write response
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		w.WriteHeader(404)
		w.Write(notFoundResponse)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(404)
	gz := gzip.NewWriter(w)
	gz.Write(notFoundResponse)
	gz.Close()
}

func healthCheckHandler(ctx *context.Context) {
	// TODO make the checks concurrently

	response := struct {
		Healthy  bool `json:"healthy"`
		Services struct {
			PostgresMaster bool   `json:"postgres_master"`
			PostgresSlaves []bool `json:"postgres_slaves"`
			RateLimiter    bool   `json:"rate_limiter"`
			AppCache       bool   `json:"app_cache"`
		} `json:"services"`
	}{
		Healthy: true,
		Services: struct {
			PostgresMaster bool   `json:"postgres_master"`
			PostgresSlaves []bool `json:"postgres_slaves"`
			RateLimiter    bool   `json:"rate_limiter"`
			AppCache       bool   `json:"app_cache"`
		}{
			PostgresMaster: true,
			PostgresSlaves: make([]bool, rawPostgresClient.SlaveCount()),
			RateLimiter:    true,
			AppCache:       true,
		},
	}

	defer func() {
		if response.Healthy {
			ctx.StatusCode = 200
		} else {
			ctx.StatusCode = 500
		}

		WriteCommonHeaders(0, ctx.W, ctx.R)
		ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
		ctx.W.WriteHeader(ctx.StatusCode)
		json.NewEncoder(ctx.W).Encode(response)
	}()

	// Check Postgres
	if _, err := rawPostgresClient.MainDatastore().Exec("SELECT 1"); err != nil {
		response.Healthy = false
		response.Services.PostgresMaster = false
	}

	// TODO add exactly the slaves
	for slave := 0; slave < rawPostgresClient.SlaveCount(); slave++ {
		if _, err := rawPostgresClient.SlaveDatastore(slave).Exec("SELECT 1"); err != nil {
			response.Healthy = false
			response.Services.PostgresSlaves[slave] = false
		} else {
			response.Services.PostgresSlaves[slave] = true
		}
	}

	// Check Rate-Limiter
	rlConn := rateLimiterPool.Get()
	if err := rlConn.Err(); err != nil {
		response.Healthy = false
		response.Services.RateLimiter = false
		ctx.LogError(err)
	} else if rlConn != nil {
		rlConn.Close()
	}

	// Check AppCache
	acConn := appCachePool.Get()
	if err := acConn.Err(); err != nil {
		response.Healthy = false
		response.Services.AppCache = false
		ctx.LogError(err)
	} else if acConn != nil {
		acConn.Close()
	}
}
