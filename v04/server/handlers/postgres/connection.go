package postgres

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/errors"
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

	p := &entity.PresentationConnection{}

	if er := json.Unmarshal(ctx.Body, p); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection := p.Connection

	connection.UserFromID = userFromID
	connection.UserToID = userToID
	if !connection.IsValidType() {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong}
	}

	existingConnection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID, connection.Type)
	if err != nil {
		return
	}
	if existingConnection == nil {
		return []errors.Error{errmsg.ErrConnectionUsersNotConnected.SetCurrentLocation()}
	}

	err = validator.UpdateConnection(conn.appUser, accountID, applicationID, existingConnection, connection)
	if err != nil {
		return
	}

	updatedConnection, err := conn.storage.Update(accountID, applicationID, existingConnection, connection, false)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, &entity.PresentationConnection{Connection: updatedConnection}, http.StatusCreated, 0)
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

	connectionType := entity.ConnectionTypeType(ctx.Vars["connectionType"])
	if connectionType != entity.ConnectionTypeFollow &&
		connectionType != entity.ConnectionTypeFriend {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.UpdateInternalMessage("got connection type: " + string(connectionType)).SetCurrentLocation()}
	}

	existingConnection, err := conn.storage.Read(accountID, applicationID, userFromID, userToID, connectionType)
	if err != nil {
		if connectionType == entity.ConnectionTypeFriend && err[0].Code() == errmsg.ErrConnectionNotFound.Code() {
			existingConnection, err = conn.storage.Read(accountID, applicationID, userToID, userFromID, connectionType)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	if existingConnection == nil {
		return []errors.Error{errmsg.ErrConnectionNotFound.SetCurrentLocation()}
	}

	if existingConnection.State == entity.ConnectionStateRejected {
		if existingConnection.UserFromID == userFromID {
			return []errors.Error{errmsg.ErrConnectionDeletionNotAllowed.SetCurrentLocation()}
		}
	}

	err = conn.storage.Delete(accountID, applicationID, existingConnection.UserFromID, existingConnection.UserToID, connectionType)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (conn *connection) Create(ctx *context.Context) []errors.Error {
	var (
		p = &entity.PresentationConnection{}

		er error
	)

	if er = json.Unmarshal(ctx.Body, p); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection := p.Connection

	connection.UserFromID = ctx.ApplicationUserID

	switch connection.Type {
	case entity.ConnectionTypeFollow, entity.ConnectionTypeFriend:
		return conn.create(ctx, connection)
	default:
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.UpdateMessage("unexpected connection type " + string(connection.Type)).SetCurrentLocation()}
	}

	return nil
}

func (conn *connection) FollowingList(ctx *context.Context) (err []errors.Error) {
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

	userIDs, err := conn.storage.Following(accountID, applicationID, userID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(accountID, applicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, userID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
		UsersCount: len(users),
	}

	status := http.StatusOK

	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) CurrentUserFollowingList(ctx *context.Context) (err []errors.Error) {
	userIDs, err := conn.storage.Following(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, ctx.ApplicationUserID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
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
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := conn.appUser.ExistsByID(ctx.OrganizationID, ctx.ApplicationID, userID)
	if err != nil {
		return
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
	}

	userIDs, err := conn.storage.FollowedBy(ctx.OrganizationID, ctx.ApplicationID, userID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, userID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
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
	userIDs, err := conn.storage.FollowedBy(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, ctx.ApplicationUserID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
		UsersCount: len(users),
	}

	status := http.StatusOK

	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return
}

func (conn *connection) CreateSocial(ctx *context.Context) (err []errors.Error) {
	request := entity.CreateSocialConnectionRequest{}

	if er := json.Unmarshal(ctx.Body, &request); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	if request.ConnectionType != entity.ConnectionTypeFriend && request.ConnectionType != entity.ConnectionTypeFollow {
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

	if request.ConnectionState == "" {
		request.ConnectionState = entity.ConnectionStateConfirmed
	}

	userIDs, err := conn.storage.SocialConnect(
		ctx.OrganizationID,
		ctx.ApplicationID,
		user,
		request.SocialPlatform,
		request.ConnectionsIDs,
		request.ConnectionType,
		request.ConnectionState)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, user.ID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
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

	userIDs, err := conn.storage.Friends(accountID, applicationID, userID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, userID, users); err != nil {
		return err
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
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
	userIDs, err := conn.storage.Friends(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID)
	if err != nil {
		return
	}

	users, err := conn.appUser.ReadMultiple(ctx.OrganizationID, ctx.ApplicationID, userIDs)
	if err != nil {
		return
	}

	response.SanitizeApplicationUsers(users)

	for idx := range users {
		relation, err := conn.storage.Relation(ctx.OrganizationID, ctx.ApplicationID, ctx.ApplicationUserID, users[idx].ID)
		if err != nil {
			ctx.LogError(err)
		} else if relation != nil {
			users[idx].Relation = *relation
		}
	}

	resp := struct {
		Users      []*entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                   `json:"users_count"`
	}{
		Users:      conn.presentationUsers(users),
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
		p = &entity.PresentationConnection{}

		er error
	)

	if er = json.Unmarshal(ctx.Body, p); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection := p.Connection

	connection.Type = entity.ConnectionTypeFriend
	connection.UserFromID = ctx.ApplicationUserID
	return conn.create(ctx, connection)
}

func (conn *connection) CreateFollow(ctx *context.Context) []errors.Error {
	var (
		p = &entity.PresentationConnection{}

		er error
	)

	if er = json.Unmarshal(ctx.Body, p); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	connection := p.Connection

	connection.Type = entity.ConnectionTypeFollow
	connection.UserFromID = ctx.ApplicationUserID
	return conn.create(ctx, connection)
}

func (conn *connection) create(
	ctx *context.Context,
	connection *entity.Connection,
) (errs []errors.Error) {
	from, to, errs := conn.getConnectionPair(ctx, connection)
	if errs != nil {
		return errs
	}

	if connection.Type == entity.ConnectionTypeFriend && from == nil && to != nil {
		from = to
	}

	if from == nil {
		if connection.State == "" {
			connection.TransferState(entity.ConnectionStateConfirmed, ctx.ApplicationUserID)
		}

		errs := validator.CreateConnection(conn.appUser, ctx.OrganizationID, ctx.ApplicationID, connection)
		if errs != nil {
			return errs
		}

		errs = conn.storage.Create(ctx.OrganizationID, ctx.ApplicationID, connection)
		if errs != nil {
			return errs
		}

		response.WriteResponse(ctx, connection, http.StatusCreated, 0)
		return nil
	}

	if connection.State == "" || connection.State == from.State {
		response.WriteResponse(ctx, connection, http.StatusOK, 0)
		return nil
	}

	errs = from.TransferState(connection.State, ctx.ApplicationUserID)
	if errs != nil {
		return errs
	}

	_, errs = conn.storage.Update(ctx.OrganizationID, ctx.ApplicationID, from, from, false)
	if errs != nil {
		return errs
	}

	response.WriteResponse(ctx, connection, http.StatusOK, 0)
	return nil
}

func (conn *connection) getConnectionPair(
	ctx *context.Context,
	connection *entity.Connection,
) (*entity.Connection, *entity.Connection, []errors.Error) {
	from, err := conn.storage.Read(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUserID,
		connection.UserToID,
		connection.Type,
	)
	if err != nil {
		if err[0].Code() != errmsg.ErrConnectionNotFound.Code() {
			return nil, nil, err
		}
	}

	to, err := conn.storage.Read(
		ctx.OrganizationID,
		ctx.ApplicationID,
		connection.UserToID,
		ctx.ApplicationUserID,
		connection.Type,
	)
	if err != nil {
		if err[0].Code() != errmsg.ErrConnectionNotFound.Code() {
			return nil, nil, err
		}
	}

	return from, to, nil
}

func (conn *connection) CurrentUserConnectionsByState(ctx *context.Context) []errors.Error {
	userID := ctx.ApplicationUserID
	connectionState := entity.ConnectionStateType(ctx.Vars["connectionState"])

	return conn.doGetUserConnectionsByState(ctx, userID, connectionState)
}

func (conn *connection) UserConnectionsByState(ctx *context.Context) []errors.Error {
	connectionState := entity.ConnectionStateType(ctx.Vars["connectionState"])

	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}

	exists, err := conn.appUser.ExistsByID(ctx.OrganizationID, ctx.ApplicationID, userID)
	if err != nil {
		return err
	}

	if !exists {
		return []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
	}

	return conn.doGetUserConnectionsByState(ctx, userID, connectionState)
}

func (conn *connection) doGetUserConnectionsByState(ctx *context.Context, userID uint64, connectionState entity.ConnectionStateType) []errors.Error {
	orgID := ctx.OrganizationID
	appID := ctx.ApplicationID

	if !entity.IsValidConectionState(connectionState) {
		return []errors.Error{errmsg.ErrConnectionStateInvalid.SetCurrentLocation()}
	}

	connections, err := conn.storage.ConnectionsByState(orgID, appID, userID, connectionState)
	if err != nil {
		return err
	}

	incomingConnections := []*entity.Connection{}
	outgoingConnections := []*entity.Connection{}
	userIDs := []uint64{}
	for idx := range connections {
		connections[idx].Enabled = entity.PFalse
		if idx > 0 {
			if connections[idx-1].UserToID == connections[idx].UserFromID {
				continue
			}
		}

		if connections[idx].UserFromID == userID {
			userIDs = append(userIDs, connections[idx].UserToID)
			outgoingConnections = append(outgoingConnections, connections[idx])
		} else {
			userIDs = append(userIDs, connections[idx].UserFromID)
			incomingConnections = append(incomingConnections, connections[idx])
		}
	}

	users, err := conn.appUser.ReadMultiple(orgID, appID, userIDs)
	if err != nil {
		return err
	}

	response.SanitizeApplicationUsers(users)

	if err := conn.addRelationInformation(ctx, userID, users); err != nil {
		return err
	}

	resp := entity.ConnectionsByStateResponse{
		IncomingConnections: conn.presentationConnections(incomingConnections),
		OutgoingConnections: conn.presentationConnections(outgoingConnections),
		Users:               conn.presentationUsers(users),
		IncomingConnectionsCount: len(incomingConnections),
		OutgoingConnectionsCount: len(outgoingConnections),
		UsersCount:               len(users),
	}

	status := http.StatusOK
	if resp.UsersCount == 0 {
		status = http.StatusNoContent
	}

	response.WriteResponse(ctx, resp, status, 10)
	return nil
}

func (conn *connection) addRelationInformation(ctx *context.Context, userID uint64, users []*entity.ApplicationUser) []errors.Error {
	for idx := range users {
		relation, err := conn.storage.Relation(ctx.OrganizationID, ctx.ApplicationID, userID, users[idx].ID)
		if err != nil {
			return err
		} else if relation != nil {
			users[idx].Relation = *relation
		}
	}
	return nil
}

func (conn *connection) presentationUsers(users []*entity.ApplicationUser) []*entity.PresentationApplicationUser {
	usrs := make([]*entity.PresentationApplicationUser, len(users))
	for idx := range users {
		usrs[idx] = &entity.PresentationApplicationUser{
			ApplicationUser: users[idx],
		}
	}

	return usrs
}

func (conn *connection) presentationConnections(connections []*entity.Connection) []*entity.PresentationConnection {
	conns := make([]*entity.PresentationConnection, len(connections))
	for idx := range connections {
		conns[idx] = &entity.PresentationConnection{
			Connection: connections[idx],
		}
	}

	return conns
}

// NewConnection initializes a new connection with an application user
func NewConnection(storage core.Connection, appUser core.ApplicationUser, evt core.Event) handlers.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
		event:   evt,
	}
}
