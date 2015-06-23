/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"net/http"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/v02/server/handlers"
	"github.com/tapglue/backend/v02/server/handlers/kinesis"
	"github.com/tapglue/backend/v02/server/handlers/postgres"
	"github.com/tapglue/backend/v02/server/response"

	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
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
	accountID           = "{accountID}"
	accountUserID       = "{accountUserID}"
	applicationID       = "{applicationID}"
	applicationUserID   = "{applicationUserID}"
	applicationUserToID = "{applicationUserToID}"
	eventID             = "{eventID}"
)

var (
	postgresAccountHandler, kinesisAccountHandler                 handlers.Account
	postgresAccountUserHandler, kinesisAccountUserHandler         handlers.AccountUser
	postgresApplicationHandler, kinesisApplicationHandler         handlers.Application
	postgresApplicationUserHandler, kinesisApplicationUserHandler handlers.ApplicationUser
	postgresConnectionHandler, kinesisConnectionHandler           handlers.Connection
	postgresEventHandler, kinesisEventHandler                     handlers.Event

	applicationUserIDPattern = "%d"
	eventIDPattern           = "%d"
)

// RoutePattern returns the full route path, inclulding the api version
func (r *Route) RoutePattern() string {
	return "/" + APIVersion + r.Path
}

func ReplaceTestApplicationUserIDPattern(pattern string) {
	applicationUserIDPattern = pattern
}

func ReplaceTestEventIDPattern(pattern string) {
	eventIDPattern = pattern
}

// TestPattern returns the pattern used during tests
func (r *Route) TestPattern() string {
	pattern := r.RoutePattern()

	pattern = strings.Replace(pattern, accountID, "%s", -1)
	pattern = strings.Replace(pattern, accountUserID, "%s", -1)
	pattern = strings.Replace(pattern, applicationID, "%s", -1)
	pattern = strings.Replace(pattern, applicationUserID, applicationUserIDPattern, -1)
	pattern = strings.Replace(pattern, applicationUserToID, applicationUserIDPattern, -1)
	pattern = strings.Replace(pattern, eventID, eventIDPattern, -1)

	return pattern
}

// Routes defined for this module
var Routes []*Route

// InitRouter initializes the router with this modules routes
func InitRouter(agent *gorelic.Agent, router *mux.Router, mainLogChan, errorLogChan chan *logger.LogMsg, env string, skipSecurity, debug bool) {
	for _, route := range Routes {
		router.
			Methods(route.Method).
			Path("/" + APIVersion + route.Path).
			Name(route.Name).
			HandlerFunc(http.HandlerFunc(agent.WrapHTTPHandler(CustomHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug)).ServeHTTP))
	}

	for _, route := range Routes {
		router.
			Methods("OPTIONS").
			Path("/" + APIVersion + route.Path).
			Name(route.Name + "-options").
			HandlerFunc(http.HandlerFunc(agent.WrapHTTPHandler(CustomOptionsHandler(route, mainLogChan, errorLogChan, env, skipSecurity, debug)).ServeHTTP))
	}
}

func InitHandlers() {
	kinesisAccountHandler = kinesis.NewAccount(kinesisAccount, postgresAccount)
	kinesisAccountUserHandler = kinesis.NewAccountUser(kinesisAccountUser, postgresAccountUser)
	kinesisApplicationHandler = kinesis.NewApplication(kinesisApplication, postgresApplication)
	kinesisApplicationUserHandler = kinesis.NewApplicationUser(kinesisApplicationUser, postgresApplicationUser)
	kinesisConnectionHandler = kinesis.NewConnectionWithApplicationUser(kinesisConnection, postgresConnection, postgresApplicationUser)
	kinesisEventHandler = kinesis.NewEventWithApplicationUser(kinesisEvent, postgresEvent, postgresApplicationUser)

	postgresAccountHandler = postgres.NewAccount(postgresAccount)
	postgresAccountUserHandler = postgres.NewAccountUser(postgresAccountUser)
	postgresApplicationHandler = postgres.NewApplication(postgresApplication)
	postgresApplicationUserHandler = postgres.NewApplicationUser(postgresApplicationUser)
	postgresConnectionHandler = postgres.NewConnectionWithApplicationUser(postgresConnection, postgresApplicationUser)
	postgresEventHandler = postgres.NewEventWithApplicationUser(postgresEvent, postgresApplicationUser)
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
