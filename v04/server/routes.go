package server

import (
	"strings"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/handlers/postgres"
	"github.com/tapglue/multiverse/v04/server/handlers/redis"
	"github.com/tapglue/multiverse/v04/server/response"

	"github.com/gorilla/mux"
)

type (
	// RouteFunc defines the pattern for a route handling function
	RouteFunc func(*context.Context) []errors.Error

	// Filter for the context
	Filter func(*context.Context) []errors.Error

	// Route holds the route pattern
	Route struct {
		Name     string
		Method   string
		Path     string
		Handlers []RouteFunc
		Filters  []Filter
	}
)

const (
	organizationID      = "{accountID}"
	memberID            = "{accountUserID}"
	applicationID       = "{applicationID}"
	applicationUserID   = "{applicationUserID:[0-9]+}"
	applicationUserToID = "{applicationUserToID:[0-9]+}"
	connectionType      = "{connectionType}"
	connectionState     = "{connectionState}"
	eventID             = "{eventID:[0-9]+}"
)

var (
	postgresOrganizationHandler                         handlers.Organization
	postgresMemberHandler                               handlers.Member
	postgresApplicationHandler, redisApplicationHandler handlers.Application
	postgresApplicationUserHandler                      handlers.ApplicationUser
	postgresConnectionHandler                           handlers.Connection
	postgresEventHandler                                handlers.Event

	applicationUserIDPattern = "%d"
	connectionTypePattern    = "%s"
	connectionStatePattern   = "%s"
	eventIDPattern           = "%d"
)

// RoutePattern returns the full route path, inclulding the api version
func (r *Route) RoutePattern() string {
	return "/" + APIVersion + r.Path
}

// ReplaceTestApplicationUserIDPattern is used in testing for replacing the application user id pattern
func ReplaceTestApplicationUserIDPattern(pattern string) {
	applicationUserIDPattern = pattern
}

// ReplaceTestEventIDPattern is used in testing for replacing the event id pattern
func ReplaceTestEventIDPattern(pattern string) {
	eventIDPattern = pattern
}

// TestPattern returns the pattern used during tests
func (r *Route) TestPattern() string {
	pattern := r.RoutePattern()

	pattern = strings.Replace(pattern, organizationID, "%s", -1)
	pattern = strings.Replace(pattern, memberID, "%s", -1)
	pattern = strings.Replace(pattern, applicationID, "%s", -1)
	pattern = strings.Replace(pattern, applicationUserID, applicationUserIDPattern, -1)
	pattern = strings.Replace(pattern, applicationUserToID, applicationUserIDPattern, -1)
	pattern = strings.Replace(pattern, connectionType, connectionTypePattern, -1)
	pattern = strings.Replace(pattern, connectionState, connectionStatePattern, -1)
	pattern = strings.Replace(pattern, eventID, eventIDPattern, -1)

	return pattern
}

// Routes defined for this module
var Routes []*Route

// InitRouter initializes the router with this modules routes
func InitRouter(
	router *mux.Router,
	metricHandler RouteMetricHandler,
	mainLogChan, errorLogChan chan *logger.LogMsg,
	env string,
	skipSecurity, debug bool,
) {

	for _, route := range Routes {
		r := router.Methods(route.Method).Path("/" + APIVersion + route.Path)

		r.Name(route.Name).HandlerFunc(
			metricHandler(
				route.Name,
				APIVersion,
				CustomHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug),
			),
		)
	}

	for _, route := range Routes {
		r := router.Methods("OPTIONS").Path("/" + APIVersion + route.Path)

		r.Name(route.Name + "-options").HandlerFunc(
			metricHandler(
				route.Name,
				APIVersion,
				CustomOptionsHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug),
			),
		)
	}
}

// InitHandlers handles the initialization of the route handlers
func InitHandlers() {
	postgresOrganizationHandler = postgres.NewOrganization(postgresOrganization)
	postgresMemberHandler = postgres.NewMember(postgresAccountUser)
	postgresApplicationHandler = postgres.NewApplication(postgresApplication)
	postgresApplicationUserHandler = postgres.NewApplicationUser(postgresApplicationUser, postgresConnection)
	postgresConnectionHandler = postgres.NewConnection(postgresConnection, postgresApplicationUser, postgresEvent)
	postgresEventHandler = postgres.NewEvent(postgresEvent, postgresApplicationUser, postgresConnection)

	redisApplicationHandler = redis.NewApplication(redisApplication, postgresApplication)
}

// VersionHandler returns the current version status
func VersionHandler(ctx *context.Context) []errors.Error {
	resp := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{APIVersion, "alpha"}
	response.WriteResponse(ctx, resp, 200, 86400)
	return nil
}
