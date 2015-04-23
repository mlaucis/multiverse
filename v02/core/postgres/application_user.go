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
	applicationUser struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

func (au *applicationUser) Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) Delete(*entity.ApplicationUser) (err errors.Error) {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err errors.Error) {
	return []*entity.ApplicationUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) CreateSession(user *entity.ApplicationUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) GetSession(user *entity.ApplicationUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) DestroySession(sessionToken string, user *entity.ApplicationUser) errors.Error {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) ExistsByUsername(accountID, applicationID int64, email string) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *applicationUser) ExistsByID(accountID, applicationID, userID int64) bool {
	panic("not implemented yet")
}

// NewApplicationUser returns a new application user handler with PostgreSQL as storage driver
func NewApplicationUser(pgsql postgres.Client) core.ApplicationUser {
	return &applicationUser{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
