/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"database/sql"

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

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) Delete(*entity.Application) errors.Error {
	return errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) List(accountID int64) ([]*entity.Application, errors.Error) {
	return []*entity.Application{}, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) Exists(accountID, applicationID int64) (bool, errors.Error) {
	panic("not implemented yet")
}

// NewApplication returns a new application handler with PostgreSQL as storage driver
func NewApplication(pgsql postgres.Client) core.Application {
	return &application{
		pg:     pgsql,
		mainPg: pgsql.MainDatastore(),
	}
}
