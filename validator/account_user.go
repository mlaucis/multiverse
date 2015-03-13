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

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.FirstName)) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.LastName)) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.Username)) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if accountUser.Email == "" || !email.Match([]byte(accountUser.Email)) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !url.Match([]byte(accountUser.URL)) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(accountUser.Image) > 0 {
		for _, image := range accountUser.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !AccountExists(accountUser.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(accountUser *entity.AccountUser) error {
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

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.FirstName)) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.LastName)) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.Username)) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.Email == "" || !email.Match([]byte(accountUser.Email)) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !url.Match([]byte(accountUser.URL)) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(accountUser.Image) > 0 {
		for _, image := range accountUser.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !AccountExists(accountUser.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}

// UserCredentialsValid checks is a certain user has the right credentials
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
func CheckAccountSession(r *http.Request) (string, error) {
	encodedSessionToken := r.Header.Get("x-tapglue-session")
	if encodedSessionToken == "" {
		return "", fmt.Errorf("missing session token")
	}

	encodedIds := r.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "", fmt.Errorf("ids not present in request")
	}

	accountID, err := strconv.ParseInt(string(decodedIds), 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed ids received")
	}

	sessionToken, err := Base64Decode(encodedSessionToken)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	splitSessionToken := strings.SplitN(string(sessionToken), ":", 4)
	if len(splitSessionToken) != 4 {
		return "", fmt.Errorf("malformed session token")
	}

	accID, err := strconv.ParseInt(splitSessionToken[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	userID, err := strconv.ParseInt(splitSessionToken[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	if accountID != accID {
		return "", fmt.Errorf("session token mismatch(1)")
	}

	sessionKey := storageClient.AccountSessionKey(accountID, userID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", fmt.Errorf("could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", fmt.Errorf("session not found")
	}

	//fmt.Printf("storedSession\t%s\nencodedSession\t%s\n", storedSessionToken, encodedSessionToken)

	if storedSessionToken != encodedSessionToken {
		return "", fmt.Errorf("session token mismatch(3)")
	}

	return encodedSessionToken, nil
}
