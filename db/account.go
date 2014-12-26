/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"github.com/gluee/backend/entity"
)

/*
 * GetAccountById returns the account matching the ID or an error
 */
func GetAccountById(accountID uint64) (account *entity.Account, err error) {
	account = &entity.Account{}

	err = GetSlave().
		QueryRowx("SELECT * FROM accounts WHERE id=?", accountID).
		StructScan(account)

	return
}
