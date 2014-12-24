/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "net/http"

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Route definitions of the API
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		home,
	},
	Route{
		"getApplication",
		"GET",
		"/app/{appId}",
		getApplication,
	},
	Route{
		"createApplication",
		"POST",
		"/app",
		createApplication,
	},
	Route{
		"getUser",
		"GET",
		"/app/{appId}/user/{userToken}",
		getUser,
	},
	Route{
		"getUserEvents",
		"GET",
		"/app/{appId}/user/{userToken}/events",
		getUserEvents,
	},
	Route{
		"getUserConnections",
		"GET",
		"/app/{appId}/user/{userToken}/connections",
		getUserConnections,
	},
	Route{
		"getUserConnectionsEvents",
		"GET",
		"/app/{appId}/user/{userToken}/connections/events",
		getUserConnectionsEvents,
	},
	Route{
		"getEvent",
		"GET",
		"/app/{appId}/event/{eventId}",
		getEvent,
	},
	Route{
		"getAccount",
		"GET",
		"/account/{accountId}",
		getAccount,
	},
	Route{
		"getAccountApplications",
		"GET",
		"/account/{accountId}/applications",
		getAccountApplications,
	},
	Route{
		"getAccountUser",
		"GET",
		"/account/{accountId}/user/{userId}",
		getAccountUser,
	},
	Route{
		"getAccountUserList",
		"GET",
		"/account/{accountId}/users",
		getAccountUserList,
	},
	Route{
		"humans.txt",
		"GET",
		"/humans.txt",
		humans,
	},
	Route{
		"robots",
		"GET",
		"/robots.txt",
		robots,
	},
}