/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	. "github.com/tapglue/backend/server/utils"
)

// Route definitions
var Routes = map[string]*Route{
	// Account
	"getAccount": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}",
		CPattern: "/account/%d",
		Scope:    "account/index",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			getAccount,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccount,
		},
	},
	/**/
	"updateAccount": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}",
		CPattern: "/account/%d",
		Scope:    "account/update",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			updateAccount,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccount,
		},
	},
	"deleteAccount": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}",
		CPattern: "/account/%d",
		Scope:    "account/delete",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			deleteAccount,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	"createAccount": &Route{
		Method:   "POST",
		Pattern:  "/accounts",
		CPattern: "/accounts",
		Scope:    "account/create",
		Handlers: []RouteFunc{
			createAccount,
		},
	},
	/**/
	// AccountUser
	"getAccountUser": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/user/%d",
		Scope:    "account/user/index",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			getAccountUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccountUserID,
			contextHasAccountUser,
		},
	},
	"updateAccountUser": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/user/%d",
		Scope:    "account/user/update",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			updateAccountUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccountUserID,
			contextHasAccountUser,
		},
	},
	"deleteAccountUser": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/user/%d",
		Scope:    "account/user/delete",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			deleteAccountUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccountUserID,
		},
	},
	"createAccountUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/users",
		CPattern: "/account/%d/users",
		Scope:    "account/user/create",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			createAccountUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	"getAccountUserList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/users",
		CPattern: "/account/%d/users",
		Scope:    "account/user/list",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			getAccountUserList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	"loginAccountUser": &Route{
		Method:   "POST",
		Pattern:  "/account/user/login",
		CPattern: "/account/user/login",
		Scope:    "account/user/login",
		Handlers: []RouteFunc{
			loginAccountUser,
		},
	},
	"refreshAccountUserSession": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/refreshSession",
		CPattern: "/account/%d/application/%d/user/refreshsession",
		Scope:    "account/user/refreshAccountUserSession",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			refreshAccountUserSession,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	"logoutAccountUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
		CPattern: "/account/%d/user/%d/logout",
		Scope:    "account/user/logout",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			logoutAccountUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasAccountUserID,
		},
	},
	// Application
	"getApplication": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d",
		Scope:    "application/index",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			getApplication,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplication,
		},
	},
	"updateApplication": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d",
		Scope:    "application/update",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			updateApplication,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplication,
		},
	},
	"deleteApplication": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d",
		Scope:    "application/delete",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			deleteApplication,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
	"createApplication": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
		CPattern: "/account/%d/applications",
		Scope:    "application/create",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			createApplication,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	"getApplications": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/applications",
		CPattern: "/account/%d/applications",
		Scope:    "account/applications/list",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			getApplicationList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
		},
	},
	// User
	"getUser": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/application/%d/user/%d",
		Scope:    "application/user/index",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
			contextHasApplicationUser,
		},
	},
	"updateUser": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/application/%d/user/%d",
		Scope:    "application/user/update",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			updateApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
			contextHasApplicationUser,
		},
	},
	"deleteUser": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}",
		CPattern: "/account/%d/application/%d/user/%d",
		Scope:    "application/user/delete",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			deleteApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"createUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
		CPattern: "/account/%d/application/%d/users",
		Scope:    "application/user/create",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			createApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
	"loginUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
		CPattern: "/account/%d/application/%d/user/login",
		Scope:    "application/user/login",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			loginApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
	"refreshUserSession": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/refreshSession",
		CPattern: "/account/%d/application/%d/user/%d/refreshSession",
		Scope:    "application/user/refreshSession",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			refreshApplicationUserSession,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
			contextHasApplicationUser,
		},
	},
	"logoutUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/logout",
		CPattern: "/account/%d/application/%d/user/%d/logout",
		Scope:    "application/user/logout",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			logoutApplicationUser,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
			contextHasApplicationUser,
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
	"createConnection": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
		CPattern: "/account/%d/application/%d/user/%d/connections",
		Scope:    "application/user/connection/create",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			createConnection,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"updateConnection": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/connection/%d",
		Scope:    "application/user/connection/update",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			updateConnection,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"deleteConnection": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/connection/{userToId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/connection/%d",
		Scope:    "application/user/connection/delete",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			deleteConnection,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"getConnectionList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections",
		CPattern: "/account/%d/application/%d/user/%d/connections",
		Scope:    "application/user/connections/list",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getConnectionList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"getFollowerList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/followers",
		CPattern: "/account/%d/application/%d/user/%d/followers",
		Scope:    "application/user/followers/list",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getFollowedByUsersList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"confirmConnection": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connection/confirm",
		CPattern: "/application/:applicationId/user/:UserID/connection/confirm",
		Scope:    "application/user/connection/confirm",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			confirmConnection,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"createSocialConnections": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[1-9]{1,20}}/connections/social/{platformName:[0-9a-zA-Z]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/connections/social/%s",
		Scope:    "application/user/connections/social",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			createSocialConnections,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
			contextHasApplicationUser,
		},
	},
	// Event
	"getEvent": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/event/%d",
		Scope:    "application/user/event/index",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getEvent,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"updateEvent": &Route{
		Method:   "PUT",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/event/%d",
		Scope:    "application/user/event/update",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			updateEvent,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"deleteEvent": &Route{
		Method:   "DELETE",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/event/{eventId:[0-9]{1,20}}",
		CPattern: "/account/%d/application/%d/user/%d/event/%d",
		Scope:    "application/user/event/delete",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			deleteEvent,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"createEvent": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
		CPattern: "/account/%d/application/%d/user/%d/events",
		Scope:    "application/user/event/create",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			createEvent,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"getEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/events",
		CPattern: "/account/%d/application/%d/user/%d/events",
		Scope:    "application/user/events/list",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getEventList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"getConnectionEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connections/events",
		CPattern: "/account/%d/application/%d/user/%d/connections/events",
		Scope:    "application/user/connection/events",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			getConnectionEventList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
			contextHasApplicationUserID,
		},
	},
	"getGeoEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/geo/{latitude:[0-9.]+}/{longitude:[0-9.]+}/{radius:[0-9.]+}",
		CPattern: "/account/%d/application/%d/events/geo/%.7f/%.7f/%.7f",
		Scope:    "application/events/geo",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			getGeoEventList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
	"getObjectEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/object/{objectKey:[0-9a-zA-Z]+}",
		CPattern: "/account/%d/application/%d/events/object/%s",
		Scope:    "application/events/object",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			getObjectEventList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
	"getLocationEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/location/{location:[0-9a-zA-Z\\-]+}",
		CPattern: "/account/%d/application/%d/events/location/%s",
		Scope:    "application/events/location",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			getLocationEventList,
		},
		Filters: []context.Filter{
			contextHasAccountID,
			contextHasApplicationID,
		},
	},
}
