/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/core"
)

// ValidateAccountRequestToken checks if the provided request token is valid or not for the specified account
func ValidateAccountRequestToken(accountID int64, requestToken string) bool {
	if len(requestToken) < 7 {
		return false
	}

	if requestToken[:7] != "Bearer " {
		return false
	}

	account, err := core.ReadAccount(accountID)
	if err != nil {
		return false
	}

	token, err := storageClient.GenerateAccountToken(account)
	if err != nil {
		return false
	}

	// TODO account token is actually the secret that we never want to transfer in plain text (or base64)
	// so this needs to actually generate a JWT signature with the account token as the secret

	requestToken = requestToken[7:]

	return token == requestToken
}
