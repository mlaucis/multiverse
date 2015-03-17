/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"

	"strings"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// updateConnection handles requests to update a user connection
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func updateConnection(ctx *context.Context) {
	var (
		userToID int64
		err      error
	)

	if userToID, err = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); err != nil {
		errorHappened(ctx, "failed to update the connection (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	existingConnection, err := core.ReadConnection(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, userToID)
	if err != nil {
		errorHappened(ctx, "failed to update the connection (2)", http.StatusInternalServerError, err)
		return
	}

	if existingConnection == nil {
		errorHappened(ctx, "failed to update the connection (3)\nusers are not connected", http.StatusBadRequest, fmt.Errorf("users are not connected"))
		return
	}

	connection := *existingConnection
	if err = json.NewDecoder(ctx.Body).Decode(&connection); err != nil {
		errorHappened(ctx, "failed to update the connection (4)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID
	connection.UserToID = userToID

	if err = validator.UpdateConnection(existingConnection, &connection); err != nil {
		errorHappened(ctx, "failed to update the connection (5)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedConnection, err := core.UpdateConnection(*existingConnection, connection, false)
	if err != nil {
		errorHappened(ctx, "failed to update the connection (6)", http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, updatedConnection, http.StatusCreated, 0)
}

// deleteConnection handles requests to delete a single connection
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func deleteConnection(ctx *context.Context) {
	var (
		userToID int64
		err      error
	)

	if userToID, err = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); err != nil {
		errorHappened(ctx, "failed to delete the connection(1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteConnection(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, userToID); err != nil {
		errorHappened(ctx, "failed to delete the connection (2)", http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createConnection handles requests to create a user connection
// Request: POST /application/:applicationId/user/:UserID/connections
func createConnection(ctx *context.Context) {
	var (
		connection = &entity.Connection{}
		err        error
	)

	if err = json.NewDecoder(ctx.Body).Decode(connection); err != nil {
		errorHappened(ctx, "failed to create the connection (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID

	if err = validator.CreateConnection(connection); err != nil {
		errorHappened(ctx, "failed to create the connection (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection, err = core.WriteConnection(connection, false); err != nil {
		errorHappened(ctx, "failed to create the connection (3)", http.StatusInternalServerError, err)
		return
	}

	if connection, err = core.ConfirmConnection(connection, true); err != nil {
		errorHappened(ctx, "failed to create the connection (4)", http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, connection, http.StatusCreated, 0)
}

// getConnectionList handles requests to list a users connections
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections
func getConnectionList(ctx *context.Context) {
	var (
		users []*entity.User
		err   error
	)

	if users, err = core.ReadConnectionList(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID); err != nil {
		errorHappened(ctx, "failed to retrieve the connections list (1)", http.StatusInternalServerError, err)
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		Users         []*entity.User
	}{
		ApplicationID: ctx.ApplicationID,
		Users:         users,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// confirmConnection handles requests to confirm a user connection
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/connection/confirm
func confirmConnection(ctx *context.Context) {
	var (
		connection = &entity.Connection{}
		userFromID int64
		err        error
	)

	if userFromID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "failed to confirm the connection (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = json.NewDecoder(ctx.Body).Decode(connection); err != nil {
		errorHappened(ctx, "failed to confirm the connection (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	connection.AccountID = ctx.AccountID
	connection.ApplicationID = ctx.ApplicationID
	connection.UserFromID = userFromID

	if err = validator.ConfirmConnection(connection); err != nil {
		errorHappened(ctx, "failed to confirm the connection (3)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection, err = core.ConfirmConnection(connection, false); err != nil {
		errorHappened(ctx, "failed to confirm the connection (4)", http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, connection, http.StatusCreated, 0)
}

var acceptedPlatforms = map[string]bool{
	"facebook": true,
	"twitter":  true,
	"gplus":    true,
	"abook":    true,
}

func createSocialConnections(ctx *context.Context) {
	platformName := strings.ToLower(ctx.Vars["platformName"])

	if _, ok := acceptedPlatforms[platformName]; !ok {
		errorHappened(ctx, "social connecting failed (1)\nunexpected social platform", http.StatusBadRequest, fmt.Errorf("expected social platform"))
		return
	}

	socialIds := []string{}
	if err := json.NewDecoder(ctx.Body).Decode(&socialIds); err != nil {
		errorHappened(ctx, "social connecting failed (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	users, err := core.SocialConnect(ctx.ApplicationUser, platformName, socialIds)
	if err != nil {
		errorHappened(ctx, "social connecting failed (3)", http.StatusInternalServerError, err)
		return
	}

	for idx := range users {
		users[idx].Password = ""
	}

	writeResponse(ctx, users, http.StatusCreated, 10)
}

// TODO: GET FOLLOWER LIST (followedbyid users)
// TODO: GET FOLLOWING USERS LIST
