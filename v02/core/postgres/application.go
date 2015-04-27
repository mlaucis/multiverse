/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/postgres"
)

type (
	application struct {
		pg     postgres.Client
		mainPg *sql.DB
	}
)

const (
	createApplicationEntryQuery           = `INSERT INTO applications (account_id, json_data) VALUES($1, $2) RETURNING id`
	selectApplicationEntryByIDQuery       = `SELECT json_data FROM applications WHERE id = $1 AND account_id = $2`
	updateApplicationEntryByIDQuery       = `UPDATE applications SET json_data = $1 WHERE id = $2 AND account_id = $3`
	deleteApplicationEntryByIDQuery       = `DELETE FROM applications WHERE id = $1`
	listApplicationsEntryByAccountIDQuery = `SELECT id, json_data FROM applications where account_id = $1`
)

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	applicationJSON, err := json.Marshal(application)
	if err != nil {
		return nil, errors.NewInternalError("error while creating the application", err.Error())
	}

	var applicationID int64
	err = app.mainPg.
		QueryRow(createApplicationEntryQuery, application.AccountID, applicationJSON).
		Scan(&applicationID)

	if !retrieve {
		return nil, nil
	}

	return app.Read(application.AccountID, applicationID)
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, errors.Error) {
	applicationJSON := &struct {
		ID       int64
		JSONData string
	}{}
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByIDQuery, applicationID, accountID).
		Scan(applicationJSON)
	if err != nil {
		return nil, errors.NewInternalError("error while reading application", err.Error())
	}

	application := &entity.Application{}
	err = json.Unmarshal([]byte(applicationJSON.JSONData), application)
	if err != nil {
		return nil, errors.NewInternalError("error while reading application", err.Error())
	}
	application.ID = applicationJSON.ID

	return application, nil
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	if updatedApplication.AuthToken == "" {
		updatedApplication.AuthToken = existingApplication.AuthToken
	}

	applicationJSON, err := json.Marshal(updatedApplication)
	if err != nil {
		return nil, errors.NewInternalError("failed to update application", err.Error())
	}

	_, err = app.mainPg.Exec(updateApplicationEntryByIDQuery, applicationJSON, existingApplication.ID, existingApplication.AccountID)
	if err != nil {
		return nil, errors.NewInternalError("failed to update application", err.Error())
	}

	if !retrieve {
		return nil, nil
	}
	return app.Read(existingApplication.AccountID, existingApplication.ID)
}

func (app *application) Delete(application *entity.Application) errors.Error {
	_, err := app.mainPg.Exec(deleteApplicationEntryByIDQuery, application.ID)
	if err != nil {
		return errors.NewInternalError("error while deleting the account", err.Error())
	}
	return nil
}

func (app *application) List(accountID int64) ([]*entity.Application, errors.Error) {
	applications := []*entity.Application{}

	rows, err := app.pg.SlaveDatastore(-1).Query(listApplicationsEntryByAccountIDQuery, accountID)
	if err != nil {
		return applications, errors.NewInternalError("failed to read the applications", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		applicationJSON := &struct {
			ID       int64
			JSONData string
		}{}
		err := rows.Scan(applicationJSON)
		if err != nil {
			return []*entity.Application{}, errors.NewInternalError("failed to read the applications", err.Error())
		}
		application := &entity.Application{}
		err = json.Unmarshal([]byte(applicationJSON.JSONData), application)
		if err != nil {
			return []*entity.Application{}, errors.NewInternalError("failed to read the applications", err.Error())
		}
		application.ID = applicationJSON.ID

		applications = append(applications, application)
	}

	return applications, nil
}

func (app *application) Exists(accountID, applicationID int64) (bool, errors.Error) {
	applicationJSON := &struct {
		ID       int
		JSONData string
	}{}
	err := app.pg.SlaveDatastore(-1).
		QueryRow(selectApplicationEntryByIDQuery, applicationID, accountID).
		Scan(applicationJSON)
	if err != nil {
		return false, errors.NewInternalError("error while reading application", err.Error())
	}
	return true, nil
}

// NewApplication returns a new application handler with PostgreSQL as storage driver
func NewApplication(pgsql postgres.Client) core.Application {
	return &application{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
