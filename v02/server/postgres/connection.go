/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
	"github.com/tapglue/backend/v02/errmsg"
)

type (
	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	var (
		userToID string
		er       error
	)

	userToID = ctx.Vars["userToId"]
	if !validator.IsValidUUID5(userToID) {
		return []errors.Error{errmsg.InvalidUserIDError}
	}

	existingConnection, err := conn.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(string),
		userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return []errors.Error{errmsg.UsersNotConnectedError}
	}

	connection := *existingConnection
	if er = json.Unmarshal(ctx.Body, &connection); er != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
	}

	if connection.UserFromID != ctx.Bag["applicationUserID"].(string) {
		return []errors.Error{errmsg.UserFromMismatchError}
	}

	if connection.UserToID != userToID {
		return []errors.Error{errmsg.UserToMismatchError}
	}

	if connection.Type != "friend" && connection.Type != "follow" {
		return []errors.Error{errmsg.WrongConnectionTypeError.UpdateMessage(fmt.Sprintf("unexpected connection type %q", connection.Type))}
	}

	if err = validator.UpdateConnection(
		conn.appUser,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		existingConnection,
		&connection); err != nil {
		return
	}

	updatedConnection, err := conn.storage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		*existingConnection,
		connection,
		false)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, updatedConnection, http.StatusCreated, 0)
	return
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	userFromID := ctx.Bag["applicationUserID"].(string)
	userToID := ctx.Vars["applicationUserToID"]
	if !validator.IsValidUUID5(userToID) {
		return []errors.Error{errmsg.InvalidUserIDError}
	}

	connection, err := conn.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), userFromID, userToID)
	if err != nil {
		return
	}

	err = conn.storage.Delete(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), connection)
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
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
	}

	if connection.Type != "friend" && connection.Type != "follow" {
		return []errors.Error{errmsg.WrongConnectionTypeError.UpdateMessage(fmt.Sprintf("unexpected connection type %q", connection.Type))}
	}

	receivedEnabled := connection.Enabled

	connection.UserFromID = ctx.Bag["applicationUserID"].(string)

	if connection.UserFromID == connection.UserToID {
		return []errors.Error{errmsg.SelfConnectingUserError}
	}

	if err = validator.CreateConnection(
		conn.appUser,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection); err != nil {
		return
	}

	if connection, err = conn.storage.Create(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection,
		true); err != nil {
		return
	}

	if receivedEnabled {
		if connection, err = conn.storage.Confirm(
			ctx.Bag["accountID"].(int64),
			ctx.Bag["applicationID"].(int64),
			connection,
			true); err != nil {
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
	if !validator.IsValidUUID5(userID) {
		return []errors.Error{errmsg.InvalidUserIDError}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ApplicationUserNotFoundError}
	}

	var users []*entity.ApplicationUser
	users, err = conn.storage.List(accountID, applicationID, userID)
	if err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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
	if !validator.IsValidUUID5(userID) {
		return []errors.Error{errmsg.InvalidUserIDError}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ApplicationUserNotFoundError}
	}

	var users []*entity.ApplicationUser
	if users, err = conn.storage.FollowedBy(accountID, applicationID, userID); err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

	response := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{
		Users:      users,
		UsersCount: len(users),
	}

	status := http.StatusOK
	computeApplicationUsersLastModified(ctx, response.Users)

	if response.UsersCount == 0 {
		status = http.StatusNoContent
	}

	server.WriteResponse(ctx, response, status, 10)
	return
}

func (conn *connection) Confirm(ctx *context.Context) (err []errors.Error) {
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
	}

	connection.UserFromID = ctx.Bag["applicationUserID"].(string)

	connection, err = conn.storage.Read(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection.UserFromID,
		connection.UserToID,
	)
	if err != nil {
		return err
	}

	if err = validator.ConfirmConnection(
		conn.appUser,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection); err != nil {
		return
	}

	if connection, err = conn.storage.Confirm(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection,
		false); err != nil {
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
		return []errors.Error{errmsg.BadJsonReceivedError.UpdateMessage(er.Error())}
	}

	if request.ConnectionType == "" || (request.ConnectionType != "friend" && request.ConnectionType != "follow") {
		return []errors.Error{errmsg.WrongConnectionTypeError}
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

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ApplicationUserNotFoundError}
	}

	var users []*entity.ApplicationUser
	if users, err = conn.storage.Friends(accountID, applicationID, userID); err != nil {
		return
	}

	computeApplicationUsersLastModified(ctx, users)

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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

	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}

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

// NewConnectionWithApplicationUser initializes a new connection with an application user
func NewConnectionWithApplicationUser(storage core.Connection, appUser core.ApplicationUser) server.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
