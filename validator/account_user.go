/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

// ValidateAccountUser will validate the passed account user
func ValidateAccountUser(accountUser *entity.AccountUser) error {
	// Check if name empty
	if accountUser.Username == "" {
		return fmt.Errorf("account name should not be empty")
	}

	// Check if account id empty
	if accountUser.AccountID == 0 {
		return fmt.Errorf("account id should not be empty")
	}

	return nil
}
