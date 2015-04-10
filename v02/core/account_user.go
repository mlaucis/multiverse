/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// AccountUser interface
	AccountUser interface {
		Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError)
		Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er tgerrors.TGError)
		Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError)
		Delete(accountID, userID int64) tgerrors.TGError
		List(accountID int64) (accountUsers []*entity.AccountUser, er tgerrors.TGError)
		CreateSession(user *entity.AccountUser) (string, tgerrors.TGError)
		RefreshSession(sessionToken string, user *entity.AccountUser) (string, tgerrors.TGError)
		DestroySession(sessionToken string, user *entity.AccountUser) tgerrors.TGError
		GetSession(user *entity.AccountUser) (string, tgerrors.TGError)
		FindByEmail(email string) (*entity.Account, *entity.AccountUser, tgerrors.TGError)
		ExistsByEmail(email string) (bool, tgerrors.TGError)
		FindByUsername(username string) (*entity.Account, *entity.AccountUser, tgerrors.TGError)
		ExistsByUsername(username string) (bool, tgerrors.TGError)
	}
)
