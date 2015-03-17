/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"fmt"

	"github.com/tapglue/backend/context"
)

// Route structure
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
	return "/" + version + r.Pattern
}

// ComposePattern returns the composed pattern for a route
func (r *Route) ComposePattern(version string) string {
	return "/" + version + r.CPattern
}

// GetRoute takes a route name and returns the route including the version
func GetRoute(routeName, apiVersion string) *Route {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version%s\n", routeName, apiVersion))
	}

	return routes[apiVersion][routeName]
}

// Route definitions
var routes = map[string]map[string]*Route{
	"0.1": {
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
		// Account
		"getAccount": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}",
			CPattern: "/account/%d",
			Scope:    "account/index",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccount,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccount,
			},
		},
		/**/
		"updateAccount": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}",
			CPattern: "/account/%d",
			Scope:    "account/update",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateAccount,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccount,
			},
		},
		"deleteAccount": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}",
			CPattern: "/account/%d",
			Scope:    "account/delete",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteAccount,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		"createAccount": &Route{
			Method:   "POST",
			Pattern:  "/accounts",
			CPattern: "/accounts",
			Scope:    "account/create",
			Handlers: []RouteFunc{
				createAccount,
			},
		},
		/**/
		// AccountUser
		"getAccountUser": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/user/%d",
			Scope:    "account/user/index",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccountUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccountUserID,
				contextHasAccountUser,
			},
		},
		"updateAccountUser": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/user/%d",
			Scope:    "account/user/update",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateAccountUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccountUserID,
				contextHasAccountUser,
			},
		},
		"deleteAccountUser": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/user/%d",
			Scope:    "account/user/delete",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteAccountUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccountUserID,
			},
		},
		"createAccountUser": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			CPattern: "/account/%d/users",
			Scope:    "account/user/create",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				createAccountUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		"getAccountUserList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			CPattern: "/account/%d/users",
			Scope:    "account/user/list",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccountUserList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		"loginAccountUser": &Route{
			Method:   "POST",
			Pattern:  "/account/user/login",
			CPattern: "/account/user/login",
			Scope:    "account/user/login",
			Handlers: []RouteFunc{
				loginAccountUser,
			},
		},
		"refreshAccountUserSession": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/user/refreshSession",
			CPattern: "/account/%d/application/%d/user/refreshsession",
			Scope:    "account/user/refreshAccountUserSession",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				refreshAccountUserSession,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		"logoutAccountUser": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
			CPattern: "/account/%d/user/%d/logout",
			Scope:    "account/user/logout",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				logoutAccountUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasAccountUserID,
			},
		},
		// Application
		"getApplication": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d",
			Scope:    "application/index",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getApplication,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplication,
			},
		},
		"updateApplication": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d",
			Scope:    "application/update",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateApplication,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplication,
			},
		},
		"deleteApplication": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d",
			Scope:    "application/delete",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteApplication,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"createApplication": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			CPattern: "/account/%d/applications",
			Scope:    "application/create",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				createApplication,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		"getApplications": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			CPattern: "/account/%d/applications",
			Scope:    "account/applications/list",
			Handlers: []RouteFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getApplicationList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
			},
		},
		// User
		"getUser": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/application/%d/user/%d",
			Scope:    "application/user/index",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		"updateUser": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/application/%d/user/%d",
			Scope:    "application/user/update",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		"deleteUser": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			CPattern: "/account/%d/application/%d/user/%d",
			Scope:    "application/user/delete",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"createUser": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
			CPattern: "/account/%d/application/%d/users",
			Scope:    "application/user/create",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				createApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"loginUser": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			CPattern: "/account/%d/application/%d/user/login",
			Scope:    "application/user/login",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				loginApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"refreshUserSession": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/refreshSession",
			CPattern: "/account/%d/application/%d/user/%d/refreshsession",
			Scope:    "application/user/refreshSession",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				refreshApplicationUserSession,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"logoutUser": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
			CPattern: "/account/%d/application/%d/user/%d/logout",
			Scope:    "application/user/logout",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				logoutApplicationUser,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		/*
			"getUserList": &route{
				"getUserList",
				"GET",
				"/application/{applicationId:[0-9]{1,20}}/users",
				getUserList,
			},
		*/
		// UserConnection
		"createConnection": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			CPattern: "/account/%d/application/%d/user/%d/connections",
			Scope:    "application/user/connection/create",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				createConnection,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"updateConnection": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[a-zA-Z0-9]+}",
			CPattern: "/account/%d/application/%d/user/%d/connection/%d",
			Scope:    "application/user/connection/update",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateConnection,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"deleteConnection": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[a-zA-Z0-9]+}",
			CPattern: "/account/%d/application/%d/user/%d/connection/%d",
			Scope:    "application/user/connection/delete",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteConnection,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getConnectionList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			CPattern: "/account/%d/application/%d/user/%d/connections",
			Scope:    "application/user/connections/list",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getConnectionList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"confirmConnection": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connection/confirm",
			CPattern: "/application/:applicationId/user/:UserID/connection/confirm",
			Scope:    "application/user/connection/confirm",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				confirmConnection,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"createSocialConnections": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{application:[0-9]{1,20}}/user/{userId:[1-9]{1,20}}/connections/social/{platformName:[0-9a-zA-Z]{1,20}}",
			CPattern: "/account/%d/application/%d/user/%d/connections/social/%s",
			Scope:    "application/user/connections/social",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				createSocialConnections,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		// Event
		"getEvent": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d/user/%d/event/%d",
			Scope:    "application/user/event/index",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getEvent,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"updateEvent": &Route{
			Method:   "PUT",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d/user/%d/event/%d",
			Scope:    "application/user/event/update",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateEvent,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"deleteEvent": &Route{
			Method:   "DELETE",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			CPattern: "/account/%d/application/%d/user/%d/event/%d",
			Scope:    "application/user/event/delete",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteEvent,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"createEvent": &Route{
			Method:   "POST",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			CPattern: "/account/%d/application/%d/user/%d/events",
			Scope:    "application/user/event/create",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				createEvent,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getEventList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			CPattern: "/account/%d/application/%d/user/%d/events",
			Scope:    "application/user/events/list",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getEventList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getConnectionEventList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections/events",
			CPattern: "/account/%d/application/%d/user/%d/connections/events",
			Scope:    "application/user/connection/events",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getConnectionEventList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getGeoEventList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/geo/{latitude:[0-9.]+}/{longitude:[0-9.]+}/{radius:[0-9.]+}",
			CPattern: "/account/%d/application/%d/events/geo/%.5f/%.5f/%.5f",
			Scope:    "application/events/geo",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getGeoEventList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"getObjectEventList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/object/{objectKey:[0-9a-zA-Z]+}",
			CPattern: "/account/%d/application/%d/events/object/%s",
			Scope:    "application/events/object",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getObjectEventList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"getLocationEventList": &Route{
			Method:   "GET",
			Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/location/{location:[0-9a-zA-Z]}",
			CPattern: "/account/%d/application/%d/events/location/%s",
			Scope:    "application/events/location",
			Handlers: []RouteFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getLocationEventList,
			},
			Filters: []context.Filter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		// Other
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
	},
}
