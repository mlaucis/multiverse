/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

type (
	// Connection holds the routes for the connections
	Connection interface {
		// Update handles requests to update a user connection
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single connection
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests to create a user connection
		Create(*context.Context) tgerrors.TGError

		// List handles requests to list a users connections
		List(*context.Context) tgerrors.TGError

		// FollowedByList handles requests to list a users list of users who follow him
		FollowedByList(*context.Context) tgerrors.TGError

		// Confirm handles requests to confirm a user connection
		Confirm(*context.Context) tgerrors.TGError

		// CreateSocialConnections creates the social connections between users of the same social network
		CreateSocial(*context.Context) tgerrors.TGError
	}

	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

var acceptedPlatforms = map[string]bool{
	"facebook": true,
	"twitter":  true,
	"gplus":    true,
	"abook":    true,
}

func (conn *connection) Update(ctx *context.Context) (err tgerrors.TGError) {
	var (
		userToID int64
		er       error
	)

	if userToID, er = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to update the connection (1)\n"+er.Error(), er.Error())
	}

	existingConnection, err := conn.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return tgerrors.NewNotFoundError("failed to update the connection (3)\nusers are not connected", "users are not connected")
	}

	connection := *existingConnection
	if er = json.Unmarshal(ctx.Body, &connection); er != nil {
		return tgerrors.NewBadRequestError("failed to update the connection (4)\n"+er.Error(), er.Error())
	}

	connection.AccountID = ctx.Bag["accountID"].(int64)
	connection.ApplicationID = ctx.Bag["applicationID"].(int64)

	if connection.UserFromID != ctx.Bag["applicationUserID"].(int64) {
		return tgerrors.NewBadRequestError("failed to update the connection (5)\nuser_from mismatch", "user_from mismatch")
	}

	if connection.UserToID != userToID {
		return tgerrors.NewBadRequestError("failed to update the connection (6)\nuser_to mismatch", "user_to mismatch")
	}

	if err = validator.UpdateConnection(conn.appUser, existingConnection, &connection); err != nil {
		return
	}

	updatedConnection, err := conn.storage.Update(*existingConnection, connection, false)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedConnection, http.StatusCreated, 0)
	return
}

func (conn *connection) Delete(ctx *context.Context) (err tgerrors.TGError) {
	var (
		userToID int64
		er       error
	)

	if userToID, er = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); er != nil {
		return tgerrors.NewBadRequestError("failed to delete the connection(1)\n"+er.Error(), er.Error())
	}

	if err = conn.storage.Delete(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64), userToID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (conn *connection) Create(ctx *context.Context) (err tgerrors.TGError) {
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return tgerrors.NewBadRequestError("failed to create the connection(1)\n"+er.Error(), er.Error())
	}

	receivedEnabled := connection.Enabled

	connection.AccountID = ctx.Bag["accountID"].(int64)
	connection.ApplicationID = ctx.Bag["applicationID"].(int64)
	connection.UserFromID = ctx.Bag["applicationUserID"].(int64)

	if connection.UserFromID == connection.UserToID {
		return tgerrors.NewBadRequestError("failed to create connection (2)\nuser is connecting with itself", "self-connecting user")
	}

	if err = validator.CreateConnection(conn.appUser, connection); err != nil {
		return
	}

	if connection, err = conn.storage.Create(connection, false); err != nil {
		return
	}

	if receivedEnabled {
		if connection, err = conn.storage.Confirm(connection, true); err != nil {
			return
		}
	}

	WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

func (conn *connection) List(ctx *context.Context) (err tgerrors.TGError) {
	var users []*entity.ApplicationUser

	if users, err = conn.storage.List(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

func (conn *connection) FollowedByList(ctx *context.Context) (err tgerrors.TGError) {
	var users []*entity.ApplicationUser

	if users, err = conn.storage.FollowedBy(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64), ctx.Bag["applicationUserID"].(int64)); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

func (conn *connection) Confirm(ctx *context.Context) (err tgerrors.TGError) {
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return tgerrors.NewBadRequestError("failed to confirm the connection (1)\n"+er.Error(), er.Error())
	}

	connection.AccountID = ctx.Bag["accountID"].(int64)
	connection.ApplicationID = ctx.Bag["applicationID"].(int64)
	connection.UserFromID = ctx.Bag["applicationUserID"].(int64)

	if err = validator.ConfirmConnection(conn.appUser, connection); err != nil {
		return
	}

	if connection, err = conn.storage.Confirm(connection, false); err != nil {
		return
	}

	WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

func (conn *connection) CreateSocial(ctx *context.Context) (err tgerrors.TGError) {
	platformName := strings.ToLower(ctx.Vars["platformName"])

	if _, ok := acceptedPlatforms[platformName]; !ok {
		return tgerrors.NewNotFoundError("social connecting failed (1)\nunexpected social platform", "platform not found")
	}

	socialConnections := struct {
		UserFromID     int64    `json:"user_from_id"`
		SocialPlatform string   `json:"social_platform"`
		ConnectionsIDs []string `json:"connection_ids"`
	}{}

	if er := json.Unmarshal(ctx.Body, &socialConnections); er != nil {
		return tgerrors.NewBadRequestError("social connecting failed (2)\n"+er.Error(), er.Error())
	}

	if ctx.Bag["applicationUserID"].(int64) != socialConnections.UserFromID {
		return tgerrors.NewBadRequestError("social connecting failed (3)\nuser mismatch", "user mismatch")
	}

	if platformName != strings.ToLower(socialConnections.SocialPlatform) {
		return tgerrors.NewBadRequestError("social connecting failed (3)\nplatform mismatch", "platform mismatch")
	}

	users, err := conn.storage.SocialConnect(ctx.Bag["applicationUser"].(*entity.ApplicationUser), platformName, socialConnections.ConnectionsIDs)
	if err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusCreated, 10)
	return
}

// NewConnection returns a new connection handler
func NewConnection(storage core.Connection) Connection {
	return &connection{
		storage: storage,
	}
}

// NewConnectionWithApplicationUser initializes a new connection with an application user
func NewConnectionWithApplicationUser(storage core.Connection, appUser core.ApplicationUser) Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
