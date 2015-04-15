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
	panic("not implemented yet")
	return nil, nil
}

func (app *application) Read(accountID, applicationID int64) (*entity.Application, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (app *application) Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (app *application) Delete(accountID, applicationID int64) tgerrors.TGError {
	panic("not implemented yet")
	return nil
}

func (app *application) List(accountID int64) ([]*entity.Application, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
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
