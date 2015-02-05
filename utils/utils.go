/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package utils holds supportive functions for tests etc.
package utils

import (
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount() (acc *entity.Account, err error) {
	account, err := core.WriteAccount(correctAccount, true)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser() (acc *entity.AccountUser, err error) {
	account, err := AddCorrectAccount()
	if err != nil {
		return nil, err
	}

	accountUserWithAccountID := correctAccountUser
	accountUserWithAccountID.AccountID = account.ID
	accountUser, err := core.WriteAccountUser(accountUserWithAccountID, true)
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
