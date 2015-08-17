package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type (
	member struct {
		storage core.Member
	}
)

func (user *member) Read(ctx *context.Context) (err []errors.Error) {
	// TODO This one read only the current account user maybe we want to have something to read any account user?
	accountUser := ctx.Bag["accountUser"].(*entity.Member)
	response.SanitizeMember(accountUser)
	response.ComputeMemberLastModified(ctx, accountUser)
	response.WriteResponse(ctx, accountUser, http.StatusOK, 10)
	return
}

func (user *member) Update(ctx *context.Context) (err []errors.Error) {
	accountUser := *(ctx.Bag["accountUser"].(*entity.Member))

	if accountUser.PublicID != ctx.Vars["accountUserID"] {
		return []errors.Error{errmsg.ErrMemberMismatchErr}
	}

	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	accountUser.ID = ctx.Bag["accountUserID"].(int64)
	accountUser.OrgID = ctx.Bag["accountID"].(int64)

	if err = validator.UpdateMember(user.storage, ctx.Bag["accountUser"].(*entity.Member), &accountUser); err != nil {
		return
	}

	updatedAccountUser, err := user.storage.Update(*(ctx.Bag["accountUser"].(*entity.Member)), accountUser, true)
	if err != nil {
		return
	}

	updatedAccountUser.Password = ""
	response.WriteResponse(ctx, updatedAccountUser, http.StatusCreated, 0)
	return
}

func (user *member) Delete(ctx *context.Context) (err []errors.Error) {
	accountUserID := ctx.Vars["accountUserID"]
	if !validator.IsValidUUID5(accountUserID) {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}
	accountUser, err := user.storage.FindByPublicID(ctx.Bag["accountID"].(int64), accountUserID)
	if err != nil {
		return
	}

	if err = user.storage.Delete(accountUser); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (user *member) Create(ctx *context.Context) (err []errors.Error) {
	var accountUser = &entity.Member{}

	if err := json.Unmarshal(ctx.Body, accountUser); err != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(err.Error())}
	}

	accountUser.OrgID = ctx.Bag["accountID"].(int64)
	accountUser.PublicAccountID = ctx.Bag["account"].(*entity.Organization).PublicID

	if err = validator.CreateMember(user.storage, accountUser); err != nil {
		return
	}

	if accountUser, err = user.storage.Create(accountUser, true); err != nil {
		return
	}

	accountUser.Password = ""

	response.WriteResponse(ctx, accountUser, http.StatusCreated, 0)
	return
}

func (user *member) List(ctx *context.Context) (err []errors.Error) {
	var (
		accountUsers []*entity.Member
	)

	if accountUsers, err = user.storage.List(ctx.Bag["accountID"].(int64)); err != nil {
		return
	}

	for _, accountUser := range accountUsers {
		response.SanitizeMember(accountUser)
	}

	resp := &struct {
		AccountUsers []*entity.Member `json:"accountUsers"`
	}{
		AccountUsers: accountUsers,
	}

	response.ComputeMembersLastModified(ctx, resp.AccountUsers)

	response.WriteResponse(ctx, resp, http.StatusOK, 10)
	return
}

func (user *member) Login(ctx *context.Context) (err []errors.Error) {
	var (
		loginPayload = &entity.LoginPayload{}
		account      *entity.Organization
		usr         *entity.Member
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if err = validator.IsValidLoginPayload(loginPayload); err != nil {
		return
	}

	if loginPayload.Email != "" {
		account, usr, err = user.storage.FindByEmail(loginPayload.Email)
		if err != nil {
			return
		}
	}

	if loginPayload.Username != "" {
		account, usr, err = user.storage.FindByUsername(loginPayload.Username)
		if err != nil {
			return
		}
	}

	if account == nil || usr == nil || !usr.Enabled {
		return []errors.Error{errmsg.ErrMemberNotFound}
	}

	if err = validator.MemberCredentialsValid(loginPayload.Password, usr); err != nil {
		return
	}

	if sessionToken, err = user.storage.CreateSession(usr); err != nil {
		return
	}

	timeNow := time.Now()
	usr.LastLogin = &timeNow
	_, err = user.storage.Update(*usr, *usr, false)

	response.WriteResponse(ctx, struct {
		ID           string `json:"id"`
		AccountID    string `json:"account_id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{
		ID:           usr.PublicID,
		AccountID:    usr.PublicAccountID,
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		AccountToken: account.AuthToken,
		Token:        sessionToken,
	}, http.StatusCreated, 0)
	return
}

func (user *member) RefreshSession(ctx *context.Context) (err []errors.Error) {
	var (
		tokenPayload struct {
			Token string `json:"token"`
		}
		sessionToken string
	)

	if er := json.Unmarshal(ctx.Body, &tokenPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if ctx.SessionToken != tokenPayload.Token {
		return []errors.Error{errmsg.ErrAuthSessionTokenMismatch}
	}

	if sessionToken, err = user.storage.RefreshSession(ctx.SessionToken, ctx.Bag["accountUser"].(*entity.Member)); err != nil {
		return
	}

	response.WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

func (user *member) Logout(ctx *context.Context) (err []errors.Error) {
	if err = user.storage.DestroySession(ctx.SessionToken, ctx.Bag["accountUser"].(*entity.Member)); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

// PopulateContext adds the accountUser to the context
func (user *member) PopulateContext(ctx *context.Context) (err []errors.Error) {
	usr, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidAccountUserCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", usr, pass))}
	}
	accountUser, err := user.storage.FindBySession(pass)
	if accountUser == nil {
		return []errors.Error{errmsg.ErrMemberNotFound}
	}
	if err == nil {
		ctx.Bag["accountUser"] = accountUser
		ctx.Bag["accountUserID"] = accountUser.ID
		ctx.SessionToken = pass
	}
	return
}

// NewMember creates a new member route handler
func NewMember(storage core.Member) handlers.Member {
	return &member{
		storage: storage,
	}
}
