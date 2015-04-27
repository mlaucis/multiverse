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
	createAccountQuery     = `INSERT INTO accounts(json_data) VALUES ($1, $2) RETURNING id`
	selectAccountByIDQuery = `SELECT id, json_data FROM accounts WHERE id = $1`
	updateAccountByIDQuery = `UPDATE accounts SET json_data = $1 WHERE id = $2`
	deleteAccountByIDQuery = `DELETE FROM accounts WHERE id = $1`
)

func (a *account) Create(account *entity.Account, retrieve bool) (*entity.Account, errors.Error) {
	// TODO we should generate the account auth key here... it would be nice to do so
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the account", err.Error())
	}

	var createdAccountID int64
	err = a.mainPg.
		QueryRow(createAccountQuery, string(accountJSON)).
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
	accountJSON := &struct {
		ID       int64
		JSONData string
	}{}
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(accountJSON)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account", err.Error())
	}

	acc := &entity.Account{}
	err = json.Unmarshal([]byte(accountJSON.JSONData), acc)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the account", err.Error())
	}
	acc.ID = accountID

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

	_, err = a.mainPg.Exec(updateAccountByIDQuery, string(accountJSON), existingAccount.ID)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the account", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return a.Read(existingAccount.ID)
}

func (a *account) Delete(account *entity.Account) errors.Error {
	_, err := a.mainPg.Exec(deleteAccountByIDQuery, account.ID)
	if err != nil {
		return errors.NewInternalError("error while deleting the account", err.Error())
	}
	return nil
}

func (a *account) Exists(accountID int64) (bool, errors.Error) {
	accountJSON := &struct {
		ID       int64
		JSONData string
	}{}
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(accountJSON)
	if err != nil {
		return false, errors.NewInternalError("error while reading the account", err.Error())
	}
	return true, nil
}

// NewAccount returns a new account handler with PostgreSQL as storage driver
func NewAccount(pgsql postgres.Client) core.Account {
	return &account{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
