/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package response

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/v02/entity"
)

func ComputeAccountLastModified(ctx *context.Context, account *entity.Account) {
	ctx.Bag["Last-Modified"] = account.UpdatedAt.Format(http.TimeFormat)
}

func ComputeAccountsLastModified(ctx *context.Context, accounts []*entity.Account) {
	if len(accounts) == 0 {
		ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
		return
	}
	highTime := accounts[0].UpdatedAt
	for idx := range accounts {
		if accounts[idx].UpdatedAt.After(*highTime) {
			highTime = accounts[idx].UpdatedAt
		}
	}
	ctx.Bag["Last-Modified"] = highTime.Format(http.TimeFormat)
}

func ComputeAccountUserLastModified(ctx *context.Context, user *entity.AccountUser) {
	ctx.Bag["Last-Modified"] = user.UpdatedAt.Format(http.TimeFormat)
}

func ComputeAccountUsersLastModified(ctx *context.Context, users []*entity.AccountUser) {
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

func ComputeApplicationLastModified(ctx *context.Context, application *entity.Application) {
	ctx.Bag["Last-Modified"] = application.UpdatedAt.Format(http.TimeFormat)
}

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

func ComputeApplicationUserLastModified(ctx *context.Context, user *entity.ApplicationUser) {
	ctx.Bag["Last-Modified"] = user.UpdatedAt.Format(http.TimeFormat)
}

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

func ComputeConnectionLastModified(ctx *context.Context, connection *entity.Connection) {
	ctx.Bag["Last-Modified"] = connection.UpdatedAt.Format(http.TimeFormat)
}

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

func ComputeEventLastModified(ctx *context.Context, event *entity.Event) {
	ctx.Bag["Last-Modified"] = event.UpdatedAt.Format(http.TimeFormat)
}

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

func ComputeLastModifiedNow(ctx *context.Context) {
	ctx.Bag["Last-Modified"] = ctx.StartTime.Format(http.TimeFormat)
}
