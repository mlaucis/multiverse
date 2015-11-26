package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/response"
	"github.com/tapglue/multiverse/v04/validator"
)

type applicationUser struct {
	storage core.ApplicationUser
	conn    core.Connection
}

func (appUser *applicationUser) Read(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	user, err := appUser.storage.Read(ctx.OrganizationID, ctx.ApplicationID, userID, false)
	if err != nil {
		return err
	}

	err = appUser.storage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, user)
	if err != nil {
		ctx.LogError(err)
	}

	user.IsFriend = entity.PFalse
	user.IsFollower = entity.PFalse
	user.IsFollowed = entity.PFalse

	// maybe not the most efficient way to do it?
	rel, err := appUser.conn.Relation(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUser.ID, user.ID)
	if err != nil {
		return err
	}
	if rel != nil {
		user.Relation = *rel
	}

	user.Password = ""
	user.Deleted = nil
	user.CreatedAt, user.UpdatedAt, user.LastLogin, user.LastRead = nil, nil, nil, nil

	response.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) ReadCurrent(ctx *context.Context) (err []errors.Error) {
	user := ctx.ApplicationUser
	user.Password = ""
	user.Deleted = nil
	user.IsFriend = nil
	user.IsFollower = nil
	user.IsFollowed = nil

	appUser.storage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, user)

	response.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) UpdateCurrent(ctx *context.Context) (err []errors.Error) {
	user := *ctx.ApplicationUser
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	user.ID = ctx.ApplicationUserID

	if err = validator.UpdateUser(
		appUser.storage,
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser,
		&user); err != nil {
		return
	}

	updatedUser, err := appUser.storage.Update(
		ctx.OrganizationID,
		ctx.ApplicationID,
		*ctx.ApplicationUser,
		user,
		false)
	if err != nil {
		return
	}
	if updatedUser == nil {
		updatedUser = &user
	}

	updatedUser.Password = ""
	appUser.storage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, updatedUser)

	response.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	if err = appUser.storage.Delete(
		ctx.OrganizationID,
		ctx.ApplicationID,
		userID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) DeleteCurrent(ctx *context.Context) (err []errors.Error) {
	if err = appUser.storage.Delete(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser.ID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) Create(ctx *context.Context) (err []errors.Error) {
	var (
		user = &entity.ApplicationUser{}
		er   error
	)

	if er = json.Unmarshal(ctx.Body, user); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	err = validator.CreateUser(appUser.storage, ctx.OrganizationID, ctx.ApplicationID, user)
	if err != nil {
		if err[0].Code() == errmsg.ErrApplicationUserEmailAlreadyExists.Code() ||
			err[0].Code() == errmsg.ErrApplicationUserUsernameInUse.Code() {
			return appUser.Login(ctx)
		}

		return
	}

	timeNow := time.Now()
	user.LastLogin = &timeNow

	user.ID, er = tgflake.FlakeNextID(ctx.ApplicationID, "users")
	if er != nil {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(er.Error()).SetCurrentLocation()}
	}

	err = appUser.storage.Create(ctx.OrganizationID, ctx.ApplicationID, user)
	if err != nil {
		return
	}

	sessionToken, err := appUser.storage.CreateSession(ctx.OrganizationID, ctx.ApplicationID, user)
	if err != nil {
		return
	}

	response.SanitizeApplicationUser(user)
	user.SessionToken = sessionToken
	appUser.storage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, user)

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/users/%d", user.ID))
	response.WriteResponse(ctx, user, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Login(ctx *context.Context) (err []errors.Error) {
	var (
		loginPayload = &entity.LoginPayload{}
		user         *entity.ApplicationUser
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		if !(err[0].Code() == errmsg.ErrAuthGotBothUsernameAndEmail.Code()) {
			return
		}
	}

	if loginPayload.EmailName != "" {
		loginPayload.Email = loginPayload.EmailName
		loginPayload.Username = loginPayload.EmailName
	}

	if loginPayload.Email != "" {
		user, err = appUser.storage.FindByEmail(ctx.OrganizationID, ctx.ApplicationID, loginPayload.Email)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err[0].Error() != "application user not found" {
			return
		}
	}

	if loginPayload.Username != "" && user == nil {
		user, err = appUser.storage.FindByUsername(ctx.OrganizationID, ctx.ApplicationID, loginPayload.Username)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err[0].Error() != "application user not found" {
			return
		}
	}

	if user == nil || !user.Enabled || user.Deleted == entity.PFalse {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = appUser.storage.CreateSession(ctx.OrganizationID, ctx.ApplicationID, user); err != nil {
		return
	}

	timeNow := time.Now()
	user.LastLogin = &timeNow
	user, err = appUser.storage.Update(ctx.OrganizationID, ctx.ApplicationID, *user, *user, true)
	if err != nil {
		return
	}

	response.SanitizeApplicationUser(user)
	user.SessionToken = sessionToken
	user.IsFriend = nil
	user.IsFollower = nil
	user.IsFollowed = nil

	appUser.storage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, user)

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/users/%d", user.ID))
	response.WriteResponse(ctx, user, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	tokenPayload := struct {
		Token string `json:"session_token"`
	}{}

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	if tokenPayload.Token != ctx.SessionToken {
		return []errors.Error{errmsg.ErrAuthSessionTokenMismatch.SetCurrentLocation()}
	}

	if tokenPayload.Token, err = appUser.storage.RefreshSession(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.SessionToken, ctx.ApplicationUser); err != nil {
		return
	}

	response.WriteResponse(ctx, tokenPayload, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Logout(ctx *context.Context) (err []errors.Error) {
	if err = appUser.storage.DestroySession(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.SessionToken,
		ctx.ApplicationUser); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

func (appUser *applicationUser) Search(ctx *context.Context) (err []errors.Error) {
	var users []*entity.ApplicationUser

	query := ctx.Query.Get("q")
	if query == "" {
		if len(ctx.Query["socialid"]) != 0 && ctx.Query.Get("social_platform") != "" {
			users, err = appUser.storage.FilterBySocialIDs(
				ctx.OrganizationID, ctx.ApplicationID,
				ctx.Query.Get("social_platform"),
				ctx.Query["socialid"])
			if err != nil {
				return err
			}
			goto userProcessing
		}

		if len(ctx.Query["email"]) != 0 {
			users, err = appUser.storage.FilterByEmail(
				ctx.OrganizationID, ctx.ApplicationID,
				ctx.Query["email"])
			if err != nil {
				return err
			}
			goto userProcessing
		}
		return
	}

	if len(query) < 3 {
		return []errors.Error{errmsg.ErrApplicationUserSearchTypeMin3Chars.SetCurrentLocation()}
	}

	users, err = appUser.storage.Search(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID, query)
	if err != nil {
		return err
	}

userProcessing:
	response.SanitizeApplicationUsers(users)

	for idx := range users {
		relation, err := appUser.conn.Relation(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID, users[idx].ID)
		if err != nil {
			return err
		} else if relation != nil {
			users[idx].Relation = *relation
		}
	}

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	if resp.UsersCount == 0 {
		response.WriteResponse(ctx, resp, http.StatusNoContent, 10)
		return
	}

	response.WriteResponse(ctx, resp, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) PopulateContext(ctx *context.Context) (err []errors.Error) {
	appToken, userToken, ok := ctx.BasicAuth()

	if ctx.TokenType == context.TokenTypeApplication {
		if !ok {
			return []errors.Error{errmsg.ErrAuthInvalidApplicationUserCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", appToken, userToken)).SetCurrentLocation()}
		}
		if userToken == "" {
			return []errors.Error{errmsg.ErrAuthUserSessionNotSet.SetCurrentLocation()}
		}
		ctx.ApplicationUser, err = appUser.storage.FindBySession(ctx.Application.OrgID, ctx.Application.ID, userToken)
		if err == nil {
			ctx.ApplicationUserID = ctx.ApplicationUser.ID
			ctx.SessionToken = userToken
		}
	} else if ctx.TokenType == context.TokenTypeBackend {
		var (
			userID uint64
			er     error
		)
		if userToken != "" {
			userID, er = strconv.ParseUint(userToken, 10, 64)
		} else if val, ok := ctx.Vars["applicationUserID"]; ok {
			userID, er = strconv.ParseUint(val, 10, 64)
		} else {
			return []errors.Error{errmsg.ErrApplicationUserIDInvalid.UpdateMessage("user ID could not be read from the request").SetCurrentLocation()}
		}
		if er != nil {
			return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
		}
		ctx.ApplicationUser, err = appUser.storage.Read(ctx.Application.OrgID, ctx.Application.ID, userID, false)
		if err == nil {
			ctx.ApplicationUserID = ctx.ApplicationUser.ID
		}
	} else {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(fmt.Sprintf("unexpected context token type, got %d", ctx.TokenType)).SetCurrentLocation()}
	}

	return
}

// NewApplicationUser returns a new application user routes handler
func NewApplicationUser(storage core.ApplicationUser, conn core.Connection) handlers.ApplicationUser {
	return &applicationUser{
		storage: storage,
		conn:    conn,
	}
}
