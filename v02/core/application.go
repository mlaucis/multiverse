/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Application interface
	Application interface {
		// Create adds an application to the database and returns the created applicaton user or an error
		Create(application *entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError)

		// Read returns the application matching the ID or an error
		Read(accountID, applicationID int64) (*entity.Application, tgerrors.TGError)

		// Update updates an application in the database and returns the created applicaton user or an error
		Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError)

		// Delete deletes the application matching the IDs or an error
		Delete(*entity.Application) tgerrors.TGError

		// List returns all applications from a certain account
		List(accountID int64) ([]*entity.Application, tgerrors.TGError)

		// Exists validates if an application exists and returns the application or an error
		Exists(accountID, applicationID int64) bool
	}
)
