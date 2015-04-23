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
	accountUser struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

func (au *accountUser) Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) Delete(*entity.AccountUser) errors.Error {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) List(accountID int64) (accountUsers []*entity.AccountUser, er errors.Error) {
	return []*entity.AccountUser{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) CreateSession(user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) RefreshSession(sessionToken string, user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) DestroySession(sessionToken string, user *entity.AccountUser) errors.Error {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) GetSession(user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) FindByEmail(email string) (*entity.Account, *entity.AccountUser, errors.Error) {
	return nil, nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) ExistsByEmail(email string) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) FindByUsername(username string) (*entity.Account, *entity.AccountUser, errors.Error) {
	return nil, nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) ExistsByUsername(username string) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) ExistsByID(accountID, accountUserID int64) bool {
	panic("not implemented yet")
}

// NewAccountUser returns a new account user handler with PostgreSQL as storage driver
func NewAccountUser(pgsql postgres.Client) core.AccountUser {
	return &accountUser{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
