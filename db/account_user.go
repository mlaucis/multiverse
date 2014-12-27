/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/gluee/backend/entity"
)

// GetAccountUser returns the user matching the account or an error
func GetAccountUserByID(accountID, userID uint64) (accountUser *entity.AccountUser, err error) {
	accountUser = &entity.AccountUser{}

	err = GetSlave().
		QueryRowx("SELECT * FROM `gluee`.`account_users` WHERE `id`=? AND `account_id`=?", userID, accountID).
		StructScan(accountUser)

	return
}

// GetAccountAllUsers returns all the users from a certain account
func GetAccountAllUsers(accountID uint64) (accountUsers []*entity.AccountUser, err error) {
	accountUsers = []*entity.AccountUser{}

	err = GetSlave().
		Select(&accountUsers, "SELECT * FROM `gluee`.`account_users` WHERE `account_id`=?", accountID)

	return
}

// AddAccountUser creates a user for an account and returns the created entry or an error
func AddAccountUser(accountID uint64, accountUser *entity.AccountUser) (*entity.AccountUser, error) {
	query := "INSERT INTO `gluee`.`account_users` (`account_id`, `name`, `password`, `email`) VALUES (?, ?, ?, ?);"
	result, err := GetMaster().Exec(query, accountID, accountUser.Name, accountUser.Password, accountUser.Email)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	var createdAccountUserID int64
	createdAccountUserID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	return GetAccountUserByID(accountID, uint64(createdAccountUserID))
}
