/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

	"encoding/json"

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

const (
	insertAccountQuery     = "INSERT INTO accounts(authkey, data) VALUES ($1, $2) RETURNING id"
	selectAccountByIDQuery = "SELECT data FROM accounts WHERE id = $1"
	updateAccountByIDQuery = "UPDATE accounts SET authkey = $1, data = $2 WHERE id = $3"
	deleteAccountByIDQuery = "DELETE FROM accounts WHERE id = $1"
)

func (a *account) Create(account *entity.Account, retrieve bool) (*entity.Account, errors.Error) {
	// TODO we should generate the account auth key here... it would be nice to do so
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the account", err.Error())
	}

	var createdAccountID int64
	err = a.mainPg.
		QueryRow(insertAccountQuery, account.AuthToken, accountJSON).
		Scan(&createdAccountID)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the account", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return a.Read(createdAccountID)
}

func (a *account) Read(accountID int64) (*entity.Account, errors.Error) {
	var accountJSON string
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(&accountJSON)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account", err.Error())
	}

	acc := &entity.Account{}
	err = json.Unmarshal([]byte(accountJSON), acc)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account", err.Error())
	}

	return acc, nil
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (*entity.Account, errors.Error) {
	if updatedAccount.AuthToken == "" {
		// TODO we should regenerate the account key here somehow?
		updatedAccount.AuthToken = existingAccount.AuthToken
	}
	accountJSON, err := json.Marshal(updatedAccount)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the account", err.Error())
	}

	_, err = a.mainPg.Exec(updateAccountByIDQuery, updatedAccount.AuthToken, accountJSON, existingAccount.ID)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the account", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return a.Read(existingAccount.ID)
}

func (a *account) Delete(acc *entity.Account) errors.Error {
	_, err := a.mainPg.Exec(deleteAccountByIDQuery, acc.ID)
	if err != nil {
		return errors.NewInternalError("error while deleting the account", err.Error())
	}
	return nil
}

func (a *account) Exists(accountID int64) bool {
	acc, err := a.Read(accountID)
	if err != nil {
		return false
	}
	return acc != nil
}

// NewAccount returns a new account handler with PostgreSQL as storage driver
func NewAccount(pgsql postgres.Client) core.Account {
	return &account{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
