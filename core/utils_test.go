/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/core/entity"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount(fetchAccount bool) (acc *entity.Account, err error) {
	account, err := WriteAccount(correctAccount, fetchAccount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser(accountID int64, fetchUser bool) (acc *entity.AccountUser, err error) {
	accountUserWithAccountID := correctAccountUser
	accountUserWithAccountID.AccountID = accountID
	accountUser, err := WriteAccountUser(accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// CorrectAccount returns a correct account
func CorrectAccount() *entity.Account {
	return correctAccount
}

// CorrectAccountUser returns a correct account user
func CorrectAccountUser() *entity.AccountUser {
	return correctAccountUser
}
