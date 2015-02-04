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
			"GET",
			"/",
			"/",
			[]http.HandlerFunc{
				home,
			},
		},
		// Account
		"getAccount": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			[]http.HandlerFunc{
				getAccount,
			},
		},
		"updateAccount": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			[]http.HandlerFunc{
				updateAccount,
			},
		},
		"deleteAccount": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}",
			"/account/%d",
			[]http.HandlerFunc{
				deleteAccount,
			},
		},
		"createAccount": &route{
			"POST",
			"/accounts",
			"/accounts",
			[]http.HandlerFunc{
				createAccount,
			},
		},
		// AccountUser
		"getAccountUser": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			[]http.HandlerFunc{
				getAccountUser,
			},
		},
		"updateAccountUser": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			[]http.HandlerFunc{
				updateAccountUser,
			},
		},
		"deleteAccountUser": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%s",
			[]http.HandlerFunc{
				deleteAccountUser,
			},
		},
		"createAccountUser": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/users",
			"/account/%d/users",
			[]http.HandlerFunc{
				createAccountUser,
			},
		},
		"getAccountUserList": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/users",
			"/account/%d/users",
			[]http.HandlerFunc{
				getAccountUserList,
			},
		},
		// Application
		"getApplication": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			[]http.HandlerFunc{
				getApplication,
			},
		},
		"updateApplication": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			[]http.HandlerFunc{
				updateApplication,
			},
		},
		"deleteApplication": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{appId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			[]http.HandlerFunc{
				deleteApplication,
			},
		},
		"createApplication": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/applications",
			"/account/%d/applications",
			[]http.HandlerFunc{
				createApplication,
			},
		},
		"getApplications": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/applications",
			"/account/%d/applications",
			[]http.HandlerFunc{
				getApplicationList,
			},
		},
		// User
		"getUser": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			[]http.HandlerFunc{
				getUser,
			},
		},
		"updateUser": &route{
			"PUT",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			[]http.HandlerFunc{
				updateUser,
			},
		},
		"deleteUser": &route{
			"DELETE",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/application/%d/user/%s",
			[]http.HandlerFunc{
				deleteUser,
			},
		},
		"createUser": &route{
			"POST",
			"/application/{appId:[0-9]{1,20}}/users",
			"/application/%d/users",
			[]http.HandlerFunc{
				createUser,
			},
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
			[]http.HandlerFunc{
				createConnection,
			},
		},
		"updateConnection": &route{
			"PUT",
			"/application/{appId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			"/application/%d/user/%d/connection/%d",
			[]http.HandlerFunc{
				updateConnection,
			},
		},
		"deleteConnection": &route{
			"DELETE",
			"/application/{appId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			"/application/%d/user/%d/connection/%d",
			[]http.HandlerFunc{
				deleteConnection,
			},
		},
		"getConnectionList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
			"/application/%d/user/%s/connections",
			[]http.HandlerFunc{
				getConnectionList,
			},
		},
		// Event
		"getEvent": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/application/%d/user/%s/event/%d",
			[]http.HandlerFunc{
				getEvent,
			},
		},
		"updateEvent": &route{
			"PUT",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/application/%d/user/%s/event/%d",
			[]http.HandlerFunc{
				updateEvent,
			},
		},
		"deleteEvent": &route{
			"DELETE",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/application/%d/user/%s/event/%d",
			[]http.HandlerFunc{
				deleteEvent,
			},
		},
		"createEvent": &route{
			"POST",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/application/%d/user/%s/events",
			[]http.HandlerFunc{
				createEvent,
			},
		},
		"getEventList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/application/%d/user/%s/events",
			[]http.HandlerFunc{
				getEventList,
			},
		},
		"getConnectionEventList": &route{
			"GET",
			"/application/{appId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections/events",
			"/application/%d/user/%s/connections/events",
			[]http.HandlerFunc{
				getConnectionEventList,
			},
		},
		// Other
		"humans.txt": &route{
			"GET",
			"/humans.txt",
			"/humans.txt",
			[]http.HandlerFunc{
				humans,
			},
		},
		"robots": &route{
			"GET",
			"/robots.txt",
			"/robots.txt",
			[]http.HandlerFunc{
				robots,
			},
		},
	},
}
