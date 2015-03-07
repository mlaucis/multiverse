/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

// Route structure
type (
	routeFunc func(*context)

	route struct {
		method   string
		pattern  string
		cPattern string
		scope    string
		handlers []routeFunc
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
				isRequestExpired,
				getAccount,
			},
		},
		"updateAccount": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/update",
			handlers: []routeFunc{
				isRequestExpired,
				updateAccount,
			},
		},
		"deleteAccount": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/delete",
			handlers: []routeFunc{
				isRequestExpired,
				deleteAccount,
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
		// AccountUser
		"getAccountUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/index",
			handlers: []routeFunc{
				isRequestExpired,
				getAccountUser,
			},
		},
		"updateAccountUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/update",
			handlers: []routeFunc{
				isRequestExpired,
				updateAccountUser,
			},
		},
		"deleteAccountUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/delete",
			handlers: []routeFunc{
				isRequestExpired,
				deleteAccountUser,
			},
		},
		"createAccountUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/create",
			handlers: []routeFunc{
				isRequestExpired,
				createAccountUser,
			},
		},
		"getAccountUserList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/list",
			handlers: []routeFunc{
				isRequestExpired,
				getAccountUserList,
			},
		},
		// Application
		"getApplication": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/index",
			handlers: []routeFunc{
				isRequestExpired,
				getApplication,
			},
		},
		"updateApplication": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/update",
			handlers: []routeFunc{
				isRequestExpired,
				updateApplication,
			},
		},
		"deleteApplication": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/delete",
			handlers: []routeFunc{
				isRequestExpired,
				deleteApplication,
			},
		},
		"createApplication": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "application/create",
			handlers: []routeFunc{
				isRequestExpired,
				createApplication,
			},
		},
		"getApplications": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "account/applications/list",
			handlers: []routeFunc{
				isRequestExpired,
				getApplicationList,
			},
		},
		// User
		"getUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/index",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				getUser,
			},
		},
		"updateUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/update",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				updateUser,
			},
		},
		"deleteUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/delete",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				deleteUser,
			},
		},
		"createUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/application/%d/users",
			scope:    "application/user/create",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				createUser,
			},
		},
		"loginUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/login",
			scope:    "application/user/login",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				loginUser,
			},
		},
		"refreshUserSession": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/refreshsession",
			scope:    "application/user/refreshToken",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				refreshUserSession,
			},
		},
		"logoutUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/logout",
			scope:    "application/user/logout",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				logoutUser,
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
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				createConnection,
			},
		},
		"confirmConnection": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connection/confirm",
			cPattern: "/application/:applicationId/user/:UserID/connection/confirm",
			scope:    "application/user/connection/confirm",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				confirmConnection,
			},
		},
		"updateConnection": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/update",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				updateConnection,
			},
		},
		"deleteConnection": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/delete",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				deleteConnection,
			},
		},
		"getConnectionList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			cPattern: "/account/%d/application/%d/user/%d/connections",
			scope:    "application/user/connections/list",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				getConnectionList,
			},
		},
		// Event
		"getEvent": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/index",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				getEvent,
			},
		},
		"updateEvent": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/update",
			handlers: []routeFunc{
				validateApplicationRequestToken,
				checkSession,
				updateEvent,
			},
		},
		"deleteEvent": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/delete",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				deleteEvent,
			},
		},
		"createEvent": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/event/create",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				createEvent,
			},
		},
		"getEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/events/list",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				getEventList,
			},
		},
		"getConnectionEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections/events",
			cPattern: "/account/%d/application/%d/user/%d/connections/events",
			scope:    "application/user/connection/events",
			handlers: []routeFunc{
				isRequestExpired,
				validateApplicationRequestToken,
				checkSession,
				getConnectionEventList,
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
