package redis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	application struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	data, er := json.Marshal(application)
	if er != nil {
		return nil, errors.NewInternalError("error while creating the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-create-%d", application.AccountID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationCreate, partitionKey, data)

	return nil, err
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, errors.Error) {
	data, er := json.Marshal(updatedApplication)
	if er != nil {
		return nil, errors.NewInternalError("error while updating the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-update-%d-%d", updatedApplication.AccountID, updatedApplication.ID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationUpdate, partitionKey, data)

	return nil, err
}

func (app *application) Delete(application *entity.Application) errors.Error {
	data, er := json.Marshal(application)
	if er != nil {
		return errors.NewInternalError("error while deleting the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-delete-%d-%d", application.AccountID, application.ID)
	_, err := app.storage.PackAndPutRecord(kinesis.StreamApplicationDelete, partitionKey, data)

	return err
}

func (app *application) List(accountID int64) ([]*entity.Application, errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Exists(accountID, applicationID int64) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) FindByKey(applicationKey string) (*entity.Application, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (app *application) FindByPublicID(publicID string) (*entity.Application, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewApplication creates a new Application
func NewApplication(storageClient kinesis.Client) core.Application {
	return &application{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
