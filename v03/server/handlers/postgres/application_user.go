package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/tgflake"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type applicationUser struct {
	storage core.ApplicationUser
	conn    core.Connection
}

func (appUser *applicationUser) Read(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	user, err := appUser.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), userID)
	if err != nil {
		return err
	}

	// maybe not the most efficient way to do it?
	relation := make(chan *entity.Relation)
	go func(relation chan *entity.Relation, accountID, applicationID int64, userFromID, userToID uint64) {
		rel, _ := appUser.conn.Relation(accountID, applicationID, userFromID, userToID)

		relation <- rel
	}(relation, ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID, user.ID)

	response.ComputeApplicationUserLastModified(ctx, user)

	user.Password = ""
	user.Deleted = nil
	user.CreatedAt, user.UpdatedAt, user.LastLogin, user.LastRead = nil, nil, nil, nil

	if rel := <-relation; rel != nil {
		user.IsFriends = rel.IsFriends
		user.IsFollower = rel.IsFollower
		user.IsFollowed = rel.IsFollowed
	} else {
		user.IsFriends = entity.PFalse
		user.IsFollower = entity.PFalse
		user.IsFollowed = entity.PFalse
	}

	response.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) ReadCurrent(ctx *context.Context) (err []errors.Error) {
	user := ctx.Bag["applicationUser"].(*entity.ApplicationUser)
	user.Password = ""
	user.Deleted = nil
	user.IsFriends = entity.PFalse
	user.IsFollower = entity.PFalse
	user.IsFollowed = entity.PFalse

	response.ComputeApplicationUserLastModified(ctx, user)

	response.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) UpdateCurrent(ctx *context.Context) (err []errors.Error) {
	user := *(ctx.Bag["applicationUser"].(*entity.ApplicationUser))
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	user.ID = ctx.Bag["applicationUserID"].(uint64)

	if err = validator.UpdateUser(
		appUser.storage,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser),
		&user); err != nil {
		return
	}

	updatedUser, err := appUser.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		*(ctx.Bag["applicationUser"].(*entity.ApplicationUser)),
		user,
		false)
	if err != nil {
		return
	}
	if updatedUser == nil {
		updatedUser = &user
	}

	updatedUser.Password = ""

	response.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	if err = appUser.storage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		userID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) DeleteCurrent(ctx *context.Context) (err []errors.Error) {
	if err = appUser.storage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID); err != nil {
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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	err = validator.CreateUser(appUser.storage, ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user)
	if err != nil {
		if err[0] == errmsg.ErrApplicationUserEmailAlreadyExists || err[0] == errmsg.ErrApplicationUserUsernameInUse {
			return appUser.Login(ctx)
		}

		return
	}

	timeNow := time.Now()
	user.LastLogin = &timeNow

	user.ID, er = tgflake.FlakeNextID(ctx.Bag["applicationID"].(int64), "users")
	if er != nil {
		return []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(er.Error())}
	}

	user, err = appUser.storage.Create(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user, true)
	if err != nil {
		return
	}

	sessionToken, err := appUser.storage.CreateSession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user)
	if err != nil {
		return
	}

	user.Password = ""
	user.SessionToken = sessionToken

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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		if !(err[0] == errmsg.ErrAuthGotBothUsernameAndEmail) {
			return
		}
	}

	if loginPayload.EmailName != "" {
		loginPayload.Email = loginPayload.EmailName
		loginPayload.Username = loginPayload.EmailName
	}

	if loginPayload.Email != "" {
		user, err = appUser.storage.FindByEmail(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Email)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err[0].Error() != "application user not found" {
			return
		}
	}

	if loginPayload.Username != "" && user == nil {
		user, err = appUser.storage.FindByUsername(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Username)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err[0].Error() != "application user not found" {
			return
		}
	}

	if user == nil || !user.Enabled {
		return []errors.Error{errmsg.ErrApplicationUserNotFound}
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = appUser.storage.CreateSession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user); err != nil {
		return
	}

	timeNow := time.Now()
	user.LastLogin = &timeNow
	_, err = appUser.storage.Update(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), *user, *user, false)
	if err != nil {
		return
	}

	user.Password = ""
	user.SessionToken = sessionToken

	ctx.W.Header().Set("Location", fmt.Sprintf("https://api.tapglue.com/0.3/users/%d", user.ID))
	response.WriteResponse(ctx, user, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	tokenPayload := struct {
		Token string `json:"session_token"`
	}{}

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if tokenPayload.Token != ctx.SessionToken {
		return []errors.Error{errmsg.ErrAuthSessionTokenMismatch}
	}

	if tokenPayload.Token, err = appUser.storage.RefreshSession(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.SessionToken, ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	response.WriteResponse(ctx, tokenPayload, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Logout(ctx *context.Context) (err []errors.Error) {
	if err = appUser.storage.DestroySession(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.SessionToken,
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

func (appUser *applicationUser) Search(ctx *context.Context) (err []errors.Error) {
	query := ctx.Query.Get("q")
	if query == "" {
		response.WriteResponse(ctx, []*entity.ApplicationUser{}, http.StatusNoContent, 10)
		return
	}

	if len(query) < 3 {
		return []errors.Error{errmsg.ErrApplicationUserSearchTypeMin3Chars}
	}

	users, err := appUser.storage.Search(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), query)
	if err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)

	for idx := range users {
		users[idx].Password = ""
		users[idx].Deleted = nil
		users[idx].SessionToken = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
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
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidApplicationUserCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", user, pass))}
	}
	if pass == "" {
		return []errors.Error{errmsg.ErrAuthUserSessionNotSet}
	}
	ctx.Bag["applicationUser"], err = appUser.storage.FindBySession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), pass)
	if err == nil {
		ctx.Bag["applicationUserID"] = ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID
		ctx.SessionToken = pass
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
