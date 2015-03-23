/**
 * @author Florin Patan <florinpatan@gmail.com>
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

// getApplicationUser handles requests to retrieve a single user
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID
func getApplicationUser(ctx *context.Context) {
	utils.WriteResponse(ctx, ctx.ApplicationUser, http.StatusOK, 10)
}

// updateApplicationUser handles requests to update a user
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserID
func updateApplicationUser(ctx *context.Context) {
	var (
		err error
	)

	user := *ctx.ApplicationUser
	if err = json.NewDecoder(ctx.Body).Decode(&user); err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	user.ID = ctx.ApplicationUserID
	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

	if err = validator.UpdateUser(ctx.ApplicationUser, &user); err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedUser, err := core.UpdateUser(*ctx.ApplicationUser, user, true)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the user (3)", http.StatusInternalServerError, err)
		return
	}

	updatedUser.Password = ""

	utils.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
}

// deleteApplicationUser handles requests to delete a single user
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserID
func deleteApplicationUser(ctx *context.Context) {
	if err := core.DeleteUser(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the user (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "", http.StatusNoContent, 10)
}

// createApplicationUser handles requests to create a user
// Request: POST account/:AccountID/application/:ApplicationID/users
func createApplicationUser(ctx *context.Context) {
	var (
		user = &entity.User{}
		err  error
	)

	if err = json.NewDecoder(ctx.Body).Decode(user); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

	if err = validator.CreateUser(user); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		utils.ErrorHappened(ctx, "failed to create the user (3)", http.StatusInternalServerError, err)
		return
	}

	user.Password = ""

	utils.WriteResponse(ctx, user, http.StatusCreated, 0)
}

// loginApplicationUser handles the requests to login the user in the system
// Request: POST account/:AccountID/application/:ApplicationID/user/login
func loginApplicationUser(ctx *context.Context) {
	var (
		loginPayload = &entity.LoginPayload{}
		user         *entity.User
		sessionToken string
		err          error
	)

	if err = json.NewDecoder(ctx.Body).Decode(loginPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err := validator.IsValidLoginPayload(loginPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (2)"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if loginPayload.Email != "" {
		user, err = core.FindApplicationUserByEmail(ctx.AccountID, ctx.ApplicationID, loginPayload.Email)
		if err != nil {
			utils.ErrorHappened(ctx, "failed to login the user (3)", http.StatusInternalServerError, err)
			return
		}
	}

	if loginPayload.Username != "" {
		user, err = core.FindApplicationUserByUsername(ctx.AccountID, ctx.ApplicationID, loginPayload.Username)
		if err != nil {
			utils.ErrorHappened(ctx, "failed to login the user (4)", http.StatusInternalServerError, err)
			return
		}
	}

	if user == nil {
		utils.ErrorHappened(ctx, "failed to login the user (5)", http.StatusInternalServerError, fmt.Errorf("user nil on login"))
		return
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (6)", http.StatusUnauthorized, err)
		return
	}

	if sessionToken, err = core.CreateApplicationUserSession(user); err != nil {
		utils.ErrorHappened(ctx, "failed to login the user (7)", http.StatusInternalServerError, err)
		return
	}

	user.LastLogin = time.Now()
	_, err = core.UpdateUser(*user, *user, false)

	utils.WriteResponse(ctx, struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{
		UserID: user.ID,
		Token:  sessionToken,
	}, http.StatusCreated, 0)
}

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: POST account/:AccountID/application/:ApplicationID/user/refreshsession
func refreshApplicationUserSession(ctx *context.Context) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		sessionToken string
		err          error
	)

	if err = json.NewDecoder(ctx.Body).Decode(&tokenPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to refresh the user session (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if tokenPayload.Token != ctx.SessionToken {
		utils.ErrorHappened(ctx, "failed to refresh the session token (2)\nsession token mismatch", http.StatusBadRequest, fmt.Errorf("session token mismatch"))
		return
	}

	if sessionToken, err = core.RefreshApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
		utils.ErrorHappened(ctx, "failed to refresh session token (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, struct {
		Token string `json:"session_token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
}

// logoutApplicationUser handles the requests to logout the user from the system
// Request: POST account/:AccountID/application/:ApplicationID/user/logout
func logoutApplicationUser(ctx *context.Context) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		err error
	)

	if err = json.NewDecoder(ctx.Body).Decode(&tokenPayload); err != nil {
		utils.ErrorHappened(ctx, "failed to logout user (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if tokenPayload.Token != ctx.SessionToken {
		utils.ErrorHappened(ctx, "failed to logout user (2)\nsession token mismatch", http.StatusBadRequest, fmt.Errorf("session token mismatch"))
		return
	}

	if err := core.DestroyApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
		utils.ErrorHappened(ctx, "failed to logout the user (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "logged out", http.StatusOK, 0)
}
