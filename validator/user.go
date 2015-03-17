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

	if !StringLengthBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !StringLengthBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !StringLengthBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.MatchString(user.FirstName) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(user.LastName) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(user.Username) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if user.Email == "" || !IsValidEmail(user.Email) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !IsValidURL(user.URL, true) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		if !checkImages(user.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if !ApplicationExists(user.AccountID, user.ApplicationID) {
		errs = append(errs, &errorApplicationDoesNotExists)
	}

	if isDuplicate, err := DuplicateApplicationUserEmail(user.AccountID, user.ApplicationID, user.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, &errorUserEmailAlreadyExists)
		} else {
			errs = append(errs, &err)
		}
	}

	if isDuplicate, err := DuplicateApplicationUserUsername(user.AccountID, user.ApplicationID, user.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, &errorUserUsernameAlreadyExists)
		} else {
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}

// UpdateUser validates a user on update
func UpdateUser(existingApplicationUser, updatedApplicationUser *entity.User) error {
	errs := []*error{}

	if !StringLengthBetween(updatedApplicationUser.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !StringLengthBetween(updatedApplicationUser.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !StringLengthBetween(updatedApplicationUser.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.FirstName) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.LastName) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.Username) {
		errs = append(errs, &errorUserUsernameType)
	}

	if updatedApplicationUser.Email == "" || !IsValidEmail(updatedApplicationUser.Email) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if updatedApplicationUser.URL != "" && !IsValidURL(updatedApplicationUser.URL, true) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(updatedApplicationUser.Image) > 0 {
		if !checkImages(updatedApplicationUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if existingApplicationUser.Email != updatedApplicationUser.Email {
		if isDuplicate, err := DuplicateApplicationUserEmail(updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, updatedApplicationUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorEmailAddressInUse)
			} else if err != nil {
				errs = append(errs, &err)
			}
		}
	}

	if existingApplicationUser.Username != updatedApplicationUser.Username {
		if isDuplicate, err := DuplicateApplicationUserUsername(updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, updatedApplicationUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorUsernameInUse)
			} else if err != nil {
				errs = append(errs, &err)
			}
		}
	}

	return packErrors(errs)
}

// ApplicationUserCredentialsValid checks is a certain user has the right credentials
func ApplicationUserCredentialsValid(password string, user *entity.User) error {
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

// CheckApplicationSession checks if the session is valid or not
func CheckApplicationSession(r *http.Request) (string, string, error) {
	encodedSessionToken := r.Header.Get("x-tapglue-session")
	if encodedSessionToken == "" {
		return "", "failed to check session token (1)\nmissing session token", fmt.Errorf("missing session token")
	}

	encodedIds := r.Header.Get("x-tapglue-id")
	decodedIds, err := Base64Decode(encodedIds)
	if err != nil {
		return "", "failed to check session token (2)", err
	}

	ids := strings.SplitN(string(decodedIds), ":", 2)
	if len(ids) != 2 {
		return "", "failed to check session token (3)", fmt.Errorf("expected %d got %d", 2, len(ids))
	}

	accountID, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return "", "failed to check session token (4)", err
	}

	applicationID, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return "", "failed to check session token (5)", err
	}

	sessionToken, err := Base64Decode(encodedSessionToken)
	if err != nil {
		return "", "failed to check session token (6)", err
	}

	splitSessionToken := strings.SplitN(string(sessionToken), ":", 5)
	if len(splitSessionToken) != 5 {
		return "", "failed to check session token (7)", fmt.Errorf("expected %d got %d", 5, len(splitSessionToken))
	}

	accID, err := strconv.ParseInt(splitSessionToken[0], 10, 64)
	if err != nil {
		return "", "failed to check session token (8)", err
	}

	appID, err := strconv.ParseInt(splitSessionToken[1], 10, 64)
	if err != nil {
		return "", "failed to check session token (9)", err
	}

	userID, err := strconv.ParseInt(splitSessionToken[2], 10, 64)
	if err != nil {
		return "", "failed to check session token (10)", err
	}

	if accountID != accID {
		return "", "failed to check session token (11)", fmt.Errorf("expected %d got %d", accountID, accID)
	}

	if applicationID != appID {
		return "", "failed to check session token (12)", fmt.Errorf("expected %d got %d", applicationID, appID)
	}

	sessionKey := storageClient.ApplicationSessionKey(accountID, applicationID, userID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", "failed to check session token (13)", err
	}

	if storedSessionToken == "" {
		return "", "failed to check session token (14)", fmt.Errorf("session not found")
	}

	if storedSessionToken != encodedSessionToken {
		return encodedSessionToken, "", nil
	}

	return "", "failed to check session token (15)", fmt.Errorf("expected %s got %s", storedSessionToken, encodedSessionToken)
}

