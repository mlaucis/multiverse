// +build postgres bench

package server

import "fmt"

// SetupRoutes returns the initialized routes
func SetupRoutes() []*Route {
	routes := []*Route{
		// Organization
		&Route{
			Name:   "getOrganization",
			Method: "GET",
			Path:   fmt.Sprintf("/organizations/%s", organizationID),
			Handlers: []RouteFunc{
				postgresOrganizationHandler.Read,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
			},
		},
		&Route{
			Name:   "updateOrganization",
			Method: "PUT",
			Path:   fmt.Sprintf("/organizations/%s", organizationID),
			Handlers: []RouteFunc{
				postgresOrganizationHandler.Update,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
			},
		},
		&Route{
			Name:   "deleteOrganization",
			Method: "DELETE",
			Path:   fmt.Sprintf("/organizations/%s", organizationID),
			Handlers: []RouteFunc{
				postgresOrganizationHandler.Delete,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
			},
		},
		&Route{
			Name:   "createOrganization",
			Method: "POST",
			Path:   "/organizations",
			Handlers: []RouteFunc{
				postgresOrganizationHandler.Create,
			},
		},
	}

	// Member
	routes = append(routes,
		&Route{
			Name:   "getMember",
			Method: "GET",
			Path:   fmt.Sprintf("/organization/%s/members/%s", organizationID, memberID),
			Handlers: []RouteFunc{
				postgresMemberHandler.Read,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "updateMember",
			Method: "PUT",
			Path:   fmt.Sprintf("/organizations/%s/members/%s", organizationID, memberID),
			Handlers: []RouteFunc{
				postgresMemberHandler.Update,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "deleteMember",
			Method: "DELETE",
			Path:   fmt.Sprintf("/organizations/%s/members/%s", organizationID, memberID),
			Handlers: []RouteFunc{
				postgresMemberHandler.Delete,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "createMember",
			Method: "POST",
			Path:   fmt.Sprintf("/organizations/%s/members", organizationID),
			Handlers: []RouteFunc{
				postgresMemberHandler.Create,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
			},
		},
		&Route{
			Name:   "getMemberList",
			Method: "GET",
			Path:   fmt.Sprintf("/organizations/%s/members", organizationID),
			Handlers: []RouteFunc{
				postgresMemberHandler.List,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "loginMember",
			Method: "POST",
			Path:   "/organizations/members/login",
			Handlers: []RouteFunc{
				postgresMemberHandler.Login,
			},
			Filters: []Filter{},
		},
		&Route{
			Name:   "refreshMemberSession",
			Method: "POST",
			Path:   fmt.Sprintf("/organization/%s/members/%s/refresh", organizationID, memberID),
			Handlers: []RouteFunc{
				postgresMemberHandler.RefreshSession,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "logoutMember",
			Method: "DELETE",
			Path:   fmt.Sprintf("/organizations/%s/members/%s/logout", organizationID, memberID),
			Handlers: []RouteFunc{
				postgresMemberHandler.Logout,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		})

	// Application
	routes = append(routes,
		&Route{
			Name:   "getApplications",
			Method: "GET",
			Path:   fmt.Sprintf("/organizations/%s/applications", organizationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.List,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
			},
		},
		&Route{
			Name:   "getApplication",
			Method: "GET",
			Path:   fmt.Sprintf("/organizations/%s/applications/%s", organizationID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Read,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
				ContextHasOrganizationApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "updateApplication",
			Method: "PUT",
			Path:   fmt.Sprintf("/organizations/%s/applications/%s", organizationID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Update,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
				ContextHasOrganizationApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "deleteApplication",
			Method: "DELETE",
			Path:   fmt.Sprintf("/organizations/%s/applications/%s", organizationID, applicationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Delete,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
				ContextHasOrganizationApplication(postgresApplicationHandler),
			},
		},
		&Route{
			Name:   "createApplication",
			Method: "POST",
			Path:   fmt.Sprintf("/organizations/%s/applications", organizationID),
			Handlers: []RouteFunc{
				postgresApplicationHandler.Create,
			},
			Filters: []Filter{
				ContextHasOrganization(postgresOrganizationHandler),
				ContextHasMember(postgresMemberHandler),
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
			Path:   "/me",
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
			Name:   "updateApplicationUser",
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s", applicationUserID),
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
			Name:   "updateCurrentApplicationUser",
			Method: "PUT",
			Path:   "/me",
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
			Name:   "deleteApplicationUser",
			Method: "DELETE",
			Path:   fmt.Sprintf("/users/%s", applicationUserID),
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Delete,
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
			Path:   "/me",
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
			Path:   "/users",
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
			Path:   "/users/login",
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Login,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
			},
		},
		&Route{
			Name:   "loginCurrentUserApplicationUser",
			Method: "POST",
			Path:   "/me/login",
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
			Path:   fmt.Sprintf("/users/%s/refresh", applicationUserID),
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
			Name:   "refreshCurrentUserApplicationUserSession",
			Method: "POST",
			Path:   "/me/refresh",
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
			Path:   fmt.Sprintf("/users/%s/logout", applicationUserID),
			Handlers: []RouteFunc{
				postgresApplicationUserHandler.Logout,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "logoutCurrentUserApplicationUser",
			Method: "DELETE",
			Path:   "/me/logout",
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
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s/connections", applicationUserID),
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
			Name:   "createCurrentUserConnection",
			Method: "PUT",
			Path:   "/me/connections",
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
			Name:   "createFriendConnectionAlias",
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s/friends", applicationUserID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFriend,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createCurrentUserFriendConnectionAlias",
			Method: "PUT",
			Path:   "/me/friends",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFriend,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createFollowConnectionAlias",
			Method: "PUT",
			Path:   fmt.Sprintf("/users/%s/follow", applicationUserID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFollow,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createCurrentUserFollowConnectionAlias",
			Method: "PUT",
			Path:   "/me/follow",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CreateFollow,
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
			Path:   fmt.Sprintf("/users/%s/connections/%s", applicationUserID, applicationUserToID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Delete,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "deleteCurrentUserConnection",
			Method: "DELETE",
			Path:   fmt.Sprintf("/me/connections/%s/%s", connectionType, applicationUserToID),
			Handlers: []RouteFunc{
				postgresConnectionHandler.Delete,
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
			Path:   fmt.Sprintf("/users/%s/connections/social", applicationUserID),
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
			Name:   "createCurrentUserSocialConnections",
			Method: "POST",
			Path:   "/me/connections/social",
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
				postgresConnectionHandler.FollowingList,
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
			Path:   "/me/follows",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFollowingList,
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
			Path:   "/me/followers",
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
			Path:   "/me/friends",
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserFriends,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getUserConnectionsByState",
			Method: "GET",
			Path:   fmt.Sprintf("/users/%s/%s", applicationUserID, connectionState),
			Handlers: []RouteFunc{
				postgresConnectionHandler.UserConnectionsByState,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserConnectionsByState",
			Method: "GET",
			Path:   fmt.Sprintf("/me/%s", connectionState),
			Handlers: []RouteFunc{
				postgresConnectionHandler.CurrentUserConnectionsByState,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
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
			Path:   fmt.Sprintf("/me/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserUpdate,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
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
			Path:   fmt.Sprintf("/me/events/%s", eventID),
			Handlers: []RouteFunc{
				postgresEventHandler.CurrentUserDelete,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "createEvent",
			Method: "POST",
			Path:   fmt.Sprintf("/users/%s/events", applicationUserID),
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
			Name:   "createCurrentUserEvent",
			Method: "POST",
			Path:   "/me/events",
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
		&Route{
			Name:   "getCurrentUserEvent",
			Method: "GET",
			Path:   fmt.Sprintf("/me/events/%s", eventID),
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
			Path:   "/me/events",
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
			Path:   fmt.Sprintf("/users/%s/feed", applicationUserID),
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
			Name:   "getCurrentUserFeed",
			Method: "GET",
			Path:   "/me/feed",
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
			Path:   fmt.Sprintf("/users/%s/feed/unread/count", applicationUserID),
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
			Name:   "getCurrentUserUnreadFeedCount",
			Method: "GET",
			Path:   "/me/feed/unread/count",
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
			Path:   fmt.Sprintf("/users/%s/feed/unread", applicationUserID),
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeed,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},
		&Route{
			Name:   "getCurrentUserUnreadFeed",
			Method: "GET",
			Path:   "/me/feed/unread",
			Handlers: []RouteFunc{
				postgresEventHandler.UnreadFeed,
			},
			Filters: []Filter{
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
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
				ContextHasApplication(postgresApplicationHandler),
				RateLimitApplication,
				ContextHasApplicationUser(postgresApplicationUserHandler),
			},
		},

		// Misc
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
