/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getApplicationUser handles requests to retrieve a single user
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID
func getApplicationUser(ctx *context.Context) {
	// Don't return the password to the users
	user := ctx.ApplicationUser
	user.Password = ""

	writeResponse(ctx, user, http.StatusOK, 10)
}

// updateApplicationUser handles requests to update a user
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserID
func updateApplicationUser(ctx *context.Context) {
	var (
		user = &entity.User{}
		err  error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	user.ID = ctx.ApplicationUserID
	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

	if err = validator.UpdateUser(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if ctx.ApplicationUser.Email != user.Email {
		if isDuplicate, err := validator.DuplicateApplicationUserEmail(user.AccountID, user.ApplicationID, user.Email); isDuplicate || err != nil {
			if isDuplicate {
				errorHappened(ctx, "email address already in use", http.StatusBadRequest, fmt.Errorf("duplicate email address on update"))
			} else if err != nil {
				errorHappened(ctx, "unexpected error", http.StatusBadRequest, err)
			}

			return
		}
	}

	if ctx.ApplicationUser.Username != user.Username {
		if isDuplicate, err := validator.DuplicateApplicationUserUsername(user.AccountID, user.ApplicationID, user.Username); isDuplicate || err != nil {
			if isDuplicate {
				errorHappened(ctx, "username already in use", http.StatusBadRequest, fmt.Errorf("duplicate username on update"))
			} else if err != nil {
				errorHappened(ctx, "unexpected error", http.StatusBadRequest, err)
			}

			return
		}
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
	if err := core.DeleteUser(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createApplicationUser handles requests to create a user
// Request: POST account/:AccountID/application/:ApplicationID/users
func createApplicationUser(ctx *context.Context) {
	var (
		user = &entity.User{}
		err  error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(user); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

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
/*func getApplicationUserList(ctx *context.Context) {
	var (
		users []*entity.User
		err   error
	)

	if users, err = core.ReadUserList(ctx.AccountID, ctx.ApplicationID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		Users         []*entity.User
	}{
		ApplicationID: ctx.ApplicationID,
		Users:         users,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}*/

// loginApplicationUser handles the requests to login the user in the system
// Request: POST account/:AccountID/application/:ApplicationID/user/login
func loginApplicationUser(ctx *context.Context) {
	var (
		loginPayload = &entity.LoginPayload{}
		user         *entity.User
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
		user, err = core.FindApplicationUserByEmail(ctx.AccountID, ctx.ApplicationID, loginPayload.Email)
		if err != nil {
			errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
			return
		}
	}

	if loginPayload.Username != "" {
		user, err = core.FindApplicationUserByUsername(ctx.AccountID, ctx.ApplicationID, loginPayload.Username)
		if err != nil {
			errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
			return
		}
	}

	if user == nil {
		errorHappened(ctx, "unexpected error happened", http.StatusInternalServerError, fmt.Errorf("user nil on login"))
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

	user.LastLogin = time.Now()
	_, err = core.UpdateUser(user, false)

	writeResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: POST account/:AccountID/application/:ApplicationID/user/refreshsession
func refreshApplicationUserSession(ctx *context.Context) {
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
		err          error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(tokenPayload); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if tokenPayload.Token != ctx.SessionToken {
		errorHappened(ctx, "session token mismatch", http.StatusBadRequest, fmt.Errorf("session token mismatch"))
		return
	}

	if sessionToken, err = core.RefreshApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
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
	if err := core.DestroyApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "logged out", http.StatusOK, 0)
}
