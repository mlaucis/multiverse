package server

import (
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/handlers/kinesis"
	"github.com/tapglue/backend/v03/server/handlers/postgres"
	"github.com/tapglue/backend/v03/server/handlers/redis"
	"github.com/tapglue/backend/v03/server/response"

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
	applicationUserID   = "{applicationUserID}"
	applicationUserToID = "{applicationUserToID}"
	eventID             = "{eventID}"
)

var (
	postgresOrganizationHandler, kinesisOrganizationHandler                        handlers.Organization
	postgresMemberHandler, kinesisAccountUserHandler                               handlers.Member
	postgresApplicationHandler, kinesisApplicationHandler, redisApplicationHandler handlers.Application
	postgresApplicationUserHandler, kinesisApplicationUserHandler                  handlers.ApplicationUser
	postgresConnectionHandler, kinesisConnectionHandler                            handlers.Connection
	postgresEventHandler, kinesisEventHandler                                      handlers.Event

	applicationUserIDPattern = "%d"
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
	pattern = strings.Replace(pattern, eventID, eventIDPattern, -1)

	return pattern
}

// Routes defined for this module
var Routes []*Route

// InitRouter initializes the router with this modules routes
func InitRouter(router *mux.Router, mainLogChan, errorLogChan chan *logger.LogMsg, env string, skipSecurity, debug bool) {
	for _, route := range Routes {
		router.
			Methods(route.Method).
			Path("/" + APIVersion + route.Path).
			Name(route.Name).
			HandlerFunc(CustomHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug))
	}

	for _, route := range Routes {
		router.
			Methods("OPTIONS").
			Path("/" + APIVersion + route.Path).
			Name(route.Name + "-options").
			HandlerFunc(CustomOptionsHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug))
	}
}

// InitHandlers handles the initialization of the route handlers
func InitHandlers() {
	kinesisOrganizationHandler = kinesis.NewOrganization(kinesisOrganization, postgresOrganization)
	kinesisAccountUserHandler = kinesis.NewMember(kinesisAccountUser, postgresAccountUser)
	kinesisApplicationHandler = kinesis.NewApplication(kinesisApplication, postgresApplication)
	kinesisApplicationUserHandler = kinesis.NewApplicationUser(kinesisApplicationUser, postgresApplicationUser)
	kinesisConnectionHandler = kinesis.NewConnectionWithApplicationUser(kinesisConnection, postgresConnection, postgresApplicationUser)
	kinesisEventHandler = kinesis.NewEventWithApplicationUser(kinesisEvent, postgresEvent, postgresApplicationUser)

	postgresOrganizationHandler = postgres.NewOrganization(postgresOrganization)
	postgresMemberHandler = postgres.NewMember(postgresAccountUser)
	postgresApplicationHandler = postgres.NewApplication(postgresApplication)
	postgresApplicationUserHandler = postgres.NewApplicationUser(postgresApplicationUser, postgresConnection)
	postgresConnectionHandler = postgres.NewConnection(postgresConnection, postgresApplicationUser)
	postgresEventHandler = postgres.NewEvent(postgresEvent, postgresApplicationUser)

	redisApplicationHandler = redis.NewApplication(redisApplication, postgresApplication)
}

// VersionHandler returns the current version status
func VersionHandler(ctx *context.Context) []errors.Error {
	resp := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{APIVersion, "current"}
	response.WriteResponse(ctx, resp, 200, 86400)
	return nil
}
