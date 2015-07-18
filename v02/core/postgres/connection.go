/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	connection struct {
		pg      postgres.Client
		mainPg  *sqlx.DB
		appUser core.ApplicationUser
	}
)

const (
	createConnectionQuery              = `INSERT INTO app_%d_%d.connections(json_data) VALUES ($1)`
	selectConnectionQuery              = `SELECT json_data FROM app_%d_%d.connections WHERE json_data @> json_build_object('user_from_id', $1::bigint, 'user_to_id', $2::bigint)::jsonb LIMIT 1`
	updateConnectionQuery              = `UPDATE app_%d_%d.connections SET json_data = $1 WHERE json_data @> json_build_object('user_from_id', $2::bigint, 'user_to_id', $3::bigint)::jsonb`
	followsQuery                       = `SELECT json_data FROM app_%d_%d.connections WHERE json_data @> json_build_object('user_from_id', $1::bigint, 'type', 'follow', 'enabled', true)::jsonb`
	followersQuery                     = `SELECT json_data FROM app_%d_%d.connections WHERE json_data @> json_build_object('user_to_id', $1::bigint, 'type', 'follow', 'enabled', true)::jsonb`
	friendConnectionsQuery             = `SELECT json_data FROM app_%d_%d.connections WHERE json_data @> json_build_object('user_to_id', $1::bigint, 'type', 'friend', 'enabled', true)::jsonb`
	friendAndFollowingConnectionsQuery = `SELECT json_data FROM app_%d_%d.connections WHERE json_data @> json_build_object('user_from_id', $1::bigint, 'enabled', true)::jsonb`
	listUsersBySocialIDQuery           = `SELECT json_data FROM app_%d_%d.users WHERE %s`
)

func (c *connection) Create(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	exists, er := c.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
	if er != nil && er[0] != errmsg.ErrConnectionNotFound {
		return nil, er
	}
	if exists != nil {
		if !exists.Enabled {
			exists.Enabled = true
			return c.Update(accountID, applicationID, *exists, *exists, true)
		}

		return nil, []errors.Error{errmsg.ErrConnectionAlreadyExists}
	}

	timeNow := time.Now()
	connection.CreatedAt, connection.UpdatedAt = &timeNow, &timeNow
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error())}
	}
	_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error())}
	}

	if connection.Type == "friend" {
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID
		connectionJSON, err = json.Marshal(connection)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error())}
		}
		_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectionCreation.UpdateInternalMessage(err.Error())}
		}
		// Switch back so we have the original IDs in place
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID
	}

	if !retrieve {
		return nil, nil
	}
	return c.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
}

func (c *connection) Read(accountID, applicationID int64, userFromID, userToID uint64) (*entity.Connection, []errors.Error) {
	var JSONData string
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectConnectionQuery, accountID, applicationID), userFromID, userToID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrConnectionNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error())}
	}

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(JSONData), connection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionRead.UpdateInternalMessage(err.Error())}
	}

	return connection, nil
}

func (c *connection) Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	timeNow := time.Now()
	updatedConnection.UpdatedAt = &timeNow
	connectionJSON, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = c.mainPg.Exec(
		appSchema(updateConnectionQuery, accountID, applicationID),
		string(connectionJSON), existingConnection.UserFromID, existingConnection.UserToID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectionUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return c.Read(accountID, applicationID, existingConnection.UserFromID, existingConnection.UserToID)
}

func (c *connection) Delete(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	connection.Enabled = false
	_, err := c.Update(accountID, applicationID, *connection, *connection, false)
	if err != nil {
		return err
	}

	if connection.Type == "friend" {
		connection.UserFromID, connection.UserToID = connection.UserToID, connection.UserFromID
		_, err = c.Update(accountID, applicationID, *connection, *connection, false)
	}

	return err
}

func (c *connection) List(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, er []errors.Error) {
	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(followsQuery, accountID, applicationID), userID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowingList.UpdateInternalMessage(err.Error())}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowingList.UpdateInternalMessage(err.Error())}
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID)
		if er != nil {
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowersList.UpdateInternalMessage(err.Error())}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFollowersList.UpdateInternalMessage(err.Error())}
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserFromID)
		if er != nil {
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error())}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error())}
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserFromID)
		if er != nil {
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()

	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error())}
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalFriendsList.UpdateInternalMessage(err.Error())}
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID)
		if er != nil {
			return nil, er
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

	if connection.Type == "friend" {
		con := *connection
		con.UserFromID, con.UserToID = connection.UserToID, connection.UserFromID
		_, err = c.Update(accountID, applicationID, con, con, retrieve)
	}

	return conn, err
}

func (c *connection) WriteEventsToList(accountID, applicationID int64, connection *entity.Connection) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID int64, userFromID, userToID uint64) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

	if len(socialFriendsIDs) == 0 {
		return users, nil
	}

	var conditions []string
	for idx := range socialFriendsIDs {
		conditions = append(conditions, fmt.Sprintf(`json_data @> '{"social_ids": {%q: %q}}'`, platform, socialFriendsIDs[idx]))
	}

	dbUsers, err := c.pg.SlaveDatastore(-1).
		Query(fmt.Sprintf(listUsersBySocialIDQuery, accountID, applicationID, strings.Join(conditions, " OR ")))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error())}
	}
	defer dbUsers.Close()
	for dbUsers.Next() {
		var JSONData string
		err := dbUsers.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error())}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalConnectingUsers.UpdateInternalMessage(err.Error())}
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

		if _, err := c.Create(accountID, applicationID, connection, false); err != nil {
			if err[0] != errmsg.ErrConnectionAlreadyExists {
				return nil, err
			}
		}

		if _, err := c.Confirm(accountID, applicationID, connection, false); err != nil {
			return nil, err
		}

		if connectionType != "friend" {
			continue
		}

		connection = &entity.Connection{
			UserFromID: ourStoredUsersIDs[idx].ID,
			UserToID:   user.ID,
			Type:       connectionType,
		}

		if _, err := c.Create(accountID, applicationID, connection, false); err != nil {
			if err[0] != errmsg.ErrConnectionAlreadyExists {
				return nil, err
			}
		}

		if _, err := c.Confirm(accountID, applicationID, connection, false); err != nil {
			return nil, err
		}
	}

	return ourStoredUsersIDs, nil
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