// CheckApplicationSimpleSession checks if the session is valid or not
func CheckApplicationSimpleSession(accountID, applicationID, applicationUserID int64, r *http.Request) (string, string, error) {
	encodedSessionToken := r.Header.Get("x-tapglue-session")
	if encodedSessionToken == "" {
		return "", "failed to check session token (1)", fmt.Errorf("missing session token")
	}

	sessionToken, err := Base64Decode(encodedSessionToken)
	if err != nil {
		return "", "failed to check session token (2)", err
	}

	splitSessionToken := strings.SplitN(string(sessionToken), ":", 5)
	if len(splitSessionToken) != 5 {
		return "", "failed to check session token (3)", fmt.Errorf("expected %d got %d", 5, len(splitSessionToken))
	}

	tokenAccountID, err := strconv.ParseInt(splitSessionToken[0], 10, 64)
	if err != nil {
		return "", "failed to check session token (4)", err
	}

	if tokenAccountID != accountID {
		return "", "failed to check session token (5)", fmt.Errorf("expected %d got %d", accountID, tokenAccountID)
	}

	tokenApplicationID, err := strconv.ParseInt(splitSessionToken[1], 10, 64)
	if err != nil {
		return "", "failed to check session token (6)", err
	}

	if tokenApplicationID != applicationID {
		return "", "failed to check session token (7)", fmt.Errorf("expected %d got %d", applicationID, tokenApplicationID)
	}

	tokenApplicationUserID, err := strconv.ParseInt(splitSessionToken[2], 10, 64)
	if err != nil {
		return "", "failed to check session token (8)", err
	}

	if tokenApplicationUserID != applicationUserID {
		return "", "failed to check session token (9)", fmt.Errorf("expected %d got %d", applicationUserID, tokenApplicationUserID)
	}

	sessionKey := storageClient.ApplicationSessionKey(accountID, applicationID, applicationUserID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", "failed to check session token (10)", err
	}

	if storedSessionToken == "" {
		return "", "failed to check session token (11)\nsession not found", fmt.Errorf("session not found")
	}

	if storedSessionToken == encodedSessionToken {
		return encodedSessionToken, "", nil
	}

	return "", "failed to check session token (12)\nsession mismatch", fmt.Errorf("expected %s got %s", storedSessionToken, encodedSessionToken)
}

// DuplicateApplicationUserEmail checks if the user email is duplicate within the application or not
func DuplicateApplicationUserEmail(accountID, applicationID int64, email string) (bool, error) {
	emailKey := storageClient.ApplicationUserByEmail(accountID, applicationID, email)
	if userExists, err := storageEngine.Exists(emailKey).Result(); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, errorUserEmailAlreadyExists
		}
	}

	return false, nil
}

// DuplicateApplicationUserUsername checks if the username is duplicate within the application or not
func DuplicateApplicationUserUsername(accountID, applicationID int64, username string) (bool, error) {
	usernameKey := storageClient.ApplicationUserByUsername(accountID, applicationID, username)
	if userExists, err := storageEngine.Exists(usernameKey).Result(); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, errorUserUsernameAlreadyExists
		}
	}

	return false, nil
}
