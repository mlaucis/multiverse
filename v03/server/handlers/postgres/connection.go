package postgres

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	userFromID := ctx.Bag["applicationUserID"].(uint64)

	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	userToID, er := strconv.ParseUint(ctx.Vars["userToID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	existingConnection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return []errors.Error{errmsg.ErrConnectionUsersNotConnected}
	}

	connection := *existingConnection
	if er := json.Unmarshal(ctx.Body, &connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	connection.UserFromID = userFromID
	connection.UserToID = userToID

	err = validator.UpdateConnection(conn.appUser, accountID, applicationID, existingConnection, &connection)
	if err != nil {
		return
	}

	updatedConnection, err := conn.storage.Update(accountID, applicationID, *existingConnection, connection, false)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, updatedConnection, http.StatusCreated, 0)
	return
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	userFromID := ctx.Bag["applicationUserID"].(uint64)

	userToID, er := strconv.ParseUint(ctx.Vars["applicationUserToID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	connection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}

	err = conn.storage.Delete(accountID, applicationID, connection)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (conn *connection) Create(ctx *context.Context) (err []errors.Error) {
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	receivedEnabled := connection.Enabled

	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	connection.UserFromID = ctx.Bag["applicationUserID"].(uint64)

	if exists, err := conn.storage.Exists(
		accountID, applicationID,
		connection.UserFromID, connection.UserToID, connection.Type); exists || err != nil {
		if exists {
			response.WriteResponse(ctx, "", http.StatusNoContent, 0)
			return nil
		}

		return err
	}

	if err = validator.CreateConnection(conn.appUser, accountID, applicationID, connection); err != nil {
		return
	}

	if connection, err = conn.storage.Create(accountID, applicationID, connection, true); err != nil {
		return
	}

	if receivedEnabled {
		if connection, err = conn.storage.Confirm(accountID, applicationID, connection, true); err != nil {
			return
		}
	}

	response.WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

func (conn *connection) List(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound}
	}

	var users []*entity.ApplicationUser
	users, err = conn.storage.List(accountID, applicationID, userID)
	if err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK

	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	var users []*entity.ApplicationUser

	if users, err = conn.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64)); err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) FollowedByList(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound}
	}

	var users []*entity.ApplicationUser
	if users, err = conn.storage.FollowedBy(accountID, applicationID, userID); err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) CurrentUserFollowedByList(ctx *context.Context) (err []errors.Error) {
	var users []*entity.ApplicationUser
	users, err = conn.storage.FollowedBy(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(uint64))
	if err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK

	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) Confirm(ctx *context.Context) (err []errors.Error) {
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	connection.UserFromID = ctx.Bag["applicationUserID"].(uint64)

	connection, err = conn.storage.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
	if err != nil {
		return err
	}

	if err = validator.ConfirmConnection(conn.appUser, accountID, applicationID, connection); err != nil {
		return
	}

	if connection, err = conn.storage.Confirm(accountID, applicationID, connection, false); err != nil {
		return
	}

	response.WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

func (conn *connection) CreateSocial(ctx *context.Context) (err []errors.Error) {
	request := struct {
		PlatformUserID string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		ConnectionType string   `json:"type"`
	}{}

	if er := json.Unmarshal(ctx.Body, &request); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	if request.ConnectionType == "" || (request.ConnectionType != "friend" && request.ConnectionType != "follow") {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong}
	}

	user := ctx.Bag["applicationUser"].(*entity.ApplicationUser)

	if _, ok := user.SocialIDs[request.SocialPlatform]; !ok {
		if len(user.SocialIDs[request.SocialPlatform]) == 0 {
			user.SocialIDs = map[string]string{}
		}
		user.SocialIDs[request.SocialPlatform] = request.PlatformUserID
		_, err = conn.appUser.Update(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			*user,
			*user,
			false)
		if err != nil {
			return err
		}
	}

	users, err := conn.storage.SocialConnect(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		user,
		request.SocialPlatform,
		request.ConnectionsIDs,
		request.ConnectionType)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	response.WriteResponse(ctx, resp, http.StatusCreated, 10)
	return
}

func (conn *connection) Friends(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound}
	}

	var users []*entity.ApplicationUser
	if users, err = conn.storage.Friends(accountID, applicationID, userID); err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) CurrentUserFriends(ctx *context.Context) (err []errors.Error) {
	users, err := conn.storage.Friends(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(uint64))
	if err != nil {
		return
	}

	response.ComputeApplicationUsersLastModified(ctx, users)
	response.SanitizeApplicationUsers(users)

	resp := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

// NewConnection initializes a new connection with an application user
func NewConnection(storage core.Connection, appUser core.ApplicationUser) handlers.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
