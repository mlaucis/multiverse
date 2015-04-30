/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"

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
	createApplicationUserQuery               = `INSERT INTO app_%d_%d.users(json_data) VALUES($1) RETURNING id`
	selectApplicationUserByIDQuery           = `SELECT json_data, enabled FROM app_%d_%d.users WHERE id = $1`
	updateApplicationUserByIDQuery           = `UPDATE app_%d_%d.users SET json_data = $1 WHERE id = $2`
	deleteApplicationUserByIDQuery           = `UPDATE app_%d_%d.users SET enabled = FALSE WHERE id = $1`
	listApplicationUsersByApplicationIDQuery = `SELECT id, json_data FROM app_%d_%d.users`
	selectApplicationUserByEmailQuery        = `SELECT id, json_data, enabled FROM app_%d_%d.users WHERE json_data->>'email' = $1`
	selectApplicationUserByUsernameQuery     = `SELECT id, json_data, enabled FROM app_%d_%d.users WHERE json_data->>'user_name' = $1`
	createApplicationUserSessionQuery        = `INSERT INTO app_%d_%d.sessions(user_id, session_id) VALUES($1, $2)`
	selectApplicationUserSessionQuery        = `SELECT session_id FROM app_%d_%d.sessions WHERE user_id = $1`
	selectApplicationUserBySessionQuery      = `SELECT user_id FROM app_%d_%d.sessions WHERE session_id = $1`
	updateApplicationUserSessionQuery        = `UPDATE app_%d_%d.sessions SET session_id = $1 WHERE user_id = $2 AND session_id = $3`
	destroyApplicationUserSessionQuery       = `DELETE FROM app_%d_%d.sessions WHERE user_id = $1 AND session_id = $2`
)

func (au *applicationUser) Create(user *entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, errors.Error) {
	connectionType := user.SocialConnectionType
	user.SocialConnectionType = ""
	user.Activated = true
	user.Password = storageHelper.EncryptPassword(user.Password)
	applicationUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the application user", err.Error())
	}

	var applicationUserID int64
	err = au.mainPg.
		QueryRow(appSchema(createApplicationUserQuery, user.AccountID, user.ApplicationID), string(applicationUserJSON)).
		Scan(&applicationUserID)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the application user", err.Error())
	}
	user.ID = applicationUserID

	for platform := range user.SocialIDs {
		_, err := au.conn.SocialConnect(user, platform, user.SocialConnectionsIDs[platform], connectionType)
		if err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return nil, nil
	}
	return au.Read(user.AccountID, user.ApplicationID, applicationUserID)
}

func (au *applicationUser) Read(accountID, applicationID, userID int64) (*entity.ApplicationUser, errors.Error) {
	var (
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData, &Enabled)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}

	applicationUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), applicationUser)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	applicationUser.ID = userID
	applicationUser.Enabled = Enabled

	return applicationUser, nil
}

func (au *applicationUser) Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, errors.Error) {
	updatedUser.SocialConnectionType = ""
	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		updatedUser.Password = storageHelper.EncryptPassword(updatedUser.Password)
	}

	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the application user", err.Error())
	}

	_, err = au.pg.SlaveDatastore(-1).
		Exec(appSchema(updateApplicationUserByIDQuery, existingUser.AccountID, existingUser.ApplicationID), string(userJSON), existingUser.ID)
	if err != nil {
		return nil, errors.NewInternalError("error while updating the application user", err.Error())
	}

	if !retrieve {
		return nil, nil
	}

	return au.Read(existingUser.AccountID, existingUser.ApplicationID, existingUser.ID)
}

func (au *applicationUser) Delete(applicationUser *entity.ApplicationUser) errors.Error {
	_, err := au.mainPg.Exec(appSchema(deleteApplicationUserByIDQuery, applicationUser.AccountID, applicationUser.ApplicationID), applicationUser.ID)
	if err != nil {
		return errors.NewInternalError("error while deleting the application user", err.Error())
	}
	return nil
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
			ID       int64
			JSONData string
			Enabled  bool
		)
		err := rows.Scan(&ID, &JSONData, &Enabled)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of application users", err.Error())
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return []*entity.ApplicationUser{}, errors.NewInternalError("error while retrieving list of application users", err.Error())
		}
		user.ID = ID
		user.Enabled = Enabled

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) CreateSession(user *entity.ApplicationUser) (string, errors.Error) {
	sessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(createApplicationUserSessionQuery, user.AccountID, user.ApplicationID), user.ID, sessionToken)
	if err != nil {
		return "", errors.NewInternalError("error while creating application user session", err.Error())
	}

	return sessionToken, nil
}

func (au *applicationUser) RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, errors.Error) {
	updatedSessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(updateApplicationUserSessionQuery, user.AccountID, user.ApplicationID), sessionToken, user.ID, updatedSessionToken)
	if err != nil {
		return "", errors.NewInternalError("error while updating application user session", err.Error())
	}

	return updatedSessionToken, nil
}

func (au *applicationUser) GetSession(user *entity.ApplicationUser) (string, errors.Error) {
	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(selectApplicationUserSessionQuery, user.AccountID, user.ApplicationID), user.ID)
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

func (au *applicationUser) DestroySession(sessionToken string, user *entity.ApplicationUser) errors.Error {
	_, err := au.mainPg.Exec(appSchema(destroyApplicationUserSessionQuery, user.AccountID, user.ApplicationID), user.ID, sessionToken)
	if err != nil {
		return errors.NewInternalError("error while deleting session", err.Error())
	}

	return nil
}

func (au *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, errors.Error) {
	var (
		ID       int64
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&ID, &JSONData, &Enabled)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user.ID = ID
	user.Enabled = Enabled

	return user, nil
}

func (au *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, errors.Error) {
	var (
		ID       int64
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&ID, &JSONData, &Enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return true, nil
}

func (au *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, errors.Error) {
	var (
		ID       int64
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&ID, &JSONData, &Enabled)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, errors.NewInternalError("error while reading the application user", err.Error())
	}
	user.ID = ID
	user.Enabled = Enabled

	return user, nil
}

func (au *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, errors.Error) {
	var (
		ID       int64
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&ID, &JSONData, &Enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}

	return true, nil
}

func (au *applicationUser) ExistsByID(accountID, applicationID, userID int64) (bool, errors.Error) {
	var (
		JSONData string
		Enabled  bool
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData, &Enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.NewInternalError("error while reading the application user", err.Error())
	}
	return true, nil
}

func (au *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, errors.Error) {
	var userID int64

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
		conn:   NewConnection(pgsql),
	}
}
