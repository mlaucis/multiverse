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
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	connection struct {
		pg      postgres.Client
		mainPg  *sql.DB
		appUser core.ApplicationUser
	}
)

const (
	createConnectionQuery     = `INSERT INTO app_%d_%d.connections(json_data) VALUES ($1, $2)`
	selectConnectionQuery     = `SELECT json_data FROM app_%d_%d.connections WHERE json_data->>'user_from_id' = $1 AND json_data->>'user_to_id' = $2`
	updateConnectionQuery     = `UPDATE app_%d_%d.connections SET json_data = $1 WHERE json_data->>'user_from_id' = $2 AND json_data->>'user_to_id' = $3`
	listConnectionQuery       = `SELECT json_data FROM app_%d_%d.connections WHERE json_data->>'user_from_id' = $1`
	followedByConnectionQuery = `SELECT json_data FROM app_%d_%d.connections WHERE json_data->>'user_to_id' = $1`
	listUsersBySocialIDQuery  = `SELECT json_data FROM app_%d_%d.users WHERE %s`
)

func (c *connection) Create(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return nil, errors.NewInternalError("error while saving the connection", err.Error())
	}
	_, err = c.mainPg.Exec(appSchema(createConnectionQuery, accountID, applicationID), string(connectionJSON))
	if err != nil {
		return nil, errors.NewInternalError("error while saving the connection", err.Error())
	}

	if !retrieve {
		return nil, nil
	}
	return c.Read(accountID, applicationID, connection.UserFromID, connection.UserToID)
}

func (c *connection) Read(accountID, applicationID, userFromID, userToID int64) (*entity.Connection, errors.Error) {
	var JSONData string
	err := c.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectConnectionQuery, accountID, applicationID), userFromID, userToID).
		Scan(&JSONData)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the connection", err.Error())
	}

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(JSONData), connection)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the connection", err.Error())
	}

	return connection, nil
}

func (c *connection) Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connectionJSON, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the connection", err.Error())
	}

	_, err = c.mainPg.Exec(appSchema(updateConnectionQuery, accountID, applicationID), string(connectionJSON), existingConnection.UserFromID, existingConnection.UserToID)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the connection", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return c.Read(accountID, applicationID, existingConnection.UserFromID, existingConnection.UserToID)
}

func (c *connection) Delete(accountID, applicationID int64, connection *entity.Connection) errors.Error {
	connection.Enabled = false
	_, err := c.Update(accountID, applicationID, *connection, *connection, false)

	return err
}

func (c *connection) List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, er errors.Error) {
	users = []*entity.ApplicationUser{}

	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(listConnectionQuery, accountID, applicationID), userID)
	if err != nil {
		return users, errors.NewInternalError("error while retrieving list of account users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of following", err.Error())
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of following", err.Error())
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID)
		if er != nil {
			return []*entity.ApplicationUser{}, er
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) FollowedBy(accountID, applicationID, userID int64) ([]*entity.ApplicationUser, errors.Error) {
	users := []*entity.ApplicationUser{}

	rows, err := c.pg.SlaveDatastore(-1).
		Query(appSchema(followedByConnectionQuery, accountID, applicationID), userID)
	if err != nil {
		return users, errors.NewInternalError("error while retrieving list of account users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var JSONData string
		err := rows.Scan(&JSONData)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of followers", err.Error())
		}
		conn := &entity.Connection{}
		err = json.Unmarshal([]byte(JSONData), conn)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of followers", err.Error())
		}
		user, er := c.appUser.Read(accountID, applicationID, conn.UserToID)
		if er != nil {
			return []*entity.ApplicationUser{}, er
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) Confirm(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connection.Enabled = true
	connection.ConfirmedAt = time.Now()
	connection.UpdatedAt = connection.ConfirmedAt

	return c.Update(accountID, applicationID, *connection, *connection, retrieve)
}

func (c *connection) WriteEventsToList(accountID, applicationID int64, connection *entity.Connection) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, errors.Error) {
	users := []*entity.ApplicationUser{}

	var conditions []string
	for idx := range socialFriendsIDs {
		conditions = append(conditions, fmt.Sprintf(`json_data @> '{"social_ids": {%q: %q}}'`, platform, socialFriendsIDs[idx]))
	}

	dbUsers, err := c.pg.SlaveDatastore(-1).
		Query(fmt.Sprintf(listUsersBySocialIDQuery, accountID, applicationID, strings.Join(conditions, " OR ")))
	if err != nil {
		return users, errors.NewInternalError("error while connecting the users", err.Error())
	}
	defer dbUsers.Close()
	for dbUsers.Next() {
		var JSONData string
		err := dbUsers.Scan(&JSONData)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while connecting the users", err.Error())
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while connecting the users", err.Error())
		}
		users = append(users, user)
	}

	return c.AutoConnectSocialFriends(accountID, applicationID, user, connectionType, users)
}

func (c *connection) AutoConnectSocialFriends(accountID, applicationID int64, user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, errors.Error) {
	if len(ourStoredUsersIDs) == 0 {
		return ourStoredUsersIDs, nil
	}

	for idx := range ourStoredUsersIDs {
		connection := &entity.Connection{
			UserFromID: user.ID,
			UserToID:   ourStoredUsersIDs[idx].ID,
		}

		if _, err := c.Create(accountID, applicationID, connection, false); err != nil {
			return nil, err
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
		}

		if _, err := c.Create(accountID, applicationID, connection, false); err != nil {
			return nil, err
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
