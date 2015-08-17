package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	storageHelper "github.com/tapglue/backend/v03/storage/helper"
	"github.com/tapglue/backend/v03/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type organization struct {
	pg     postgres.Client
	mainPg *sqlx.DB
}

const (
	createAccountQuery           = `INSERT INTO tg.accounts(json_data) VALUES ($1) RETURNING id`
	selectAccountByIDQuery       = `SELECT json_data FROM tg.accounts WHERE id = $1 LIMIT 1`
	selectAccountByKeyQuery      = `SELECT id, json_data FROM tg.accounts WHERE json_data @> json_build_object('token', $1::text)::jsonb LIMIT 1`
	selectAccountByPublicIDQuery = `SELECT id, json_data FROM tg.accounts WHERE json_data @> json_build_object('id', $1::text)::jsonb LIMIT 1`
	updateAccountByIDQuery       = `UPDATE tg.accounts SET json_data = $1 WHERE id = $2`
	deleteAccountByIDQuery       = `DELETE FROM tg.accounts WHERE id = $1`
)

func (org *organization) Create(account *entity.Organization, retrieve bool) (*entity.Organization, []errors.Error) {
	account.PublicID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	account.Enabled = true
	timeNow := time.Now()
	account.CreatedAt, account.UpdatedAt = &timeNow, &timeNow
	account.AuthToken = storageHelper.GenerateAccountSecretKey(account)

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountCreation.SetCurrentLocation().UpdateInternalMessage(err.Error())}
	}

	var createdAccountID int64
	err = org.mainPg.
		QueryRow(createAccountQuery, string(accountJSON)).
		Scan(&createdAccountID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return org.Read(createdAccountID)
}

func (org *organization) Read(accountID int64) (*entity.Organization, []errors.Error) {
	var JSONData string
	err := org.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}

	acc := &entity.Organization{}
	err = json.Unmarshal([]byte(JSONData), acc)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	acc.ID = accountID

	return acc, nil
}

func (org *organization) Update(existingAccount, updatedAccount entity.Organization, retrieve bool) (*entity.Organization, []errors.Error) {
	if updatedAccount.AuthToken == "" {
		updatedAccount.AuthToken = existingAccount.AuthToken
	}
	timeNow := time.Now()
	updatedAccount.UpdatedAt = &timeNow
	accountJSON, err := json.Marshal(updatedAccount)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = org.mainPg.Exec(updateAccountByIDQuery, string(accountJSON), existingAccount.ID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return org.Read(existingAccount.ID)
}

func (org *organization) Delete(account *entity.Organization) []errors.Error {
	_, err := org.mainPg.Exec(deleteAccountByIDQuery, account.ID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalAccountDelete.UpdateInternalMessage(err.Error())}
	}
	return nil
}

func (org *organization) Exists(accountID int64) (bool, []errors.Error) {
	var JSONData string
	err := org.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByIDQuery, accountID).
		Scan(&JSONData)
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (org *organization) FindByKey(authKey string) (*entity.Organization, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := org.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByKeyQuery, authKey).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	account := &entity.Organization{}
	err = json.Unmarshal([]byte(JSONData), account)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	account.ID = ID

	return account, nil
}

func (org *organization) ReadByPublicID(id string) (*entity.Organization, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := org.pg.SlaveDatastore(-1).
		QueryRow(selectAccountByPublicIDQuery, id).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrAccountNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	account := &entity.Organization{}
	err = json.Unmarshal([]byte(JSONData), account)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountRead.UpdateInternalMessage(err.Error())}
	}
	account.ID = ID

	return account, nil
}

// NewOrganization returns a new account handler with PostgreSQL as storage driver
func NewOrganization(pgsql postgres.Client) core.Organization {
	return &organization{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
