package redis

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"

	"github.com/garyburd/redigo/redis"
)

type application struct {
	storage *redis.Pool
}

const (
	redisAppKey = "applications:bykey:"
)

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	conn := app.storage.Get()
	defer conn.Close()

	a := struct {
		entity.OrgAppIDs
		entity.Application
	}{}
	a.Application = *application
	a.OrgAppIDs.OrgID = application.OrgID
	a.AppID = application.ID

	ap, err := json.Marshal(a)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}

	if _, err := conn.Do("SET", redisAppKey+application.AuthToken, ap, "NX"); err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return app.Read(application.OrgID, application.ID)
}

func (app *application) Read(orgID, applicationID int64) (*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	conn := app.storage.Get()
	defer conn.Close()

	if updatedApplication.AuthToken == "" {
		updatedApplication.AuthToken = existingApplication.AuthToken
	}

	a := struct {
		entity.OrgAppIDs
		entity.Application
	}{}
	a.Application = updatedApplication
	a.OrgAppIDs.OrgID = updatedApplication.OrgID
	a.AppID = updatedApplication.ID

	ap, err := json.Marshal(a)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}

	if _, err := conn.Do("SET", redisAppKey+updatedApplication.AuthToken, ap); err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationCreation.UpdateInternalMessage(err.Error())}
	}

	if !retrieve {
		return nil, nil
	}

	return app.Read(updatedApplication.OrgID, updatedApplication.ID)
}

func (app *application) Delete(application *entity.Application) []errors.Error {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) List(orgID int64) ([]*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Exists(orgID, applicationID int64) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) FindByKey(applicationKey string) (*entity.Application, []errors.Error) {
	conn := app.storage.Get()
	defer conn.Close()

	application, err := conn.Do("GET", redisAppKey+applicationKey)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}

	if application == nil {
		return nil, []errors.Error{errmsg.ErrApplicationNotFound}
	}

	ap := &struct {
		entity.OrgAppIDs
		entity.Application
	}{}
	if err := json.Unmarshal(application.([]byte), ap); err != nil {
		return nil, []errors.Error{errmsg.ErrInternalApplicationRead.UpdateInternalMessage(err.Error())}
	}

	ap.Application.ID = ap.OrgAppIDs.AppID
	ap.Application.OrgID = ap.OrgAppIDs.OrgID

	return &ap.Application, nil
}

func (app *application) FindByPublicID(publicID string) (*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewApplication returns a new application handler with Redis as storage driver
func NewApplication(driver *redis.Pool) core.Application {
	return &application{
		storage: driver,
	}
}
