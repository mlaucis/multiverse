package postgres

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/response"
	"github.com/tapglue/multiverse/v04/validator"
)

type connection struct {
	appUser core.ApplicationUser
	storage core.Connection
	event   core.Event
}

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	userFromID := ctx.ApplicationUserID

	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID

	userToID, er := strconv.ParseUint(ctx.Vars["userToID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	existingConnection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return []errors.Error{errmsg.ErrConnectionUsersNotConnected.SetCurrentLocation()}
	}

	connection := *existingConnection
	if er := json.Unmarshal(ctx.Body, &connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
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
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID

	userFromID := ctx.ApplicationUserID

	userToID, er := strconv.ParseUint(ctx.Vars["applicationUserToID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
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

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	return conn.doCreateConnection(ctx, connection)
}

func (conn *connection) List(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
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
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID); err != nil {
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
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
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
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID)
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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID

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
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	if request.ConnectionType == "" || (request.ConnectionType != "friend" && request.ConnectionType != "follow") {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.SetCurrentLocation()}
	}

	user := ctx.ApplicationUser

	if _, ok := user.SocialIDs[request.SocialPlatform]; !ok {
		if len(user.SocialIDs[request.SocialPlatform]) == 0 {
			user.SocialIDs = map[string]string{}
		}
		user.SocialIDs[request.SocialPlatform] = request.PlatformUserID
		_, err = conn.appUser.Update(
			ctx.OrganizationID,
			ctx.ApplicationID,
			*user,
			*user,
			false)
		if err != nil {
			return err
		}
	}

	users, err := conn.storage.SocialConnect(
		ctx.OrganizationID,
		ctx.ApplicationID,
		user,
		request.SocialPlatform,
		request.ConnectionsIDs,
		request.ConnectionType)
	if err != nil {
		return
	}

	if ctx.Query.Get("with_event") == "true" {
		_, err := conn.CreateAutoConnectionEvents(ctx, user, users, request.ConnectionType)
		if err != nil {
			ctx.LogError(err)
		}
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
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := conn.appUser.ExistsByID(accountID, applicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
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
	users, err := conn.storage.Friends(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID)
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

func (conn *connection) CreateFriend(ctx *context.Context) []errors.Error {
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection.Type = "friend"
	return conn.doCreateConnection(ctx, connection)
}

func (conn *connection) CreateFollow(ctx *context.Context) []errors.Error {
	var (
		connection = &entity.Connection{}
		er         error
	)
	connection.Enabled = true

	if er = json.Unmarshal(ctx.Body, connection); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection.Type = "follow"
	return conn.doCreateConnection(ctx, connection)
}

func (conn *connection) doCreateConnection(ctx *context.Context, connection *entity.Connection) []errors.Error {
	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID
	connection.UserFromID = ctx.ApplicationUserID

	receivedEnabled := connection.Enabled
	connection.Enabled = true

	if exists, err := conn.storage.Exists(
		accountID, applicationID,
		connection.UserFromID, connection.UserToID, connection.Type); exists || err != nil {
		if exists {
			response.WriteResponse(ctx, "", http.StatusNoContent, 0)
			return nil
		}

		return err
	}

	err := validator.CreateConnection(conn.appUser, accountID, applicationID, connection)
	if err != nil {
		return err
	}

	err = conn.storage.Create(accountID, applicationID, connection)
	if err != nil {
		return err
	}

	if ctx.Query.Get("with_event") == "true" {
		_, err := conn.CreateAutoConnectionEvent(ctx, connection)
		if err != nil {
			ctx.LogError(err)
		}
	}

	if receivedEnabled {
		connection, err = conn.storage.Confirm(accountID, applicationID, connection, true)
		if err != nil {
			return err
		}
	}

	response.WriteResponse(ctx, connection, http.StatusCreated, 0)
	return nil
}

func (conn *connection) CreateAutoConnectionEvent(ctx *context.Context, connection *entity.Connection) (*entity.Event, []errors.Error) {
	event := &entity.Event{
		UserID:     connection.UserFromID,
		Type:       "tg_" + connection.Type,
		Visibility: entity.EventPrivate,
		Target: &entity.Object{
			ID:   strconv.FormatUint(connection.UserToID, 10),
			Type: "tg_user",
		},
	}

	var err error
	event.ID, err = tgflake.FlakeNextID(ctx.ApplicationID, "events")
	if err != nil {
		return nil, []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	accountID := ctx.OrganizationID
	applicationID := ctx.ApplicationID

	er := conn.event.Create(accountID, applicationID, connection.UserFromID, event)
	return event, er
}

func (conn *connection) CreateAutoConnectionEvents(
	ctx *context.Context,
	user *entity.ApplicationUser, users []*entity.ApplicationUser,
	connectionType string) ([]*entity.Event, []errors.Error) {

	events := []*entity.Event{}
	errs := []errors.Error{}
	for idx := range users {
		connection := &entity.Connection{
			UserFromID: user.ID,
			UserToID:   users[idx].ID,
			Type:       connectionType,
			Common: entity.Common{
				Enabled: true,
			},
		}

		evt, err := conn.CreateAutoConnectionEvent(ctx, connection)

		events = append(events, evt)
		errs = append(errs, err...)
	}

	return events, errs
}

// NewConnection initializes a new connection with an application user
func NewConnection(storage core.Connection, appUser core.ApplicationUser, evt core.Event) handlers.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
		event:   evt,
	}
}
