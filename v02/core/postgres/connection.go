/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

	"encoding/json"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	connection struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

const (
	createConnectionQuery     = `INSERT INTO app_$1_$2.connections(json_data, enabled) VALUES ($3, $4)`
	selectConnectionQuery     = `SELECT json_data, enabled FROM app_$1_$2 WHERE json_data @> '{"user_from_id": $3, "user_to_id": $4}'`
	updateConnectionQuery     = `UPDATE app_$1_$2.connections SET json_data = $3 WHERE json_data @> '{"user_from_id": $3, "user_to_id": $4}'`
	deleteConnectionQuery     = `UPDATE app_$1_$2.connections SET enabled = 0 WHERE json_data @> '{"user_from_id": $3, "user_to_id": $4}'`
	listConnectionQuery       = `SELECT json_data, enabled FROM app_$1_$2.connections WHERE json_data @> '{"user_from_id": $3}'`
	followedByConnectionQuery = `SELECT json_data, enabled FROM app_$1_$2.connections WHERE json_data @> '{"user_to_id": $3}'`
)

func (c *connection) Create(connection *entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connectionJSON, err := json.Marshal(connection)
	if err != nil {
		return nil, errors.NewInternalError("error while saving the connection", err.Error())
	}
	_, err = c.mainPg.Exec(createConnectionQuery, connection.AccountID, connection.ApplicationID, string(connectionJSON), connection.Enabled)
	if err != nil {
		return nil, errors.NewInternalError("error while saving the connection", err.Error())
	}

	if !retrieve {
		return nil, nil
	}
	return c.Read(connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)
}

func (c *connection) Read(accountID, applicationID, userFromID, userToID int64) (*entity.Connection, errors.Error) {
	var (
		JSONData string
		Enabled  bool
	)
	err := c.pg.SlaveDatastore(-1).
		QueryRow(selectConnectionQuery, accountID, applicationID, userFromID, userToID).
		Scan(&JSONData, &Enabled)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the connection", err.Error())
	}

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(JSONData), connection)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the connection", err.Error())
	}
	connection.Enabled = Enabled

	return connection, nil
}

func (c *connection) Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connectionJSON, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the connection", err.Error())
	}

	_, err = c.mainPg.Exec(updateConnectionQuery, existingConnection.AccountID, existingConnection.ApplicationID, string(connectionJSON))
	if err != nil {
		return nil, errors.NewInternalError("error while updating the connection", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return c.Read(existingConnection.AccountID, existingConnection.ApplicationID, existingConnection.UserFromID, existingConnection.UserToID)
}

func (c *connection) Delete(connection *entity.Connection) errors.Error {
	_, err := c.mainPg.Exec(deleteConnectionQuery, connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)
	if err != nil {
		return errors.NewInternalError("error while deleting the connection", err.Error())
	}

	return nil
}

func (c *connection) List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, er errors.Error) {
	users = []*entity.ApplicationUser{}

	rows, err := c.pg.SlaveDatastore(-1).
		Query(listConnectionQuery, accountID, applicationID, userID)
	if err != nil {
		return users, errors.NewInternalError("error while retrieving list of account users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ID       int64
			JSONData string
			Enabled  bool
		)
		err := rows.Scan(&ID, &JSONData, &Enabled)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of connections", err.Error())
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of connections", err.Error())
		}
		user.ID = ID
		user.Enabled = Enabled

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) FollowedBy(accountID, applicationID, userID int64) ([]*entity.ApplicationUser, errors.Error) {
	users := []*entity.ApplicationUser{}

	rows, err := c.pg.SlaveDatastore(-1).
		Query(followedByConnectionQuery, accountID, applicationID, userID)
	if err != nil {
		return users, errors.NewInternalError("error while retrieving list of account users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ID       int64
			JSONData string
			Enabled  bool
		)
		err := rows.Scan(&ID, &JSONData, &Enabled)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of followers", err.Error())
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of followers", err.Error())
		}
		user.ID = ID
		user.Enabled = Enabled

		users = append(users, user)
	}

	return users, nil
}

func (c *connection) Confirm(connection *entity.Connection, retrieve bool) (*entity.Connection, errors.Error) {
	connection.Enabled = true
	connection.ConfirmedAt = time.Now()
	connection.UpdatedAt = connection.ConfirmedAt

	return c.Update(*connection, *connection, retrieve)
}

func (c *connection) WriteEventsToList(connection *entity.Connection) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string) ([]*entity.ApplicationUser, errors.Error) {
	return []*entity.ApplicationUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) AutoConnectSocialFriends(user *entity.ApplicationUser, ourStoredUsersIDs []interface{}) (users []*entity.ApplicationUser, err errors.Error) {
	return []*entity.ApplicationUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewConnection returns a new connection handler with PostgreSQL as storage driver
func NewConnection(pgsql postgres.Client) core.Connection {
	return &connection{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
