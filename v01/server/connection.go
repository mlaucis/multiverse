/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v01/context"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// updateConnection handles requests to update a user connection
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func updateConnection(ctx *context.Context) (err errors.Error) {
	var (
		userToID int64
		er       error
	)

	if userToID, er = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); er != nil {
		return errors.NewBadRequestError("failed to update the connection (1)\n"+er.Error(), er.Error())
	}

	existingConnection, err := core.ReadConnection(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, userToID)
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

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID

	if connection.UserFromID != ctx.ApplicationUserID {
		return errors.NewBadRequestError("failed to update the connection (5)\nuser_from mismatch", "user_from mismatch")
	}

	if connection.UserToID != userToID {
		return errors.NewBadRequestError("failed to update the connection (6)\nuser_to mismatch", "user_to mismatch")
	}

	if err = validator.UpdateConnection(existingConnection, &connection); err != nil {
		return
	}

	updatedConnection, err := core.UpdateConnection(*existingConnection, connection, false)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedConnection, http.StatusCreated, 0)
	return
}

// deleteConnection handles requests to delete a single connection
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func deleteConnection(ctx *context.Context) (err errors.Error) {
	var (
		userToID int64
		er       error
	)

	if userToID, er = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); er != nil {
		return errors.NewBadRequestError("failed to delete the connection(1)\n"+er.Error(), er.Error())
	}

	if err = core.DeleteConnection(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, userToID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// createConnection handles requests to create a user connection
// Request: POST /application/:applicationId/user/:UserID/connections
func createConnection(ctx *context.Context) (err errors.Error) {
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return errors.NewBadRequestError("failed to create the connection(1)\n"+er.Error(), er.Error())
	}

	receivedEnabled := connection.Enabled

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID

	if connection.UserFromID == connection.UserToID {
		return errors.NewBadRequestError("failed to create connection (2)\nuser is connecting with itself", "self-connecting user")
	}

	if err = validator.CreateConnection(connection); err != nil {
		return
	}

	if connection, err = core.WriteConnection(connection, false); err != nil {
		return
	}

	if receivedEnabled {
		if connection, err = core.ConfirmConnection(connection, true); err != nil {
			return
		}
	}

	WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

// getConnectionList handles requests to list a users connections
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections
func getConnectionList(ctx *context.Context) (err errors.Error) {
	var users []*entity.User

	if users, err = core.ReadConnectionList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

// getFollowedByUsersList handles requests to list a users list of users who follow him
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/followers
func getFollowedByUsersList(ctx *context.Context) (err errors.Error) {
	var users []*entity.User

	if users, err = core.ReadFollowedByList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusOK, 10)
	return
}

// confirmConnection handles requests to confirm a user connection
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/connection/confirm
func confirmConnection(ctx *context.Context) (err errors.Error) {
	var connection = &entity.Connection{}

	if er := json.Unmarshal(ctx.Body, connection); er != nil {
		return errors.NewBadRequestError("failed to confirm the connection (1)\n"+er.Error(), er.Error())
	}

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID

	if err = validator.ConfirmConnection(connection); err != nil {
		return
	}

	if connection, err = core.ConfirmConnection(connection, false); err != nil {
		return
	}

	WriteResponse(ctx, connection, http.StatusCreated, 0)
	return
}

var acceptedPlatforms = map[string]bool{
	"facebook": true,
	"twitter":  true,
	"gplus":    true,
	"abook":    true,
}

func createSocialConnections(ctx *context.Context) (err errors.Error) {
	platformName := strings.ToLower(ctx.Vars["platformName"])

	if _, ok := acceptedPlatforms[platformName]; !ok {
		return errors.NewNotFoundError("social connecting failed (1)\nunexpected social platform", "platform not found")
	}

	socialConnections := struct {
		UserFromID     int64    `json:"user_from_id"`
		SocialPlatform string   `json:"social_platform"`
		ConnectionsIDs []string `json:"connection_ids"`
	}{}

	if er := json.Unmarshal(ctx.Body, &socialConnections); er != nil {
		return errors.NewBadRequestError("social connecting failed (2)\n"+er.Error(), er.Error())
	}

	if ctx.ApplicationUserID != socialConnections.UserFromID {
		return errors.NewBadRequestError("social connecting failed (3)\nuser mismatch", "user mismatch")
	}

	if platformName != strings.ToLower(socialConnections.SocialPlatform) {
		return errors.NewBadRequestError("social connecting failed (3)\nplatform mismatch", "platform mismatch")
	}

	users, err := core.SocialConnect(ctx.ApplicationUser, platformName, socialConnections.ConnectionsIDs)
	if err != nil {
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	WriteResponse(ctx, users, http.StatusCreated, 10)
	return
}
