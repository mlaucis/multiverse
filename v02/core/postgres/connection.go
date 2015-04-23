/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

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

func (c *connection) Create(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) Delete(*entity.Connection) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error) {
	return []*entity.ApplicationUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error) {
	return []*entity.ApplicationUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (c *connection) Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
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
