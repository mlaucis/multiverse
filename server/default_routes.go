package server

import (
	"github.com/tapglue/backend/context"
	. "github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/tgerrors"
)

var defaultRoutes = map[string]*Route{
	// General
	"index": &Route{
		Method:   "GET",
		Pattern:  "/",
		CPattern: "/",
		Scope:    "/",
		Handlers: []RouteFunc{
			home,
		},
	},
	"humans": &Route{
		Method:   "GET",
		Pattern:  "/humans.txt",
		CPattern: "/humans.txt",
		Scope:    "humans",
		Handlers: []RouteFunc{
			humans,
		},
	},
	"robots": &Route{
		Method:   "GET",
		Pattern:  "/robots.txt",
		CPattern: "/robots.txt",
		Scope:    "robots",
		Handlers: []RouteFunc{
			robots,
		},
	},
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(ctx *context.Context) (err *tgerrors.TGError) {
	WriteCommonHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Header().Set("Refresh", "3; url=https://tapglue.com")
	ctx.W.Write([]byte(`these aren't the droids you're looking for`))
	ctx.StatusCode = 200

	return
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(ctx *context.Context) (err *tgerrors.TGError) {
	WriteCommonHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
Lead developer: Florin Patan
http://tapglue.com
Location: Berlin, Germany.

/* SITE */
Last update: 2015/03/15
Standards: HTML5
Components: None
Software: Go, Redis`))
	ctx.StatusCode = 200

	return
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(ctx *context.Context) (err *tgerrors.TGError) {
	WriteCommonHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`User-agent: *
Disallow: /`))
	ctx.StatusCode = 200

	return
}
