package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	storageHelper "github.com/tapglue/backend/v03/storage/helper"
	"github.com/tapglue/backend/v03/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type member struct {
	pg     postgres.Client
	mainPg *sqlx.DB
	org    core.Organization
}

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

func (mem *member) Create(accountUser *entity.Member, retrieve bool) (*entity.Member, []errors.Error) {
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
	err = mem.mainPg.
		QueryRow(createAccountUserQuery, accountUser.OrgID, string(accountUserJSON)).
		Scan(&accountUserID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}
	return mem.Read(accountUser.OrgID, accountUserID)
}

func (mem *member) Read(accountID, accountUserID int64) (accountUser *entity.Member, er []errors.Error) {
	var JSONData string
	err := mem.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByIDQuery, accountUserID, accountID).
		Scan(&JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrMemberNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser = &entity.Member{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = accountUserID
	accountUser.OrgID = accountID

	return accountUser, nil
}

func (mem *member) Update(existingAccountUser, updatedAccountUser entity.Member, retrieve bool) (*entity.Member, []errors.Error) {
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

	_, err = mem.mainPg.
		Exec(updateAccountUserByIDQuery, string(accountUserJSON), existingAccountUser.ID, existingAccountUser.OrgID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return mem.Read(existingAccountUser.OrgID, existingAccountUser.ID)
}

func (mem *member) Delete(accountUser *entity.Member) []errors.Error {
	user, err := mem.Read(accountUser.OrgID, accountUser.ID)
	if err != nil {
		return err
	}
	user.Enabled = false

	_, err = mem.Update(*user, *user, false)
	if err != nil {
		return err
	}

	_, er := mem.mainPg.Exec(destroyAccountUserSessionsQuery, user.OrgID, user.ID)
	if er != nil {
		return []errors.Error{errmsg.ErrInternalAccountUserSessionDelete.UpdateInternalMessage(er.Error())}
	}

	return nil
}

func (mem *member) List(accountID int64) (accountUsers []*entity.Member, er []errors.Error) {
	accountUsers = []*entity.Member{}

	rows, err := mem.pg.SlaveDatastore(-1).
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
		accountUser := &entity.Member{}
		err = json.Unmarshal([]byte(JSONData), accountUser)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalAccountUserList.UpdateInternalMessage(err.Error())}
		}
		accountUser.ID = ID

		accountUsers = append(accountUsers, accountUser)
	}

	return accountUsers, nil
}

func (mem *member) CreateSession(user *entity.Member) (string, []errors.Error) {
	sessionToken := storageHelper.GenerateAccountSessionID(user)
	_, err := mem.mainPg.Exec(createAccountUserSessionQuery, user.OrgID, user.ID, sessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalAccountUserSessionCreation.UpdateInternalMessage(err.Error())}
	}

	return sessionToken, nil
}

func (mem *member) RefreshSession(sessionToken string, user *entity.Member) (string, []errors.Error) {
	updatedSessionToken := storageHelper.GenerateAccountSessionID(user)
	_, err := mem.mainPg.Exec(updateAccountUserSessionQuery, sessionToken, user.OrgID, user.ID, updatedSessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalAccountUserSessionUpdate.UpdateInternalMessage(err.Error())}
	}

	return updatedSessionToken, nil
}

func (mem *member) DestroySession(sessionToken string, user *entity.Member) []errors.Error {
	_, err := mem.mainPg.Exec(destroyAccountUserSessionQuery, user.OrgID, user.ID, sessionToken)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalAccountUserSessionDelete.UpdateInternalMessage(err.Error())}
	}

	return nil
}

func (mem *member) GetSession(user *entity.Member) (string, []errors.Error) {
	rows, err := mem.pg.SlaveDatastore(-1).
		Query(selectAccountUserSessionQuery, user.OrgID, user.ID)
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

func (mem *member) FindByEmail(email string) (*entity.Organization, *entity.Member, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := mem.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByEmailQuery, email).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser := &entity.Member{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = ID

	account, er := mem.org.ReadByPublicID(accountUser.PublicAccountID)
	if er != nil {
		return nil, nil, er
	}
	accountUser.OrgID = account.ID

	return account, accountUser, nil
}

func (mem *member) ExistsByEmail(email string) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := mem.pg.SlaveDatastore(-1).
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

func (mem *member) FindByUsername(username string) (*entity.Organization, *entity.Member, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := mem.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByUsernameQuery, username).
		Scan(&ID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser := &entity.Member{}
	err = json.Unmarshal([]byte(JSONData), accountUser)
	if err != nil {
		return nil, nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = ID

	account, er := mem.org.ReadByPublicID(accountUser.PublicAccountID)
	if er != nil {
		return nil, nil, er
	}
	accountUser.OrgID = account.ID

	return account, accountUser, nil
}

func (mem *member) ExistsByUsername(username string) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := mem.pg.SlaveDatastore(-1).
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

func (mem *member) ExistsByID(accountID, accountUserID int64) (bool, []errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := mem.pg.SlaveDatastore(-1).
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

func (mem *member) FindBySession(sessionKey string) (*entity.Member, []errors.Error) {
	var accountID, accountUserID int64

	err := mem.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserBySessionQuery, sessionKey).
		Scan(&accountID, &accountUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrMemberNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser, er := mem.Read(accountID, accountUserID)
	if er != nil {
		return nil, er
	}
	if accountUser == nil || accountUser.Enabled == false {
		return nil, []errors.Error{errmsg.ErrMemberNotFound}
	}

	return accountUser, nil
}

func (mem *member) FindByPublicID(accountID int64, publicID string) (*entity.Member, []errors.Error) {
	var (
		accountUserID int64
		JSONData      string
	)

	err := mem.pg.SlaveDatastore(-1).
		QueryRow(selectAccountUserByPublicIDQuery, accountID, publicID).
		Scan(&accountUserID, &JSONData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrMemberNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}

	accountUser := &entity.Member{}
	if err := json.Unmarshal([]byte(JSONData), accountUser); err != nil {
		return nil, []errors.Error{errmsg.ErrInternalAccountUserRead.UpdateInternalMessage(err.Error())}
	}
	accountUser.ID = accountUserID
	accountUser.OrgID = accountID

	return accountUser, nil
}

// NewMember returns a new account user handler with PostgreSQL as storage driver
func NewMember(pgsql postgres.Client) core.Member {
	return &member{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		org:    NewOrganization(pgsql),
	}
}

// NewMemeberWithOrganization returns a new account user handler with PostgreSQL as storage driver
func NewMemeberWithOrganization(pgsql postgres.Client, account core.Organization) core.Member {
	return &member{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		org:    account,
	}
}
