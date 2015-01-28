/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

func ValidateAccount(account *entity.Account) error {
	// Check if name empty
	if account.Name == "" {
		return fmt.Errorf("account name should not be empty")
	}

	return nil
}
