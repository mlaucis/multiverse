package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type connection struct {
	pg      postgres.Client
	mainPg  *sqlx.DB
	appUser core.ApplicationUser
}

const (
	createConnectionQuery = `INSERT INTO app_%d_%d.connections(json_data) VALUES ($1)`

	selectConnectionQuery = `SELECT json_data
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
			AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT
			AND json_data->>'type' = $3::TEXT
			AND (json_data->>'enabled')::BOOL = true
		LIMIT 1`

	updateConnectionQuery = `UPDATE
	app_%d_%d.connections
	SET json_data = $1
	WHERE (json_data->>'user_from_id')::BIGINT = $2::BIGINT
		AND (json_data->>'user_to_id')::BIGINT = $3::BIGINT
		AND json_data->>'type' = $4::TEXT
		AND (json_data->>'enabled')::BOOL = true`

	followsQuery = `SELECT
		json_data
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
			AND json_data->>'type' = '` + string(entity.ConnectionTypeFollow) + `'
			AND (json_data->>'enabled')::BOOL = true`

	followersQuery = `SELECT
		json_data
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_to_id')::BIGINT = $1::BIGINT
			AND json_data->>'type' = '` + string(entity.ConnectionTypeFollow) + `'
			AND (json_data->>'enabled')::BOOL = true`

	friendConnectionsQuery = `SELECT
		json_data->>'user_from_id', json_data->>'user_to_id'
		FROM app_%d_%d.connections
		WHERE ((json_data->>'user_to_id')::BIGINT = $1::BIGINT OR (json_data->>'user_from_id')::BIGINT = $1::BIGINT)
			AND json_data->>'type' = '` + string(entity.ConnectionTypeFriend) + `'
			AND json_data->>'state' = '` + string(entity.ConnectionStateConfirmed) + `'
			AND (json_data->>'enabled')::BOOL = true`

	friendAndFollowingConnectionsQuery = `SELECT
		json_data
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
			AND (json_data->>'enabled')::BOOL = true`

	friendAndFollowingConnectionsIDsQuery = `SELECT
		json_data->>'user_to_id' as "user_id"
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
			AND (json_data->>'enabled')::BOOL = true`

	listUserIDssBySocialIDQuery = `SELECT
		json_data->>'id'
		FROM app_%d_%d.users
		WHERE (json_data->>'enabled')::BOOL = true
			AND (json_data->>'deleted')::BOOL = false
			AND json_data->'social_ids'->>'%s' IN (?)`

	getUsersRelationQuery = `SELECT
		json_data ->> 'user_from_id' AS "from",
		json_data ->> 'user_to_id'   AS "to",
		json_data ->> 'type'         AS "type"
		FROM app_%d_%d.connections
		WHERE json_data ->>'state' = '` + string(entity.ConnectionStateConfirmed) + `'
			AND (((json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT) OR
				((json_data->>'user_from_id')::BIGINT = $2::BIGINT AND (json_data->>'user_to_id')::BIGINT = $1::BIGINT))
			AND (json_data->>'enabled')::BOOL = true`

	connectionExistsQuery = `SELECT
		(count(*) > 0) :: BOOL AS "exists"
		FROM app_%d_%d.connections
		WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT
			AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT
			AND json_data->>'type' = $3::TEXT
			AND (json_data->>'enabled')::BOOL = true`

	userConnectionsByStateQuery = `SELECT
		json_data
		FROM app_%d_%d.connections
		WHERE ((json_data->>'user_to_id')::BIGINT = $1::BIGINT OR (json_data->>'user_from_id')::BIGINT = $1::BIGINT)
			AND json_data->>'state' = $2::TEXT
			AND (json_data->>'enabled')::BOOL = true`
)

func (c *connection) Create(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	if !connection.IsValidState() {
		return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage("got connection state: " + string(connection.State)).SetCurrentLocation()}
	}

	timeNow := time.Now()
	connection.CreatedAt, connection.UpdatedAt = &timeNow, &timeNow
	connection.Enabled = entity.PTrue
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
	if err != nil {
		return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return nil
}

func (c *connection) Read(accountID, applicationID int64, userFromID, userToID uint64, connectionType entity.ConnectionTypeType) (*entity.Connection, []errors.Error) {
	var JSONData string
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectConnectionQuery, accountID, applicationID), userFromID, userToID, string(connectionType)).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrConnectionNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(JSONData), connection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return connection, nil
}

func (c *connection) Update(
	accountID, applicationID int64,
	existingConnection, updatedConnection entity.Connection,
	retrieve bool,
) (*entity.Connection, []errors.Error) {

	timeNow := time.Now()
	updatedConnection.UpdatedAt = &timeNow
	connectionJSON, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = c.mainPg.Exec(
		appSchema(updateConnectionQuery, accountID, applicationID),
		string(connectionJSON), existingConnection.UserFromID, existingConnection.UserToID, string(existingConnection.Type))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	if !retrieve {
		return nil, nil
	}

	return &updatedConnection, nil
}

func (c *connection) Delete(accountID, applicationID int64, userFromID, userToID uint64, connectionType entity.ConnectionTypeType) []errors.Error {
	existingConnection, err := c.Read(accountID, applicationID, userFromID, userToID, connectionType)
	if err != nil {
		return err
	}

	existingConnection.Enabled = entity.PFalse
	_, err = c.Update(accountID, applicationID, *existingConnection, *existingConnection, false)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) Following(accountID, applicationID int64, userID uint64) (userIDs []uint64, er []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(followsQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowingList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowingList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		userIDs = append(userIDs, conn.UserToID)
	}

	return userIDs, nil
}

