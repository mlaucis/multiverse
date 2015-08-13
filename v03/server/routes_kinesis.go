// +build kinesis
// +build !postgres

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
				kinesisAccountHandler.Delete,
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
				kinesisApplicationHandler.Delete,
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentApplicationUser",
			Method: "GET",
			Path:   "/me",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.ReadCurrent,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "updateCurrentApplicationUser",
			Method: "PUT",
			Path:   "/me",
			Handlers: []RouteFunc{
				kinesisApplicationUserHandler.UpdateCurrent,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteCurrentApplicationUser",
			Method: "DELETE",
			Path:   "/me",
			Handlers: []RouteFunc{
				kinesisApplicationUserHandler.DeleteCurrent,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createApplicationUser",
			Method: "POST",
			Path:   "/users",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Create,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "loginApplicationUser",
			Method: "POST",
			Path:   "/me/login",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Login,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "refreshApplicationUserSession",
			Method: "POST",
			Path:   "/me/refresh",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.RefreshSession,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "logoutApplicationUser",
			Method: "DELETE",
			Path:   "/me/logout",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Logout,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
	)

	// UserConnection
	routes = append(routes,
		&Route{
			Name:   "createConnection",
			Method: "PUT",
			Path:   "/me/connections",
			Handlers: []RouteFunc{
				postgresConnectionHandler.Create,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createFriendConnectionAlias",
			Method: "PUT",
			Path:   "/me/friends",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFriend,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createFollowConnectionAlias",
			Method: "PUT",
			Path:   "/me/follow",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFollow,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteConnection",
			Method: "DELETE",
			Path:   fmt.Sprintf("/me/connections/%s", applicationUserToID),
			Handlers: []RouteFunc{
				kinesisConnectionHandler.Delete,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "confirmConnection",
			Method: "POST",
			Path:   fmt.Sprintf("/me/connections/%s/confirm", applicationUserToID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Confirm,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createSocialConnections",
			Method: "POST",
			Path:   "/me/connections/social",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateSocial,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFollows",
			Method: "GET",
			Path:   "/me/follows",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserList,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFollowers",
			Method: "GET",
			Path:   "/me/followers",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFollowedByList,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserFriends",
			Method: "GET",
			Path:   "/me/friends",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFriends,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		})

	// Event
	routes = append(routes,
		&Route{
			Name:   "updateUserEvent",
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s/events/%s", applicationUserID, eventID),
			Handlers: []RouteFunc{
				kinesisEventHandler.Update,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "updateCurrentUserEvent",
			Method: "PUT",
			Path:   fmt.Sprintf("/me/events/%s", eventID),
			Handlers: []RouteFunc{
				kinesisEventHandler.CurrentUserUpdate,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteEvent",
			Method: "DELETE",
			Path:   fmt.Sprintf("/users/%s/events/%s", applicationUserID, eventID),
			Handlers: []RouteFunc{
				kinesisEventHandler.Delete,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteCurrentUserEvent",
			Method: "DELETE",
			Path:   fmt.Sprintf("/me/events/%s", eventID),
			Handlers: []RouteFunc{
				kinesisEventHandler.Delete,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createEvent",
			Method: "POST",
			Path:   fmt.Sprintf("/users/%s/events", applicationUserID),
			Handlers: []RouteFunc{
				kinesisEventHandler.Create,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createCurrentUserEvent",
			Method: "POST",
			Path:   "/me/events",
			Handlers: []RouteFunc{
				kinesisEventHandler.CurrentUserCreate,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserEvent",
			Method: "GET",
			Path:   fmt.Sprintf("/me/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserRead,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserEventList",
			Method: "GET",
			Path:   "/me/events",
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserList,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getFeed",
			Method: "GET",
			Path:   "/me/feed",
			Handlers: []RouteFunc{
				postgresEventHandler.Feed,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUnreadFeedCount",
			Method: "GET",
			Path:   "/me/feed/unread/count",
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeedCount,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUnreadFeed",
			Method: "GET",
			Path:   "/me/feed/unread",
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeed,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "searchEvents",
			Method: "GET",
			Path:   "/events",
			Handlers: []RouteFunc{
				postgresEventHandler.Search,
			},
			Filters: []Filter{
				RateLimitApplication,
				ContextHasApplication(postgresApplicationHandler),
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
