/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	userFromID := ctx.Bag["applicationUserID"].(string)
	userToID := ctx.Vars["userToId"]

	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	userToID, err = conn.determineTGUserID(accountID, applicationID, userToID)
	if err != nil {
		return
	}

	existingConnection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return []errors.Error{errmsg.ErrApplicationUsersNotConnected}
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

	server.WriteResponse(ctx, updatedConnection, http.StatusCreated, 0)
	return
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	userFromID := ctx.Bag["applicationUserID"].(string)
	userToID := ctx.Vars["applicationUserToID"]

	userToID, err = conn.determineTGUserID(accountID, applicationID, userToID)
	if err != nil {
		return
	}

	connection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}

	err = conn.storage.Delete(accountID, applicationID, connection)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
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

	connection.UserFromID = ctx.Bag["applicationUserID"].(string)
	connection.UserToID, err = conn.determineTGUserID(accountID, applicationID, connection.UserToID)
	if err != nil {
		return
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

	server.WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

func (conn *connection) List(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Vars["applicationUserID"]
	userID, err = conn.determineTGUserID(accountID, applicationID, userID)
	if err != nil {
		return
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

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK

	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	var users []*entity.ApplicationUser

	if users, err = conn.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string)); err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) FollowedByList(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Vars["applicationUserID"]
	userID, err = conn.determineTGUserID(accountID, applicationID, userID)
	if err != nil {
		return
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

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) CurrentUserFollowedByList(ctx *context.Context) (err []errors.Error) {
	var users []*entity.ApplicationUser
	users, err = conn.storage.FollowedBy(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string))
	if err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK

	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) Confirm(ctx *context.Context) (err []errors.Error) {
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	connection.UserFromID = ctx.Bag["applicationUserID"].(string)
	connection.UserToID, err = conn.determineTGUserID(accountID, applicationID, connection.UserToID)
	if err != nil {
		return
	}

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

	server.WriteResponse(ctx, connection, http.StatusCreated, 0)
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

	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	server.WriteResponse(ctx, response, http.StatusCreated, 10)
	return
}

func (conn *connection) Friends(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userID := ctx.Vars["applicationUserID"]
	userID, err = conn.determineTGUserID(accountID, applicationID, userID)
	if err != nil {
		return
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

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) CurrentUserFriends(ctx *context.Context) (err []errors.Error) {
	users, err := conn.storage.Friends(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(string))
	if err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)
	sanitizeApplicationUsers(users)

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) determineTGUserID(accountID, applicationID int64, userID string) (string, []errors.Error) {
	if validator.IsValidUUID5(userID) {
		return userID, nil
	}

	user, err := conn.appUser.FindByCustomID(accountID, applicationID, userID)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

// NewConnectionWithApplicationUser initializes a new connection with an application user
func NewConnectionWithApplicationUser(storage core.Connection, appUser core.ApplicationUser) server.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
