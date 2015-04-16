package redis

import (
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
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (app *application) Delete(accountID, applicationID int64) tgerrors.TGError {
	return tgerrors.NewNotFoundError("not found", "invalid handler specified")
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
