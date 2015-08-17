package response

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/v03/entity"
)

// ComputeOrganizationLastModified computes the last-modified information for the account
func ComputeOrganizationLastModified(ctx *context.Context, organization *entity.Organization) {
	ctx.Bag["Last-Modified"] = organization.UpdatedAt.Format(http.TimeFormat)
}

// ComputeOrganizationsLastModified computes the last-modified information for a list of accounts
func ComputeOrganizationsLastModified(ctx *context.Context, organizations []*entity.Organization) {
	if len(organizations) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := organizations[0].UpdatedAt
	for idx := range organizations {
		if organizations[idx].UpdatedAt.After(*highTime) {
			highTime = organizations[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeMemberLastModified computes the last-modified information for the account user
func ComputeMemberLastModified(ctx *context.Context, member *entity.Member) {
	ctx.Bag["Last-Modified"] = member.UpdatedAt.Format(http.TimeFormat)
}

// ComputeMembersLastModified computes the last-modified information for a list of account users
func ComputeMembersLastModified(ctx *context.Context, members []*entity.Member) {
	if len(members) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := members[0].UpdatedAt
	for idx := range members {
		if members[idx].UpdatedAt.After(*highTime) {
			highTime = members[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeApplicationLastModified computes the last-modified information for an application
func ComputeApplicationLastModified(ctx *context.Context, application *entity.Application) {
	ctx.Bag["Last-Modified"] = application.UpdatedAt.Format(http.TimeFormat)
}

// ComputeApplicationsLastModified computes the last-modified information for a list of applications
func ComputeApplicationsLastModified(ctx *context.Context, applications []*entity.Application) {
	if len(applications) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := applications[0].UpdatedAt
	for idx := range applications {
		if applications[idx].UpdatedAt.After(*highTime) {
			highTime = applications[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeApplicationUserLastModified computes the last-modified information for an application user
func ComputeApplicationUserLastModified(ctx *context.Context, user *entity.ApplicationUser) {
	ctx.Bag["Last-Modified"] = user.UpdatedAt.Format(http.TimeFormat)
}

// ComputeApplicationUsersLastModified computes the last-modified information for a list of application users
func ComputeApplicationUsersLastModified(ctx *context.Context, users []*entity.ApplicationUser) {
	if len(users) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := users[0].UpdatedAt
	for idx := range users {
		if users[idx].UpdatedAt.After(*highTime) {
			highTime = users[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeConnectionLastModified computes the last-modified information for a connection
func ComputeConnectionLastModified(ctx *context.Context, connection *entity.Connection) {
	ctx.Bag["Last-Modified"] = connection.UpdatedAt.Format(http.TimeFormat)
}

// ComputeConnectionsLastModified computes the last-modified information for a list of connections
func ComputeConnectionsLastModified(ctx *context.Context, connections []*entity.Connection) {
	if len(connections) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := connections[0].UpdatedAt
	for idx := range connections {
		if connections[idx].UpdatedAt.After(*highTime) {
			highTime = connections[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeEventLastModified computes the last-modified information for an event
func ComputeEventLastModified(ctx *context.Context, event *entity.Event) {
	ctx.Bag["Last-Modified"] = event.UpdatedAt.Format(http.TimeFormat)
}

// ComputeEventsLastModified computes the last-modified information for a list of events
func ComputeEventsLastModified(ctx *context.Context, events []*entity.Event) {
	if len(events) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := events[0].UpdatedAt
	for idx := range events {
		if events[idx].UpdatedAt.After(*highTime) {
			highTime = events[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

// ComputeLastModifiedNow computes the last-modified to be the request start time
func ComputeLastModifiedNow(ctx *context.Context) {
	ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
}
