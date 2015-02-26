/**
 * @author Florin Patan <florinpatan@gmail.com>
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

// getUser handles requests to retrieve a single user
// Request: GET /application/:applicationId/user/:ID
// Test with: curl -i localhost/0.1/application/:applicationId/user/:ID
func getUser(ctx *context) {
	var (
		user          *entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
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

// updateUser handles requests to update a user
// Request: PUT /application/:applicationId/user/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}'  -X PUT localhost/0.1/application/:applicationId/user/:ID
func updateUser(ctx *context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
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
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if user, err = core.UpdateUser(user, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(ctx, user, http.StatusCreated, 0)
}

// deleteUser handles requests to delete a single user
// Request: DELETE /application/:applicationId/user/:ID
// Test with: curl -i -X DELETE localhost/0.1/application/:applicationId/user/:ID
func deleteUser(ctx *context) {
	var (
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if err = core.DeleteUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// getUserList handles requests to retrieve all users of an app
// THIS ROUTE IS NOT YET ACTIVATED
// Request: GET /application/:applicationId/users
// Test with: curl -i localhost/0.1/application/:applicationId/users
func getUserList(ctx *context) {
	var (
		accountID     int64
		applicationID int64
		users         []*entity.User
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if users, err = core.ReadUserList(accountID, applicationID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
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

// createUser handles requests to create a user
// Request: POST /application/:applicationId/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost/0.1/application/:applicationId/users
func createUser(ctx *context) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	user.AccountID = accountID
	user.ApplicationID = applicationID

	if err = validator.CreateUser(user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(ctx, user, http.StatusCreated, 0)
}

// loginUser handles the requests to login the user in the system
func loginUser(ctx *context) {
	var (
		loginPayload struct {
			Password string `json:"password"`
		}
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		sessionToken  string
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(loginPayload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if err = validator.UserCredentialsValid(loginPayload.Password, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	if sessionToken, err = core.CreateUserSession(user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// refreshUserSession handles the requests to refresh the user session token
func refreshUserSession(ctx *context) {
	var (
		payload struct {
			Token string `json:"token"`
		}
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		sessionToken  string
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(payload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if sessionToken, err = core.RefreshUserSession(payload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// logoutUser handles the requests to logout the user from the system
func logoutUser(ctx *context) {
	var (
		logoutPayload struct {
			Token string `json:"token"`
		}
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if userID, err = strconv.ParseInt(ctx.vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(logoutPayload); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if err = validator.UserCredentialsValid(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusUnauthorized)
		return
	}

	if err = core.DestroyUserSession(logoutPayload.Token, user); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
