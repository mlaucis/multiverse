/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "github.com/tapglue/backend/context"

// Route structure
type (
	routeFunc func(*context.Context)

	route struct {
		method         string
		pattern        string
		cPattern       string
		scope          string
		handlers       []routeFunc
		contextFilters []context.ContextFilter
	}
)

func (r *route) routePattern(version string) string {
	return "/" + version + r.pattern
}

func (r *route) composePattern(version string) string {
	return "/" + version + r.cPattern
}

// Route definitions
var routes = map[string]map[string]*route{
	"0.1": {
		// General
		"index": &route{
			method:   "GET",
			pattern:  "/",
			cPattern: "/",
			scope:    "/",
			handlers: []routeFunc{
				home,
			},
		},
		// Account
		"getAccount": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/index",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccount,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccount,
			},
		},
		/**/
		"updateAccount": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/update",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateAccount,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccount,
			},
		},
		"deleteAccount": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/delete",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteAccount,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		"createAccount": &route{
			method:   "POST",
			pattern:  "/accounts",
			cPattern: "/accounts",
			scope:    "account/create",
			handlers: []routeFunc{
				createAccount,
			},
		},
		/**/
		// AccountUser
		"getAccountUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/index",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccountUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccountUserID,
				contextHasAccountUser,
			},
		},
		"updateAccountUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/update",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateAccountUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccountUserID,
				contextHasAccountUser,
			},
		},
		"deleteAccountUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/delete",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteAccountUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccountUserID,
			},
		},
		"createAccountUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/create",
			handlers: []routeFunc{
				validateAccountRequestToken,
				createAccountUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		"getAccountUserList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/list",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getAccountUserList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		"loginAccountUser": &route{
			method:   "POST",
			pattern:  "/account/user/login",
			cPattern: "/account/user/login",
			scope:    "account/user/login",
			handlers: []routeFunc{
				loginAccountUser,
			},
		},
		"refreshAccountUserSession": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/refreshSession",
			cPattern: "/account/%d/application/%d/user/refreshsession",
			scope:    "account/user/refreshAccountUserSession",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				refreshAccountUserSession,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		"logoutAccountUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
			cPattern: "/account/%d/user/%d/logout",
			scope:    "account/user/logout",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				logoutAccountUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasAccountUserID,
			},
		},
		// Application
		"getApplication": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/index",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getApplication,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplication,
			},
		},
		"updateApplication": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/update",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				updateApplication,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplication,
			},
		},
		"deleteApplication": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/delete",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				deleteApplication,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"createApplication": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "application/create",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				createApplication,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		"getApplications": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "account/applications/list",
			handlers: []routeFunc{
				validateAccountRequestToken,
				checkAccountSession,
				getApplicationList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
			},
		},
		// User
		"getUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/index",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getApplicationUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		"updateUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/update",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateApplicationUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
				contextHasApplicationUser,
			},
		},
		"deleteUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/delete",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteApplicationUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"createUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/application/%d/users",
			scope:    "application/user/create",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				createApplicationUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"loginUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/login",
			scope:    "application/user/login",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				loginApplicationUser,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"refreshUserSession": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/refreshSession",
			cPattern: "/account/%d/application/%d/user/%d/refreshsession",
			scope:    "application/user/refreshSession",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				refreshApplicationUserSession,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"logoutUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
			cPattern: "/account/%d/application/%d/user/%d/logout",
			scope:    "application/user/logout",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				logoutApplicationUser,
			},
			contextFilters: []context.ContextFilter{
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
		"createConnection": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			cPattern: "/account/%d/application/%d/user/%d/connections",
			scope:    "application/user/connection/create",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				createConnection,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"updateConnection": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/update",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateConnection,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"deleteConnection": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/delete",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteConnection,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getConnectionList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			cPattern: "/account/%d/application/%d/user/%d/connections",
			scope:    "application/user/connections/list",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getConnectionList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"confirmConnection": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connection/confirm",
			cPattern: "/application/:applicationId/user/:UserID/connection/confirm",
			scope:    "application/user/connection/confirm",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				confirmConnection,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		// Event
		"getEvent": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/index",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getEvent,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"updateEvent": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/update",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				updateEvent,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"deleteEvent": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/delete",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				deleteEvent,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"createEvent": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/event/create",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				createEvent,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/events/list",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getEventList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getConnectionEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections/events",
			cPattern: "/account/%d/application/%d/user/%d/connections/events",
			scope:    "application/user/connection/events",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getConnectionEventList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
				contextHasApplicationUserID,
			},
		},
		"getGeoEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/geo/{latitude:[0-9.]+}/{longitude:[0-9.]+}/{radius:[0-9.]+}",
			cPattern: "/account/%d/application/%d/events/geo/%.5f/%.5f/%.5f",
			scope:    "application/events/geo",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getGeoEventList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"getObjectEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/object/{objectKey:[0-9a-zA-Z]+}",
			cPattern: "/account/%d/application/%d/events/object/%s",
			scope:    "application/events/object",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getObjectEventList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		"getLocationEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/location/{location:[0-9a-zA-Z]}",
			cPattern: "/account/%d/application/%d/events/location/%s",
			scope:    "application/events/location",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkApplicationSession,
				getLocationEventList,
			},
			contextFilters: []context.ContextFilter{
				contextHasAccountID,
				contextHasApplicationID,
			},
		},
		// Other
		"humans": &route{
			method:   "GET",
			pattern:  "/humans.txt",
			cPattern: "/humans.txt",
			scope:    "humans",
			handlers: []routeFunc{
				humans,
			},
		},
		"robots": &route{
			method:   "GET",
			pattern:  "/robots.txt",
			cPattern: "/robots.txt",
			scope:    "robots",
			handlers: []routeFunc{
				robots,
			},
		},
	},
}
