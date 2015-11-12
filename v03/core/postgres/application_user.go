package postgres

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	storageHelper "github.com/tapglue/multiverse/v03/storage/helper"
	"github.com/tapglue/multiverse/v03/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type applicationUser struct {
	pg     postgres.Client
	mainPg *sqlx.DB
	conn   core.Connection
}

const (
	createApplicationUserQuery = `INSERT INTO app_%d_%d.users(json_data) VALUES($1)`

	selectApplicationUserByIDQuery = `SELECT json_data, last_read
	FROM app_%d_%d.users
	WHERE (json_data->>'id')::BIGINT = $1::BIGINT
  		AND (json_data->>'enabled')::BOOL = TRUE
  		AND (json_data->>'deleted')::BOOL = FALSE
	LIMIT 1`

	selectApplicationUserExistsByIDQuery = `SELECT true
	FROM app_%d_%d.users
	WHERE (json_data->>'id')::BIGINT = $1::BIGINT
  		AND (json_data->>'enabled')::BOOL = TRUE
  		AND (json_data->>'deleted')::BOOL = FALSE
	LIMIT 1`

	updateApplicationUserByIDQuery = `UPDATE app_%d_%d.users
	SET json_data = $1
	WHERE (json_data->>'id')::BIGINT = $2::BIGINT`

	listApplicationUsersByApplicationIDQuery = `SELECT json_data
	FROM app_%d_%d.users
	WHERE (json_data->>'enabled')::BOOL = TRUE
  		AND (json_data->>'deleted')::BOOL = FALSE
  		LIMIT 1`

	listApplicationUsersByUserIDsQuery = `SELECT json_data
	FROM app_%d_%d.users
	WHERE (json_data->>'id')::BIGINT = ANY(%s)
		AND (json_data->>'enabled')::BOOL = true
		AND (json_data->>'deleted')::BOOL = false`

	selectApplicationUserByEmailQuery = `SELECT json_data, last_read
	FROM app_%d_%d.users
	WHERE (json_data->>'email')::TEXT = $1::text
		AND (json_data->>'enabled')::BOOL = true
		AND (json_data->>'deleted')::BOOL = false
	LIMIT 1`

	selectApplicationUserByUsernameQuery = `SELECT json_data, last_read
	FROM app_%d_%d.users
	WHERE (json_data->>'user_name')::TEXT = $1::TEXT
		AND (json_data->>'enabled')::BOOL = true
		AND (json_data->>'deleted')::BOOL = false
	LIMIT 1`

	createApplicationUserSessionQuery = `INSERT INTO app_%d_%d.sessions(user_id, session_id) VALUES($1, $2)`

	selectApplicationUserSessionQuery = `SELECT session_id FROM app_%d_%d.sessions WHERE user_id = $1 AND enabled = TRUE LIMIT 1`

	selectApplicationUserBySessionQuery = `SELECT user_id FROM app_%d_%d.sessions WHERE session_id = $1 AND enabled = TRUE LIMIT 1`

	updateApplicationUserSessionQuery = `UPDATE app_%d_%d.sessions SET session_id = $1 WHERE user_id = $2 AND session_id = $3`

	destroyApplicationUserSessionQuery = `UPDATE app_%d_%d.sessions SET enabled = FALSE WHERE user_id = $1 AND session_id = $2`

	destroyAllApplicationUserSessionQuery = `UPDATE app_%d_%d.sessions SET enabled = FALSE WHERE user_id = $1`

	searchApplicationUsersQuery = `SELECT json_data
	FROM app_%d_%d.users
	WHERE ((json_data->>'user_name' ILIKE $1)
		OR (json_data->>'email' ILIKE $1)
		OR (json_data->>'first_name' ILIKE $1)
		OR (json_data->>'last_name' ILIKE $1))
		AND (json_data->>'id')::BIGINT != $2
		AND (json_data->>'enabled')::BOOL = true
		AND (json_data->>'deleted')::BOOL = false
	LIMIT 50`

	selectApplicationUserCountsQuery = `SELECT
  (SELECT count(*) FROM app_%d_%d.connections
    WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'enabled')::BOOL = true AND json_data->>'type' = '` + entity.ConnectionTypeFriend + `') AS "friends",
  (SELECT count(*) FROM app_%d_%d.connections
    WHERE (json_data->>'user_to_id')::BIGINT = $1::BIGINT AND (json_data->>'enabled')::BOOL = true AND json_data->>'type' = '` + entity.ConnectionTypeFollow + `') AS "follower",
  (SELECT count(*) FROM app_%d_%d.connections
    WHERE (json_data->>'user_from_id')::BIGINT = $1::BIGINT AND (json_data->>'enabled')::BOOL = true AND json_data->>'type' = '` + entity.ConnectionTypeFollow + `') AS "followed"`
)

