/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package utils holds supportive functions for tests etc.
package utils

import (
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

// Create a correct account
func AddCorrectAccount() (acc *entity.Account, err error) {
	savedAccount, err := core.WriteAccount(correctAccount, true)
	if err != nil {
		return nil, err
	}

	return savedAccount, nil
}

// EmptyAccount returns an empty account
func EmptyAccount() *entity.Account {
	return emtpyAccount
}

// CorrectAccount returns a correct account
func CorrectAccount() *entity.Account {
	return correctAccount
}
