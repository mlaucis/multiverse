/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
func getAccountUser(ctx *context.Context) {
	var (
		accountID   int64
		userID      int64
		accountUser *entity.AccountUser
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if accountUser, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, accountUser, http.StatusOK, 10)
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:UserID
func updateAccountUser(ctx *context.Context) {
	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		userID      int64
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	if accountUser.ID == 0 {
		accountUser.ID = userID
	}
	if accountUser.AccountID == 0 {
		accountUser.AccountID = accountID
	}

	if err = validator.UpdateAccountUser(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
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
	var (
		accountID int64
		userID    int64
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteAccountUser(accountID, userID); err != nil {
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
		accountID   int64
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, fmt.Sprintf("accountId is not set or the value is incorrect %v", ctx.Vars["accountId"]), http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest, err)
		return
	}

	accountUser.AccountID = accountID

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
		accountID    int64
		account      *entity.Account
		accountUsers []*entity.AccountUser
		err          error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
		return
	}

	if accountUsers, err = core.ReadAccountUserList(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError, err)
		return
	}

	response := &struct {
		entity.Account
		AccountUsers []*entity.AccountUser `json:"accountUsers"`
	}{
		Account:      *account,
		AccountUsers: accountUsers,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// loginAccountUser handles the requests to login the user in the system
// Request: POST /account/user/login
func loginAccountUser(ctx *context.Context) {
	var (
		loginPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		sessionToken string
		err          error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&loginPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if !validator.IsValidEmail(loginPayload.Email) {
		errorHappened(ctx, "invalid e-mail", http.StatusBadRequest, nil)
		return
	}

	account, user, err := core.FindAccountAndUserByEmail(loginPayload.Email)
	if err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
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
		user         = &entity.AccountUser{}
		accountID    int64
		userID       int64
		sessionToken string
		err          error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&payload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if sessionToken, err = core.RefreshAccountUserSession(payload.Token, user); err != nil {
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
		user      = &entity.AccountUser{}
		accountID int64
		userID    int64
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&logoutPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = validator.AccountUserCredentialsValid(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusUnauthorized, err)
		return
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
