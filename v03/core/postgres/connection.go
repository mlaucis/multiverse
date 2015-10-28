package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type connection struct {
	pg      postgres.Client
	mainPg  *sqlx.DB
	appUser core.ApplicationUser
}

const (
	createConnectionQuery                 = `INSERT INTO app_%d_%d.connections(json_data) VALUES ($1)`
	selectConnectionQuery                 = `SELECT json_data FROM app_%d_%d.connections WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT AND (json_data->>'enabled')::BOOL = true LIMIT 1`
	updateConnectionQuery                 = `UPDATE app_%d_%d.connections SET json_data = $1 WHERE (json_data->>'user_from_id')::BIGINT = $2::BIGINT AND (json_data->>'user_to_id')::BIGINT = $3::BIGINT`
	followsQuery                          = `SELECT json_data FROM app_%d_%d.connections WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND json_data->>'type' = '` + entity.ConnectionTypeFollow + `' AND (json_data->>'enabled')::BOOL = true`
	followersQuery                        = `SELECT json_data FROM app_%d_%d.connections WHERE (json_data->>'user_to_id')::BIGINT = $1::BIGINT AND json_data->>'type' = '` + entity.ConnectionTypeFollow + `' AND (json_data->>'enabled')::BOOL = true`
	friendConnectionsQuery                = `SELECT json_data FROM app_%d_%d.connections WHERE (json_data->>'user_to_id')::BIGINT = $1::BIGINT AND json_data->>'type' = '` + entity.ConnectionTypeFriend + `' AND (json_data->>'enabled')::BOOL = true`
	friendAndFollowingConnectionsQuery    = `SELECT json_data FROM app_%d_%d.connections WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'enabled')::BOOL = true`
	friendAndFollowingConnectionsIDsQuery = `SELECT json_data->>'user_to_id' as "user_id" FROM app_%d_%d.connections WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'enabled')::BOOL = true`
	listUsersBySocialIDQuery              = `SELECT json_data FROM app_%d_%d.users WHERE (json_data->>'enabled')::BOOL = true AND (json_data->>'deleted')::BOOL = false AND json_data->'social_ids'->>'%s' IN (?)`

	getUsersRelationQuery = `SELECT
  json_data ->> 'user_from_id' AS "from",
  json_data ->> 'user_to_id'   AS "to",
  json_data ->> 'type'         AS "type"
FROM app_%d_%d.connections
WHERE json_data @> '{"enabled": true}'
      AND (((json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT) OR
           (json_data->>'user_from_id')::BIGINT = $2::BIGINT AND (json_data->>'user_to_id')::BIGINT = $1::BIGINT)`

	connectionExistsQuery = `SELECT
  (count(*) > 0) :: BOOL AS "exists"
FROM app_%d_%d.connections
WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'user_to_id')::BIGINT = $2::BIGINT AND json_data @> json_build_object('type', $3::TEXT, 'enabled', true)::JSONB;`
)

