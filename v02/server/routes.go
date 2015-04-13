/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "github.com/tapglue/backend/context"

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
			GetAccount,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccount,
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
			UpdateAccount,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccount,
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
			DeleteAccount,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
		},
	},
	"createAccount": &Route{
		Method:   "POST",
		Pattern:  "/accounts",
		CPattern: "/accounts",
		Scope:    "account/create",
		Handlers: []RouteFunc{
			CreateAccount,
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
			GetAccountUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccountUserID,
			ContextHasAccountUser,
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
			UpdateAccountUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccountUserID,
			ContextHasAccountUser,
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
			DeleteAccountUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccountUserID,
		},
	},
	"createAccountUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/users",
		CPattern: "/account/%d/users",
		Scope:    "account/user/create",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CreateAccountUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
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
			GetAccountUserList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
		},
	},
	"loginAccountUser": &Route{
		Method:   "POST",
		Pattern:  "/account/user/login",
		CPattern: "/account/user/login",
		Scope:    "account/user/login",
		Handlers: []RouteFunc{
			LoginAccountUser,
		},
	},
	"refreshAccountUserSession": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/user/{userId:[0-9]{1,20}}/refreshSession",
		CPattern: "/account/%d/user/%d/refreshSession",
		Scope:    "account/user/refreshAccountUserSession",
		Handlers: []RouteFunc{
			ValidateAccountRequestToken,
			CheckAccountSession,
			RefreshAccountUserSession,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccountUserID,
			ContextHasAccountUser,
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
			LogoutAccountUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasAccountUserID,
			ContextHasAccountUser,
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
			GetApplication,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplication,
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
			UpdateApplication,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplication,
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
			DeleteApplication,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
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
			CreateApplication,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
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
			GetApplicationList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
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
			GetApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
			ContextHasApplicationUser,
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
			UpdateApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
			ContextHasApplicationUser,
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
			DeleteApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
		},
	},
	"createUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/users",
		CPattern: "/account/%d/application/%d/users",
		Scope:    "application/user/create",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CreateApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
		},
	},
	"loginUser": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/login",
		CPattern: "/account/%d/application/%d/user/login",
		Scope:    "application/user/login",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			LoginApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
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
			RefreshApplicationUserSession,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
			ContextHasApplicationUser,
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
			LogoutApplicationUser,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
			ContextHasApplicationUser,
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
			CreateConnection,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			UpdateConnection,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			DeleteConnection,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			GetConnectionList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			GetFollowedByUsersList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
		},
	},
	"confirmConnection": &Route{
		Method:   "POST",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/user/{userId:[0-9]+}/connection/confirm",
		CPattern: "/account/%d/application/%d/user/%d/connection/confirm",
		Scope:    "application/user/connection/confirm",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			CheckApplicationSession,
			ConfirmConnection,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			CreateSocialConnections,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
			ContextHasApplicationUser,
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
			GetEvent,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			UpdateEvent,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			DeleteEvent,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			CreateEvent,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			GetEventList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
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
			GetConnectionEventList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
			ContextHasApplicationUserID,
		},
	},
	"getGeoEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/geo/{latitude:[0-9.\\-]+}/{longitude:[0-9.\\-]+}/{radius:[0-9.\\-]+}",
		CPattern: "/account/%d/application/%d/events/geo/%.7f/%.7f/%.7f",
		Scope:    "application/events/geo",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			GetGeoEventList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
		},
	},
	"getObjectEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/object/{objectKey:[0-9a-zA-Z\\-]+}",
		CPattern: "/account/%d/application/%d/events/object/%s",
		Scope:    "application/events/object",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			GetObjectEventList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
		},
	},
	"getLocationEventList": &Route{
		Method:   "GET",
		Pattern:  "/account/{accountId:[0-9]{1,20}}/application/{applicationId:[0-9]{1,20}}/events/location/{location:[0-9a-zA-Z\\-]+}",
		CPattern: "/account/%d/application/%d/events/location/%s",
		Scope:    "application/events/location",
		Handlers: []RouteFunc{
			ValidateApplicationRequestToken,
			GetLocationEventList,
		},
		Filters: []context.Filter{
			ContextHasAccountID,
			ContextHasApplicationID,
		},
	},
	"version": &Route{
		Method:   "GET",
		Pattern:  "/version",
		CPattern: "/version",
		Scope:    "version",
		Handlers: []RouteFunc{
			VersionHandler,
		},
	},
}
