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
	ErrAccountDescriptionSize                 = errors.NewBadRequestError(0, "account description must be between %d and %d characters", "")
	ErrAccountDescriptionType                 = errors.NewBadRequestError(0, "account description is not a valid alphanumeric sequence", "")
	ErrAccountIDIsAlreadySet                  = errors.NewBadRequestError(0, "account id is already set", "")
	ErrAccountIDMismatch                      = errors.New(http.StatusConflict, 0, "account ID mismatch", "", false)
	ErrAccountIDZero                          = errors.NewBadRequestError(0, "account id can't be 0", "")
	ErrAccountMismatch                        = errors.NewBadRequestError(0, "account mismatch", "")
	ErrAccountNameSize                        = errors.NewBadRequestError(0, "account name must be between %d and %d characters", "")
	ErrAccountNameType                        = errors.NewBadRequestError(0, "account name is not a valid alphanumeric sequence", "")
	ErrAccountNotFound                        = errors.NewNotFoundError(0, "account not found", "")
	ErrAccountTokenAlreadySet                 = errors.NewBadRequestError(0, "account token is already set", "")
	ErrAccountUserEmailInvalid                = errors.NewBadRequestError(0, "user email is not valid", "")
	ErrAccountUserFirstNameSize               = errors.NewBadRequestError(0, "user first name must be between %d and %d characters", "")
	ErrAccountUserFirstNameType               = errors.NewBadRequestError(0, "user first name is not a valid alphanumeric sequence", "")
	ErrAccountUserLastNameSize                = errors.NewBadRequestError(0, "user last name must be between %d and %d characters", "")
	ErrAccountUserLastNameType                = errors.NewBadRequestError(0, "user last name is not a valid alphanumeric sequence", "")
	ErrAccountUserMismatchErr                 = errors.New(http.StatusConflict, 0, "account user mismatch", "", false)
	ErrAccountUserNotFound                    = errors.NewNotFoundError(0, "account user not found", "")
	ErrAccountUserPasswordSize                = errors.NewBadRequestError(0, "user password must be between %d and %d characters", "")
	ErrAccountUserURLInvalid                  = errors.NewBadRequestError(0, "user url is not a valid url", "")
	ErrAccountUserUsernameSize                = errors.NewBadRequestError(0, "user username must be between %d and %d characters", "")
	ErrAccountUserUsernameType                = errors.NewBadRequestError(0, "user username is not a valid alphanumeric sequence", "")
	ErrApplicationAuthTokenUpdateNotAllowed   = errors.NewBadRequestError(0, "not allowed to update the application token", "")
	ErrApplicationDescriptionSize             = errors.NewBadRequestError(0, "application description must be between %d and %d characters", "")
	ErrApplicationDescriptionType             = errors.NewBadRequestError(0, "application description is not a valid alphanumeric sequence", "")
	ErrApplicationIDIsAlreadySet              = errors.NewBadRequestError(0, "application id is already set", "")
	ErrApplicationNameSize                    = errors.NewBadRequestError(0, "application name must be between %d and %d characters", "")
	ErrApplicationNameType                    = errors.NewBadRequestError(0, "application name is not a valid alphanumeric sequence", "")
	ErrApplicationNotFound                    = errors.NewNotFoundError(0, "application not found", "application not found")
	ErrApplicationUserNotActivated            = errors.NewInternalError(0, "user not activated", "")
	ErrApplicationUserNotFound                = errors.NewNotFoundError(0, "application user not found", "user not found")
	ErrApplicationUserURLInvalid              = errors.NewBadRequestError(0, "application url is not a valid url", "")
	ErrAuthMethodNotSupported                 = errors.NewBadRequestError(0, "authorization method not supported", "auth method not supported")
	ErrBadJSONReceived                        = errors.NewBadRequestError(0, "malformed json received", "")
	ErrBadUserAgent                           = errors.NewBadRequestError(0, "User-Agent header must be set (1)", "missing ua header")
	ErrConnectionAlreadyExists                = errors.NewBadRequestError(0, "connection already exists", "")
	ErrConnectionNotFound                     = errors.NewNotFoundError(0, "connection not found", "")
	ErrContentLengthInvalid                   = errors.NewBadRequestError(0, "Content-Length header is invalid", "content-length header is not an int")
	ErrContentLengthMissing                   = errors.NewBadRequestError(0, "Content-Length header missing", "missing content-length header")
	ErrContentLengthSizeMismatch              = errors.NewBadRequestError(0, "Content-Length header size mismatch", "content-length header size mismatch")
	ErrContentTypeMismatch                    = errors.NewBadRequestError(0, "Content-Type header mismatch", "content-type header mismatch")
	ErrContentTypeMissing                     = errors.NewBadRequestError(0, "Content-Type header empty", "missing content-type header")
	ErrEmailAddressInUse                      = errors.NewBadRequestError(0, "email address already in use", "")
	ErrEventIDIsAlreadySet                    = errors.NewBadRequestError(0, "event id is already set", "")
	ErrEventInvalidVisiblity                  = errors.NewBadRequestError(0, "event visibility is invalid", "")
	ErrEventMissingVisiblity                  = errors.NewBadRequestError(0, "event visibility is missing", "")
	ErrEventNearestNotInBounds                = errors.NewBadRequestError(0, "near events limits not within accepted bounds", "nearest not within bounds")
	ErrEventNotFound                          = errors.NewNotFoundError(0, "event not found", "")
	ErrGenericAuthentication                  = errors.NewBadRequestError(0, "authentication error", "")
	ErrGotBothUsernameAndEmail                = errors.NewBadRequestError(0, "both username and email are specified", "")
	ErrGotNoUsernameOrEmail                   = errors.NewBadRequestError(0, "both username and email are empty", "")
	ErrInternalAccountCreation                = errors.NewInternalError(0, "error while creating the account", "")
	ErrInternalAccountDelete                  = errors.NewInternalError(0, "error while deleting the account", "")
	ErrInternalAccountRead                    = errors.NewInternalError(0, "error while reading the account", "")
	ErrInternalAccountUpdate                  = errors.NewInternalError(0, "error while updating the account", "")
	ErrInternalAccountUserCreation            = errors.NewInternalError(0, "error while creating the account user", "")
	ErrInternalAccountUserList                = errors.NewInternalError(0, "error while retrieving list of account users", "")
	ErrInternalAccountUserRead                = errors.NewInternalError(0, "error while reading the account user", "")
	ErrInternalAccountUserSessionCreation     = errors.NewInternalError(0, "error while creating the account user session", "")
	ErrInternalAccountUserSessionDelete       = errors.NewInternalError(0, "error while deleting the account user session", "")
	ErrInternalAccountUserSessionRead         = errors.NewInternalError(0, "error while reading the account user session", "")
	ErrInternalAccountUserSessionUpdate       = errors.NewInternalError(0, "error while updating the account user session", "")
	ErrInternalAccountUserUpdate              = errors.NewInternalError(0, "error while updating the account user", "")
	ErrInternalApplicationCreation            = errors.NewInternalError(0, "error while creating the application", "")
	ErrInternalApplicationDelete              = errors.NewInternalError(0, "error while deleting the application", "")
	ErrInternalApplicationRead                = errors.NewInternalError(0, "error while reading the application", "")
	ErrInternalApplicationUpdate              = errors.NewInternalError(0, "error while updating the application", "")
	ErrInternalApplicationUserCreation        = errors.NewInternalError(0, "error while creating the application user", "")
	ErrInternalApplicationUserList            = errors.NewInternalError(0, "error while retrieving list of application users", "")
	ErrInternalApplicationUserRead            = errors.NewInternalError(0, "error while reading the application user", "")
	ErrInternalApplicationUserSessionCreation = errors.NewInternalError(0, "error while creating the application user session", "")
	ErrInternalApplicationUserSessionDelete   = errors.NewInternalError(0, "error while deleting session", "")
	ErrInternalApplicationUserSessionRead     = errors.NewInternalError(0, "error while reading the application user session", "")
	ErrInternalApplicationUserSessionsDelete  = errors.NewInternalError(0, "error while deleting all session", "")
	ErrInternalApplicationUserSessionUpdate   = errors.NewInternalError(0, "error while updating the application user session", "")
	ErrInternalApplicationUserUpdate          = errors.NewInternalError(0, "error while updating the application user", "")
	ErrInternalConnectingUsers                = errors.NewInternalError(0, "error while connecting the users", "")
	ErrInternalConnectionCreation             = errors.NewInternalError(0, "error while creating the connection", "")
	ErrInternalConnectionRead                 = errors.NewInternalError(0, "error while reading the connection", "")
	ErrInternalConnectionUpdate               = errors.NewInternalError(0, "error while updating the connection", "")
	ErrInternalEventCreation                  = errors.NewInternalError(0, "error while saving the event", "")
	ErrInternalEventRead                      = errors.NewInternalError(0, "error while reading the event", "")
	ErrInternalEventsList                     = errors.NewInternalError(0, "failed to read the events", "")
	ErrInternalEventUpdate                    = errors.NewInternalError(0, "error while updating the event", "")
	ErrInternalFollowersList                  = errors.NewInternalError(0, "error while retrieving list of followers", "")
	ErrInternalFollowingList                  = errors.NewInternalError(0, "error while retrieving list of following", "")
	ErrInternalFriendsList                    = errors.NewInternalError(0, "error while retrieving list of friends", "")
	ErrInvalidAccountCredentials              = errors.NewBadRequestError(0, "error while reading account credentials", "")
	ErrInvalidAccountUserCredentials          = errors.NewBadRequestError(0, "error while reading account user credentials", "")
	ErrInvalidAppID                           = errors.NewBadRequestError(0, "application id is not valid", "")
	ErrInvalidApplicationCredentials          = errors.NewBadRequestError(0, "error while reading application credentials", "")
	ErrInvalidApplicationUserCredentials      = errors.NewBadRequestError(0, "error while reading user credentials", "")
	ErrInvalidEmailAddress                    = errors.NewBadRequestError(0, "invalid email address", "")
	ErrInvalidEventID                         = errors.NewBadRequestError(0, "event id is not valid", "")
	ErrInvalidImageURL                        = errors.NewBadRequestError(0, "image url is not valid", "")
	ErrInvalidUserID                          = errors.NewBadRequestError(0, "user id is not valid", "")
	ErrMissingAccountInContext                = errors.NewInternalError(0, "missing account context", "")
	ErrMissingJarvisID                        = errors.NewNotFoundError(0, "not found", "request does not contain a correct Jarvis auth")
	ErrMissingLastModifiedHeader              = errors.NewInternalError(0, "something went wrong", "")
	ErrNoKnownSearchTermsSupplied             = errors.NewBadRequestError(0, "no known search terms supplied", "no known search terms supplied")
	ErrNotImplementedYet                      = errors.NewInternalError(0, "not implemented yet", "")
	ErrParseFloat                             = errors.NewBadRequestError(0, "", "parse float error")
	ErrPasswordEmpty                          = errors.NewBadRequestError(0, "password is empty", "")
	ErrPasswordMismatch                       = errors.NewBadRequestError(0, "different passwords", "")
	ErrPayloadTooBig                          = errors.New(http.StatusRequestEntityTooLarge, 0, "payload too big", "fat payload detected", false)
	ErrRadiusAndNearestMissing                = errors.NewBadRequestError(0, "you must specify either a radius or a how many nearest events you want", "invalid radius and nearest")
	ErrRadiusUnder2M                          = errors.NewBadRequestError(0, "Location radius can't be smaller than 2 meters", "radius smaller than 2")
	ErrRequestBodyEmpty                       = errors.NewBadRequestError(0, "Empty request body", "empty request body")
	ErrSelfConnectingUser                     = errors.NewBadRequestError(0, "user is connecting with itself", "self-connecting user")
	ErrSessionTokenMismatch                   = errors.NewBadRequestError(0, "session token mismatch", "")
	ErrTypeMin3Chars                          = errors.NewBadRequestError(0, "type at least 3 characters to search", "")
	ErrUserEmailAlreadyExists                 = errors.NewBadRequestError(0, "user already exists", "")
	ErrUserEmailInvalid                       = errors.NewBadRequestError(0, "user email is not valid", "")
	ErrUserFirstNameSize                      = errors.NewBadRequestError(0, "user first name must be between %d and %d characters", "")
	ErrUserFirstNameType                      = errors.NewBadRequestError(0, "user first name is not a valid alphanumeric sequence", "")
	ErrUserLastNameSize                       = errors.NewBadRequestError(0, "user last name must be between %d and %d characters", "")
	ErrUserLastNameType                       = errors.NewBadRequestError(0, "user last name is not a valid alphanumeric sequence", "")
	ErrUsernameAndEmailAreEmpty               = errors.NewBadRequestError(0, "user email and username are both empty", "")
	ErrUsernameInUse                          = errors.NewBadRequestError(0, "username already in use", "")
	ErrUsersNotConnected                      = errors.NewNotFoundError(0, "users are not connected", "")
	ErrUserURLInvalid                         = errors.NewBadRequestError(0, "user url is not a valid url", "")
	ErrUserUsernameAlreadyExists              = errors.NewBadRequestError(0, "user already exists", "")
	ErrUserUsernameSize                       = errors.NewBadRequestError(0, "user username must be between %d and %d characters", "")
	ErrUserUsernameType                       = errors.NewBadRequestError(0, "user username is not a valid alphanumeric sequence", "")
	ErrVerbSize                               = errors.NewBadRequestError(0, "verb must be between %d and %d characters", "")
	ErrVerbType                               = errors.NewBadRequestError(0, "verb is not a valid alphanumeric sequence", "")
	ErrWrongConnectionType                    = errors.NewBadRequestError(0, "unexpected connection type", "")
)
