/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "net/http"

// Route structure
type route struct {
	method   string
	pattern  string
	cPattern string
	scope    string
	handlers []http.HandlerFunc
}

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
			handlers: []http.HandlerFunc{
				home,
			},
		},
		// Account
		"getAccount": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/index",
			handlers: []http.HandlerFunc{
				getAccount,
			},
		},
		"updateAccount": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/update",
			handlers: []http.HandlerFunc{
				updateAccount,
			},
		},
		"deleteAccount": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}",
			cPattern: "/account/%d",
			scope:    "account/delete",
			handlers: []http.HandlerFunc{
				deleteAccount,
			},
		},
		"createAccount": &route{
			method:   "POST",
			pattern:  "/accounts",
			cPattern: "/accounts",
			scope:    "account/create",
			handlers: []http.HandlerFunc{
				createAccount,
			},
		},
		// AccountUser
		"getAccountUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/index",
			handlers: []http.HandlerFunc{
				getAccountUser,
			},
		},
		"updateAccountUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/update",
			handlers: []http.HandlerFunc{
				updateAccountUser,
			},
		},
		"deleteAccountUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/user/%d",
			scope:    "account/user/delete",
			handlers: []http.HandlerFunc{
				deleteAccountUser,
			},
		},
		"createAccountUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/create",
			handlers: []http.HandlerFunc{
				createAccountUser,
			},
		},
		"getAccountUserList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/users",
			scope:    "account/user/list",
			handlers: []http.HandlerFunc{
				getAccountUserList,
			},
		},
		// Application
		"getApplication": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/index",
			handlers: []http.HandlerFunc{
				getApplication,
			},
		},
		"updateApplication": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/update",
			handlers: []http.HandlerFunc{
				updateApplication,
			},
		},
		"deleteApplication": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d",
			scope:    "application/delete",
			handlers: []http.HandlerFunc{
				deleteApplication,
			},
		},
		"createApplication": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "application/create",
			handlers: []http.HandlerFunc{
				createApplication,
			},
		},
		"getApplications": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
			cPattern: "/account/%d/applications",
			scope:    "account/applications/list",
			handlers: []http.HandlerFunc{
				getApplicationList,
			},
		},
		// User
		"getUser": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/index",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/index", "0.1", w, r)
				},
				isSessionValid,
				getUser,
			},
		},
		"updateUser": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/update",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/update", "0.1", w, r)
				},
				isSessionValid,
				updateUser,
			},
		},
		"deleteUser": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d",
			scope:    "application/user/delete",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/delete", "0.1", w, r)
				},
				isSessionValid,
				deleteUser,
			},
		},
		"createUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
			cPattern: "/account/%d/application/%d/users",
			scope:    "application/user/create",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/create", "0.1", w, r)
				},
				createUser,
			},
		},
		"loginUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/login",
			scope:    "application/user/login",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/login", "0.1", w, r)
				},
				loginUser,
			},
		},
		"refreshUserSession": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/refreshsession",
			scope:    "application/user/refreshToken",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/refreshToken", "0.1", w, r)
				},
				isSessionValid,
				refreshUserSession,
			},
		},
		"logoutUser": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
			cPattern: "/account/%d/application/%d/user/logout",
			scope:    "application/user/logout",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/logout", "0.1", w, r)
				},
				isSessionValid,
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
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/connection/create", "0.1", w, r)
				},
				isSessionValid,
				createConnection,
			},
		},
		"updateConnection": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/update",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/connection/update", "0.1", w, r)
				},
				isSessionValid,
				updateConnection,
			},
		},
		"deleteConnection": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			cPattern: "/account/%d/application/%d/user/%d/connection/%d",
			scope:    "application/user/connection/delete",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/connection/delete", "0.1", w, r)
				},
				isSessionValid,
				deleteConnection,
			},
		},
		"getConnectionList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
			cPattern: "/account/%d/application/%d/user/%d/connections",
			scope:    "application/user/connections/list",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/connections/list", "0.1", w, r)
				},
				isSessionValid,
				getConnectionList,
			},
		},
		// Event
		"getEvent": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/index",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/event/index", "0.1", w, r)
				},
				isSessionValid,
				getEvent,
			},
		},
		"updateEvent": &route{
			method:   "PUT",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/update",
			handlers: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/event/update", "0.1", w, r)
				},
				isSessionValid,
				updateEvent,
			},
		},
		"deleteEvent": &route{
			method:   "DELETE",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
			cPattern: "/account/%d/application/%d/user/%d/event/%d",
			scope:    "application/user/event/delete",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/event/delete", "0.1", w, r)
				},
				isSessionValid,
				deleteEvent,
			},
		},
		"createEvent": &route{
			method:   "POST",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/event/create",
			handlers: []http.HandlerFunc{
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/event/create", "0.1", w, r)
				},
				isSessionValid,
				createEvent,
			},
		},
		"getEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
			cPattern: "/account/%d/application/%d/user/%d/events",
			scope:    "application/user/events/list",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/events/list", "0.1", w, r)
				},
				isSessionValid,
				getEventList,
			},
		},
		"getConnectionEventList": &route{
			method:   "GET",
			pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections/events",
			cPattern: "/account/%d/application/%d/user/%d/connections/events",
			scope:    "application/user/connection/events",
			handlers: []http.HandlerFunc{
				isRequestExpired,
				func(w http.ResponseWriter, r *http.Request) {
					validateApplicationRequestToken("application/user/connection/events", "0.1", w, r)
				},
				isSessionValid,
				getConnectionEventList,
			},
		},
		// Other
		"humans": &route{
			method:   "GET",
			pattern:  "/humans.txt",
			cPattern: "/humans.txt",
			scope:    "humans",
			handlers: []http.HandlerFunc{
				humans,
			},
		},
		"robots": &route{
			method:   "GET",
			pattern:  "/robots.txt",
			cPattern: "/robots.txt",
			scope:    "robots",
			handlers: []http.HandlerFunc{
				robots,
			},
		},
	},
}
