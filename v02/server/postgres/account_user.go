/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	accountUser struct {
		storage core.AccountUser
	}
)

func (accUser *accountUser) Read(ctx *context.Context) (err []errors.Error) {
	// TODO This one read only the current account user maybe we want to have something to read any account user?
	computeAccountUserLastModified(ctx, ctx.Bag["accountUser"].(*entity.AccountUser))
	server.WriteResponse(ctx, ctx.Bag["accountUser"].(*entity.AccountUser), http.StatusOK, 10)
	return
}

func (accUser *accountUser) Update(ctx *context.Context) (err []errors.Error) {
	accountUser := *(ctx.Bag["accountUser"].(*entity.AccountUser))

	if accountUser.PublicID != ctx.Vars["accountUserID"] {
		return []errors.Error{errmsg.AccountUserMismatchError}
	}

	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
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
	if ctx.R.Header.Get("X-Jarvis-Auth") != "ZTBmZjI3MGE2M2YzYzAzOWI1MjhiYTNi" {
		return []errors.Error{errmsg.MissingJarvisIDError}
	}

	accountUserID := ctx.Vars["accountUserID"]
	if !validator.IsValidUUID5(accountUserID) {
		return []errors.Error{errmsg.InvalidUserIDError}
	}
	accountUser, err := accUser.storage.FindByPublicID(ctx.Bag["accountID"].(int64), accountUserID)
	if err != nil {
		return
	}

	if err = accUser.storage.Delete(accountUser); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (accUser *accountUser) Create(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("X-Jarvis-Auth") != "ZTBmZjI3MGE2M2YzYzAzOWI1MjhiYTNi" {
		return []errors.Error{errmsg.MissingJarvisIDError}
	}

	var accountUser = &entity.AccountUser{}

	if err := json.Unmarshal(ctx.Body, accountUser); err != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(err.Error())}
	}

	accountUser.AccountID = ctx.Bag["accountID"].(int64)
	accountUser.PublicAccountID = ctx.Bag["account"].(*entity.Account).PublicID

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
	var (
		accountUsers []*entity.AccountUser
	)

	if accountUsers, err = accUser.storage.List(ctx.Bag["accountID"].(int64)); err != nil {
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

	computeAccountUsersLastModified(ctx, response.AccountUsers)

	server.WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

func (accUser *accountUser) Login(ctx *context.Context) (err []errors.Error) {
	var (
		loginPayload = &entity.LoginPayload{}
		account      *entity.Account
		user         *entity.AccountUser
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
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

	if account == nil || user == nil {
		return []errors.Error{errmsg.AccountUserNotFoundError}
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
		ID           string `json:"id"`
		AccountID    string `json:"account_id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{
		ID:           user.PublicID,
		AccountID:    user.PublicAccountID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccountToken: account.AuthToken,
		Token:        sessionToken,
	}, http.StatusCreated, 0)
	return
}

func (accUser *accountUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
	)

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
	}

	if ctx.SessionToken != tokenPayload.Token {
		return []errors.Error{errmsg.SessionTokenMismatchError}
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
	if err = accUser.storage.DestroySession(ctx.SessionToken, ctx.Bag["accountUser"].(*entity.AccountUser)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

// PopulateContext adds the accountUser to the context
func (accUser *accountUser) PopulateContext(ctx *context.Context) (err []errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.InvalidAccountUserCredentialsError.UpdateInternalMessage(fmt.Sprintf("got %s:%s", user, pass))}
	}
	accountUser, err := accUser.storage.FindBySession(pass)
	if accountUser == nil {
		return []errors.Error{errmsg.AccountUserNotFoundError}
	}
	if err == nil {
		ctx.Bag["accountUser"] = accountUser
		ctx.Bag["accountUserID"] = accountUser.ID
		ctx.SessionToken = pass
	}
	return
}

// NewAccountUser creates a new Account Route handler
func NewAccountUser(storage core.AccountUser) server.AccountUser {
	return &accountUser{
		storage: storage,
	}
}
