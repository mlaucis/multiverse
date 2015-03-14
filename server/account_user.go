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
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
func getAccountUser(ctx *context.Context) {
	writeResponse(ctx, ctx.AccountUser, http.StatusOK, 10)
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:UserID
func updateAccountUser(ctx *context.Context) {
	var (
		accountUser = &entity.AccountUser{}
		err         error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	accountUser.ID = ctx.AccountUserID
	accountUser.AccountID = ctx.AccountID

	if err = validator.UpdateAccountUser(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	if ctx.AccountUser.Email != accountUser.Email {
		if isDuplicate, err := validator.DuplicateAccountUserEmail(accountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errorHappened(ctx, "email address already in use", http.StatusBadRequest, fmt.Errorf("duplicate email address on update"))
			} else if err != nil {
				errorHappened(ctx, "unexpected error", http.StatusBadRequest, err)
			}

			return
		}
	}

	if ctx.AccountUser.Username != accountUser.Username {
		if isDuplicate, err := validator.DuplicateAccountUserUsername(accountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errorHappened(ctx, "username already in use", http.StatusBadRequest, fmt.Errorf("duplicate username on update"))
			} else if err != nil {
				errorHappened(ctx, "unexpected error", http.StatusBadRequest, err)
			}

			return
		}
	}

	if accountUser, err = core.UpdateAccountUser(accountUser, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, accountUser, http.StatusCreated, 0)
}

// deleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:UserID
func deleteAccountUser(ctx *context.Context) {
	if err := core.DeleteAccountUser(ctx.AccountID, ctx.AccountUserID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
func createAccountUser(ctx *context.Context) {
	var (
		accountUser = &entity.AccountUser{}
		err         error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	accountUser.AccountID = ctx.AccountID

	if err = validator.CreateAccountUser(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, accountUser, http.StatusCreated, 0)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
func getAccountUserList(ctx *context.Context) {
	var (
		accountUsers []*entity.AccountUser
		err          error
	)

	if accountUsers, err = core.ReadAccountUserList(ctx.AccountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
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

	writeResponse(ctx, response, http.StatusOK, 10)
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

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(loginPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if err := validator.IsValidLoginPayload(loginPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if loginPayload.Email != "" {
		account, user, err = core.FindAccountAndUserByEmail(loginPayload.Email)
		if err != nil {
			errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
			return
		}
	}

	if loginPayload.Username != "" {
		account, user, err = core.FindAccountAndUserByUsername(loginPayload.Username)
		if err != nil {
			errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
			return
		}
	}

	if account == nil || user == nil {
		errorHappened(ctx, "unexpected error happened", http.StatusInternalServerError, fmt.Errorf("account or user nil on login"))
		return
	}

	if err = validator.AccountUserCredentialsValid(loginPayload.Password, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusUnauthorized, err)
		return
	}

	if sessionToken, err = core.CreateAccountUserSession(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	user.LastLogin = time.Now()
	_, err = core.UpdateAccountUser(user, false)

	writeResponse(ctx, struct {
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
		payload struct {
			Token string `json:"token"`
		}
		sessionToken string
		err          error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&payload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if sessionToken, err = core.RefreshAccountUserSession(payload.Token, ctx.AccountUser); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, struct {
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

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&logoutPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = validator.AccountUserCredentialsValid(logoutPayload.Token, ctx.AccountUser); err != nil {
		errorHappened(ctx, err.Error(), http.StatusUnauthorized, err)
		return
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, ctx.AccountUser); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
