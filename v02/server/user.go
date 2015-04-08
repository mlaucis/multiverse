/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/context"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

// getApplicationUser handles requests to retrieve a single user
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID
func GetApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	WriteResponse(ctx, ctx.ApplicationUser, http.StatusOK, 10)
	return
}

// updateApplicationUser handles requests to update a user
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserID
func UpdateApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	user := *ctx.ApplicationUser
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return tgerrors.NewBadRequestError("failed to update the user (1)\n"+er.Error(), er.Error())
	}

	user.ID = ctx.ApplicationUserID
	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

	if err = validator.UpdateUser(ctx.ApplicationUser, &user); err != nil {
		return
	}

	updatedUser, err := core.UpdateUser(*ctx.ApplicationUser, user, true)
	if err != nil {
		return
	}

	updatedUser.Password = ""

	WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

// deleteApplicationUser handles requests to delete a single user
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserID
func DeleteApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	if err = core.DeleteUser(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// createApplicationUser handles requests to create a user
// Request: POST account/:AccountID/application/:ApplicationID/users
func CreateApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		user = &entity.User{}
		er   error
	)

	if er = json.Unmarshal(ctx.Body, user); er != nil {
		return tgerrors.NewBadRequestError("failed to create the application user (1)\n"+er.Error(), er.Error())
	}

	user.AccountID = ctx.AccountID
	user.ApplicationID = ctx.ApplicationID

	if err = validator.CreateUser(user); err != nil {
		return
	}

	if user, err = core.WriteUser(user, true); err != nil {
		return
	}

	user.Password = ""

	WriteResponse(ctx, user, http.StatusCreated, 0)
	return
}

// loginApplicationUser handles the requests to login the user in the system
// Request: POST account/:AccountID/application/:ApplicationID/user/login
func LoginApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		loginPayload = &entity.LoginPayload{}
		user         *entity.User
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
		user, err = core.FindApplicationUserByEmail(ctx.AccountID, ctx.ApplicationID, loginPayload.Email)
		if err != nil {
			return
		}
	}

	if loginPayload.Username != "" {
		user, err = core.FindApplicationUserByUsername(ctx.AccountID, ctx.ApplicationID, loginPayload.Username)
		if err != nil {
			return
		}
	}

	if user == nil {
		return tgerrors.NewInternalError("failed to login the application user (2)\n", "user is nil")
	}

	if !user.Enabled {
		return tgerrors.NewNotFoundError("failed to login the user (3)\nuser is disabled", "user is disabled")
	}

	if err = validator.ApplicationUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = core.CreateApplicationUserSession(user); err != nil {
		return
	}

	user.LastLogin = time.Now()
	_, err = core.UpdateUser(*user, *user, false)
	if err != nil {
		return
	}

	WriteResponse(ctx, struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{
		UserID: user.ID,
		Token:  sessionToken,
	}, http.StatusCreated, 0)
	return
}

// refreshApplicationUserSession handles the requests to refresh the user session token
// Request: POST account/:AccountID/application/:ApplicationID/user/refreshsession
func RefreshApplicationUserSession(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to refresh the session token (1)\n"+er.Error(), er.Error())
	}

	if tokenPayload.Token != ctx.SessionToken {
		return tgerrors.NewBadRequestError("failed to refresh the session token (2)\nsession token mismatch", "session token mismatch")
	}

	if sessionToken, err = core.RefreshApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
		return
	}

	WriteResponse(ctx, struct {
		Token string `json:"session_token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

// logoutApplicationUser handles the requests to logout the user from the system
// Request: POST account/:AccountID/application/:ApplicationID/user/logout
func LogoutApplicationUser(ctx *context.Context) (err *tgerrors.TGError) {
	var (
		tokenPayload struct {
			Token string `json:"session_token"`
		}
		er error
	)

	if er = json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return tgerrors.NewBadRequestError("failed to logout the user (1)\n"+er.Error(), er.Error())
	}

	if tokenPayload.Token != ctx.SessionToken {
		return tgerrors.NewBadRequestError("failed to logout the user (2)\nsession token mismatch", "session token mismatch")
	}

	if err = core.DestroyApplicationUserSession(ctx.SessionToken, ctx.ApplicationUser); err != nil {
		return
	}

	WriteResponse(ctx, "logged out", http.StatusOK, 0)
	return
}
