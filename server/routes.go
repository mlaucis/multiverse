/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "net/http"

// Route structure
type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

// Route definitions
var routes = []*route{
	// General
	&route{
		"Index",
		"GET",
		"/",
		home,
	},
	// Account
	&route{
		"getAccount",
		"GET",
		"/account/{accountId:[0-9]{1,20}}",
		getAccount,
	},
	&route{
		"createAccount",
		"POST",
		"/account",
		createAccount,
	},
	// AccountUser
	&route{
		"getAccountUser",
		"GET",
		"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
		getAccountUser,
	},
	&route{
		"createAccountUser",
		"POST",
		"/account/{accountId:[0-9]{1,20}}/user",
		createAccountUser,
	},
	&route{
		"getAccountUserList",
		"GET",
		"/account/{accountId:[0-9]{1,20}}/users",
		getAccountUserList,
	},
	// Application
	&route{
		"getApplication",
		"GET",
		"/app/{appId:[0-9]{1,20}}",
		getAccountApplication,
	},
	&route{
		"createAccountApplication",
		"POST",
		"/account/{accountId:[0-9]{1,20}}/app",
		createAccountApplication,
	},
	&route{
		"getAccountApplications",
		"GET",
		"/account/{accountId:[0-9]{1,20}}/applications",
		getAccountApplicationList,
	},
	// User
	&route{
		"getApplicationUser",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}",
		getApplicationUser,
	},
	&route{
		"createApplicationUser",
		"POST",
		"/app/{appId:[0-9]{1,20}}/user",
		createApplicationUser,
	},
	&route{
		"getApplicationUserList",
		"GET",
		"/app/{appId:[0-9]{1,20}}/users",
		getApplicationUserList,
	},
	// UserConnection
	&route{
		"createUserConnection",
		"POST",
		"/app/{appId:[0-9]{1,20}}/connection",
		createUserConnection,
	},
	&route{
		"getUserConnections",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/connections",
		getUserConnections,
	},
	// Session
	&route{
		"getUserSession",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/session/{sessionId:[0-9]{1,20}}",
		getUserSession,
	},
	&route{
		"createUserSession",
		"POST",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/session",
		createUserSession,
	},
	&route{
		"getUserSessionList",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/sessions",
		getUserSessionList,
	},
	// Event
	&route{
		"getApplicationEvent",
		"GET",
		"/app/{appId:[0-9]{1,20}}/event/{eventId:[0-9]{1,20}}",
		getApplicationEvent,
	},
	&route{
		"createApplicationEvent",
		"POST",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/session/{sessionId:[0-9]{1,20}}/event/{eventId:[0-9]{1,20}}",
		createApplicationEvent,
	},
	&route{
		"getApplicationUserEvents",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/events",
		getApplicationUserEvents,
	},
	&route{
		"getSessionEvents",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/session/{sessionId:[0-9]{1,20}}/events",
		getSessionEvents,
	},
	&route{
		"getUserConnectionsEvents",
		"GET",
		"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/connections/events",
		getUserConnectionsEvents,
	},
	// Other
	&route{
		"humans.txt",
		"GET",
		"/humans.txt",
		humans,
	},
	&route{
		"robots",
		"GET",
		"/robots.txt",
		robots,
	},
}
