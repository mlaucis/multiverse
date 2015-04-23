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
	account struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

func (a *account) Create(account *entity.Account, retrieve bool) (*entity.Account, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) Read(accountID int64) (*entity.Account, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (*entity.Account, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) Delete(*entity.Account) errors.Error {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) Exists(accountID int64) bool {
	panic("not implemented yet")
}

// NewAccount returns a new account handler with PostgreSQL as storage driver
func NewAccount(pgsql postgres.Client) core.Account {
	return &account{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
