package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v02/core"
	"github.com/tapglue/multiverse/v02/entity"
	"github.com/tapglue/multiverse/v02/errmsg"
	storageHelper "github.com/tapglue/multiverse/v02/storage/helper"
	"github.com/tapglue/multiverse/v02/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	accountUser struct {
		pg     postgres.Client
		mainPg *sqlx.DB
		a      core.Account
	}
)

const (
	createAccountUserQuery           = `INSERT INTO tg.account_users(account_id, json_data) VALUES($1, $2) RETURNING id`
	selectAccountUserByIDQuery       = `SELECT json_data FROM tg.account_users WHERE id = $1 AND account_id = $2 LIMIT 1`
	updateAccountUserByIDQuery       = `UPDATE tg.account_users SET json_data = $1 WHERE id = $2 AND account_id = $3`
	deleteAccountUserByIDQuery       = `UPDATE tg.account_users SET enabled = 0 WHERE id = $1`
	listAccountUsersByAccountIDQuery = `SELECT id, json_data FROM tg.account_users WHERE account_id = $1 ORDER BY json_data->>'created_at' DESC`
	selectAccountUserByEmailQuery    = `SELECT id, json_data FROM tg.account_users WHERE json_data @> json_build_object('email', $1::text)::jsonb LIMIT 1`
	selectAccountUserByUsernameQuery = `SELECT id, json_data FROM tg.account_users WHERE json_data @> json_build_object('user_name', $1::text)::jsonb LIMIT 1`
	createAccountUserSessionQuery    = `INSERT INTO tg.account_user_sessions(account_id, account_user_id, session_id) VALUES($1, $2, $3)`
	selectAccountUserSessionQuery    = `SELECT session_id FROM tg.account_user_sessions WHERE account_id = $1 AND account_user_id = $2 LIMIT 1`
	selectAccountUserBySessionQuery  = `SELECT account_id, account_user_id FROM tg.account_user_sessions WHERE session_id = $1 LIMIT 1`
	selectAccountUserByPublicIDQuery = `SELECT id, json_data FROM tg.account_users WHERE account_id = $1 AND json_data @> json_build_object('id', $2::text)::jsonb LIMIT 1`
	updateAccountUserSessionQuery    = `UPDATE tg.account_user_sessions SET session_id = $1 WHERE account_id = $2 AND account_user_id = $3 AND session_id = $4`
	destroyAccountUserSessionQuery   = `DELETE FROM tg.account_user_sessions WHERE account_id = $1 AND account_user_id = $2 AND session_id = $3`
	destroyAccountUserSessionsQuery  = `DELETE FROM tg.account_user_sessions WHERE account_id = $1 AND account_user_id = $2`
)

func (au *accountUser) Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, []errors.Error) {
	accountUser.PublicID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	accountUser.Password = storageHelper.EncryptPassword(accountUser.Password)
	accountUser.Enabled = true
	timeNow := time.Now()
	accountUser.CreatedAt, accountUser.UpdatedAt, accountUser.LastLogin = &timeNow, &timeNow, &timeNow

	accountUserJSON, err := json.Marshal(accountUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserCreation.UpdateInternalMessage(err.Error())}
	}

	var accountUserID int64
	err = au.mainPg.
		QueryRow(createAccountUserQuery, accountUser.AccountID, string(accountUserJSON)).
		Scan(&accountUserID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}
	return au.Read(accountUser.AccountID, accountUserID)
}

func (au *accountUser) Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er []errors.Error) {
	var JSONData string
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByIDQuery, accountUserID, accountID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrAccountUserNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser = &entity.AccountUser{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = accountUserID
	accountUser.AccountID = accountID

	return accountUser, nil
}

func (au *accountUser) Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, []errors.Error) {
	if updatedAccountUser.Password == "" {
		updatedAccountUser.Password = existingAccountUser.Password
	} else if updatedAccountUser.Password != existingAccountUser.Password {
		updatedAccountUser.Password = storageHelper.EncryptPassword(updatedAccountUser.Password)
	}
	timeNow := time.Now()
	updatedAccountUser.UpdatedAt = &timeNow

	accountUserJSON, err := json.Marshal(updatedAccountUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = au.mainPg.
		Exec(updateAccountUserByIDQuery, string(accountUserJSON), existingAccountUser.ID, existingAccountUser.AccountID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return au.Read(existingAccountUser.AccountID, existingAccountUser.ID)
}

func (au *accountUser) Delete(accountUser *entity.AccountUser) []errors.Error {
	user, err := au.Read(accountUser.AccountID, accountUser.ID)
	if err != nil {
		return err
	}
	user.Enabled = false

	_, err = au.Update(*user, *user, false)
	if err != nil {
		return err
	}

	_, er := au.mainPg.Exec(destroyAccountUserSessionsQuery, user.AccountID, user.ID)
	if er != nil {
		return []errors.Error{errmsg.ErrInternalAccountUserSessionDelete.UpdateInternalMessage(er.Error())}
	}

	return nil
}

func (au *accountUser) List(accountID int64) (accountUsers []*entity.AccountUser, er []errors.Error) {
	accountUsers = []*entity.AccountUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(listAccountUsersByAccountIDQuery, accountID)
	if err != nil {
		return accountUsers, []errors.Error{errmsg.ErrInternalAccountUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ID       int64
			JSONData string
		)
		err := rows.Scan(&ID, &JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalAccountUserList.UpdateInternalMessage(err.Error())}
		}
		accountUser := &entity.AccountUser{}
		err = json.Unmarshal([]byte(JSONData), accountUser)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalAccountUserList.UpdateInternalMessage(err.Error())}
		}
		accountUser.ID = ID

		accountUsers = append(accountUsers, accountUser)
	}

	return accountUsers, nil
}

