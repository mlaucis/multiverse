package http

import (
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/member"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	ctxKeyApp       = "app"
	ctxKeyDeviceID  = "deviceID"
	ctxKeyMember    = "member"
	ctxKeyOrg       = "org"
	ctxKeyRoute     = "route"
	ctxKeyToken     = "token"
	ctxKeyTokenType = "tokenType"
	ctxKeyUser      = "user"
	ctxKeyVersion   = "version"

	tokenApplication = "application"
	tokenBackend     = "backend"
)

func appFromContext(ctx context.Context) *app.App {
	return ctx.Value(ctxKeyApp).(*app.App)
}

func appInContext(ctx context.Context, app *app.App) context.Context {
	return context.WithValue(ctx, ctxKeyApp, app)
}

func deviceIDFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyDeviceID).(string)
}

func deviceIDInContext(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, ctxKeyDeviceID, deviceID)
}

func memberFromContext(ctx context.Context) *member.Member {
	return ctx.Value(ctxKeyMember).(*member.Member)
}

func memberInContext(ctx context.Context, app *member.Member) context.Context {
	return context.WithValue(ctx, ctxKeyMember, app)
}

func originFromContext(ctx context.Context) controller.Origin {
	var (
		currentUser = userFromContext(ctx)
		deviceID    = deviceIDFromContext(ctx)
		tokenType   = tokenTypeFromContext(ctx)
	)

	return createOrigin(deviceID, tokenType, currentUser.ID)
}

func orgFromContext(ctx context.Context) *v04_entity.Organization {
	return ctx.Value(ctxKeyOrg).(*v04_entity.Organization)
}

func orgInContext(ctx context.Context, app *v04_entity.Organization) context.Context {
	return context.WithValue(ctx, ctxKeyOrg, app)
}

func routeFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyRoute).(string)
}

func routeInContext(ctx context.Context, route string) context.Context {
	return context.WithValue(ctx, ctxKeyRoute, route)
}

func tokenFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyToken).(string)
}

func tokenInContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxKeyToken, token)
}

func tokenTypeFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyTokenType).(string)
}

func tokenTypeInContext(ctx context.Context, tokenType string) context.Context {
	return context.WithValue(ctx, ctxKeyTokenType, tokenType)
}

func userFromContext(ctx context.Context) *user.User {
	return ctx.Value(ctxKeyUser).(*user.User)
}

func userInContext(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, user)
}

func versionFromContext(ctx context.Context) string {
	return ctx.Value(ctxKeyVersion).(string)
}

func versionInContext(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, ctxKeyVersion, version)
}
