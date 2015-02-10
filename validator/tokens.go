/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"encoding/base64"
	"strconv"

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

	requestToken = requestToken[7:]

	if requestToken == "" {
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

	return token == requestToken
}

// ValidateApplicationRequestToken checks if the provided request token is valid or not for the specified app
func ValidateApplicationRequestToken(accountID, applicationID int64, requestToken string) bool {
	if len(requestToken) < 7 {
		return false
	}

	if requestToken[:7] != "Bearer " {
		return false
	}

	requestToken = requestToken[7:]

	if requestToken == "" {
		return false
	}

	// Store the token details in redis
	storedToken, err := storageEngine.HMGet(
		"tokens:"+base64.URLEncoding.EncodeToString([]byte(requestToken)),
		"acc",
		"app",
	).Result()
	if err != nil {
		return false
	}

	if storedToken == nil {
		return false
	}

	var acc, app int64

	switch storedToken[0].(type) {
	case nil:
		return false
	case string:
		acc, err = strconv.ParseInt(storedToken[0].(string), 10, 64)
		if err != nil {
			return false
		}
	}

	switch storedToken[1].(type) {
	case nil:
		return false
	case string:
		app, err = strconv.ParseInt(storedToken[1].(string), 10, 64)
		if err != nil {
			return false
		}
	}

	if acc != accountID || app != applicationID {
		return false
	}

	return true
}
