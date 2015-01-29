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
var routes = map[string]*route{
	// General
	"Index": &route{
		"Index",
		"GET",
		"/",
		home,
	},
	// Account
	"getAccount": &route{
		"getAccount",
		"GET",
		"/account/{accountId:[0-9]{1,20}}",
		getAccount,
	},
	"createAccount": &route{
		"createAccount",
		"POST",
		"/accounts",
		createAccount,
	},
	// AccountUser
	"getAccountUser": &route{
		"getAccountUser",
		"GET",
		"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
		getAccountUser,
	},
	"createAccountUser": &route{
		"createAccountUser",
		"POST",
		"/account/{accountId:[0-9]{1,20}}/users",
		createAccountUser,
	},
	/*
		// AccountUser
		"getAccountUser": &route{
			"getAccountUser",
			"GET",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			getAccountUser,
		},
		"createAccountUser": &route{
			"createAccountUser",
			"POST",
			"/account/{accountId:[0-9]{1,20}}/user",
			createAccountUser,
		},
		"getAccountUserList": &route{
			"getAccountUserList",
			"GET",
			"/account/{accountId:[0-9]{1,20}}/users",
			getAccountUserList,
		},
		// Application
		"getApplication": &route{
			"getApplication",
			"GET",
			"/app/{appId:[0-9]{1,20}}",
			getAccountApplication,
		},
		"createAccountApplication": &route{
			"createAccountApplication",
			"POST",
			"/account/{accountId:[0-9]{1,20}}/app",
			createAccountApplication,
		},
		"getAccountApplications": &route{
			"getAccountApplications",
			"GET",
			"/account/{accountId:[0-9]{1,20}}/applications",
			getAccountApplicationList,
		},
		// User
		"getApplicationUser": &route{
			"getApplicationUser",
			"GET",
			"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}",
			getApplicationUser,
		},
		"createApplicationUser": &route{
			"createApplicationUser",
			"POST",
			"/app/{appId:[0-9]{1,20}}/user",
			createApplicationUser,
		},
		"getApplicationUserList": &route{
			"getApplicationUserList",
			"GET",
			"/app/{appId:[0-9]{1,20}}/users",
			getApplicationUserList,
		},
		// UserConnection
		"createUserConnection": &route{
			"createUserConnection",
			"POST",
			"/app/{appId:[0-9]{1,20}}/connection",
			createUserConnection,
		},
		"getUserConnections": &route{
			"getUserConnections",
			"GET",
			"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/connections",
			getUserConnections,
		},
		// Event
		"getApplicationEvent": &route{
			"getApplicationEvent",
			"GET",
			"/app/{appId:[0-9]{1,20}}/event/{eventId:[0-9]{1,20}}",
			getApplicationEvent,
		},
		"createApplicationEvent": &route{
			"createApplicationEvent",
			"POST",
			"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/event",
			createApplicationEvent,
		},
		"getApplicationUserEvents": &route{
			"getApplicationUserEvents",
			"GET",
			"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/events",
			getApplicationUserEvents,
		},
		"getUserConnectionsEvents": &route{
			"getUserConnectionsEvents",
			"GET",
			"/app/{appId:[0-9]{1,20}}/user/{userToken:[a-zA-Z0-9]+}/connections/events",
			getUserConnectionsEvents,
		},
	*/
	// Other
	"humans.txt": &route{
		"humans.txt",
		"GET",
		"/humans.txt",
		humans,
	},
	"robots": &route{
		"robots",
		"GET",
		"/robots.txt",
		robots,
	},
}
