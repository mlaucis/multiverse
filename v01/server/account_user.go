/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
func getAccountUser(ctx *context.Context) {
	utils.WriteResponse(ctx, ctx.AccountUser, http.StatusOK, 10)
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:UserID
func updateAccountUser(ctx *context.Context) {
	var err error

	accountUser := *ctx.AccountUser
	if err = json.NewDecoder(ctx.Body).Decode(&accountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (1)"+err.Error(), http.StatusBadRequest, err)
		return
	}

	accountUser.ID = ctx.AccountUserID
	accountUser.AccountID = ctx.AccountID

	if err = validator.UpdateAccountUser(ctx.AccountUser, &accountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedAccountUser, err := core.UpdateAccountUser(*ctx.AccountUser, accountUser, true)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (3)", http.StatusInternalServerError, err)
		return
	}

	updatedAccountUser.Password = ""
	utils.WriteResponse(ctx, updatedAccountUser, http.StatusCreated, 0)
}

// deleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:UserID
func deleteAccountUser(ctx *context.Context) {
	if err := core.DeleteAccountUser(ctx.AccountID, ctx.AccountUserID); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the user (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "", http.StatusNoContent, 10)
}

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
func createAccountUser(ctx *context.Context) {
	var (
		accountUser = &entity.AccountUser{}
		err         error
	)

	if err = json.NewDecoder(ctx.Body).Decode(accountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	accountUser.AccountID = ctx.AccountID

	if err = validator.CreateAccountUser(accountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (3)", http.StatusInternalServerError, err)
		return
	}

	accountUser.Password = ""

	utils.WriteResponse(ctx, accountUser, http.StatusCreated, 0)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
func getAccountUserList(ctx *context.Context) {
	var (
		accountUsers []*entity.AccountUser
		err          error
	)

	if accountUsers, err = core.ReadAccountUserList(ctx.AccountID); err != nil {
		utils.ErrorHappened(ctx, "failed to retrieve the users (1)", http.StatusInternalServerError, err)
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

	utils.WriteResponse(ctx, response, http.StatusOK, 10)
}

// loginAccountUser handles the requests to login the user in the system
// Request: POST /account/user/login
func loginAccountUser(ctx *context.Context) {
	var (
		loginPayload = &entity.LoginPayload{}
		account      *entity.Account
		user         *entity.AccountUser
		sessionToken string
		err          error
	)

	if err = json.NewDecoder(ctx.Body).Decode(loginPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err := validator.IsValidLoginPayload(loginPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if loginPayload.Email != "" {
		account, user, err = core.FindAccountAndUserByEmail(loginPayload.Email)
		if err != nil {
			utils.ErrorHappened(ctx, "failed to login the user (3)\n"+err.Error(), http.StatusBadRequest, err)
			return
		}
	}

	if loginPayload.Username != "" {
		account, user, err = core.FindAccountAndUserByUsername(loginPayload.Username)
		if err != nil {
			utils.ErrorHappened(ctx, "failed to login the user (4)\n"+err.Error(), http.StatusBadRequest, err)
			return
		}
	}

	if account == nil || user == nil {
		utils.ErrorHappened(ctx, "failed to login the user (5)", http.StatusInternalServerError, fmt.Errorf("account or user nil on login"))
		return
	}

	if err = validator.AccountUserCredentialsValid(loginPayload.Password, user); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (6)", http.StatusUnauthorized, err)
		return
	}

	if sessionToken, err = core.CreateAccountUserSession(user); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (7)", http.StatusInternalServerError, err)
		return
	}

	user.LastLogin = time.Now()
	_, err = core.UpdateAccountUser(*user, *user, false)

	utils.WriteResponse(ctx, struct {
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
	}{
		AccountToken: account.AuthToken,
		Token:        sessionToken,
	}, http.StatusCreated, 0)
}

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: Post /account/:AccountID/application/:ApplicationID/user/refreshsession
func refreshAccountUserSession(ctx *context.Context) {
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
		err          error
	)

	if err = json.NewDecoder(ctx.Body).Decode(&tokenPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to refresh session token (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if ctx.SessionToken != tokenPayload.Token {
		utils.ErrorHappened(ctx, "failed to refresh session token (2) \nsession token mismatch", http.StatusBadRequest, fmt.Errorf("session token mismatch"))
		return
	}

	if sessionToken, err = core.RefreshAccountUserSession(ctx.SessionToken, ctx.AccountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to refresh session token (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// logoutApplicationUser handles the requests to logout the user from the system
// Request: Post /account/:AccountID/application/:ApplicationID/user/logout
func logoutAccountUser(ctx *context.Context) {
	var (
		logoutPayload struct {
			Token string `json:"token"`
		}
		err error
	)

	if err = json.NewDecoder(ctx.Body).Decode(&logoutPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to logout the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = validator.AccountUserCredentialsValid(logoutPayload.Token, ctx.AccountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to logout the user (2)", http.StatusUnauthorized, err)
		return
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, ctx.AccountUser); err != nil {
		utils.ErrorHappened(ctx, "failed to logout the user (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "logged out", http.StatusOK, 0)
}
