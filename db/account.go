/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/gluee/backend/entity"
)

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID uint64) (account *entity.Account, err error) {
	account = &entity.Account{}

	err = GetSlave().
		QueryRowx("SELECT * FROM `gluee`.`accounts` WHERE `id`=?", accountID).
		StructScan(account)

	return
}

// AddAccount adds a new account to the database and returns the created account or an error
func AddAccount(account *entity.Account) (*entity.Account, error) {
	query := "INSERT INTO `gluee`.`accounts` (`name`) VALUES (?)"
	result, err := GetMaster().Exec(query, account.Name)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	var createdAccountID int64
	createdAccountID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	return GetAccountByID(uint64(createdAccountID))
}