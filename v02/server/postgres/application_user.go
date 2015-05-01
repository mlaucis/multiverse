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
	user := ctx.Bag["applicationUser"].(*entity.ApplicationUser)
	user.Password = ""
	server.WriteResponse(ctx, user, http.StatusOK, 10)
	return
}

func (appUser *applicationUser) Update(ctx *context.Context) (err errors.Error) {
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
		true)
	if err != nil {
		return
	}

	updatedUser.Password = ""

	server.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err errors.Error) {
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
		user = &entity.ApplicationUser{}
		er   error
	)

	if er = json.Unmarshal(ctx.Body, user); er != nil {
		return errors.NewBadRequestError("failed to create the application user (1)\n"+er.Error(), er.Error())
	}

	if err = validator.CreateUser(appUser.storage, ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user); err != nil {
		return
	}

	if user, err = appUser.storage.Create(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user, true); err != nil {
		return
	}

	user.Password = ""

	server.WriteResponse(ctx, user, http.StatusCreated, 0)
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

	if loginPayload.Email != "" {
		user, err = appUser.storage.FindByEmail(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Email)
		if err != nil {
			return
		}
	}

	if loginPayload.Username != "" {
		user, err = appUser.storage.FindByUsername(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), loginPayload.Username)
		if err != nil {
			return
		}
	}

	if user == nil {
		return errors.NewInternalError("failed to login the application user (2)\n", "user is nil")
	}

	if !user.Enabled {
		return errors.NewNotFoundError("failed to login the user (3)\nuser is disabled", "user is disabled")
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = appUser.storage.CreateSession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), user); err != nil {
		return
	}

	user.LastLogin = time.Now()
	_, err = appUser.storage.Update(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), *user, *user, false)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, struct {
		UserID string `json:"id"`
		Token  string `json:"session_token"`
	}{
		UserID: user.ID,
		Token:  sessionToken,
	}, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err errors.Error) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
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

func (appUser *applicationUser) PopulateContext(ctx *context.Context) (err errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return errors.NewBadRequestError("error while reading user credentials", fmt.Sprintf("got %s:%s", user, pass))
	}
	ctx.Bag["applicationUser"], err = appUser.storage.FindBySession(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), pass)
	if err == nil {
		ctx.Bag["applicationUserID"] = ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID
	}
	return
}

// NewApplicationUser returns a new application user routes handler
func NewApplicationUser(storage core.ApplicationUser) server.ApplicationUser {
	return &applicationUser{
		storage: storage,
	}
}
