/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tapglue/backend/core/entity"
	. "github.com/tapglue/backend/utils"
)

const (
	accountUserNameMin = 2
	accountUserNameMax = 40

	accountUserPasswordMin = 4
	accountUserPasswordMax = 60
)

var (
	errorAccountUserFirstNameSize = fmt.Errorf("user first name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserFirstNameType = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorAccountUserLastNameSize = fmt.Errorf("user last name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserLastNameType = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorAccountUserUsernameSize = fmt.Errorf("user username must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserUsernameType = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorAccountUserPasswordSize = fmt.Errorf("user password must be between %d and %d characters", accountUserPasswordMin, accountUserPasswordMax)

	errorAccountIDZero = fmt.Errorf("account id can't be 0")
	errorAccountIDType = fmt.Errorf("account id is not a valid integer")

	errorAccountUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorAccountUserEmailInvalid = fmt.Errorf("user email is not valid")
)

// CreateAccountUser validates an account user on create
func CreateAccountUser(accountUser *entity.AccountUser) error {
	errs := []*error{}

	if !StringLengthBetween(accountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserUsernameSize)
	}

	if !StringLengthBetween(accountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
		errs = append(errs, &errorAccountUserPasswordSize)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.FirstName) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.LastName) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.Username) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if accountUser.Email == "" || !IsValidEmail(accountUser.Email) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !IsValidURL(accountUser.URL, false) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(accountUser.Image) > 0 {
		if !checkImages(accountUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateAccountUserEmail(accountUser.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, &errorUserEmailAlreadyExists)
		} else {
			errs = append(errs, &err)
		}
	}

	if isDuplicate, err := DuplicateAccountUserUsername(accountUser.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, &errorUserEmailAlreadyExists)
		} else {
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(existingAccountUser, updatedAccountUser *entity.AccountUser) error {
	errs := []*error{}

	if !StringLengthBetween(updatedAccountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserUsernameSize)
	}

	if updatedAccountUser.Password != "" {
		if !StringLengthBetween(updatedAccountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
			errs = append(errs, &errorAccountUserPasswordSize)
		}
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.FirstName) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.LastName) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.Username) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars
	if updatedAccountUser.Email == "" || !IsValidEmail(updatedAccountUser.Email) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if updatedAccountUser.URL != "" && !IsValidURL(updatedAccountUser.URL, true) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(updatedAccountUser.Image) > 0 {
		if !checkImages(updatedAccountUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if !AccountExists(updatedAccountUser.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	if existingAccountUser.Email != updatedAccountUser.Email {
		if isDuplicate, err := DuplicateAccountUserEmail(updatedAccountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorEmailAddressInUse)
			} else if err != nil {
				errs = append(errs, &err)
			}
		}
	}

	if existingAccountUser.Username != updatedAccountUser.Username {
		if isDuplicate, err := DuplicateAccountUserUsername(updatedAccountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorUsernameInUse)
			} else if err != nil {
				errs = append(errs, &err)
			}
		}
	}

	return packErrors(errs)
}

// AccountUserCredentialsValid checks is a certain user has the right credentials
func AccountUserCredentialsValid(password string, user *entity.AccountUser) error {
	pass, err := Base64Decode(user.Password)
	if err != nil {
		return err
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return fmt.Errorf("invalid password parts")
	}

	salt, err := Base64Decode(passwordParts[0])
	if err != nil {
		return err
	}

	timestamp, err := Base64Decode(passwordParts[1])
	if err != nil {
		return err
	}

	encryptedPassword := storageClient.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return fmt.Errorf("invalid user credentials")
	}

	return nil
}

// CheckAccountSession checks if the session is valid or not
func CheckAccountSession(r *http.Request) (string, string, error) {
	encodedSessionToken := r.Header.Get("x-tapglue-session")
	if encodedSessionToken == "" {
		return "", "failed to check session token (1)", fmt.Errorf("missing session token")
	}

	encodedIds := r.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "", "failed to check session token (2)\nmalformed token received", err
	}

	accountID, err := strconv.ParseInt(string(decodedIds), 10, 64)
	if err != nil {
		return "", "failed to check session token (3)\nmalformed token received", err
	}

	sessionToken, err := Base64Decode(encodedSessionToken)
	if err != nil {
		return "", "failed to check session token (4)\nmalformed token received", err
	}

	splitSessionToken := strings.SplitN(string(sessionToken), ":", 4)
	if len(splitSessionToken) != 4 {
		return "", "failed to check session token (5)\nmalformed token received", fmt.Errorf("malformed session token parts expected %d got %d", 4, len(splitSessionToken))
	}

	accID, err := strconv.ParseInt(splitSessionToken[0], 10, 64)
	if err != nil {
		return "", "failed to check session token (6)\nmalformed token received", err
	}

	userID, err := strconv.ParseInt(splitSessionToken[1], 10, 64)
	if err != nil {
		return "", "failed to check session token (7)\nmalformed token received", err
	}

	if accountID != accID {
		return "", "failed to check session token (8)\nmalformed token received", fmt.Errorf("account id mismatch expected %d got %d", accountID, accID)
	}

	sessionKey := storageClient.AccountSessionKey(accountID, userID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", "failed to check session token (9)\nmalformed token received", err
	}

	if storedSessionToken == "" {
		return "", "failed to check session token (10)\nsession not found", fmt.Errorf("session not found")
	}

	//fmt.Printf("storedSession\t%s\nencodedSession\t%s\n", storedSessionToken, encodedSessionToken)

	if storedSessionToken == encodedSessionToken {
		return encodedSessionToken, "", nil
	}

	return "", "failed to check session token (11)\nsession token mismatch", fmt.Errorf("session tokens mismatch expected %s got %s", storedSessionToken, encodedSessionToken)
}

// DuplicateAccountUserEmail checks if the user e-mail is duplicate within the provided account
func DuplicateAccountUserEmail(email string) (bool, error) {
	emailKey := storageClient.AccountUserByEmail(email)
	if userExists, err := storageEngine.Exists(emailKey).Result(); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, errorUserEmailAlreadyExists
		}
	}

	return false, nil
}

// DuplicateAccountUserUsername checks if the username is duplicate within the provided account
func DuplicateAccountUserUsername(username string) (bool, error) {
	usernameKey := storageClient.AccountUserByUsername(username)
	if userExists, err := storageEngine.Exists(usernameKey).Result(); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, errorUserUsernameAlreadyExists
		}
	}

	return false, nil
}
