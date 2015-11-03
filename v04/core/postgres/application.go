package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	storageHelper "github.com/tapglue/multiverse/v04/storage/helper"
	"github.com/tapglue/multiverse/v04/storage/postgres"

	"github.com/jmoiron/sqlx"
)

type application struct {
	pg     postgres.Client
	mainPg *sqlx.DB
	redis  core.Application
}

const (
	createApplicationEntryQuery                   = `INSERT INTO tg.applications (account_id, json_data) VALUES($1, $2) RETURNING id`
	selectApplicationEntryByIDQuery               = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE id = $1 AND account_id = $2 and enabled = 1`
	checkApplicationExistsByIDQuery               = `SELECT id FROM tg.applications WHERE id = $1 AND account_id = $2 and enabled = 1`
	selectApplicationEntryByPublicIDsQuery        = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE json_data->>'id' = $1::text LIMIT 1`
	selectApplicationEntryByApplicationTokenQuery = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE json_data->>'token' = $1::text LIMIT 1`
	selectApplicationEntryByBackendTokenQuery     = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE json_data->>'backend_token' = $1::text LIMIT 1`
	updateApplicationEntryByIDQuery               = `UPDATE tg.applications SET json_data = $1 WHERE id = $2 AND account_id = $3`
	deleteApplicationEntryByIDQuery               = `UPDATE tg.applications SET enabled = 0 WHERE id = $1 AND account_id = $2`
	listApplicationsEntryByAccountIDQuery         = `SELECT id, json_data, enabled FROM tg.applications where account_id = $1 and enabled = 1`
)

