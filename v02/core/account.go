/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Account interface
	Account interface {
		Create(account *entity.Account, retrieve bool) (acc *entity.Account, err tgerrors.TGError)
		Read(accountID int64) (account *entity.Account, err tgerrors.TGError)
		Update(existingAccount, updatedAccount entity.Account, retrieve bool) (acc *entity.Account, err tgerrors.TGError)
		Delete(accountID int64) (err tgerrors.TGError)
	}
)
