/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "net/http"

// Route structure
type route struct {
	method      string
	pattern     string
	cPattern    string
	handlerFunc http.HandlerFunc
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
			"GET",
			"/",
			"/",
			home,
		},
		// Account
		"getAccount": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			getAccount,
		},
		"updateAccount": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			updateAccount,
		},
		"deleteAccount": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			deleteAccount,
		},
		"createAccount": &route{
			"POST",
			"/accounts",
			"/accounts",
			createAccount,
		},
		// AccountUser
		"getAccountUser": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			getAccountUser,
		},
		"updateAccountUser": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			updateAccountUser,
		},
		"deleteAccountUser": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			deleteAccountUser,
		},
		"createAccountUser": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/users",
			"/account/%d/users",
			createAccountUser,
		},
		"getAccountUserList": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/users",
			"/account/%d/users",
			getAccountUserList,
		},
		// Application
		"getApplication": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			getApplication,
		},
		"updateApplication": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			updateApplication,
		},
		"deleteApplication": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			deleteApplication,
		},
		"createApplication": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/applications",
			"/account/%d/applications",
			createApplication,
		},
		"getApplications": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/applications",
			"/account/%d/applications",
			getApplicationList,
		},
		// User
		"getUser": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			getUser,
		},
		"updateUser": &route{
			"PUT",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			updateUser,
		},
		"deleteUser": &route{
			"DELETE",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			deleteUser,
		},
		"createUser": &route{
			"POST",
			"/application/{appId:[0-9]{1,20}}/users",
			"/application/%d/users",
			createUser,
		},
		/*
			"getUserList": &route{
				"getUserList",
				"GET",
				"/application/{appId:[0-9]{1,20}}/users",
				getUserList,
			},
		*/
		// UserConnection
		"createConnection": &route{
			"POST",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
			"/application/%d/connections",
			createConnection,
		},
		// "updateConnection": &route{
		// 	"PUT",
		// 	"/application/{appId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
		// 	"/application/%d/user/%d/connection/%d",
		// 	updateConnection,
		// },
		// "deleteConnection": &route{
		// 	"DELETE",
		// 	"/application/{appId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
		// 	"/application/%d/user/%d/connection/%d",
		// 	deleteConnection,
		// },
		"getConnectionList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
			"/application/%d/user/%s/connections",
			getConnectionList,
		},
		// Event
		"getEvent": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/application/%d/user/%s/event/%d",
			getEvent,
		},
		// "updateEvent": &route{
		// 	"PUT",
		// 	"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
		// 	"/application/%d/user/%s/event/%d",
		// 	updateEvent,
		// },
		// "deleteEvent": &route{
		// 	"DELETE",
		// 	"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
		// 	"/application/%d/user/%s/event/%d",
		// 	deleteEvent,
		// },
		"createEvent": &route{
			"POST",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/application/%d/user/%s/events",
			createEvent,
		},
		"getEventList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/application/%d/user/%s/events",
			getEventList,
		},
		"getConnectionEventList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections/events",
			"/application/%d/user/%s/connections/events",
			getConnectionEventList,
		},
		// Other
		"humans.txt": &route{
			"GET",
			"/humans.txt",
			"/humans.txt",
			humans,
		},
		"robots": &route{
			"GET",
			"/robots.txt",
			"/robots.txt",
			robots,
		},
	},
}
