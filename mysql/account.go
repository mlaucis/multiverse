/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package mysql

import (
	"fmt"

	"github.com/tapglue/backend/entity"
)

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID int64) (account *entity.Account, err error) {
	account = &entity.Account{}

	// Execute query to get account
	err = GetSlave().
		QueryRowx("SELECT * FROM `accounts` WHERE `id`=?", accountID).
		StructScan(account)

	return
}

// AddAccount adds a new account to the database and returns the created account or an error
func AddAccount(account *entity.Account) (*entity.Account, error) {
	// Check if name empty
	if account.Name == "" {
		return nil, fmt.Errorf("account name should not be empty")
	}

	// Write to db
	query := "INSERT INTO `accounts` (`name`) VALUES (?)"
	result, err := GetMaster().Exec(query, account.Name)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	// Retrieve account
	var createdAccountID int64
	createdAccountID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	// Return account
	return GetAccountByID(createdAccountID)
}
