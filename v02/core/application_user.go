/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// ApplicationUser interface
	ApplicationUser interface {
		Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError)
		Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err tgerrors.TGError)
		Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError)
		Delete(accountID, applicationID, userID int64) (err tgerrors.TGError)
		List(accountID, applicationID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)
		CreateSession(user *entity.ApplicationUser) (string, tgerrors.TGError)
		RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, tgerrors.TGError)
		GetSession(user *entity.ApplicationUser) (string, tgerrors.TGError)
		DestroySession(sessionToken string, user *entity.ApplicationUser) tgerrors.TGError
		FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, tgerrors.TGError)
		ExistsByEmail(accountID, applicationID int64, email string) (bool, tgerrors.TGError)
		FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, tgerrors.TGError)
		ExistsByUsername(accountID, applicationID int64, email string) (bool, tgerrors.TGError)
	}
)
