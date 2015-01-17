/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package mysql

import (
	"fmt"

	"github.com/tapglue/backend/entity"
)

// GetAccountUserByID returns the user matching the account or an error
func GetAccountUserByID(accountID int64, userID uint64) (accountUser *entity.AccountUser, err error) {
	accountUser = &entity.AccountUser{}

	// Execute query to get account user
	err = GetSlave().
		QueryRowx("SELECT * FROM `account_users` WHERE `id`=? AND `account_id`=?", userID, accountID).
		StructScan(accountUser)

	return
}

// GetAccountAllUsers returns all the users from a certain account
func GetAccountAllUsers(accountID int64) (accountUsers []*entity.AccountUser, err error) {
	accountUsers = []*entity.AccountUser{}

	// Execute query to get account users
	err = GetSlave().
		Select(&accountUsers, "SELECT * FROM `account_users` WHERE `account_id`=?", accountID)

	return
}

// AddAccountUser creates a user for an account and returns the created entry or an error
func AddAccountUser(accountID int64, accountUser *entity.AccountUser) (*entity.AccountUser, error) {
	// Check if name empty
	if accountUser.Username == "" {
		return nil, fmt.Errorf("empty account user username is not allowed")
	}
	// Check if password empty
	if accountUser.Password == "" {
		return nil, fmt.Errorf("empty account user password is not allowed")
	}
	// Check if email empty
	if accountUser.Email == "" {
		return nil, fmt.Errorf("empty account user email is not allowed")
	}

	// Write to db
	query := "INSERT INTO `account_users` (`account_id`, `name`, `password`, `email`) VALUES (?, ?, ?, ?)"
	result, err := GetMaster().Exec(query, accountID, accountUser.Username, accountUser.Password, accountUser.Email)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	// Retrieve account user
	var createdAccountUserID int64
	createdAccountUserID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	// Return account user
	return GetAccountUserByID(accountID, uint64(createdAccountUserID))
}
