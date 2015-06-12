/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	application struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

const (
	createApplicationEntryQuery            = `INSERT INTO tg.applications (account_id, json_data) VALUES($1, $2) RETURNING id`
	selectApplicationEntryByIDQuery        = `SELECT json_data, enabled FROM tg.applications WHERE id = $1 AND account_id = $2 and enabled = 1`
	selectApplicationEntryByPublicIDsQuery = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE json_data @> json_build_object('id', $1::text)::jsonb LIMIT 1`
	selectApplicationEntryByKeyQuery       = `SELECT id, account_id, json_data, enabled FROM tg.applications WHERE json_data @> json_build_object('token', $1::text)::jsonb LIMIT 1`
	updateApplicationEntryByIDQuery        = `UPDATE tg.applications SET json_data = $1 WHERE id = $2 AND account_id = $3`
	deleteApplicationEntryByIDQuery        = `UPDATE tg.applications SET enabled = 0 WHERE id = $1 AND account_id = $2`
	listApplicationsEntryByAccountIDQuery  = `SELECT id, json_data, enabled FROM tg.applications where account_id = $1 and enabled = 1`
)

var (
	createApplicationNamespaceQuery = []string{
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
		user_id CHAR(36) NOT NULL,
		session_id CHAR(40) NOT NULL,
		created_at TIMESTAMP DEFAULT now() NOT NULL,
		enabled BOOL DEFAULT TRUE NOT NULL
	)`,

		`CREATE INDEX ON app_%d_%d.users USING GIN (json_data jsonb_path_ops)`,
		`CREATE INDEX ON app_%d_%d.connections USING GIN (json_data jsonb_path_ops)`,
		`CREATE INDEX ON app_%d_%d.events USING GIN (json_data jsonb_path_ops)`,
		`CREATE INDEX ON app_%d_%d.events USING GIST (geo)`,
	}
)

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	application.PublicID = storageHelper.GenerateUUIDV5(storageHelper.OIDUUIDNamespace, storageHelper.GenerateRandomString(20))
	application.Enabled = true
	timeNow := time.Now()
	application.CreatedAt, application.UpdatedAt = &timeNow, &timeNow
	application.AuthToken = storageHelper.GenerateApplicationSecretKey(application)

	applicationJSON, err := json.Marshal(application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}

	var applicationID int64
	err = app.mainPg.
		QueryRow(createApplicationEntryQuery, application.AccountID, applicationJSON).
		Scan(&applicationID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}
	application.ID = applicationID

	for idx := range createApplicationNamespaceQuery {
		_, err = app.mainPg.Exec(fmt.Sprintf(createApplicationNamespaceQuery[idx], application.AccountID, application.ID))
		if err != nil {
			// TODO rollback the creation from the field if we fail to create all the stuff here
			// TODO learn transactions :)
			return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
		}
	}

	if !retrieve {
		return nil, nil
	}

	return app.Read(application.AccountID, applicationID)
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, []errors.Error) {
	var (
		JSONData string
		enabled  bool
	)
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByIDQuery, applicationID, accountID).
		Scan(&JSONData, &enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}

	application := &entity.Application{}
	err = json.Unmarshal([]byte(JSONData), application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	application.ID = applicationID
	application.Enabled = enabled

	return application, nil
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	if updatedApplication.AuthToken == "" {
		updatedApplication.AuthToken = existingApplication.AuthToken
	}
	timeNow := time.Now()
	updatedApplication.UpdatedAt = &timeNow

	applicationJSON, err := json.Marshal(updatedApplication)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error())}
	}

	_, err = app.mainPg.Exec(updateApplicationEntryByIDQuery, applicationJSON, existingApplication.ID, existingApplication.AccountID)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationUpdate.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}
	return app.Read(existingApplication.AccountID, existingApplication.ID)
}

func (app *application) Delete(application *entity.Application) []errors.Error {
	_, err := app.mainPg.Exec(deleteApplicationEntryByIDQuery, application.ID, application.AccountID)
	if err != nil {
		return []errors.Error{errmsg.ErrInternalApplicationDelete.UpdateInternalMessage(err.Error())}
	}
	return nil
}

func (app *application) List(accountID int64) ([]*entity.Application, []errors.Error) {
	applications := []*entity.Application{}

	rows, err := app.pg.SlaveDatastore(-1).Query(listApplicationsEntryByAccountIDQuery, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return applications, nil
		}
		return applications, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
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
			return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
		}
		application := &entity.Application{}
		err = json.Unmarshal([]byte(JSONData), application)
		if err != nil {
			return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
		}
		application.ID = ID
		application.Enabled = Enabled

		applications = append(applications, application)
	}

	return applications, nil
}

func (app *application) Exists(accountID, applicationID int64) (bool, []errors.Error) {
	var (
		ID       int
		JSONData string
	)
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByIDQuery, applicationID, accountID).
		Scan(&ID, &JSONData)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	return true, nil
}

func (app *application) FindByKey(applicationKey string) (*entity.Application, []errors.Error) {
	var (
		ID, accountID int64
		JSONData      string
		Enabled       bool
	)
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByKeyQuery, applicationKey).
		Scan(&ID, &accountID, &JSONData, &Enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	application := &entity.Application{}
	err = json.Unmarshal([]byte(JSONData), application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	application.ID = ID
	application.AccountID = accountID
	application.Enabled = Enabled

	return application, nil
}

func (app *application) FindByPublicID(publicID string) (*entity.Application, []errors.Error) {
	var (
		ID, accountID int64
		JSONData      string
		Enabled       bool
	)
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByPublicIDsQuery, publicID).
		Scan(&ID, &accountID, &JSONData, &Enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, []errors.Error{errmsg.ErrApplicationNotFound}
		}
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	application := &entity.Application{}
	err = json.Unmarshal([]byte(JSONData), application)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}
	application.ID = ID
	application.AccountID = accountID
	application.Enabled = Enabled

	return application, nil
}

// NewApplication returns a new application handler with PostgreSQL as storage driver
func NewApplication(pgsql postgres.Client) core.Application {
	return &application{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
