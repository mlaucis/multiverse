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

// Route definitions of the API
var routes = []*route{
	&route{
		"getApplication",
		"GET",
		"/app/{appId}",
		getApplication,
	},
	&route{
		"createApplication",
		"POST",
		"/app",
		createApplication,
	},

	&route{
		"getUser",
		"GET",
		"/app/{appId}/user/{userToken}",
		getUser,
	},
	&route{
		"getUserEvents",
		"GET",
		"/app/{appId}/user/{userToken}/events",
		getUserEvents,
	},
	&route{
		"getUserConnections",
		"GET",
		"/app/{appId}/user/{userToken}/connections",
		getUserConnections,
	},
	&route{
		"getUserConnectionsEvents",
		"GET",
		"/app/{appId}/user/{userToken}/connections/events",
		getUserConnectionsEvents,
	},
	&route{
		"getEvent",
		"GET",
		"/app/{appId}/event/{eventId}",
		getEvent,
	},
	&route{
		"getAccount",
		"GET",
		"/account/{accountId}",
		getAccount,
	},
	&route{
		"getAccountApplications",
		"GET",
		"/account/{accountId}/applications",
		getAccountApplications,
	},
	&route{
		"getAccountUser",
		"GET",
		"/account/{accountId}/user/{userId}",
		getAccountUser,
	},
	&route{
		"getAccountUserList",
		"GET",
		"/account/{accountId}/users",
		getAccountUserList,
	},

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

	&route{
		"Index",
		"GET",
		"/",
		home,
	},
}
