/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

// GetAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
func GetAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	WriteResponse(ctx, ctx.Bag["accountUser"].(*entity.AccountUser), http.StatusOK, 10)
	return
}

// UpdateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:UserID
func UpdateAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	accountUser := *(ctx.Bag["accountUser"].(*entity.AccountUser))
	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return tgerrors.NewBadRequestError("failed to update the account user (1)\n"+er.Error(), er.Error())
	}

	accountUser.ID = ctx.Bag["accountUserID"].(int64)
	accountUser.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.UpdateAccountUser(ctx.Bag["accountUser"].(*entity.AccountUser), &accountUser); err != nil {
		return
	}

	updatedAccountUser, err := core.UpdateAccountUser(*(ctx.Bag["accountUser"].(*entity.AccountUser)), accountUser, true)
	if err != nil {
		return
	}

	updatedAccountUser.Password = ""
	WriteResponse(ctx, updatedAccountUser, http.StatusCreated, 0)
	return
}

// DeleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:UserID
func DeleteAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	if err = core.DeleteAccountUser(ctx.Bag["accountID"].(int64), ctx.Bag["accountUserID"].(int64)); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// CreateAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
func CreateAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		accountUser = &entity.AccountUser{}
	)

	if err := json.Unmarshal(ctx.Body, accountUser); err != nil {
		return tgerrors.NewBadRequestError("failed to create the account user (1)"+err.Error(), err.Error())
	}

	accountUser.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.CreateAccountUser(accountUser); err != nil {
		return
	}

	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		return
	}

	accountUser.Password = ""

	WriteResponse(ctx, accountUser, http.StatusCreated, 0)
	return
}

// GetAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
func GetAccountUserList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		accountUsers []*entity.AccountUser
	)

	if accountUsers, err = core.ReadAccountUserList(ctx.Bag["accountID"].(int64)); err != nil {
		//		utils.ErrorHappened(ctx, "failed to retrieve the users (1)", http.StatusInternalServerError, err)
		return
	}

	for idx := range accountUsers {
		accountUsers[idx].Password = ""
	}

	response := &struct {
		AccountUsers []*entity.AccountUser `json:"accountUsers"`
	}{
		AccountUsers: accountUsers,
	}

	WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

// LoginAccountUser handles the requests to login the user in the system
// Request: POST /account/user/login
func LoginAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		loginPayload = &entity.LoginPayload{}
		account      *entity.Account
		user         *entity.AccountUser
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to login the user (1)\n"+er.Error(), er.Error())
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		return
	}

	if loginPayload.Email != "" {
		account, user, err = core.FindAccountAndUserByEmail(loginPayload.Email)
		if err != nil {
			return
		}
	}

	if loginPayload.Username != "" {
		account, user, err = core.FindAccountAndUserByUsername(loginPayload.Username)
		if err != nil {
			return
		}
	}

	if err = validator.AccountUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = core.CreateAccountUserSession(user); err != nil {
		return
	}

	user.LastLogin = time.Now()
	_, err = core.UpdateAccountUser(*user, *user, false)

	WriteResponse(ctx, struct {
		ID           int64  `json:"id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccountToken: account.AuthToken,
		Token:        sessionToken,
	}, http.StatusCreated, 0)
	return
}

// RefreshAccountUserSession handles the requests to refresh the account user session token
// Request: Post /account/:AccountID/application/:ApplicationID/user/refreshsession
func RefreshAccountUserSession(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
	)

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to refresh session token (1)\n"+er.Error(), er.Error())
	}

	if ctx.SessionToken != tokenPayload.Token {
		return tgerrors.NewBadRequestError("failed to refresh session token (2) \nsession token mismatch", "session token mismatch")
	}

	if sessionToken, err = core.RefreshAccountUserSession(ctx.SessionToken, ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

// LogoutAccountUser handles the requests to logout the account user from the system
// Request: Post /account/:AccountID/application/:ApplicationID/user/logout
func LogoutAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	var logoutPayload struct {
		Token string `json:"token"`
	}

	if er := json.Unmarshal(ctx.Body, &logoutPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to logout the user (1)\n"+er.Error(), er.Error())
	}

	if ctx.SessionToken != logoutPayload.Token {
		return tgerrors.NewBadRequestError("failed to logout the user (2) \nsession token mismatch", "session token mismatch")
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	WriteResponse(ctx, "logged out", http.StatusOK, 0)
	return
}
