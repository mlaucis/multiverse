package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/response"
	"github.com/tapglue/multiverse/v04/validator"
)

type member struct {
	storage core.Member
}

func (user *member) Read(ctx *context.Context) (err []errors.Error) {
	// TODO This one read only the current account user maybe we want to have something to read any account user?
	accountUser := ctx.Member
	response.SanitizeMember(accountUser)
	response.ComputeMemberLastModified(ctx, accountUser)
	response.WriteResponse(ctx, accountUser, http.StatusOK, 10)
	return
}

func (user *member) Update(ctx *context.Context) (err []errors.Error) {
	accountUser := *ctx.Member

	if accountUser.PublicID != ctx.Vars["accountUserID"] {
		return []errors.Error{errmsg.ErrMemberMismatchErr.SetCurrentLocation()}
	}

	if er := json.Unmarshal(ctx.Body, &accountUser); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	accountUser.ID = ctx.MemberID
	accountUser.OrgID = ctx.OrganizationID

	if err = validator.UpdateMember(user.storage, ctx.Member, &accountUser); err != nil {
		return
	}

	updatedAccountUser, err := user.storage.Update(*ctx.Member, accountUser, true)
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
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}
	accountUser, err := user.storage.FindByPublicID(ctx.OrganizationID, accountUserID)
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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(err.Error()).SetCurrentLocation()}
	}

	accountUser.OrgID = ctx.OrganizationID
	accountUser.PublicAccountID = ctx.Organization.PublicID

	if err = validator.CreateMember(user.storage, accountUser); err != nil {
		return
	}

	if accountUser, err = user.storage.Create(accountUser, true); err != nil {
		return
	}

	sessionToken := ""
	if sessionToken, err = user.storage.CreateSession(accountUser); err != nil {
		return
	}

	response.SanitizeMember(accountUser)

	rsp := struct {
		entity.Member
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
	}{
		Member:       *accountUser,
		Token:        sessionToken,
		AccountToken: ctx.Organization.AuthToken,
	}

	response.WriteResponse(ctx, rsp, http.StatusCreated, 0)
	return
}

func (user *member) List(ctx *context.Context) (err []errors.Error) {
	var (
		accountUsers []*entity.Member
	)

	if accountUsers, err = user.storage.List(ctx.OrganizationID); err != nil {
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
		usr          *entity.Member
		sessionToken string
		er           error
	)

	if er = json.Unmarshal(ctx.Body, loginPayload); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
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
		return []errors.Error{errmsg.ErrMemberNotFound.SetCurrentLocation()}
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

	response.SanitizeMember(usr)

	response.WriteResponse(ctx, struct {
		entity.Member
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
	}{
		Member:       *usr,
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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	if ctx.SessionToken != tokenPayload.Token {
		return []errors.Error{errmsg.ErrAuthSessionTokenMismatch.SetCurrentLocation()}
	}

	if sessionToken, err = user.storage.RefreshSession(ctx.SessionToken, ctx.Member); err != nil {
		return
	}

	response.WriteResponse(ctx, struct {
		Token string `json:"token"`
	}{Token: sessionToken}, http.StatusCreated, 0)
	return
}

func (user *member) Logout(ctx *context.Context) (err []errors.Error) {
	if err = user.storage.DestroySession(ctx.SessionToken, ctx.Member); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 0)
	return
}

// PopulateContext adds the accountUser to the context
func (user *member) PopulateContext(ctx *context.Context) (err []errors.Error) {
	usr, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidAccountUserCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", usr, pass)).SetCurrentLocation()}
	}
	accountUser, err := user.storage.FindBySession(pass)
	if accountUser == nil {
		return []errors.Error{errmsg.ErrMemberNotFound.SetCurrentLocation()}
	}
	if err == nil {
		ctx.Member = accountUser
		ctx.MemberID = accountUser.ID
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
