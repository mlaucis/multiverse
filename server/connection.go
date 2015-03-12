/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// updateConnection handles requests to update a user connection
// Request: PUT account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func updateConnection(ctx *context.Context) {
	var (
		connection    = &entity.Connection{}
		accountID     int64
		applicationID int64
		userFromID    int64
		userToID      int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userFromID, err = strconv.ParseInt(ctx.Vars["userFromId"], 10, 64); err != nil {
		errorHappened(ctx, "userFromId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userToID, err = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); err != nil {
		errorHappened(ctx, "userToId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection.AccountID == 0 {
		connection.AccountID = accountID
	}
	if connection.ApplicationID == 0 {
		connection.ApplicationID = applicationID
	}
	if connection.UserFromID == 0 {
		connection.UserFromID = userFromID
	}
	if connection.UserToID == 0 {
		connection.UserToID = userToID
	}

	if err = validator.UpdateConnection(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection, err = core.UpdateConnection(connection, false); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, connection, http.StatusCreated, 0)
}

// deleteConnection handles requests to delete a single connection
// Request: DELETE account/:AccountID/application/:ApplicationID/user/:UserFromID/connection/:UserToID
func deleteConnection(ctx *context.Context) {
	var (
		accountID     int64
		applicationID int64
		userFromID    int64
		userToID      int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userFromID, err = strconv.ParseInt(ctx.Vars["userFromId"], 10, 64); err != nil {
		errorHappened(ctx, "userFromId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userToID, err = strconv.ParseInt(ctx.Vars["userToId"], 10, 64); err != nil {
		errorHappened(ctx, "userToId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = core.DeleteConnection(accountID, applicationID, userFromID, userToID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createConnection handles requests to create a user connection
// Request: POST /application/:applicationId/user/:UserID/connections
func createConnection(ctx *context.Context) {
	var (
		connection    = &entity.Connection{}
		accountID     int64
		applicationID int64
		userFromID    int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userFromID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = json.NewDecoder(ctx.Body).Decode(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	connection.AccountID = accountID
	connection.ApplicationID = applicationID
	connection.UserFromID = userFromID

	if err = validator.CreateConnection(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection, err = core.WriteConnection(connection, false); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, connection, http.StatusCreated, 0)
}

// getConnectionList handles requests to list a users connections
// Request: GET account/:AccountID/application/:ApplicationID/user/:UserID/connections
func getConnectionList(ctx *context.Context) {
	var (
		users         []*entity.User
		accountID     int64
		applicationID int64
		userID        int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if users, err = core.ReadConnectionList(accountID, applicationID, userID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response := &struct {
		ApplicationID int64 `json:"applicationId"`
		Users         []*entity.User
	}{
		ApplicationID: applicationID,
		Users:         users,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}

// confirmConnection handles requests to confirm a user connection
// Request: POST account/:AccountID/application/:ApplicationID/user/:UserID/connection/confirm
func confirmConnection(ctx *context.Context) {
	var (
		connection    = &entity.Connection{}
		accountID     int64
		applicationID int64
		userFromID    int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.Vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.Vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if userFromID, err = strconv.ParseInt(ctx.Vars["userId"], 10, 64); err != nil {
		errorHappened(ctx, "userId is not set or the value is incorrect", http.StatusBadRequest, err)
		return
	}

	if err = json.NewDecoder(ctx.Body).Decode(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	connection.AccountID = accountID
	connection.ApplicationID = applicationID
	connection.UserFromID = userFromID

	if err = validator.ConfirmConnection(connection); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if connection, err = core.ConfirmConnection(connection, false); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, connection, http.StatusCreated, 0)
}

// TODO: GET FOLLOWER LIST (followedbyid users)
// TODO: GET FOLLOWING USERS LIST
