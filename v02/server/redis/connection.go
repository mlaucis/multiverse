/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	var (
		userToID int64
		er       error
	)

	if userToID, er = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); er != nil {
		return errors.NewBadRequestError("failed to update the connection (1)\n"+er.Error(), er.Error())
	}

	existingConnection, err := conn.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return errors.NewNotFoundError("failed to update the connection (3)\nusers are not connected", "users are not connected")
	}

	connection := *existingConnection
	if er = json.Unmarshal(ctx.Body, &connection); er != nil {
		return errors.NewBadRequestError("failed to update the connection (4)\n"+er.Error(), er.Error())
	}

	if connection.UserFromID != ctx.Bag["applicationUserID"].(int64) {
		return errors.NewBadRequestError("failed to update the connection (5)\nuser_from mismatch", "user_from mismatch")
	}

	if connection.UserToID != userToID {
		return errors.NewBadRequestError("failed to update the connection (6)\nuser_to mismatch", "user_to mismatch")
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

func (conn *connection) Delete(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	connection := &entity.Connection{}
	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return errors.NewBadRequestError(er.Error(), er.Error())
	}

	if err = conn.storage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		connection); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (conn *connection) Create(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return errors.NewBadRequestError("failed to create the connection(1)\n"+er.Error(), er.Error())
	}

	receivedEnabled := connection.Enabled

	connection.UserFromID = ctx.Bag["applicationUserID"].(int64)

	if connection.UserFromID == connection.UserToID {
		return errors.NewBadRequestError("failed to create connection (2)\nuser is connecting with itself", "self-connecting user")
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
		false); err != nil {
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

func (conn *connection) List(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	var users []*entity.ApplicationUser

	if users, err = conn.storage.List(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	server.WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

func (conn *connection) FollowedByList(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	var users []*entity.ApplicationUser

	if users, err = conn.storage.FollowedBy(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	server.WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

func (conn *connection) Confirm(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return errors.NewBadRequestError("failed to confirm the connection (1)\n"+er.Error(), er.Error())
	}

	connection.UserFromID = ctx.Bag["applicationUserID"].(int64)

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

func (conn *connection) CreateSocial(ctx *context.Context) (err errors.Error) {
	return errors.NewInternalError("deprecated storage used", "redis storage is deprecated")
	platformName := strings.ToLower(ctx.Vars["platformName"])

	socialConnections := struct {
		UserFromID     int64    `json:"user_from_id"`
		SocialPlatform string   `json:"social_platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		ConnectionType string   `json:"type"`
	}{}

	if er := json.Unmarshal(ctx.Body, &socialConnections); er != nil {
		return errors.NewBadRequestError("social connecting failed (2)\n"+er.Error(), er.Error())
	}

	if ctx.Bag["applicationUserID"].(int64) != socialConnections.UserFromID {
		return errors.NewBadRequestError("social connecting failed (3)\nuser mismatch", "user mismatch")
	}

	if platformName != strings.ToLower(socialConnections.SocialPlatform) {
		return errors.NewBadRequestError("social connecting failed (3)\nplatform mismatch", "platform mismatch")
	}

	users, err := conn.storage.SocialConnect(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser),
		platformName,
		socialConnections.ConnectionsIDs,
		socialConnections.ConnectionType)
	if err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	server.WriteResponse(ctx, users, http.StatusCreated, 10)
	return
}

// NewConnection returns a new connection handler
func NewConnection(storage core.Connection) server.Connection {
	return &connection{
		storage: storage,
	}
}

// NewConnectionWithApplicationUser initializes a new connection with an application user
func NewConnectionWithApplicationUser(storage core.Connection, appUser core.ApplicationUser) server.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
