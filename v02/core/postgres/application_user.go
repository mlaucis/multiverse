/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	applicationUser struct {
		pg     postgres.Client
		mainPg *sql.DB
		conn   core.Connection
	}
)

const (
	createApplicationUserQuery               = `INSERT INTO app_%d_%d.users(json_data) VALUES($1)`
	selectApplicationUserByIDQuery           = `SELECT json_data FROM app_%d_%d.users WHERE json_data->>'id' = $1`
	updateApplicationUserByIDQuery           = `UPDATE app_%d_%d.users SET json_data = $1 WHERE json_data->>'id' = $2`
	listApplicationUsersByApplicationIDQuery = `SELECT json_data FROM app_%d_%d.users`
	selectApplicationUserByEmailQuery        = `SELECT json_data FROM app_%d_%d.users WHERE json_data->>'email' = $1`
	selectApplicationUserByUsernameQuery     = `SELECT json_data FROM app_%d_%d.users WHERE json_data->>'user_name' = $1`
	createApplicationUserSessionQuery        = `INSERT INTO app_%d_%d.sessions(user_id, session_id) VALUES($1, $2)`
	selectApplicationUserSessionQuery        = `SELECT session_id FROM app_%d_%d.sessions WHERE user_id = $1`
	selectApplicationUserBySessionQuery      = `SELECT user_id FROM app_%d_%d.sessions WHERE session_id = $1`
	updateApplicationUserSessionQuery        = `UPDATE app_%d_%d.sessions SET session_id = $1 WHERE user_id = $2 AND session_id = $3`
	destroyApplicationUserSessionQuery       = `DELETE FROM app_%d_%d.sessions WHERE user_id = $1 AND session_id = $2`
)

func (au *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, errors.Error) {
	connectionType := user.SocialConnectionType

	user.ID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	user.SocialConnectionType = ""
	user.Enabled = true
	user.Password = storageHelper.EncryptPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	applicationUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the application user", err.Error())
	}

	_, err = au.mainPg.
		Exec(appSchema(createApplicationUserQuery, accountID, applicationID), string(applicationUserJSON))
	if err != nil {
		return nil, errors.NewInternalError("error while creating the application user", err.Error())
	}

	for platform := range user.SocialIDs {
		_, err := au.conn.SocialConnect(accountID, applicationID, user, platform, user.SocialConnectionsIDs[platform], connectionType)
		if err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return nil, nil
	}
	return au.Read(accountID, applicationID, user.ID)
}

func (au *applicationUser) Read(accountID, applicationID int64, userID string) (*entity.ApplicationUser, errors.Error) {
	var JSONData string
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}

	applicationUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), applicationUser)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return applicationUser, nil
}

func (au *applicationUser) Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, errors.Error) {
	updatedUser.SocialConnectionType = ""
	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		updatedUser.Password = storageHelper.EncryptPassword(updatedUser.Password)
	}
	updatedUser.UpdatedAt = time.Now()

	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the application user", err.Error())
	}

	_, err = au.pg.SlaveDatastore(-1).
		Exec(appSchema(updateApplicationUserByIDQuery, accountID, applicationID), string(userJSON), existingUser.ID)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the application user", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return au.Read(accountID, applicationID, existingUser.ID)
}

func (au *applicationUser) Delete(accountID, applicationID int64, applicationUser *entity.ApplicationUser) errors.Error {
	applicationUser.Enabled = false
	_, err := au.Update(accountID, applicationID, *applicationUser, *applicationUser, false)

	return err
}

func (au *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, er errors.Error) {
	users = []*entity.ApplicationUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(listApplicationUsersByApplicationIDQuery, accountID, applicationID))
	if err != nil {
		return users, errors.NewInternalError("error while retrieving list of application users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of application users", err.Error())
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of application users", err.Error())
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, errors.Error) {
	sessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(createApplicationUserSessionQuery, accountID, applicationID), user.ID, sessionToken)
	if err != nil {
		return "", errors.NewInternalError("error while creating application user session", err.Error())
	}

	return sessionToken, nil
}

func (au *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, errors.Error) {
	updatedSessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(updateApplicationUserSessionQuery, accountID, applicationID), sessionToken, user.ID, updatedSessionToken)
	if err != nil {
		return "", errors.NewInternalError("error while updating application user session", err.Error())
	}

	return updatedSessionToken, nil
}

func (au *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, errors.Error) {
	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(selectApplicationUserSessionQuery, accountID, applicationID), user.ID)
	if err != nil {
		return "", errors.NewInternalError("error while reading session from the database", err.Error())
	}
	sessions := []string{}
	for rows.Next() {
		session := ""
		if err := rows.Scan(&session); err != nil {
			return "", errors.NewInternalError("error while reading session from the database", err.Error())
		}
		if session == user.SessionToken {
			return session, nil
		}
		sessions = append(sessions, session)
	}

	// TODO think about what we have to do when we have multiple sessions and we are asking for one...
	return sessions[rand.Intn(len(sessions))], nil
}

func (au *applicationUser) DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) errors.Error {
	_, err := au.mainPg.Exec(appSchema(destroyApplicationUserSessionQuery, accountID, applicationID), user.ID, sessionToken)
	if err != nil {
		return errors.NewInternalError("error while deleting session", err.Error())
	}

	return nil
}

func (au *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, errors.Error) {
	var JSONData string
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&JSONData)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return user, nil
}

func (au *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return true, nil
}

func (au *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, errors.Error) {
	var JSONData string
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&JSONData)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return user, nil
}

func (au *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, errors.Error) {
	var (
		ID       int64
		JSONData string
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return true, nil
}

func (au *applicationUser) ExistsByID(accountID, applicationID int64, userID string) (bool, errors.Error) {
	var JSONData string
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}
	return true, nil
}

func (au *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, errors.Error) {
	var userID string

	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserBySessionQuery, accountID, applicationID), sessionKey).
		Scan(&userID)
	if err != nil {
		return nil, errors.NewInternalError("error while loading the application user", err.Error())
	}

	return au.Read(accountID, applicationID, userID)
}

// NewApplicationUser returns a new application user handler with PostgreSQL as storage driver
func NewApplicationUser(pgsql postgres.Client) core.ApplicationUser {
	return &applicationUser{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		conn: &connection{
			pg:     pgsql,
			mainPg: pgsql.MainDatastore(),
		},
	}
}