func (au *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error {
	if user.ID == 0 {
		return []errors.Error{errmsg.ErrInternalApplicationUserIDMissing.SetCurrentLocation()}
	}
	connectionType := user.SocialConnectionType
	user.SocialConnectionType = ""
	user.Enabled = true
	user.Deleted = entity.PFalse

	var err error
	user.Password, err = storageHelper.StrongEncryptPassword(user.Password)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	timeNow := time.Now()
	user.CreatedAt, user.UpdatedAt = &timeNow, &timeNow

	applicationUserJSON, err := json.Marshal(user)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = au.mainPg.
		Exec(appSchema(createApplicationUserQuery, accountID, applicationID), string(applicationUserJSON))
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	for platform := range user.SocialIDs {
		_, err := au.conn.SocialConnect(accountID, applicationID, user, platform, user.SocialConnectionsIDs[platform], connectionType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (au *applicationUser) Read(accountID, applicationID int64, userID uint64, withStatistics bool) (*entity.ApplicationUser, []errors.Error) {
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	appUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), appUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	appUser.LastRead = &lastRead
	if withStatistics {
		au.FriendStatistics(accountID, applicationID, appUser)
	}

	return appUser, nil
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
		return users, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(er.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
	}
	timeNow := time.Now()
	updatedUser.UpdatedAt = &timeNow
	if updatedUser.Deleted == nil {
		updatedUser.Deleted = entity.PFalse
	}

	userJSON, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = au.mainPg.
		Exec(appSchema(updateApplicationUserByIDQuery, accountID, applicationID), string(userJSON), existingUser.ID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	if !retrieve {
		return nil, nil
	}

	return &updatedUser, nil
}

func (au *applicationUser) Delete(accountID, applicationID int64, userID uint64) []errors.Error {
	user, err := au.Read(accountID, applicationID, userID, false)
	if err != nil {
		return err
	}

	user.Enabled = false
	*user.Deleted = true
	_, err = au.Update(accountID, applicationID, *user, *user, false)

	if err != nil {
		return err
	}

	return au.destroyAllUserSession(accountID, applicationID, user)
}

func (au *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, er []errors.Error) {
	users = []*entity.ApplicationUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(listApplicationUsersByApplicationIDQuery, accountID, applicationID))
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	sessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(createApplicationUserSessionQuery, accountID, applicationID), user.ID, sessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return sessionToken, nil
}

func (au *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, []errors.Error) {
	updatedSessionToken := storageHelper.GenerateApplicationSessionID(user)
	_, err := au.mainPg.Exec(appSchema(updateApplicationUserSessionQuery, accountID, applicationID), sessionToken, user.ID, updatedSessionToken)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return updatedSessionToken, nil
}

func (au *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(selectApplicationUserSessionQuery, accountID, applicationID), user.ID)
	if err != nil {
		return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	sessions := []string{}
	for rows.Next() {
		session := ""
		if err := rows.Scan(&session); err != nil {
			return "", []errors.Error{errmsg.ErrInternalApplicationUserSessionRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return []errors.Error{errmsg.ErrInternalApplicationUserSessionDelete.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	user := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(JSONData), user)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return true, nil
}

func (au *applicationUser) ExistsByID(accountID, applicationID int64, userID uint64) (bool, []errors.Error) {
	var exists bool

	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchema(selectApplicationUserExistsByIDQuery, accountID, applicationID), userID).
		Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
		return nil, []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return au.Read(accountID, applicationID, userID, false)
}

func (au *applicationUser) Search(accountID, applicationID int64, userID uint64, searchTerm string) ([]*entity.ApplicationUser, []errors.Error) {
	users := []*entity.ApplicationUser{}

	rows, err := au.pg.SlaveDatastore(-1).
		Query(appSchema(searchApplicationUsersQuery, accountID, applicationID), "%"+searchTerm+"%", userID)
	if err != nil {
		return users, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	defer rows.Close()
	for rows.Next() {
		var (
			JSONData string
		)
		err := rows.Scan(&JSONData)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		user := &entity.ApplicationUser{}
		err = json.Unmarshal([]byte(JSONData), user)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationUserList.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}

		users = append(users, user)
	}

	return users, nil
}

func (au *applicationUser) destroyAllUserSession(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error {
	_, err := au.mainPg.Exec(appSchema(destroyAllApplicationUserSessionQuery, accountID, applicationID), user.ID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserSessionsDelete.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	return nil
}

func (au *applicationUser) FriendStatistics(accountID, applicationID int64, appUser *entity.ApplicationUser) []errors.Error {
	err := au.pg.SlaveDatastore(-1).
		QueryRow(appSchemaWithParams(selectApplicationUserCountsQuery, accountID, applicationID, accountID, applicationID, accountID, applicationID), appUser.ID).
		Scan(&appUser.FriendCount, &appUser.FollowerCount, &appUser.FollowedCount)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationUserRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
