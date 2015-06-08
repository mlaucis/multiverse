/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	accountUser struct {
		storage core.AccountUser
	}
)

func (accUser *accountUser) Read(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	server.WriteResponse(ctx, ctx.Bag["accountUser"].(*entity.AccountUser), http.StatusOK, 10)
	return
}

func (accUser *accountUser) Update(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	accountUser := *(ctx.Bag["accountUser"].(*entity.AccountUser))
	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return []errors.Error{errors.NewBadRequestError("failed to update the account user (1)\n"+er.Error(), er.Error())}
	}

	accountUser.ID = ctx.Bag["accountUserID"].(int64)
	accountUser.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.UpdateAccountUser(accUser.storage, ctx.Bag["accountUser"].(*entity.AccountUser), &accountUser); err != nil {
		return
	}

	updatedAccountUser, err := accUser.storage.Update(*(ctx.Bag["accountUser"].(*entity.AccountUser)), accountUser, true)
	if err != nil {
		return
	}

	updatedAccountUser.Password = ""
	server.WriteResponse(ctx, updatedAccountUser, http.StatusCreated, 0)
	return
}

func (accUser *accountUser) Delete(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	if err = accUser.storage.Delete(ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (accUser *accountUser) Create(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	var accountUser = &entity.AccountUser{}

	if err := json.Unmarshal(ctx.Body, accountUser); err != nil {
		return []errors.Error{errors.NewBadRequestError("failed to create the account user (1)"+err.Error(), err.Error())}
	}

	accountUser.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.CreateAccountUser(accUser.storage, accountUser); err != nil {
		return
	}

	if accountUser, err = accUser.storage.Create(accountUser, true); err != nil {
		return
	}

	accountUser.Password = ""

	server.WriteResponse(ctx, accountUser, http.StatusCreated, 0)
	return
}

func (accUser *accountUser) List(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	var (
		accountUsers []*entity.AccountUser
	)

	if accountUsers, err = accUser.storage.List(ctx.Bag["accountID"].(int64)); err != nil {
		//		utils.ErrorHappened(ctx, "failed to retrieve the users (1)", http.StatusInternalServerError, err)
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

	server.WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

func (accUser *accountUser) Login(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	var (
		loginPayload = &entity.LoginPayload{}
		account      *entity.Account
		user         *entity.AccountUser
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return []errors.Error{errors.NewBadRequestError("failed to login the user (1)\n"+er.Error(), er.Error())}
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		return
	}

	if loginPayload.Email != "" {
		account, user, err = accUser.storage.FindByEmail(loginPayload.Email)
		if err != nil {
			return
		}
	}

	if loginPayload.Username != "" {
		account, user, err = accUser.storage.FindByUsername(loginPayload.Username)
		if err != nil {
			return
		}
	}

	if err = validator.AccountUserCredentialsValid(loginPayload.Password, user); err != nil {
		return
	}

	if sessionToken, err = accUser.storage.CreateSession(user); err != nil {
		return
	}

	timeNow := time.Now()
	user.LastLogin = &timeNow
	_, err = accUser.storage.Update(*user, *user, false)

	server.WriteResponse(ctx, struct {
		ID           int64  `json:"id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccountToken: account.AuthToken,
		Token:        sessionToken,
	}, http.StatusCreated, 0)
	return
}

func (accUser *accountUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
	)

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return []errors.Error{errors.NewBadRequestError("failed to refresh session token (1)\n"+er.Error(), er.Error())}
	}

	if ctx.SessionToken != tokenPayload.Token {
		return []errors.Error{errors.NewBadRequestError("failed to refresh session token (2) \nsession token mismatch", "session token mismatch")}
	}

	if sessionToken, err = accUser.storage.RefreshSession(ctx.SessionToken, ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

func (accUser *accountUser) Logout(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	var logoutPayload struct {
		Token string `json:"token"`
	}

	if er := json.Unmarshal(ctx.Body, &logoutPayload); er != nil {
		return []errors.Error{errors.NewBadRequestError("failed to logout the user (1)\n"+er.Error(), er.Error())}
	}

	if ctx.SessionToken != logoutPayload.Token {
		return []errors.Error{errors.NewBadRequestError("failed to logout the user (2) \nsession token mismatch", "session token mismatch")}
	}

	if err = accUser.storage.DestroySession(logoutPayload.Token, ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, "logged out", http.StatusOK, 0)
	return
}

// PopulateContext adds the accountUser to the context
func (accUser *accountUser) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return deprecatedStorageError
	ctx.Bag["accountUser"], err = accUser.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["accountUserID"].(int64))
	return
}

// NewAccountUser creates a new Account Route handler
func NewAccountUser(storage core.AccountUser) server.AccountUser {
	return &accountUser{
		storage: storage,
	}
}
