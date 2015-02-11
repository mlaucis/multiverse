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

	"github.com/gorilla/mux"
)

// getUser handles requests to retrieve a single user
// Request: GET /application/:applicationId/user/:ID
// Test with: curl -i localhost/0.1/application/:applicationId/user/:ID
func getUser(w http.ResponseWriter, r *http.Request) {
	var (
		user          *entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
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

	writeResponse(response, http.StatusOK, 10, w, r)
}

// updateUser handles requests to update a user
// Request: PUT /application/:applicationId/user/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}'  -X PUT localhost/0.1/application/:applicationId/user/:ID
func updateUser(w http.ResponseWriter, r *http.Request) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
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
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.UpdateUser(user, true); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(user, http.StatusCreated, 0, w, r)
}

// deleteUser handles requests to delete a single user
// Request: DELETE /application/:applicationId/user/:ID
// Test with: curl -i -X DELETE localhost/0.1/application/:applicationId/user/:ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteUser(accountID, applicationID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getUserList handles requests to retrieve all users of an app
// THIS ROUTE IS NOT YET ACTIVATED
// Request: GET /application/:applicationId/users
// Test with: curl -i localhost/0.1/application/:applicationId/users
func getUserList(w http.ResponseWriter, r *http.Request) {
	var (
		accountID     int64
		applicationID int64
		users         []*entity.User
		err           error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if users, err = core.ReadUserList(accountID, applicationID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
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

	writeResponse(response, http.StatusOK, 10, w, r)
}

// createUser handles requests to create a user
// Request: POST /application/:applicationId/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost/0.1/application/:applicationId/users
func createUser(w http.ResponseWriter, r *http.Request) {
	var (
		user          = &entity.User{}
		accountID     int64
		applicationID int64
		err           error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	user.AccountID = accountID
	user.ApplicationID = applicationID

	if err = validator.CreateUser(user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	// Don't return the password to the users
	user.Password = ""

	writeResponse(user, http.StatusCreated, 0, w, r)
}

// loginUser handles the requests to login the user in the system
func loginUser(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(loginPayload); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if err = validator.UserCredentialsValid(loginPayload.Password, user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusUnauthorized, r, w)
		return
	}

	if sessionToken, err = core.CreateUserSession(user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0, w, r)
}

// refreshUserSession handles the requests to refresh the user session token
func refreshUserSession(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(payload); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if sessionToken, err = core.RefreshUserSession(payload.Token, user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0, w, r)
}

// logoutUser handles the requests to logout the user from the system
func logoutUser(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("applicationId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	if user, err = core.ReadUser(accountID, applicationID, userID); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(logoutPayload); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusBadRequest, r, w)
		return
	}

	if err = validator.UserCredentialsValid(logoutPayload.Token, user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusUnauthorized, r, w)
		return
	}

	if err = core.DestroyUserSession(logoutPayload.Token, user); err != nil {
		errorHappened(fmt.Sprintf("%s", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("logged out", http.StatusOK, 0, w, r)
}
