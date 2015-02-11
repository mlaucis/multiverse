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
			"/account/%d/user/%d",
			[]http.HandlerFunc{
				getAccountUser,
			},
		},
		"updateAccountUser": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%d",
			[]http.HandlerFunc{
				updateAccountUser,
			},
		},
		"deleteAccountUser": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/user/%d",
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
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			[]http.HandlerFunc{
				getApplication,
			},
		},
		"updateApplication": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
			"/account/%d/application/%d",
			[]http.HandlerFunc{
				updateApplication,
			},
		},
		"deleteApplication": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
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
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/application/%d/user/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				getUser,
			},
		},
		"updateUser": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/application/%d/user/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				updateUser,
			},
		},
		"deleteUser": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}",
			"/account/%d/application/%d/user/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				deleteUser,
			},
		},
		"createUser": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
			"/account/%d/application/%d/users",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				createUser,
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
			"POST",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
			"/account/%d/application/%d/user/%d/connections",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				createConnection,
			},
		},
		"updateConnection": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			"/account/%d/application/%d/user/%d/connection/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				updateConnection,
			},
		},
		"deleteConnection": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userFromId:[a-zA-Z0-9]+}/connection/{userToId:[a-zA-Z0-9]+}",
			"/account/%d/application/%d/user/%d/connection/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				deleteConnection,
			},
		},
		"getConnectionList": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections",
			"/account/%d/application/%d/user/%d/connections",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				getConnectionList,
			},
		},
		// Event
		"getEvent": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/account/%d/application/%d/user/%d/event/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				getEvent,
			},
		},
		"updateEvent": &route{
			"PUT",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/account/%d/application/%d/user/%d/event/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				updateEvent,
			},
		},
		"deleteEvent": &route{
			"DELETE",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/event/{eventId:[0-9]{1,20}}",
			"/account/%d/application/%d/user/%d/event/%d",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				deleteEvent,
			},
		},
		"createEvent": &route{
			"POST",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/account/%d/application/%d/user/%d/events",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				createEvent,
			},
		},
		"getEventList": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/events",
			"/account/%d/application/%d/user/%d/events",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				getEventList,
			},
		},
		"getConnectionEventList": &route{
			"GET",
			"/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[a-zA-Z0-9]+}/connections/events",
			"/account/%d/application/%d/user/%d/connections/events",
			[]http.HandlerFunc{
				validateApplicationRequestToken,
				getConnectionEventList,
			},
		},
		// Other
		"humans": &route{
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
