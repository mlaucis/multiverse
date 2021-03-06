package core

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/entity"
)

// Application interface
type Application interface {
	// Create adds an application to the database and returns the created applicaton user or an error
	Create(application *entity.Application, retrieve bool) (*entity.Application, []errors.Error)

	// Read returns the application matching the ID or an error
	Read(accountID, applicationID int64) (*entity.Application, []errors.Error)

	// Update updates an application in the database and returns the created applicaton user or an error
	Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, []errors.Error)

	// Delete deletes the application matching the IDs or an error
	Delete(*entity.Application) []errors.Error

	// List returns all applications from a certain account
	List(accountID int64) ([]*entity.Application, []errors.Error)

	// Exists validates if an application exists and returns the application or an error
	Exists(accountID, applicationID int64) (bool, []errors.Error)

	// FindByToken will load an application by the application token, if it exists
	FindByApplicationToken(applicationToken string) (*entity.Application, []errors.Error)

	// FindByBackendToken will load an application by the backend token, if it exists
	FindByBackendToken(applicationToken string) (*entity.Application, []errors.Error)

	// FindByPublicID finds an application by it's public ID
	FindByPublicID(publicID string) (*entity.Application, []errors.Error)
}
