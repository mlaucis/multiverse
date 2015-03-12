/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getApplicationUser handles requests to retrieve a single user
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID
func getApplicationUser(ctx *context.Context) {
	var (
		user          *entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadApplicationUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		*entity.User
	}{
		ApplicationID: applicationID,
		User:          user,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// updateApplicationUser handles requests to update a user
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserID
func updateApplicationUser(ctx *context.Context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if user.ID == 0 {
		user.ID = userID
	}
	if user.AccountID == 0 {
		user.AccountID = accountID
	}
	if user.ApplicationID == 0 {
		user.ApplicationID = applicationID
	}

	if err = validator.UpdateUser(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if user, err = core.UpdateUser(user, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(ctx, user, http.StatusCreated, 0)
}

// deleteApplicationUser handles requests to delete a single user
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserID
func deleteApplicationUser(ctx *context.Context) {
	var (
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createApplicationUser handles requests to create a user
// Request: POST account/:AccountID/application/:ApplicationID/users
func createApplicationUser(ctx *context.Context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	user.AccountID = accountID
	user.ApplicationID = applicationID

	if err = validator.CreateUser(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(ctx, user, http.StatusCreated, 0)
}

// getApplicationUserList handles requests to retrieve all users of an app
// THIS ROUTE IS NOT YET ACTIVATED
// Request: GET account/:AccountID/application/:ApplicationID/users
func getApplicationUserList(ctx *context.Context) {
	var (
		accountID     int64
		applicationID int64
		users         []*entity.User
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if users, err = core.ReadUserList(accountID, applicationID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	// TODO iterate the users and strip their password

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		Users         []*entity.User
	}{
		ApplicationID: applicationID,
		Users:         users,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// loginApplicationUser handles the requests to login the user in the system
// Request: POST account/:AccountID/application/:ApplicationID/user/login
func loginApplicationUser(ctx *context.Context) {
	var (
		loginPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		accountID     int64
		applicationID int64
		sessionToken  string
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(&loginPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if !validator.IsValidEmail(loginPayload.Email) {
		errorHappened(ctx, "invalid e-mail", http.StatusBadRequest, err)
		return
	}

	user, err := core.FindApplicationUserByEmail(accountID, applicationID, loginPayload.Email)
	if err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusUnauthorized, err)
		return
	}

	if sessionToken, err = core.CreateApplicationUserSession(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: POST account/:AccountID/application/:ApplicationID/user/refreshsession
func refreshApplicationUserSession(ctx *context.Context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		sessionToken  string
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadApplicationUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	if sessionToken, err = core.RefreshApplicationUserSession(ctx.SessionToken, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// logoutApplicationUser handles the requests to logout the user from the system
// Request: POST account/:AccountID/application/:ApplicationID/user/logout
func logoutApplicationUser(ctx *context.Context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if user, err = core.ReadApplicationUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	if err = core.DestroyApplicationUserSession(ctx.SessionToken, user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
