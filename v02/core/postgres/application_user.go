package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	applicationUser struct {
		pg     postgres.Client
		mainPg *sqlx.DB
		conn   core.Connection
	}
)

const (
	createApplicationUserQuery               = `INSERT INTO app_%d_%d.users(json_data) VALUES($1)`
	selectApplicationUserByIDQuery           = `SELECT json_data, last_read FROM app_%d_%d.users WHERE json_data @> json_build_object('id', $1::bigint, 'enabled', true, 'deleted', false)::jsonb LIMIT 1`
	updateApplicationUserByIDQuery           = `UPDATE app_%d_%d.users SET json_data = $1 WHERE json_data @> json_build_object('id', $2::bigint)::jsonb`
	listApplicationUsersByApplicationIDQuery = `SELECT json_data FROM app_%d_%d.users WHERE json_data @> '{"enabled": true, "deleted": false}' LIMIT 1`
	listApplicationUsersByUserIDsQuery       = `SELECT json_data FROM app_%d_%d.users WHERE (json_data->>'id')::BIGINT = ANY(%s) AND json_data @> '{"enabled": true, "deleted": false}'`
	selectApplicationUserByEmailQuery        = `SELECT json_data, last_read FROM app_%d_%d.users WHERE json_data @> json_build_object('email', $1::text, 'enabled', true, 'deleted', false)::jsonb LIMIT 1`
	selectApplicationUserByUsernameQuery     = `SELECT json_data, last_read FROM app_%d_%d.users WHERE json_data @> json_build_object('user_name', $1::text, 'enabled', true, 'deleted', false)::jsonb LIMIT 1`
	selectApplicationUserByCustomIDQuery     = `SELECT json_data, last_read FROM app_%d_%d.users WHERE json_data @> json_build_object('custom_id', $1::text, 'enabled', true, 'deleted', false)::jsonb LIMIT 1`
	createApplicationUserSessionQuery        = `INSERT INTO app_%d_%d.sessions(user_id, session_id) VALUES($1, $2)`
	selectApplicationUserSessionQuery        = `SELECT session_id FROM app_%d_%d.sessions WHERE user_id = $1 AND enabled = TRUE LIMIT 1`
	selectApplicationUserBySessionQuery      = `SELECT user_id FROM app_%d_%d.sessions WHERE session_id = $1 AND enabled = TRUE LIMIT 1`
	updateApplicationUserSessionQuery        = `UPDATE app_%d_%d.sessions SET session_id = $1 WHERE user_id = $2 AND session_id = $3`
	destroyApplicationUserSessionQuery       = `UPDATE app_%d_%d.sessions SET enabled = FALSE WHERE user_id = $1 AND session_id = $2`
	destroyAllApplicationUserSessionQuery    = `UPDATE app_%d_%d.sessions SET enabled = FALSE WHERE user_id = $1`
	searchApplicationUsersQuery              = `SELECT json_data FROM app_%d_%d.users WHERE ((json_data->>'user_name' ILIKE $1) OR (json_data->>'email' ILIKE $1) OR (json_data->>'first_name' ILIKE $1) OR (json_data->>'last_name' ILIKE $1)) AND json_data @> '{"enabled": true, "deleted": false}' LIMIT 50`
)

func (au *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, []errors.Error) {
	if user.ID == 0 {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserIDMissing}
	}
	connectionType := user.SocialConnectionType
	user.SocialConnectionType = ""
	user.Activated = true
	user.Enabled = true
	user.Deleted = entity.PFalse

	var err error
	user.Password, err = storageHelper.StrongEncryptPassword(user.Password)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error())}
	}

	timeNow := time.Now()
	user.CreatedAt, user.UpdatedAt = &timeNow, &timeNow

	applicationUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error())}
	}

	_, err = au.mainPg.
		Exec(appSchema(createApplicationUserQuery, accountID, applicationID), string(applicationUserJSON))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error())}
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

func (au *applicationUser) Read(accountID, applicationID int64, userID uint64) (*entity.ApplicationUser, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)

	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData, &lastRead)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}

	applicationUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), applicationUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	applicationUser.LastRead = &lastRead

	return applicationUser, nil
}

