/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/context"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
func getAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	WriteResponse(ctx, ctx.AccountUser, http.StatusOK, 10)
	return
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:UserID
func updateAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	accountUser := *ctx.AccountUser
	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return tgerrors.NewBadRequestError("failed to update the account user (1)\n"+er.Error(), er.Error())
	}

	accountUser.ID = ctx.AccountUserID
	accountUser.AccountID = ctx.AccountID

	if err = validator.UpdateAccountUser(ctx.AccountUser, &accountUser); err != nil {
		return
	}

	updatedAccountUser, err := core.UpdateAccountUser(*ctx.AccountUser, accountUser, true)
	if err != nil {
		return
	}

	updatedAccountUser.Password = ""
	WriteResponse(ctx, updatedAccountUser, http.StatusCreated, 0)
	return
}

// deleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:UserID
func deleteAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	if err = core.DeleteAccountUser(ctx.AccountID, ctx.AccountUserID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
func createAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		accountUser = &entity.AccountUser{}
	)

	if err := json.Unmarshal(ctx.Body, accountUser); err != nil {
		return tgerrors.NewBadRequestError("failed to create the account user (1)"+err.Error(), err.Error())
	}

	accountUser.AccountID = ctx.AccountID

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

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
func getAccountUserList(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		accountUsers []*entity.AccountUser
	)

	if accountUsers, err = core.ReadAccountUserList(ctx.AccountID); err != nil {
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

// loginAccountUser handles the requests to login the user in the system
// Request: POST /account/user/login
func loginAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
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

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: Post /account/:AccountID/application/:ApplicationID/user/refreshsession
func refreshAccountUserSession(ctx *context.Context) (err *tgerrors.TGError) {
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

	if sessionToken, err = core.RefreshAccountUserSession(ctx.SessionToken, ctx.AccountUser); err != nil {
		return
	}

	WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

// logoutApplicationUser handles the requests to logout the user from the system
// Request: Post /account/:AccountID/application/:ApplicationID/user/logout
func logoutAccountUser(ctx *context.Context) (err *tgerrors.TGError) {
	var logoutPayload struct {
		Token string `json:"token"`
	}

	if er := json.Unmarshal(ctx.Body, &logoutPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to logout the user (1)\n"+er.Error(), er.Error())
	}

	if ctx.SessionToken != logoutPayload.Token {
		return tgerrors.NewBadRequestError("failed to logout the user (2) \nsession token mismatch", "session token mismatch")
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, ctx.AccountUser); err != nil {
		return
	}

	WriteResponse(ctx, "logged out", http.StatusOK, 0)
	return
}