func (c *connection) Create(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	// Check if the connection already exists between users
	exists, er := c.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
	if er != nil && er[0].Code() != errmsg.ErrConnectionNotFound.Code() {
		return er
	}

	// If it exists and it's not enabled then enable it
	if exists == nil {
		timeNow := time.Now()
		connection.CreatedAt, connection.UpdatedAt = &timeNow, &timeNow
		connection.Enabled = true
		connection.State = "confirmed"
		connectionJSON, err := json.Marshal(connection)
		if err != nil {
			return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
		if err != nil {
			return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
	}

	// if the connection is of type friend then reverse the roles and create the other connection
	if connection.Type == entity.ConnectionTypeFriend {
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID

		// Check if the connection exists
		exists, er := c.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
		if er != nil && er[0].Code() != errmsg.ErrConnectionNotFound.Code() {
			return er
		}

		// If it doesn't exists then create it
		if exists == nil {
			connectionJSON, err := json.Marshal(connection)
			if err != nil {
				return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
			}
			_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
			if err != nil {
				return []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
			}
		}

		// Switch back so we have the original IDs in place
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID
	}

	return nil
}

func (c *connection) Read(accountID, applicationID int64, userFromID, userToID uint64) (*entity.Connection, []errors.Error) {
	var JSONData string
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectConnectionQuery, accountID, applicationID), userFromID, userToID).
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

func (c *connection) Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	timeNow := time.Now()
	updatedConnection.UpdatedAt = &timeNow
	updatedConnection.State = "confirmed"
	connectionJSON, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = c.mainPg.Exec(
		appSchema(updateConnectionQuery, accountID, applicationID),
		string(connectionJSON), existingConnection.UserFromID, existingConnection.UserToID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	if !retrieve {
		return nil, nil
	}

	return &updatedConnection, nil
}

func (c *connection) Delete(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	connection.Enabled = false
	_, err := c.Update(accountID, applicationID, *connection, *connection, false)
	if err != nil {
		return err
	}

	if connection.Type == entity.ConnectionTypeFriend {
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID
		_, err = c.Update(accountID, applicationID, *connection, *connection, false)
	}

	return err
}

func (c *connection) List(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, er []errors.Error) {
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
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID, false)
		if er != nil {
			if er[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, er
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) FollowedBy(accountID, applicationID int64, userID uint64) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

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
		user, er := c.appUser.Read(accountID, applicationID, conn.UserFromID, false)
		if er != nil {
			if er[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, er
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) Friends(accountID, applicationID int64, userID uint64) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(friendConnectionsQuery, accountID, applicationID), userID)
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
		user, er := c.appUser.Read(accountID, applicationID, conn.UserFromID, false)
		if er != nil {
			if er[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, er
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) FriendsAndFollowing(accountID, applicationID int64, userID uint64) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

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
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID, false)
		if er != nil {
			if er[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, er
		}

		users = append(users, user)
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

func (c *connection) Confirm(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	connection.Enabled = true
	timeNow := time.Now()
	connection.ConfirmedAt, connection.UpdatedAt = &timeNow, &timeNow

	conn, err := c.Update(accountID, applicationID, *connection, *connection, retrieve)
	if err != nil {
		return conn, err
	}

	if connection.Type == entity.ConnectionTypeFriend {
		con := *connection
		con.UserFromID, con.UserToID = connection.UserToID, connection.UserFromID
		_, err = c.Update(accountID, applicationID, con, con, retrieve)
	}

	return conn, err
}

func (c *connection) WriteEventsToList(accountID, applicationID int64, connection *entity.Connection) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID int64, userFromID, userToID uint64) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

	if len(socialFriendsIDs) == 0 {
		return users, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(listUsersBySocialIDQuery, accountID, applicationID, platform), socialFriendsIDs)
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
		var JSONData string
		err := dbUsers.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		users = append(users, user)
	}

	return c.AutoConnectSocialFriends(accountID, applicationID, user, connectionType, users)
}

func (c *connection) AutoConnectSocialFriends(accountID, applicationID int64, user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, []errors.Error) {
	if len(ourStoredUsersIDs) == 0 {
		return ourStoredUsersIDs, nil
	}

	for idx := range ourStoredUsersIDs {
		connection := &entity.Connection{
			UserFromID: user.ID,
			UserToID:   ourStoredUsersIDs[idx].ID,
			Type:       connectionType,
		}
		connection.Enabled = true

		if err := c.Create(accountID, applicationID, connection); err != nil {
			if err[0].Code() != errmsg.ErrConnectionAlreadyExists.Code() {
				return nil, err
			}
		}

		if connectionType != entity.ConnectionTypeFriend {
			continue
		}

		connection = &entity.Connection{
			UserFromID: ourStoredUsersIDs[idx].ID,
			UserToID:   user.ID,
			Type:       connectionType,
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
	if userFromID == userToID {
		return &entity.Relation{}, nil
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

	rel := &entity.Relation{
		IsFriend:   entity.PFalse,
		IsFollowed: entity.PFalse,
		IsFollower: entity.PFalse,
	}
	var (
		relationFrom, relationTo uint64
		relationType             string
	)
	for relations.Next() {
		err := relations.Scan(&relationFrom, &relationTo, &relationType)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		if relationType == entity.ConnectionTypeFriend {
			rel.IsFriend = entity.PTrue
		}

		if relationFrom == userFromID && relationTo == userToID && relationType == entity.ConnectionTypeFollow {
			rel.IsFollowed = entity.PTrue
		}

		if relationFrom == userToID && relationTo == userFromID && relationType == entity.ConnectionTypeFollow {
			rel.IsFollower = entity.PTrue
		}
	}

	return rel, nil
}

func (c *connection) Exists(accountID, applicationID int64, userFromID, userToID uint64, connType string) (bool, []errors.Error) {
	exists := false
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(connectionExistsQuery, accountID, applicationID), userFromID, userToID, connType).
		Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, []errors.Error{errmsg.ErrConnectionNotFound.SetCurrentLocation()}
		}
		return false, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return exists, nil
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
