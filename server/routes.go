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
		"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
		getApplication,
	},
	"createApplication": &route{
		"createApplication",
		"POST",
		"/account/{accountId:[0-9]{1,20}}/applications",
		createApplication,
	},
	"getApplications": &route{
		"getApplications",
		"GET",
		"/account/{accountId:[0-9]{1,20}}/applications",
		getApplicationList,
	},
	// User
	"getUser": &route{
		"getUser",
		"GET",
		"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
		getUser,
	},
	"createUser": &route{
		"createUser",
		"POST",
		"/application/{appId:[0-9]{1,20}}/users",
		createUser,
	},
	// "getUserList": &route{
	// 	"getUserList",
	// 	"GET",
	// 	"/application/{appId:[0-9]{1,20}}/users",
	// 	getUserList,
	// },
	// UserConnection
	"createConnection": &route{
		"createConnection",
		"POST",
		"/application/{appId:[0-9]{1,20}}/connections",
		createConnection,
	},
	"getConnectionList": &route{
		"getConnectionList",
		"GET",
		"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
		getConnectionList,
	},
	/*
		// Event
		"getApplicationEvent": &route{
			"getApplicationEvent",
			"GET",
			"/application/{appId:[0-9]{1,20}}/event/{eventId:[0-9]{1,20}}",
			getApplicationEvent,
		},
		"createApplicationEvent": &route{
			"createApplicationEvent",
			"POST",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event",
			createApplicationEvent,
		},
		"getApplicationUserEvents": &route{
			"getApplicationUserEvents",
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			getApplicationUserEvents,
		},
		"getUserConnectionsEvents": &route{
			"getUserConnectionsEvents",
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections/events",
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
