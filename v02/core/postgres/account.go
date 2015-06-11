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
	"github.com/tapglue/backend/v02/errmsg"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	account struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

const (
	createAccountQuery           = `INSERT INTO tg.accounts(json_data) VALUES ($1) RETURNING id`
	selectAccountByIDQuery       = `SELECT json_data FROM tg.accounts WHERE id = $1 LIMIT 1`
	selectAccountByKeyQuery      = `SELECT id, json_data FROM tg.accounts WHERE json_data @> json_build_object('token', $1::text)::jsonb LIMIT 1`
	selectAccountByPublicIDQuery = `SELECT id, json_data FROM tg.accounts WHERE json_data @> json_build_object('id', $1::text)::jsonb LIMIT 1`
	updateAccountByIDQuery       = `UPDATE tg.accounts SET json_data = $1 WHERE id = $2`
	deleteAccountByIDQuery       = `DELETE FROM tg.accounts WHERE id = $1`
)

func (a *account) Create(account *entity.Account, retrieve bool) (*entity.Account, []errors.Error) {
	account.PublicID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	account.Enabled = true
	timeNow := time.Now()
	account.CreatedAt, account.UpdatedAt = &timeNow, &timeNow
	account.AuthToken = storageHelper.GenerateAccountSecretKey(account)

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountCreationError.SetCurrentLocation().UpdateInternalMessage(err.Error())}
	}

	var createdAccountID int64
	err = a.mainPg.
		QueryRow(createAccountQuery, string(accountJSON)).
		Scan(&createdAccountID)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountCreationError.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return a.Read(createdAccountID)
}

func (a *account) Read(accountID int64) (*entity.Account, []errors.Error) {
	var JSONData string
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}

	acc := &entity.Account{}
	err = json.Unmarshal([]byte(JSONData), acc)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	acc.ID = accountID

	return acc, nil
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (*entity.Account, []errors.Error) {
	if updatedAccount.AuthToken == "" {
		updatedAccount.AuthToken = existingAccount.AuthToken
	}
	timeNow := time.Now()
	updatedAccount.UpdatedAt = &timeNow
	accountJSON, err := json.Marshal(updatedAccount)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountUpdateError.UpdateInternalMessage(err.Error())}
	}

	_, err = a.mainPg.Exec(updateAccountByIDQuery, string(accountJSON), existingAccount.ID)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountUpdateError.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return a.Read(existingAccount.ID)
}

func (a *account) Delete(account *entity.Account) []errors.Error {
	_, err := a.mainPg.Exec(deleteAccountByIDQuery, account.ID)
	if err != nil {
		return []errors.Error{errmsg.InternalAccountDeleteError.UpdateInternalMessage(err.Error())}
	}
	return nil
}

func (a *account) Exists(accountID int64) (bool, []errors.Error) {
	var JSONData string
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(&JSONData)
	if err != nil {
		return false, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (a *account) FindByKey(authKey string) (*entity.Account, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByKeyQuery, authKey).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	account := &entity.Account{}
	err = json.Unmarshal([]byte(JSONData), account)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	account.ID = ID

	return account, nil
}

func (a *account) ReadByPublicID(id string) (*entity.Account, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := a.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByPublicIDQuery, id).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.AccountNotFoundError}
		}
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	account := &entity.Account{}
	err = json.Unmarshal([]byte(JSONData), account)
	if err != nil {
		return nil, []errors.Error{errmsg.InternalAccountReadError.UpdateInternalMessage(err.Error())}
	}
	account.ID = ID

	return account, nil
}

// NewAccount returns a new account handler with PostgreSQL as storage driver
func NewAccount(pgsql postgres.Client) core.Account {
	return &account{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
