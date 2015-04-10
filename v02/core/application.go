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
		Create(application *entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError)
		Read(accountID, applicationID int64) (*entity.Application, tgerrors.TGError)
		Update(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, tgerrors.TGError)
		Delete(accountID, applicationID int64) tgerrors.TGError
		List(accountID int64) ([]*entity.Application, tgerrors.TGError)
	}
)
