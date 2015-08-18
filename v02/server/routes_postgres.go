// +build !kinesis
// +build postgres redis

package server

import "fmt"

// SetupRoutes returns the initialized routes
func SetupRoutes() []*Route {
	routes := []*Route{
		// Account
		&Route{
			Name:   "getAccount",
			Method: "GET",
			Path:   fmt.Sprintf("/accounts/%s", accountID),
			Handlers: []RouteFunc{
				postgresAccountHandler.Read,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
			},
		},
		&Route{
			Name:   "updateAccount",
			Method: "PUT",
			Path:   fmt.Sprintf("/accounts/%s", accountID),
			Handlers: []RouteFunc{
				postgresAccountHandler.Update,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
			},
		},
		&Route{
			Name:   "deleteAccount",
			Method: "DELETE",
			Path:   fmt.Sprintf("/accounts/%s", accountID),
			Handlers: []RouteFunc{
				postgresAccountHandler.Delete,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
			},
		},
		&Route{
			Name:   "createAccount",
			Method: "POST",
			Path:   "/accounts",
			Handlers: []RouteFunc{
				postgresAccountHandler.Create,
			},
		},
	}

	// AccountUser
	routes = append(routes,
		&Route{
			Name:   "getAccountUser",
			Method: "GET",
			Path:   fmt.Sprintf("/accounts/%s/users/%s", accountID, accountUserID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Read,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "updateAccountUser",
			Method: "PUT",
			Path:   fmt.Sprintf("/accounts/%s/users/%s", accountID, accountUserID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Update,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "deleteAccountUser",
			Method: "DELETE",
			Path:   fmt.Sprintf("/accounts/%s/users/%s", accountID, accountUserID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Delete,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "createAccountUser",
			Method: "POST",
			Path:   fmt.Sprintf("/accounts/%s/users", accountID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Create,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
			},
		},
		&Route{
			Name:   "getAccountUserList",
			Method: "GET",
			Path:   fmt.Sprintf("/accounts/%s/users", accountID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.List,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "loginAccountUser",
			Method: "POST",
			Path:   "/accounts/users/login",
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Login,
			},
			Filters: []Filter{},
		},
		&Route{
			Name:   "refreshAccountUserSession",
			Method: "POST",
			Path:   fmt.Sprintf("/accounts/%s/users/%s/refresh", accountID, accountUserID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.RefreshSession,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "logoutAccountUser",
			Method: "DELETE",
			Path:   fmt.Sprintf("/accounts/%s/users/%s/logout", accountID, accountUserID),
			Handlers: []RouteFunc{
				postgresAccountUserHandler.Logout,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		})

	// Application
	routes = append(routes,
		&Route{
			Name:   "getApplications",
			Method: "GET",
			Path:   fmt.Sprintf("/accounts/%s/applications", accountID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.List,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		},
		&Route{
			Name:   "getApplication",
			Method: "GET",
			Path:   fmt.Sprintf("/accounts/%s/applications/%s", accountID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Read,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
				ContextHasAccountApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "updateApplication",
			Method: "PUT",
			Path:   fmt.Sprintf("/accounts/%s/applications/%s", accountID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Update,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
				ContextHasAccountApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "deleteApplication",
			Method: "DELETE",
			Path:   fmt.Sprintf("/accounts/%s/applications/%s", accountID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Delete,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
				ContextHasAccountApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "createApplication",
			Method: "POST",
			Path:   fmt.Sprintf("/accounts/%s/applications", accountID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Create,
			},
			Filters: []Filter{
				ContextHasAccount(postgresAccountHandler),
				ContextHasAccountUser(postgresAccountUserHandler),
			},
		})

	// User
	routes = append(routes,
		&Route{
			Name:   "searchApplicationUser",
			Method: "GET",
			Path:   "/users/search",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Search,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getApplicationUser",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s", applicationUserID),
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Read,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentApplicationUser",
			Method: "GET",
			Path:   "/user",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.ReadCurrent,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "updateCurrentApplicationUser",
			Method: "PUT",
			Path:   "/user",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.UpdateCurrent,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteCurrentApplicationUser",
			Method: "DELETE",
			Path:   "/user",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.DeleteCurrent,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createApplicationUser",
			Method: "POST",
			Path:   fmt.Sprintf("/users"),
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Create,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
			},
		},
		&Route{
			Name:   "loginApplicationUser",
			Method: "POST",
			Path:   "/user/login",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Login,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
			},
		},
		&Route{
			Name:   "refreshApplicationUserSession",
			Method: "POST",
			Path:   "/user/refresh",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.RefreshSession,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "logoutApplicationUser",
			Method: "DELETE",
			Path:   "/user/logout",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Logout,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
	)

	// UserConnection
	routes = append(routes,
		&Route{
			Name:   "createConnection",
			Method: "POST",
			Path:   "/user/connections",
			Handlers: []RouteFunc{
				postgresConnectionHandler.Create,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteConnection",
			Method: "DELETE",
			Path:   fmt.Sprintf("/user/connections/%s", applicationUserToID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Delete,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		// TODO NOT DOCUMENTED
		&Route{
			Name:   "confirmConnection",
			Method: "POST",
			Path:   fmt.Sprintf("/user/connections/%s/confirm", applicationUserToID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Confirm,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createSocialConnections",
			Method: "POST",
			Path:   "/user/connections/social",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateSocial,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUserFollows",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/follows", applicationUserID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.List,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFollows",
			Method: "GET",
			Path:   "/user/follows",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserList,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUserFollowers",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/followers", applicationUserID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.FollowedByList,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFollowers",
			Method: "GET",
			Path:   "/user/followers",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFollowedByList,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUserFriends",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/friends", applicationUserID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Friends,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFriends",
			Method: "GET",
			Path:   "/user/friends",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFriends,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		})

	// Event
	routes = append(routes,
		// TODO NOT DOCUMENTED, NOT IMPLEMENTED
		&Route{
			Name:   "updateUserEvent",
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s/events/%s", applicationUserID, eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.Update,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "updateCurrentUserEvent",
			Method: "PUT",
			Path:   fmt.Sprintf("/user/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserUpdate,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		// TODO NOT DOCUMENTED, NOT IMPLEMENTED
		&Route{
			Name:   "deleteEvent",
			Method: "DELETE",
			Path:   fmt.Sprintf("/users/%s/events/%s", applicationUserID, eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.Delete,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteCurrentUserEvent",
			Method: "DELETE",
			Path:   fmt.Sprintf("/user/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.Delete,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},

		// TODO NOT DOCUMENTED, NOT IMPLEMENTED
		&Route{
			Name:   "createEvent",
			Method: "POST",
			Path:   fmt.Sprintf("/users/%s/events", applicationUserID),
			Handlers: []RouteFunc{
				postgresEventHandler.Create,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createCurrentUserEvent",
			Method: "POST",
			Path:   "/user/events",
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserCreate,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getEvent",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/events/%s", applicationUserID, eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.Read,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		// TODO NOT DOCUMENTED, NOT IMPLEMENTED
		&Route{
			Name:   "getCurrentUserEvent",
			Method: "GET",
			Path:   fmt.Sprintf("/user/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserRead,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getEventList",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/events", applicationUserID),
			Handlers: []RouteFunc{
				postgresEventHandler.List,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserEventList",
			Method: "GET",
			Path:   "/user/events",
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserList,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getFeed",
			Method: "GET",
			Path:   "/user/feed",
			Handlers: []RouteFunc{
				postgresEventHandler.Feed,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUnreadFeedCount",
			Method: "GET",
			Path:   "/user/feed/unread/count",
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeedCount,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUnreadFeed",
			Method: "GET",
			Path:   "/user/feed/unread",
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeed,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		// TODO NOT DOCUMENTED, NOT IMPLEMENTED
		&Route{
			Name:   "searchEvents",
			Method: "GET",
			Path:   "/events",
			Handlers: []RouteFunc{
				postgresEventHandler.Search,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "version",
			Method: "GET",
			Path:   "/version",
			Handlers: []RouteFunc{
				VersionHandler,
			},
		})

	return routes
}
