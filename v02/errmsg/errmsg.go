/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package errmsg

import (
	"net/http"

	"github.com/tapglue/backend/errors"
)

// These are all the errors in the application, sorted alphabetically by their name
var (
	// Account errors
	ErrAccountDescriptionSize  = errors.New(http.StatusBadRequest, 0, "account description must be between 0 and 100 characters", "", false)
	ErrAccountDescriptionType  = errors.New(http.StatusBadRequest, 0, "account description is not a valid alphanumeric sequence", "", false)
	ErrAccountIDIsAlreadySet   = errors.New(http.StatusBadRequest, 0, "account id is already set", "", false)
	ErrAccountIDMismatch       = errors.New(http.StatusConflict, 0, "account ID mismatch", "", false)
	ErrAccountIDZero           = errors.New(http.StatusBadRequest, 0, "account id can't be 0", "", false)
	ErrAccountMismatch         = errors.New(http.StatusBadRequest, 0, "account mismatch", "", false)
	ErrAccountMissingInContext = errors.New(http.StatusInternalServerError, 0, "missing account context", "", false)
	ErrAccountNameSize         = errors.New(http.StatusBadRequest, 0, "account name must be between 3 and 40 characters", "", false)
	ErrAccountNameType         = errors.New(http.StatusBadRequest, 0, "account name is not a valid alphanumeric sequence", "", false)
	ErrAccountNotFound         = errors.New(http.StatusNotFound, 0, "account not found", "", false)
	ErrAccountTokenAlreadySet  = errors.New(http.StatusBadRequest, 0, "account token is already set", "", false)

	// Account user errors
	ErrAccountUserEmailInvalid  = errors.New(http.StatusBadRequest, 0, "user email is not valid", "", false)
	ErrAccountUserFirstNameSize = errors.New(http.StatusBadRequest, 0, "user first name must be between 2 and 40 characters", "", false)
	ErrAccountUserFirstNameType = errors.New(http.StatusBadRequest, 0, "user first name is not a valid alphanumeric sequence", "", false)
	ErrAccountUserLastNameSize  = errors.New(http.StatusBadRequest, 0, "user last name must be between 2 and 40 characters", "", false)
	ErrAccountUserLastNameType  = errors.New(http.StatusBadRequest, 0, "user last name is not a valid alphanumeric sequence", "", false)
	ErrAccountUserMismatchErr   = errors.New(http.StatusConflict, 0, "account user mismatch", "", false)
	ErrAccountUserNotFound      = errors.New(http.StatusNotFound, 0, "account user not found", "", false)
	ErrAccountUserPasswordSize  = errors.New(http.StatusBadRequest, 0, "user password must be between 4 and 60 characters", "", false)
	ErrAccountUserURLInvalid    = errors.New(http.StatusBadRequest, 0, "user url is not a valid url", "", false)
	ErrAccountUserUsernameSize  = errors.New(http.StatusBadRequest, 0, "user username must be between 2 and 40 characters", "", false)
	ErrAccountUserUsernameType  = errors.New(http.StatusBadRequest, 0, "user username is not a valid alphanumeric sequence", "", false)

	// Application errors
	ErrApplicationAuthTokenUpdateNotAllowed = errors.New(http.StatusBadRequest, 0, "not allowed to update the application token", "", false)
	ErrApplicationDescriptionSize           = errors.New(http.StatusBadRequest, 0, "application description must be between 0 and 100 characters", "", false)
	ErrApplicationDescriptionType           = errors.New(http.StatusBadRequest, 0, "application description is not a valid alphanumeric sequence", "", false)
	ErrApplicationIDInvalid                 = errors.New(http.StatusBadRequest, 0, "application id is not valid", "", false)
	ErrApplicationIDIsAlreadySet            = errors.New(http.StatusBadRequest, 0, "application id is already set", "", false)
	ErrApplicationNameSize                  = errors.New(http.StatusBadRequest, 0, "application name must be between 2 and 40 characters", "", false)
	ErrApplicationNameType                  = errors.New(http.StatusBadRequest, 0, "application name is not a valid alphanumeric sequence", "", false)
	ErrApplicationNotFound                  = errors.New(http.StatusNotFound, 0, "application not found", "application not found", false)
	ErrApplicationURLInvalid                = errors.New(http.StatusBadRequest, 0, "application url is not a valid url", "", false)

	// Application user errors
	ErrApplicationUserNotActivated          = errors.New(http.StatusNotAcceptable, 0, "user not activated", "", false)
	ErrApplicationUserNotFound              = errors.New(http.StatusNotFound, 0, "application user not found", "user not found", false)
	ErrApplicationUserEmailAlreadyExists    = errors.New(http.StatusBadRequest, 0, "email address already in use", "", false)
	ErrApplicationUserEmailInvalid          = errors.New(http.StatusBadRequest, 0, "user email is not valid", "", false)
	ErrApplicationUserFirstNameSize         = errors.New(http.StatusBadRequest, 0, "user first name must be between 2 and 40 characters", "", false)
	ErrApplicationUserFirstNameType         = errors.New(http.StatusBadRequest, 0, "user first name is not a valid alphanumeric sequence", "", false)
	ErrApplicationUserIDInvalid             = errors.New(http.StatusBadRequest, 0, "user id is not valid", "", false)
	ErrApplicationUserLastNameSize          = errors.New(http.StatusBadRequest, 0, "user last name must be between 2 and 40 characters", "", false)
	ErrApplicationUserLastNameType          = errors.New(http.StatusBadRequest, 0, "user last name is not a valid alphanumeric sequence", "", false)
	ErrApplicationUsernameAndEmailAreEmpty  = errors.New(http.StatusBadRequest, 0, "user email and username are both empty", "", false)
	ErrApplicationUserUsernameInUse         = errors.New(http.StatusBadRequest, 0, "username already in use", "", false)
	ErrApplicationUserSearchTypeMin3Chars   = errors.New(http.StatusBadRequest, 0, "type at least 3 characters to search", "", false)
	ErrApplicationUsersNotConnected         = errors.New(http.StatusNotFound, 0, "users are not connected", "", false)
	ErrApplicationUserURLInvalid            = errors.New(http.StatusBadRequest, 0, "user url is not a valid url", "", false)
	ErrApplicationUserUsernameAlreadyExists = errors.New(http.StatusBadRequest, 0, "user already exists", "", false)
	ErrApplicationUserUsernameSize          = errors.New(http.StatusBadRequest, 0, "user username must be between 2 and 40 characters", "", false)
	ErrApplicationUserUsernameType          = errors.New(http.StatusBadRequest, 0, "user username is not a valid alphanumeric sequence", "", false)

	// Authentication errors
	ErrAuthGeneric                           = errors.New(http.StatusBadRequest, 0, "authentication error", "", false)
	ErrAuthGotBothUsernameAndEmail           = errors.New(http.StatusBadRequest, 0, "both username and email are specified", "", false)
	ErrAuthGotNoUsernameOrEmail              = errors.New(http.StatusBadRequest, 0, "both username and email are empty", "", false)
	ErrAuthInvalidAccountCredentials         = errors.New(http.StatusBadRequest, 0, "error while reading account credentials", "", false)
	ErrAuthInvalidAccountUserCredentials     = errors.New(http.StatusBadRequest, 0, "error while reading account user credentials", "", false)
	ErrAuthInvalidApplicationCredentials     = errors.New(http.StatusBadRequest, 0, "error while reading application credentials", "", false)
	ErrAuthInvalidApplicationUserCredentials = errors.New(http.StatusBadRequest, 0, "error while reading user credentials", "", false)
	ErrAuthInvalidEmailAddress               = errors.New(http.StatusBadRequest, 0, "invalid email address", "", false)
	ErrAuthMethodNotSupported                = errors.New(http.StatusBadRequest, 0, "authorization method not supported", "auth method not supported", false)
	ErrAuthPasswordEmpty                     = errors.New(http.StatusBadRequest, 0, "password is empty", "", false)
	ErrAuthPasswordMismatch                  = errors.New(http.StatusBadRequest, 0, "different passwords", "", false)
	ErrAuthSessionTokenMismatch              = errors.New(http.StatusBadRequest, 0, "session token mismatch", "", false)

	// Connection errors
	ErrConnectionAlreadyExists      = errors.New(http.StatusBadRequest, 0, "connection already exists", "", false)
	ErrConnectionNotFound           = errors.New(http.StatusNotFound, 0, "connection not found", "", false)
	ErrConnectionTypeIsWrong        = errors.New(http.StatusBadRequest, 0, "unexpected connection type", "", false)
	ErrConnectionSelfConnectingUser = errors.New(http.StatusBadRequest, 0, "user is connecting with itself", "self-connecting user", false)

	// Event errors
	ErrEventGeoRadiusAndNearestMissing = errors.New(http.StatusBadRequest, 0, "you must specify either a radius or a how many nearest events you want", "invalid radius and nearest", false)
	ErrEventGeoRadiusUnder2M           = errors.New(http.StatusBadRequest, 0, "Location radius can't be smaller than 2 meters", "radius smaller than 2", false)
	ErrEventIDInvalid                  = errors.New(http.StatusBadRequest, 0, "event id is not valid", "", false)
	ErrEventIDIsAlreadySet             = errors.New(http.StatusBadRequest, 0, "event id is already set", "", false)
	ErrEventInvalidVisiblity           = errors.New(http.StatusBadRequest, 0, "event visibility is invalid", "", false)
	ErrEventMissingVisiblity           = errors.New(http.StatusBadRequest, 0, "event visibility is missing", "", false)
	ErrEventNearestNotInBounds         = errors.New(http.StatusBadRequest, 0, "near events limits not within accepted bounds", "nearest not within bounds", false)
	ErrEventNotFound                   = errors.New(http.StatusNotFound, 0, "event not found", "", false)
	ErrEventTypeSize                   = errors.New(http.StatusBadRequest, 0, "type must be between 1 and 30 characters", "", false)
	ErrEventTypeType                   = errors.New(http.StatusBadRequest, 0, "type is not a valid alphanumeric sequence", "", false)

	// Internal account errors
	ErrInternalAccountCreation = errors.New(http.StatusInternalServerError, 0, "error while creating the account", "", false)
	ErrInternalAccountDelete   = errors.New(http.StatusInternalServerError, 0, "error while deleting the account", "", false)
	ErrInternalAccountRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the account", "", false)
	ErrInternalAccountUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the account", "", false)

	// Internal account user errors
	ErrInternalAccountUserCreation        = errors.New(http.StatusInternalServerError, 0, "error while creating the account user", "", false)
	ErrInternalAccountUserList            = errors.New(http.StatusInternalServerError, 0, "error while retrieving list of account users", "", false)
	ErrInternalAccountUserRead            = errors.New(http.StatusInternalServerError, 0, "error while reading the account user", "", false)
	ErrInternalAccountUserSessionCreation = errors.New(http.StatusInternalServerError, 0, "error while creating the account user session", "", false)
	ErrInternalAccountUserSessionDelete   = errors.New(http.StatusInternalServerError, 0, "error while deleting the account user session", "", false)
	ErrInternalAccountUserSessionRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the account user session", "", false)
	ErrInternalAccountUserSessionUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the account user session", "", false)
	ErrInternalAccountUserUpdate          = errors.New(http.StatusInternalServerError, 0, "error while updating the account user", "", false)

	// Internal application errors
	ErrInternalApplicationCreation = errors.New(http.StatusInternalServerError, 0, "error while creating the application", "", false)
	ErrInternalApplicationDelete   = errors.New(http.StatusInternalServerError, 0, "error while deleting the application", "", false)
	ErrInternalApplicationRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the application", "", false)
	ErrInternalApplicationUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the application", "", false)

	// Internal application user errors
	ErrInternalApplicationUserCreation        = errors.New(http.StatusInternalServerError, 0, "error while creating the application user", "", false)
	ErrInternalApplicationUserList            = errors.New(http.StatusInternalServerError, 0, "error while retrieving list of application users", "", false)
	ErrInternalApplicationUserRead            = errors.New(http.StatusInternalServerError, 0, "error while reading the application user", "", false)
	ErrInternalApplicationUserSessionCreation = errors.New(http.StatusInternalServerError, 0, "error while creating the application user session", "", false)
	ErrInternalApplicationUserSessionDelete   = errors.New(http.StatusInternalServerError, 0, "error while deleting session", "", false)
	ErrInternalApplicationUserSessionRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the application user session", "", false)
	ErrInternalApplicationUserSessionsDelete  = errors.New(http.StatusInternalServerError, 0, "error while deleting all session", "", false)
	ErrInternalApplicationUserSessionUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the application user session", "", false)
	ErrInternalApplicationUserUpdate          = errors.New(http.StatusInternalServerError, 0, "error while updating the application user", "", false)

	// Internal connection errors
	ErrInternalConnectingUsers    = errors.New(http.StatusInternalServerError, 0, "error while connecting the users", "", false)
	ErrInternalConnectionCreation = errors.New(http.StatusInternalServerError, 0, "error while creating the connection", "", false)
	ErrInternalConnectionRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the connection", "", false)
	ErrInternalConnectionUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the connection", "", false)

	// Internal event errors
	ErrInternalEventCreation = errors.New(http.StatusInternalServerError, 0, "error while saving the event", "", false)
	ErrInternalEventRead     = errors.New(http.StatusInternalServerError, 0, "error while reading the event", "", false)
	ErrInternalEventsList    = errors.New(http.StatusInternalServerError, 0, "failed to read the events", "", false)
	ErrInternalEventUpdate   = errors.New(http.StatusInternalServerError, 0, "error while updating the event", "", false)
	ErrInternalFollowersList = errors.New(http.StatusInternalServerError, 0, "error while retrieving list of followers", "", false)
	ErrInternalFollowingList = errors.New(http.StatusInternalServerError, 0, "error while retrieving list of following", "", false)
	ErrInternalFriendsList   = errors.New(http.StatusInternalServerError, 0, "error while retrieving list of friends", "", false)

	// Server errors
	ErrServerNotImplementedYet = errors.New(http.StatusInternalServerError, 0, "not implemented yet", "", false)

	// Server request errors
	ErrServerReqBadJSONReceived = errors.New(http.StatusBadRequest, 0, "malformed json received", "", false)
	ErrServerReqBadUserAgent = errors.New(http.StatusBadRequest, 0, "User-Agent header must be set (1)", "missing ua header", false)
	ErrServerReqContentLengthInvalid = errors.New(http.StatusBadRequest, 0, "Content-Length header is invalid", "content-length header is not an int", false)
	ErrServerReqContentLengthMissing = errors.New(http.StatusBadRequest, 0, "Content-Length header missing", "missing content-length header", false)
	ErrServerReqContentLengthSizeMismatch = errors.New(http.StatusBadRequest, 0, "Content-Length header size mismatch", "content-length header size mismatch", false)
	ErrServerReqContentTypeMismatch = errors.New(http.StatusBadRequest, 0, "Content-Type header mismatch", "content-type header mismatch", false)
	ErrServerReqContentTypeMissing = errors.New(http.StatusBadRequest, 0, "Content-Type header empty", "missing content-type header", false)
	ErrServerReqMissingJarvisID = errors.New(http.StatusNotFound, 0, "not found", "request does not contain a correct Jarvis auth", false)
	ErrServerReqNoKnownSearchTermsSupplied = errors.New(http.StatusBadRequest, 0, "no known search terms supplied", "no known search terms supplied", false)
	ErrServerReqParseFloat = errors.New(http.StatusBadRequest, 0, "", "parse float error", false)
	ErrServerReqPayloadTooBig = errors.New(http.StatusRequestEntityTooLarge, 0, "payload too big", "fat payload detected", false)
	ErrServerReqBodyEmpty = errors.New(http.StatusBadRequest, 0, "Empty request body", "empty request body", false)

	// Server response errors
	ErrServerRespMissingLastModifiedHeader = errors.New(http.StatusInternalServerError, 0, "something went wrong", "", false)

	// Misc errors
	ErrInvalidImageURL = errors.New(http.StatusBadRequest, 0, "image url is not valid", "", false)
)