func (au *applicationUser) ReadMultiple(accountID, applicationID int64, userIDs []uint64) (users []*entity.ApplicationUser, err []errors.Error) {
	users = []*entity.ApplicationUser{}
	if len(userIDs) == 0 {
		return
	}

	ids := ""
	for idx := 0; idx < len(userIDs)-1; idx++ {
		ids += strconv.FormatUint(userIDs[idx], 10) + ", "
	}
	ids += strconv.FormatUint(userIDs[len(userIDs)-1], 10)
	condition := `ARRAY[` + ids + `]`

	rows, er := au.pg.SlaveDatastore(-1).
		Query(appSchemaWithParams(listApplicationUsersByUserIDsQuery, accountID, applicationID, condition))
	if er != nil {
		return users, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(er.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, []errors.Error) {
	updatedUser.SocialConnectionType = ""
	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		var err error
		updatedUser.Password, err = storageHelper.StrongEncryptPassword(updatedUser.Password)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error())}
		}
	}
	timeNow := time.Now()
	updatedUser.UpdatedAt = &timeNow

	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = au.mainPg.
		Exec(appSchema(updateApplicationUserByIDQuery, accountID, applicationID), string(userJSON), existingUser.ID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return au.Read(accountID, applicationID, existingUser.ID)
}

func (au *applicationUser) Delete(accountID, applicationID int64, applicationUser *entity.ApplicationUser) []errors.Error {
	applicationUser.Enabled = false
	*applicationUser.Deleted = true
	_, err := au.Update(accountID, applicationID, *applicationUser, *applicationUser, false)

	go au.destroyAllUserSession(accountID, applicationID, applicationUser)

	return err
}

func (au *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, er []errors.Error) {
	users = []*entity.ApplicationUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(listApplicationUsersByApplicationIDQuery, accountID, applicationID))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	sessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(createApplicationUserSessionQuery, accountID, applicationID), user.ID, sessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionCreation.UpdateInternalMessage(err.Error())}
	}

	return sessionToken, nil
}

func (au *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, []errors.Error) {
	updatedSessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(updateApplicationUserSessionQuery, accountID, applicationID), sessionToken, user.ID, updatedSessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionUpdate.UpdateInternalMessage(err.Error())}
	}

	return updatedSessionToken, nil
}

func (au *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(selectApplicationUserSessionQuery, accountID, applicationID), user.ID)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionRead.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	sessions := []string{}
	for rows.Next() {
		session := ""
		if err := rows.Scan(&session); err != nil {
			return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionRead.UpdateInternalMessage(err.Error())}
		}
		if session == user.SessionToken {
			return session, nil
		}
		sessions = append(sessions, session)
	}

	// TODO think about what we have to do when we have multiple sessions and we are asking for one...
	return sessions[rand.Intn(len(sessions))], nil
}

func (au *applicationUser) DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) []errors.Error {
	_, err := au.mainPg.Exec(appSchema(destroyApplicationUserSessionQuery, accountID, applicationID), user.ID, sessionToken)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserSessionDelete.UpdateInternalMessage(err.Error())}
	}

	return nil
}

func (au *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&JSONData, &lastRead)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user.LastRead = &lastRead

	return user, nil
}

func (au *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByEmailQuery, accountID, applicationID), email).
		Scan(&JSONData, &lastRead)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}

	return true, nil
}

func (au *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&JSONData, &lastRead)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user.LastRead = &lastRead

	return user, nil
}

func (au *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByUsernameQuery, accountID, applicationID), username).
		Scan(&JSONData, &lastRead)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}

	return true, nil
}

func (au *applicationUser) ExistsByID(accountID, applicationID int64, userID uint64) (bool, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)

	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByIDQuery, accountID, applicationID), userID).
		Scan(&JSONData, &lastRead)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (au *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, []errors.Error) {
	var userID uint64

	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserBySessionQuery, accountID, applicationID), sessionKey).
		Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}

	return au.Read(accountID, applicationID, userID)
}

func (au *applicationUser) Search(accountID, applicationID int64, searchTerm string) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(searchApplicationUsersQuery, accountID, applicationID), "%"+searchTerm+"%")
	if err != nil {
		return users, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error())}
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) FindByCustomID(accountID, applicationID int64, customID string) (*entity.ApplicationUser, []errors.Error) {
	var (
		JSONData string
		lastRead time.Time
	)
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserByCustomIDQuery, accountID, applicationID), customID).
		Scan(&JSONData, &lastRead)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationUserNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error())}
	}
	user.LastRead = &lastRead

	return user, nil
}

func (au *applicationUser) destroyAllUserSession(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error {
	_, err := au.mainPg.Exec(appSchema(destroyAllApplicationUserSessionQuery, accountID, applicationID), user.ID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserSessionsDelete.UpdateInternalMessage(err.Error())}
	}

	return nil
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
