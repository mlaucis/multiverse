/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	applicationUser struct {
		storage core.ApplicationUser
	}
)

func (appUser *applicationUser) Read(ctx *context.Context) (err errors.Error) {
	user, err := appUser.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Vars["applicationUserID"])
	if err != nil {
		return err
	}
	user.Password = ""
	user.Enabled = false
	user.CreatedAt, user.UpdatedAt, user.LastLogin, user.LastRead = nil, nil, nil, nil
	server.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) ReadCurrent(ctx *context.Context) (err errors.Error) {
	user := ctx.Bag["applicationUser"].(*entity.ApplicationUser)
	user.Password = ""
	user.Enabled = false
	server.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) UpdateCurrent(ctx *context.Context) (err errors.Error) {
	user := *(ctx.Bag["applicationUser"].(*entity.ApplicationUser))
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return errors.NewBadRequestError("failed to update the user (1)\n"+er.Error(), er.Error())
	}

	user.ID = ctx.Bag["applicationUserID"].(string)

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
	updatedUser.Enabled = false

	server.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) DeleteCurrent(ctx *context.Context) (err errors.Error) {
	if err = appUser.storage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) Create(ctx *context.Context) (err errors.Error) {
	var (
		user      = &entity.ApplicationUser{}
		er        error
		withLogin = ctx.R.URL.Query().Get("withLogin") == "true"
	)

	if er = json.Unmarshal(ctx.Body, user); er != nil {
		return errors.NewBadRequestError("failed to create the application user (1)\n"+er.Error(), er.Error())
	}

	if err = validator.CreateUser(appUser.storage, ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user); err != nil {
		return
	}

	if withLogin {
		timeNow := time.Now()
		user.LastLogin = &timeNow
	}
	if user, err = appUser.storage.Create(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user, true); err != nil {
		return
	}

	sessionToken := ""
	if withLogin {
		if sessionToken, err = appUser.storage.CreateSession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user); err != nil {
			return
		}
	}

	user.Password = ""
	user.Enabled = false

	result := struct {
		entity.ApplicationUser
		SessionToken string `json:"session_token,omitempty"`
	}{
		ApplicationUser: *user,
		SessionToken:    sessionToken,
	}

	server.WriteResponse(ctx, result, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Login(ctx *context.Context) (err errors.Error) {
	var (
		loginPayload = &entity.LoginPayload{}
		user         *entity.ApplicationUser
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return errors.NewBadRequestError("failed to login the user (1)\n"+er.Error(), er.Error())
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		return
	}

	if loginPayload.EmailName != "" {
		loginPayload.Email = loginPayload.EmailName
		loginPayload.Username = loginPayload.EmailName
	}

	if loginPayload.Email != "" {
		user, err = appUser.storage.FindByEmail(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Email)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err.Error() != "application user not found" {
			return
		}
	}

	if loginPayload.Username != "" && user == nil {
		user, err = appUser.storage.FindByUsername(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Username)
		// TODO This is horrible and I should change it when we have constant errors
		if err != nil && err.Error() != "application user not found" {
			return
		}
	}

	if user == nil || !user.Enabled {
		return errors.NewNotFoundError("application user not found", "user not found")
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

	response := struct {
		entity.ApplicationUser
		UserID string `json:"id"`
		Token  string `json:"session_token"`
	}{
		UserID: user.ID,
		Token:  sessionToken,
	}

	if ctx.R.URL.Query().Get("withUserDetails") == "true" {
		user.Password = ""
		user.Enabled = false
		response.ApplicationUser = *user
	}

	server.WriteResponse(ctx, response, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err errors.Error) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		sessionToken string
	)

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return errors.NewBadRequestError("failed to refresh the session token (1)\n"+er.Error(), er.Error())
	}

	if tokenPayload.Token != ctx.SessionToken {
		return errors.NewBadRequestError("failed to refresh the session token (2)\nsession token mismatch", "session token mismatch")
	}

	if sessionToken, err = appUser.storage.RefreshSession(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.SessionToken, ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, struct {
		Token string `json:"session_token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Logout(ctx *context.Context) (err errors.Error) {
	if err = appUser.storage.DestroySession(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.SessionToken,
		ctx.Bag["applicationUser"].(*entity.ApplicationUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

func (appUser *applicationUser) Search(ctx *context.Context) (err errors.Error) {
	query := ctx.Query.Get("q")
	if query == "" {
		server.WriteResponse(ctx, []*entity.ApplicationUser{}, http.StatusNoContent, 10)
		return
	}

	if len(query) < 3 {
		return errors.NewBadRequestError("type at least 3 characters to search", "less than 3 chars for search")
	}

	users, err := appUser.storage.Search(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), query)
	if err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	if response.UsersCount == 0 {
		server.WriteResponse(ctx, response, http.StatusNoContent, 10)
		return
	}

	server.WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) PopulateContext(ctx *context.Context) (err errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return errors.NewBadRequestError("error while reading user credentials", fmt.Sprintf("got %s:%s", user, pass))
	}
	ctx.Bag["applicationUser"], err = appUser.storage.FindBySession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), pass)
	if err == nil {
		ctx.Bag["applicationUserID"] = ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID
		ctx.SessionToken = pass
	}
	return
}

// NewApplicationUser returns a new application user routes handler
func NewApplicationUser(storage core.ApplicationUser) server.ApplicationUser {
	return &applicationUser{
		storage: storage,
	}
}
