package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type application struct {
	storage kinesis.Client
	ksis    *ksis.Kinesis
}

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	data, er := json.Marshal(application)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the application (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-create-%d", application.OrgID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationCreate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return application, nil
	}

	return nil, nil
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, []errors.Error) {
	data, er := json.Marshal(updatedApplication)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the application (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-update-%d-%d", updatedApplication.OrgID, updatedApplication.ID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationUpdate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return &updatedApplication, nil
	}

	return nil, nil
}

func (app *application) Delete(application *entity.Application) []errors.Error {
	data, er := json.Marshal(application)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the application (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-delete-%d-%d", application.OrgID, application.ID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationDelete, partitionKey, data)
	if err != nil {
		return []errors.Error{err}
	}

	return nil
}

func (app *application) List(accountID int64) ([]*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (app *application) Exists(accountID, applicationID int64) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (app *application) FindByKey(applicationKey string) (*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (app *application) FindByPublicID(publicID string) (*entity.Application, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

// NewApplication creates a new Application
func NewApplication(storageClient kinesis.Client) core.Application {
	return &application{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