func (c *connection) FollowedBy(accountID, applicationID int64, userID uint64) (userIDs []uint64, er []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(followersQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowersList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowersList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		userIDs = append(userIDs, conn.UserFromID)
	}

	return userIDs, nil
}

func (c *connection) Friends(accountID, applicationID int64, userID uint64) (userIDs []uint64, er []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(friendConnectionsQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			userFromID uint64
			userToID   uint64
		)
		err := rows.Scan(&userFromID, &userToID)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		if userID != userFromID {
			userIDs = append(userIDs, userFromID)
		} else if userID != userToID {
			userIDs = append(userIDs, userToID)
		}
	}

	return userIDs, nil
}

func (c *connection) FriendsAndFollowing(accountID, applicationID int64, userID uint64) (users []uint64, er []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(friendAndFollowingConnectionsQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		users = append(users, conn.UserToID)
	}

	return users, nil
}

func (c *connection) FriendsAndFollowingIDs(accountID, applicationID int64, userID uint64) ([]uint64, []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(friendAndFollowingConnectionsIDsQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()

	var users []uint64
	for rows.Next() {
		var user uint64
		err := rows.Scan(&user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType entity.ConnectionTypeType, connectionState entity.ConnectionStateType) (users []uint64, er []errors.Error) {
	if len(socialFriendsIDs) == 0 {
		return
	}

	query, args, err := sqlx.In(fmt.Sprintf(listUserIDssBySocialIDQuery, accountID, applicationID, platform), socialFriendsIDs)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrServerInternalError.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	dbUsers, err := c.pg.SlaveDatastore(-1).
		Query(query, args...)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer dbUsers.Close()
	for dbUsers.Next() {
		var userID uint64
		err := dbUsers.Scan(&userID)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		users = append(users, userID)
	}

	return c.CreateMultiple(accountID, applicationID, user, connectionType, connectionState, users)
}

func (c *connection) CreateMultiple(accountID, applicationID int64, user *entity.ApplicationUser, connectionType entity.ConnectionTypeType, connectionState entity.ConnectionStateType, ourStoredUsersIDs []uint64) ([]uint64, []errors.Error) {
	if len(ourStoredUsersIDs) == 0 {
		return ourStoredUsersIDs, nil
	}

	for idx := range ourStoredUsersIDs {
		connection := &entity.Connection{
			UserFromID: user.ID,
			UserToID:   ourStoredUsersIDs[idx],
			Type:       connectionType,
			State:      connectionState,
		}

		if err := c.Create(accountID, applicationID, connection); err != nil {
			if err[0].Code() != errmsg.ErrConnectionAlreadyExists.Code() {
				return nil, err
			}
		}
	}

	return ourStoredUsersIDs, nil
}

func (c *connection) Relation(accountID, applicationID int64, userFromID, userToID uint64) (*entity.Relation, []errors.Error) {
	rel := &entity.Relation{
		IsFriend:   entity.PFalse,
		IsFollowed: entity.PFalse,
		IsFollower: entity.PFalse,
	}

	if userFromID == userToID {
		return rel, nil
	}

	relations, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(getUsersRelationQuery, accountID, applicationID), userFromID, userToID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrConnectionNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer relations.Close()

	var (
		relationFrom, relationTo uint64
		relationType             string
	)
	for relations.Next() {
		err := relations.Scan(&relationFrom, &relationTo, &relationType)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		if entity.ConnectionTypeType(relationType) == entity.ConnectionTypeFriend {
			rel.IsFriend = entity.PTrue
		}

		if relationFrom == userFromID && relationTo == userToID && entity.ConnectionTypeType(relationType) == entity.ConnectionTypeFollow {
			rel.IsFollowed = entity.PTrue
		}

		if relationFrom == userToID && relationTo == userFromID && entity.ConnectionTypeType(relationType) == entity.ConnectionTypeFollow {
			rel.IsFollower = entity.PTrue
		}
	}

	return rel, nil
}

func (c *connection) Exists(accountID, applicationID int64, userFromID, userToID uint64, connType entity.ConnectionTypeType) (bool, []errors.Error) {
	exists := false
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(connectionExistsQuery, accountID, applicationID), userFromID, userToID, string(connType)).
		Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, []errors.Error{errmsg.ErrConnectionNotFound.SetCurrentLocation()}
		}
		return false, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return exists, nil
}

func (c *connection) ConnectionsByState(accountID, applicationID int64, userID uint64, state entity.ConnectionStateType) ([]*entity.Connection, []errors.Error) {
	var connections []*entity.Connection
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(userConnectionsByStateQuery, accountID, applicationID), userID, string(state))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		connection := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), connection)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage("got error: " + err.Error()).SetCurrentLocation()}
		}

		connections = append(connections, connection)
	}

	return connections, nil
}

// NewConnection returns a new connection handler with PostgreSQL as storage driver
func NewConnection(pgsql postgres.Client) core.Connection {
	return &connection{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		appUser: &applicationUser{
			pg:     pgsql,
			mainPg: pgsql.MainDatastore(),
		},
	}
}
