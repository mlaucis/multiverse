/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
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
	userNameMin = 2
	userNameMax = 40
)

var (
	errorUserFirstNameSize = fmt.Errorf("user first name must be between %d and %d characters", userNameMin, userNameMax)
	errorUserFirstNameType = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorUserLastNameSize = fmt.Errorf("user last name must be between %d and %d characters", userNameMin, userNameMax)
	errorUserLastNameType = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorUserUsernameSize = fmt.Errorf("user username must be between %d and %d characters", userNameMin, userNameMax)
	errorUserUsernameType = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorApplicationIDZero = fmt.Errorf("application id can't be 0")
	errorApplicationIDType = fmt.Errorf("application id is not a valid integer")

	errorAuthTokenInvalid = fmt.Errorf("auth token is invalid")
	errorUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorUserEmailInvalid = fmt.Errorf("user email is not valid")

	errorUserIDIsAlreadySet = fmt.Errorf("user id is already set")
)

// CreateUser validates a user on create
func CreateUser(user *entity.User) error {
	errs := []*error{}

	if !stringLenghtBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !stringLenghtBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !stringLenghtBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.FirstName)) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.LastName)) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.Username)) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if user.AuthToken == "" {
		errs = append(errs, &errorAuthTokenInvalid)
	}

	if user.Email == "" || !email.Match([]byte(user.Email)) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !url.Match([]byte(user.URL)) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		for _, image := range user.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !applicationExists(user.AccountID, user.ApplicationID) {
		errs = append(errs, &errorApplicationDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateUser validates a user on update
func UpdateUser(user *entity.User) error {
	errs := []*error{}

	if !stringLenghtBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !stringLenghtBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !stringLenghtBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.FirstName)) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.LastName)) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.Username)) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.AuthToken == "" {
		errs = append(errs, &errorAuthTokenInvalid)
	}

	if user.Email == "" || !email.Match([]byte(user.Email)) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !url.Match([]byte(user.URL)) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		for _, image := range user.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !applicationExists(user.AccountID, user.ApplicationID) {
		errs = append(errs, &errorApplicationDoesNotExists)
	}

	return packErrors(errs)
}

// UserCredentialsValid checks is a certain user has the right credentials
func ApplicationUserCredentialsValid(password string, user *entity.User) error {
	// TODO improve this with a salt and stuff
	if Base64Encode(Sha256String([]byte(password))) != user.Password {
		return fmt.Errorf("invalid user credentials")
	}

	return nil
}

// CheckApplicationSession checks if the session is valid or not
func CheckApplicationSession(r *http.Request) (string, error) {
	encodedSessionToken := r.Header.Get("x-tapglue-session")
	if encodedSessionToken == "" {
		return "", fmt.Errorf("missing session token")
	}

	encodedIds := r.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "", fmt.Errorf("ids not present in request")
	}

	ids := strings.SplitN(string(decodedIds), ":", 2)
	if len(ids) != 2 {
		return "", fmt.Errorf("malformed ids received")
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed ids received")
	}

	applicationID, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed ids received")
	}

	sessionToken, err := Base64Decode(encodedSessionToken)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	splitSessionToken := strings.SplitN(string(sessionToken), ":", 5)
	if len(splitSessionToken) != 5 {
		return "", fmt.Errorf("malformed session token")
	}

	accID, err := strconv.ParseInt(splitSessionToken[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	appID, err := strconv.ParseInt(splitSessionToken[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	userID, err := strconv.ParseInt(splitSessionToken[2], 10, 64)
	if err != nil {
		return "", fmt.Errorf("malformed session token")
	}

	if accountID != accID {
		return "", fmt.Errorf("session token mismatch(1)")
	}

	if applicationID != appID {
		return "", fmt.Errorf("session token mismatch(2)")
	}

	sessionKey := storageClient.ApplicationSessionKey(accountID, applicationID, userID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", fmt.Errorf("could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", fmt.Errorf("session not found")
	}

	if storedSessionToken != encodedSessionToken {
		return "", fmt.Errorf("session token mismatch(3)")
	}

	return encodedSessionToken, nil
}