var createApplicationNamespaceQuery = []string{
	`CREATE SCHEMA app_%d_%d`,
	`CREATE TABLE app_%d_%d.users
	(
		json_data JSONB NOT NULL,
		last_read TIMESTAMP DEFAULT '2015-05-01 01:23:45' NOT NULL
	)`,
	`CREATE TABLE app_%d_%d.events
	(
		json_data JSONB NOT NULL,
		geo GEOMETRY(POINT, 4326)
	)`,
	`CREATE TABLE app_%d_%d.connections
	(
		json_data JSONB NOT NULL
	)`,
	`CREATE TABLE app_%d_%d.sessions
	(
		user_id BIGINT NOT NULL,
		session_id VARCHAR(40) NOT NULL,
		created_at TIMESTAMP DEFAULT now() NOT NULL,
		enabled BOOL DEFAULT TRUE NOT NULL
	)`,

	`CREATE INDEX user_id ON app_%d_%d.users USING btree (((json_data ->> 'id') :: BIGINT))`,
	`CREATE INDEX user_email ON app_%d_%d.users USING btree (((json_data ->> 'email') :: TEXT))`,
	`CREATE INDEX user_username ON app_%d_%d.users USING btree (((json_data ->> 'user_name') :: TEXT))`,

	`CREATE INDEX ON app_%d_%d.sessions (session_id, user_id)`,

	`CREATE INDEX conection_from_id ON app_%d_%d.connections USING btree ((((json_data ->> 'user_from_id'::text))::bigint))`,
	`CREATE INDEX conection_to_id ON app_%d_%d.connections USING btree ((((json_data ->> 'user_to_id'::text))::bigint))`,

	`CREATE INDEX event_id ON app_%d_%d.events USING btree (((json_data ->> 'id') :: BIGINT))`,
	`CREATE INDEX event_user_id ON app_%d_%d.events USING btree (((json_data ->> 'user_id') :: BIGINT))`,
	`CREATE INDEX event_visibility ON app_%d_%d.events USING btree (((json_data ->> 'visibility') :: INT))`,
	`CREATE INDEX event_target_id ON app_%d_%d.events USING BTREE (((json_data -> 'target'->>'id') :: TEXT))`,
	`CREATE INDEX event_created_at ON app_%d_%d.events USING btree ((json_data ->> 'created_at'))`,
	`CREATE INDEX ON app_%d_%d.events USING GIST (geo)`,
}

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	application.PublicID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	application.Enabled = true
	timeNow := time.Now()
	application.CreatedAt, application.UpdatedAt = &timeNow, &timeNow
	application.AuthToken = storageHelper.GenerateApplicationSecretKey(application)
	application.BackendToken = storageHelper.GenerateBackendApplicationSecretKey(application)

	applicationJSON, err := json.Marshal(application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	var applicationID int64
	err = app.mainPg.
		QueryRow(createApplicationEntryQuery, application.OrgID, applicationJSON).
		Scan(&applicationID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	application.ID = applicationID

	for idx := range createApplicationNamespaceQuery {
		_, err = app.mainPg.Exec(fmt.Sprintf(createApplicationNamespaceQuery[idx], application.OrgID, application.ID))
		if err != nil {
			// TODO rollback the creation from the field if we fail to create all the stuff here
			// TODO learn transactions :)
			return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
	}

	if !retrieve {
		return nil, nil
	}

	return application, nil
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, []errors.Error) {
	return app.findByQuery(selectApplicationEntryByIDQuery, applicationID, accountID)
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	if updatedApplication.AuthToken == "" {
		updatedApplication.AuthToken = existingApplication.AuthToken
	}
	if updatedApplication.BackendToken == "" {
		updatedApplication.BackendToken = existingApplication.BackendToken
	}
	timeNow := time.Now()
	updatedApplication.UpdatedAt = &timeNow

	applicationJSON, err := json.Marshal(updatedApplication)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, err = app.mainPg.Exec(updateApplicationEntryByIDQuery, applicationJSON, existingApplication.ID, existingApplication.OrgID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	_, er := app.redis.Update(existingApplication, updatedApplication, false)
	if er != nil {
		return nil, er
	}

	if !retrieve {
		return nil, nil
	}
	return &updatedApplication, nil
}

func (app *application) Delete(application *entity.Application) []errors.Error {
	_, err := app.mainPg.Exec(deleteApplicationEntryByIDQuery, application.ID, application.OrgID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationDelete.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	return nil
}

func (app *application) List(accountID int64) ([]*entity.Application, []errors.Error) {
	applications := []*entity.Application{}

	rows, err := app.pg.SlaveDatastore(-1).
		Query(listApplicationsEntryByAccountIDQuery, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return applications, nil
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
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
			return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		application := &entity.Application{}
		err = json.Unmarshal([]byte(JSONData), application)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
		}
		application.ID = ID
		application.Enabled = Enabled

		applications = append(applications, application)
	}

	return applications, nil
}

func (app *application) Exists(accountID, applicationID int64) (bool, []errors.Error) {
	var ID int
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByIDQuery, applicationID, accountID).
		Scan(&ID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	return true, nil
}

func (app *application) FindByApplicationToken(applicationToken string) (*entity.Application, []errors.Error) {
	return app.findByQuery(selectApplicationEntryByApplicationTokenQuery, applicationToken)
}

func (app *application) FindByBackendToken(backendToken string) (*entity.Application, []errors.Error) {
	return app.findByQuery(selectApplicationEntryByBackendTokenQuery, backendToken)
}

func (app *application) FindByPublicID(publicID string) (*entity.Application, []errors.Error) {
	return app.findByQuery(selectApplicationEntryByPublicIDsQuery, publicID)
}

func (app *application) findByQuery(query string, params ...interface{}) (*entity.Application, []errors.Error) {
	var (
		ID, accountID int64
		jsonData      string
		enabled       bool
	)
	err := app.pg.SlaveDatastore(-1).
		QueryRow(query, params...).
		Scan(&ID, &accountID, &jsonData, &enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationNotFound.SetCurrentLocation()}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	application := &entity.Application{}
	err = json.Unmarshal([]byte(jsonData), application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	application.ID = ID
	application.OrgID = accountID
	application.Enabled = enabled

	return application, nil
}

// NewApplication returns a new application handler with PostgreSQL as storage driver
func NewApplication(pgsql postgres.Client, redis core.Application) core.Application {
	return &application{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
		redis:  redis,
	}
}
