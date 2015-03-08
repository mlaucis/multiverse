/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:ID
// Test with: curl -i localhost/0.1/account/:AccountID/user/:ID
func getAccountUser(ctx *context) {
	var (
		accountID   int64
		userID      int64
		accountUser *entity.AccountUser
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if accountUser, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, accountUser, http.StatusOK, 10)
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' -X PUT localhost/0.1/account/:AccountID/user/:ID
func updateAccountUser(ctx *context) {
	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		userID      int64
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest)
		return
	}

	if accountUser.ID == 0 {
		accountUser.ID = userID
	}
	if accountUser.AccountID == 0 {
		accountUser.AccountID = accountID
	}

	if err = validator.UpdateAccountUser(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest)
		return
	}

	if accountUser, err = core.UpdateAccountUser(accountUser, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, accountUser, http.StatusCreated, 0)
}

// deleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:ID
// Test with: curl -i -X DELETE localhost/0.1/account/:AccountID/user/:ID
func deleteAccountUser(ctx *context) {
	var (
		accountID int64
		userID    int64
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if err = core.DeleteAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
// Test with: curl -i localhost/0.1/account/:AccountID/users
func getAccountUserList(ctx *context) {
	var (
		accountID    int64
		account      *entity.Account
		accountUsers []*entity.AccountUser
		err          error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError)
		return
	}

	if accountUsers, err = core.ReadAccountUserList(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError)
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

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/0.1/account/:AccountID/users
func createAccountUser(ctx *context) {
	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, fmt.Sprintf("accountId is not set or the value is incorrect %v", ctx.vars["accountId"]), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest)
		return
	}

	accountUser.AccountID = accountID

	if err = validator.CreateAccountUser(accountUser); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusBadRequest)
		return
	}

	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%q", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, accountUser, http.StatusCreated, 0)
}

// loginAccountUser handles the requests to login the user in the system
func loginAccountUser(ctx *context) {
	var (
		loginPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		sessionToken string
		err          error
	)

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(loginPayload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if !validator.IsValidEmail(loginPayload.Email) {
		errorHappened(ctx, "invalid e-mail", http.StatusBadRequest)
		return
	}

	_, user, err := core.FindAccountAndUserByEmail(loginPayload.Email)
	if err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validator.AccountUserCredentialsValid(loginPayload.Password, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	if sessionToken, err = core.CreateAccountUserSession(user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// refreshApplicationUserSession handles the requests to refresh the user session token
func refreshAccountUserSession(ctx *context) {
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

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(payload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if sessionToken, err = core.RefreshAccountUserSession(payload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// logoutApplicationUser handles the requests to logout the user from the system
func logoutAccountUser(ctx *context) {
	var (
		logoutPayload struct {
			Token string `json:"token"`
		}
		user      = &entity.AccountUser{}
		accountID int64
		userID    int64
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(logoutPayload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if err = validator.AccountUserCredentialsValid(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	if err = core.DestroyAccountUserSession(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
