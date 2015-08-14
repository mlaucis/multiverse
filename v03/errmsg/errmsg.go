// Package errmsg holds the error messages accros the application
package errmsg

import (
	"net/http"

	"github.com/tapglue/backend/errors"
)

// These are all the errors in the application, sorted alphabetically by their name
var (
	// Application user errors

	ErrApplicationUserNotActivated         = errors.New(http.StatusNotAcceptable, 1000, "user not activated", "", false)
	ErrApplicationUserNotFound             = errors.New(http.StatusNotFound, 1001, "application user not found", "", false)
	ErrApplicationUserEmailAlreadyExists   = errors.New(http.StatusBadRequest, 1002, "email address already in use", "", false)
	ErrApplicationUserEmailInvalid         = errors.New(http.StatusBadRequest, 1003, "user email is not valid", "", false)
	ErrApplicationUserFirstNameSize        = errors.New(http.StatusBadRequest, 1004, "user first name must be between 2 and 40 characters", "", false)
	ErrApplicationUserIDInvalid            = errors.New(http.StatusBadRequest, 1005, "user id is not valid", "", false)
	ErrApplicationUserLastNameSize         = errors.New(http.StatusBadRequest, 1006, "user last name must be between 2 and 40 characters", "", false)
	ErrApplicationUsernameAndEmailAreEmpty = errors.New(http.StatusBadRequest, 1007, "user email and username are both empty", "", false)
	ErrApplicationUserUsernameInUse        = errors.New(http.StatusBadRequest, 1008, "username already in use", "", false)
	ErrApplicationUserSearchTypeMin3Chars  = errors.New(http.StatusBadRequest, 1009, "type at least 3 characters to search", "", false)
	ErrApplicationUserURLInvalid           = errors.New(http.StatusBadRequest, 1010, "user url is not a valid url", "", false)
	ErrApplicationUserUsernameSize         = errors.New(http.StatusBadRequest, 1012, "user username must be between 2 and 40 characters", "", false)

	// Internal application user errors

	ErrInternalApplicationUserCreation        = errors.New(http.StatusInternalServerError, 1500, "error while creating the application user", "", false)
	ErrInternalApplicationUserList            = errors.New(http.StatusInternalServerError, 1501, "error while retrieving list of application users", "", false)
	ErrInternalApplicationUserRead            = errors.New(http.StatusInternalServerError, 1502, "error while reading the application user", "", false)
	ErrInternalApplicationUserSessionCreation = errors.New(http.StatusInternalServerError, 1503, "error while creating the application user session", "", false)
	ErrInternalApplicationUserSessionDelete   = errors.New(http.StatusInternalServerError, 1504, "error while deleting session", "", false)
	ErrInternalApplicationUserSessionRead     = errors.New(http.StatusInternalServerError, 1505, "error while reading the application user session", "", false)
	ErrInternalApplicationUserSessionsDelete  = errors.New(http.StatusInternalServerError, 1506, "error while deleting all session", "", false)
	ErrInternalApplicationUserSessionUpdate   = errors.New(http.StatusInternalServerError, 1507, "error while updating the application user session", "", false)
	ErrInternalApplicationUserUpdate          = errors.New(http.StatusInternalServerError, 1508, "error while updating the application user", "", false)
	ErrInternalApplicationUserIDMissing       = errors.New(http.StatusInternalServerError, 1509, "user ID is missing", "", false)

	// Connection errors

	ErrConnectionAlreadyExists      = errors.New(http.StatusBadRequest, 2000, "connection already exists", "", false)
	ErrConnectionNotFound           = errors.New(http.StatusNotFound, 2001, "connection not found", "", false)
	ErrConnectionTypeIsWrong        = errors.New(http.StatusBadRequest, 2002, "unexpected connection type", "", false)
	ErrConnectionSelfConnectingUser = errors.New(http.StatusBadRequest, 2003, "user is connecting with itself", "self-connecting user", false)
	ErrConnectionUsersNotConnected  = errors.New(http.StatusNotFound, 2004, "users are not connected", "", false)

	// Internal connection errors

	ErrInternalConnectingUsers    = errors.New(http.StatusInternalServerError, 2500, "error while connecting the users", "", false)
	ErrInternalConnectionCreation = errors.New(http.StatusInternalServerError, 2501, "error while creating the connection", "", false)
	ErrInternalConnectionRead     = errors.New(http.StatusInternalServerError, 2502, "error while reading the connection", "", false)
	ErrInternalConnectionUpdate   = errors.New(http.StatusInternalServerError, 2503, "error while updating the connection", "", false)
	ErrInternalFollowersList      = errors.New(http.StatusInternalServerError, 2504, "error while retrieving list of followers", "", false)
	ErrInternalFollowingList      = errors.New(http.StatusInternalServerError, 2505, "error while retrieving list of following", "", false)
	ErrInternalFriendsList        = errors.New(http.StatusInternalServerError, 2506, "error while retrieving list of friends", "", false)

	// Event errors

	ErrEventGeoRadiusAndNearestMissing = errors.New(http.StatusBadRequest, 3000, "you must specify either a radius or a how many nearest events you want", "invalid radius and nearest", false)
	ErrEventGeoRadiusUnder2M           = errors.New(http.StatusBadRequest, 3001, "Location radius can't be smaller than 2 meters", "radius smaller than 2", false)
	ErrEventIDInvalid                  = errors.New(http.StatusBadRequest, 3002, "event id is not valid", "", false)
	ErrEventIDIsAlreadySet             = errors.New(http.StatusBadRequest, 3003, "event id is already set", "", false)
	ErrEventInvalidVisiblity           = errors.New(http.StatusBadRequest, 3004, "event visibility is invalid", "", false)
	ErrEventMissingVisiblity           = errors.New(http.StatusBadRequest, 3005, "event visibility is missing", "", false)
	ErrEventNearestNotInBounds         = errors.New(http.StatusBadRequest, 3006, "near events limits not within accepted bounds", "nearest not within bounds", false)
	ErrEventNotFound                   = errors.New(http.StatusNotFound, 3007, "event not found", "", false)
	ErrEventTypeSize                   = errors.New(http.StatusBadRequest, 3008, "type must be between 1 and 30 characters", "", false)

	// Internal event errors

	ErrInternalEventCreation  = errors.New(http.StatusInternalServerError, 3500, "error while saving the event", "", false)
	ErrInternalEventRead      = errors.New(http.StatusInternalServerError, 3501, "error while reading the event", "", false)
	ErrInternalEventsList     = errors.New(http.StatusInternalServerError, 3502, "failed to read the events", "", false)
	ErrInternalEventUpdate    = errors.New(http.StatusInternalServerError, 3503, "error while updating the event", "", false)
	ErrInternalEventMissingID = errors.New(http.StatusInternalServerError, 3504, "event is missing ID", "", false)

	// Authentication errors

	ErrAuthGeneric                           = errors.New(http.StatusBadRequest, 4001, "authentication error", "", false)
	ErrAuthGotBothUsernameAndEmail           = errors.New(http.StatusBadRequest, 4002, "both username and email are specified", "", false)
	ErrAuthGotNoUsernameOrEmail              = errors.New(http.StatusBadRequest, 4003, "both username and email are empty", "", false)
	ErrAuthInvalidAccountCredentials         = errors.New(http.StatusBadRequest, 4004, "error while reading account credentials", "", false)
	ErrAuthInvalidAccountUserCredentials     = errors.New(http.StatusBadRequest, 4005, "error while reading account user credentials", "", false)
	ErrAuthInvalidApplicationCredentials     = errors.New(http.StatusBadRequest, 4006, "error while reading application credentials", "", false)
	ErrAuthInvalidApplicationUserCredentials = errors.New(http.StatusBadRequest, 4007, "error while reading user credentials", "", false)
	ErrAuthInvalidEmailAddress               = errors.New(http.StatusBadRequest, 4008, "invalid email address", "", false)
	ErrAuthMethodNotSupported                = errors.New(http.StatusBadRequest, 4009, "authorization method not supported", "auth method not supported", false)
	ErrAuthPasswordEmpty                     = errors.New(http.StatusBadRequest, 4010, "password is empty", "", false)
	ErrAuthPasswordMismatch                  = errors.New(http.StatusBadRequest, 4011, "different passwords", "", false)
	ErrAuthSessionTokenMismatch              = errors.New(http.StatusBadRequest, 4012, "session token mismatch", "", false)
	ErrAuthUserSessionNotSet                 = errors.New(http.StatusBadRequest, 4013, "session token missing from request", "", false)

	// Server errors

	ErrInvalidImageURL = errors.New(http.StatusBadRequest, 5000, "image url is not valid", "", false)

	// Server request errors

	ErrServerReqBadJSONReceived            = errors.New(http.StatusBadRequest, 5001, "malformed json received", "", false)
	ErrServerReqBadUserAgent               = errors.New(http.StatusBadRequest, 5002, "User-Agent header must be set (1)", "missing ua header", false)
	ErrServerReqContentLengthInvalid       = errors.New(http.StatusBadRequest, 5003, "Content-Length header is invalid", "content-length header is not an int", false)
	ErrServerReqContentLengthMissing       = errors.New(http.StatusBadRequest, 5004, "Content-Length header missing", "missing content-length header", false)
	ErrServerReqContentLengthSizeMismatch  = errors.New(http.StatusBadRequest, 5005, "Content-Length header size mismatch", "content-length header size mismatch", false)
	ErrServerReqContentTypeMismatch        = errors.New(http.StatusBadRequest, 5006, "Content-Type header mismatch", "content-type header mismatch", false)
	ErrServerReqContentTypeMissing         = errors.New(http.StatusBadRequest, 5007, "Content-Type header empty", "missing content-type header", false)
	ErrServerReqNoKnownSearchTermsSupplied = errors.New(http.StatusBadRequest, 5009, "no known search terms supplied", "no known search terms supplied", false)
	ErrServerReqParseFloat                 = errors.New(http.StatusBadRequest, 5010, "", "parse float error", false)
	ErrServerReqPayloadTooBig              = errors.New(http.StatusRequestEntityTooLarge, 5011, "payload too big", "fat payload detected", false)
	ErrServerReqBodyEmpty                  = errors.New(http.StatusBadRequest, 5012, "Empty request body", "empty request body", false)

	// Server response errors

	ErrServerRespMissingLastModifiedHeader = errors.New(http.StatusInternalServerError, 5013, "something went wrong", "", false)

	// Misc errors

	ErrServerNotImplementedYet         = errors.New(http.StatusInternalServerError, 5500, "not implemented yet", "", false)
	ErrServerDeprecatedStorage         = errors.New(http.StatusInternalServerError, 5501, "deprecated storage", "", false)
	ErrServerInvalidHandler            = errors.New(http.StatusInternalServerError, 5502, "something went wrong", "handler used in wrong context", false)
	ErrServerUnsportedHandlerOperation = errors.New(http.StatusInternalServerError, 5503, "something went wrong", "handler does not support operation", false)
	ErrServerInternalError             = errors.New(http.StatusInternalServerError, 5504, "something went wrong", "", false)

	// Account errors

	ErrAccountDescriptionSize  = errors.New(http.StatusBadRequest, 6000, "account description must be between 0 and 100 characters", "", false)
	ErrOrgIDIsAlreadySet       = errors.New(http.StatusBadRequest, 6001, "account id is already set", "", false)
	ErrOrgIDZero               = errors.New(http.StatusBadRequest, 6002, "account id can't be 0", "", false)
	ErrAccountMismatch         = errors.New(http.StatusBadRequest, 6003, "organization mismatch", "", false)
	ErrAccountMissingInContext = errors.New(http.StatusInternalServerError, 6004, "missing account context", "", false)
	ErrAccountNameSize         = errors.New(http.StatusBadRequest, 6005, "account name must be between 3 and 40 characters", "", false)
	ErrAccountNotFound         = errors.New(http.StatusNotFound, 6006, "account not found", "", false)
	ErrOrgTokenAlreadySet      = errors.New(http.StatusBadRequest, 6007, "account token is already set", "", false)

	// Internal account errors

	ErrInternalAccountCreation = errors.New(http.StatusInternalServerError, 6500, "error while creating the account", "", false)
	ErrInternalAccountDelete   = errors.New(http.StatusInternalServerError, 6501, "error while deleting the account", "", false)
	ErrInternalAccountRead     = errors.New(http.StatusInternalServerError, 6502, "error while reading the account", "", false)
	ErrInternalAccountUpdate   = errors.New(http.StatusInternalServerError, 6503, "error while updating the account", "", false)

	// Account user errors

	ErrMemberEmailInvalid  = errors.New(http.StatusBadRequest, 7000, "user email is not valid", "", false)
	ErrMemberFirstNameSize = errors.New(http.StatusBadRequest, 7001, "user first name must be between 2 and 40 characters", "", false)
	ErrMemberLastNameSize  = errors.New(http.StatusBadRequest, 7002, "user last name must be between 2 and 40 characters", "", false)
	ErrMemberMismatchErr   = errors.New(http.StatusConflict, 7003, "member mismatch", "", false)
	ErrMemberNotFound      = errors.New(http.StatusNotFound, 7004, "member not found", "", false)
	ErrMemberPasswordSize  = errors.New(http.StatusBadRequest, 7005, "user password must be between 4 and 60 characters", "", false)
	ErrMemberURLInvalid    = errors.New(http.StatusBadRequest, 7006, "user url is not a valid url", "", false)
	ErrMemberUsernameSize  = errors.New(http.StatusBadRequest, 7007, "user username must be between 2 and 40 characters", "", false)

	// Internal account user errors

	ErrInternalAccountUserCreation        = errors.New(http.StatusInternalServerError, 7500, "error while creating the account user", "", false)
	ErrInternalAccountUserList            = errors.New(http.StatusInternalServerError, 7501, "error while retrieving list of account users", "", false)
	ErrInternalAccountUserRead            = errors.New(http.StatusInternalServerError, 7502, "error while reading the account user", "", false)
	ErrInternalAccountUserSessionCreation = errors.New(http.StatusInternalServerError, 7503, "error while creating the account user session", "", false)
	ErrInternalAccountUserSessionDelete   = errors.New(http.StatusInternalServerError, 7504, "error while deleting the account user session", "", false)
	ErrInternalAccountUserSessionRead     = errors.New(http.StatusInternalServerError, 7505, "error while reading the account user session", "", false)
	ErrInternalAccountUserSessionUpdate   = errors.New(http.StatusInternalServerError, 7506, "error while updating the account user session", "", false)
	ErrInternalAccountUserUpdate          = errors.New(http.StatusInternalServerError, 7507, "error while updating the account user", "", false)

	// Application errors

	ErrApplicationAuthTokenUpdateNotAllowed = errors.New(http.StatusBadRequest, 8000, "not allowed to update the application token", "", false)
	ErrApplicationDescriptionSize           = errors.New(http.StatusBadRequest, 8001, "application description must be between 0 and 100 characters", "", false)
	ErrApplicationIDInvalid                 = errors.New(http.StatusBadRequest, 8002, "application id is not valid", "", false)
	ErrApplicationIDIsAlreadySet            = errors.New(http.StatusBadRequest, 8003, "application id is already set", "", false)
	ErrApplicationNameSize                  = errors.New(http.StatusBadRequest, 8004, "application name must be between 2 and 40 characters", "", false)
	ErrApplicationNotFound                  = errors.New(http.StatusNotFound, 8005, "application not found", "application not found", false)
	ErrApplicationURLInvalid                = errors.New(http.StatusBadRequest, 8006, "application url is not a valid url", "", false)

	// Internal application errors

	ErrInternalApplicationCreation = errors.New(http.StatusInternalServerError, 8500, "error while creating the application", "", false)
	ErrInternalApplicationDelete   = errors.New(http.StatusInternalServerError, 8501, "error while deleting the application", "", false)
	ErrInternalApplicationRead     = errors.New(http.StatusInternalServerError, 8502, "error while reading the application", "", false)
	ErrInternalApplicationUpdate   = errors.New(http.StatusInternalServerError, 8503, "error while updating the application", "", false)
)
