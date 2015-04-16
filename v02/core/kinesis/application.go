package redis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/tgerrors"
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

func (app *application) Create(application *entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError) {
	data, er := json.Marshal(application)
	if er != nil {
		return tgerrors.NewInternalError("error while creating the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-create-%d", application.AccountID)
	_, err := app.storage.PutRecord("application_create", partitionKey, data)

	return nil, err
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError) {
	data, er := json.Marshal(updatedApplication)
	if er != nil {
		return tgerrors.NewInternalError("error while updating the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-update-%d-%d", updatedApplication.AccountID, updatedApplication.ID)
	_, err := app.storage.PutRecord("application_update", partitionKey, data)

	return nil, err
}

func (app *application) Delete(accountUser *entity.AccountUser) tgerrors.TGError {
	data, er := json.Marshal(accountUser)
	if er != nil {
		return tgerrors.NewInternalError("error while deleting the application (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-delete-%d-%d", accountUser.AccountID, accountUser.ID)
	_, err := app.storage.PutRecord("application_delete", partitionKey, data)

	return err
}

func (app *application) List(accountID int64) ([]*entity.Application, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Exists(accountID, applicationID int64) bool {
	panic("not implemented yet")
	return false
}

// NewApplication creates a new Application
func NewApplication(storageClient kinesis.Client) core.Application {
	return &application{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
