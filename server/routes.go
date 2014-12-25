/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package server holds all the server related logic
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
		"/account/{accountId}",
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
		"/account/{accountId}/user/{userId}",
		getAccountUser,
	},
	&route{
		"createAccountUser",
		"POST",
		"/account/{accountId}/user",
		createAccountUser,
	},
	&route{
		"getAccountUserList",
		"GET",
		"/account/{accountId}/users",
		getAccountUserList,
	},
	// Application
	&route{
		"getApplication",
		"GET",
		"/app/{appId}",
		getAccountApplication,
	},
	&route{
		"createAccountApplication",
		"POST",
		"/account/{accountId}/app",
		createAccountApplication,
	},
	&route{
		"getAccountApplications",
		"GET",
		"/account/{accountId}/applications",
		getAccountApplicationList,
	},
	// User
	&route{
		"getApplicationUser",
		"GET",
		"/app/{appId}/user/{userToken}",
		getApplicationUser,
	},
	&route{
		"createApplicationUser",
		"POST",
		"/app/{appId}/user",
		createApplicationUser,
	},
	&route{
		"getApplicationUserList",
		"GET",
		"/app/{appId}/users",
		getApplicationUserList,
	},
	// UserConnection
	&route{
		"createUserConnection",
		"POST",
		"/app/{appId}/connection",
		createUserConnection,
	},
	&route{
		"getUserConnections",
		"GET",
		"/app/{appId}/user/{userToken}/connections",
		getUserConnections,
	},
	// Session
	&route{
		"getUserSession",
		"GET",
		"/app/{appId}/user/{userToken}/session/{sessionId}",
		getUserSession,
	},
	&route{
		"createUserSession",
		"POST",
		"/app/{appId}/user/{userToken}/session",
		createUserSession,
	},
	&route{
		"getUserSessionList",
		"GET",
		"/app/{appId}/user/{userToken}/sessions",
		getUserSessionList,
	},
	// Event
	&route{
		"getApplicationEvent",
		"GET",
		"/app/{appId}/event/{eventId}",
		getApplicationEvent,
	},
	&route{
		"createApplicationEvent",
		"POST",
		"/app/{appId}/user/{userToken}/session/{sessionId}/event/{eventId}",
		createApplicationEvent,
	},
	&route{
		"getApplicationUserEvents",
		"GET",
		"/app/{appId}/user/{userToken}/events",
		getApplicationUserEvents,
	},
	&route{
		"getSessionEvents",
		"GET",
		"/app/{appId}/user/{userToken}/session/{sessionId}/events",
		getSessionEvents,
	},
	&route{
		"getUserConnectionsEvents",
		"GET",
		"/app/{appId}/user/{userToken}/connections/events",
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