func (au *accountUser) CreateSession(user *entity.AccountUser) (string, []errors.Error) {
	sessionToken := storageHelper.GenerateAccountSessionID(user)
	_, err := au.mainPg.Exec(createAccountUserSessionQuery, user.AccountID, user.ID, sessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalAccountUserSessionCreation.UpdateInternalMessage(err.Error())}
	}

	return sessionToken, nil
}

func (au *accountUser) RefreshSession(sessionToken string, user *entity.AccountUser) (string, []errors.Error) {
	updatedSessionToken := storageHelper.GenerateAccountSessionID(user)
	_, err := au.mainPg.Exec(updateAccountUserSessionQuery, sessionToken, user.AccountID, user.ID, updatedSessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalAccountUserSessionUpdate.UpdateInternalMessage(err.Error())}
	}

	return updatedSessionToken, nil
}

func (au *accountUser) DestroySession(sessionToken string, user *entity.AccountUser) []errors.Error {
	_, err := au.mainPg.Exec(destroyAccountUserSessionQuery, user.AccountID, user.ID, sessionToken)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalAccountUserSessionDelete.UpdateInternalMessage(err.Error())}
	}

	return nil
}

func (au *accountUser) GetSession(user *entity.AccountUser) (string, []errors.Error) {
	rows, err := au.pg.SlaveDatastore(-1).
		Query(selectAccountUserSessionQuery, user.AccountID, user.ID)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalAccountUserSessionRead.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	sessions := []string{}
	for rows.Next() {
		session := ""
		if err := rows.Scan(&session); err != nil {
			return "", []errors.Error{errmsg.ErrInternalAccountUserSessionRead.UpdateInternalMessage(err.Error())}
		}
		if session == user.SessionToken {
			return session, nil
		}
		sessions = append(sessions, session)
	}

	// TODO think about what we have to do when we have multiple sessions and we are asking for one...
	return sessions[rand.Intn(len(sessions))], nil
}

func (au *accountUser) FindByEmail(email string) (*entity.Account, *entity.AccountUser, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByEmailQuery, email).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = ID

	account, er := au.a.ReadByPublicID(accountUser.PublicAccountID)
	if er != nil {
		return nil, nil, er
	}
	accountUser.AccountID = account.ID

	return account, accountUser, nil
}

func (au *accountUser) ExistsByEmail(email string) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByEmailQuery, email).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (au *accountUser) FindByUsername(username string) (*entity.Account, *entity.AccountUser, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByUsernameQuery, username).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = ID

	account, er := au.a.ReadByPublicID(accountUser.PublicAccountID)
	if er != nil {
		return nil, nil, er
	}
	accountUser.AccountID = account.ID

	return account, accountUser, nil
}

func (au *accountUser) ExistsByUsername(username string) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByUsernameQuery, username).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (au *accountUser) ExistsByID(accountID, accountUserID int64) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByIDQuery, accountUserID, accountID).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (au *accountUser) FindBySession(sessionKey string) (*entity.AccountUser, []errors.Error) {
	var accountID, accountUserID int64

	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserBySessionQuery, sessionKey).
		Scan(&accountID, &accountUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrAccountUserNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser, er := au.Read(accountID, accountUserID)
	if er != nil {
		return nil, er
	}
	if accountUser == nil || accountUser.Enabled == false {
		return nil, []errors.Error{errmsg.ErrAccountUserNotFound}
	}

	return accountUser, nil
}

func (au *accountUser) FindByPublicID(accountID int64, publicID string) (*entity.AccountUser, []errors.Error) {
	var (
		accountUserID int64
		JSONData      string
	)

	err := au.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByPublicIDQuery, accountID, publicID).
		Scan(&accountUserID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrAccountUserNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser := &entity.AccountUser{}
	if err := json.Unmarshal([]byte(JSONData), accountUser); err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = accountUserID
	accountUser.AccountID = accountID

	return accountUser, nil
}

// NewAccountUser returns a new account user handler with PostgreSQL as storage driver
func NewAccountUser(pgsql postgres.Client) core.AccountUser {
	return &accountUser{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		a:      NewAccount(pgsql),
	}
}

// NewAccountUserWithAccount returns a new account user handler with PostgreSQL as storage driver
func NewAccountUserWithAccount(pgsql postgres.Client, account core.Account) core.AccountUser {
	return &accountUser{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		a:      account,
	}
}
