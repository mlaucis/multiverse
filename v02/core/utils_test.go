/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/v02/entity"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount(fetchAccount bool) (acc *entity.Account, err error) {
	account, err := CreateAccount(correctAccount, fetchAccount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser(accountID int64, fetchUser bool) (accUsr *entity.AccountUser, err error) {
	accountUserWithAccountID := correctAccountUser
	accountUserWithAccountID.AccountID = accountID
	accountUser, err := WriteAccountUser(accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// AddCorrectApplication creates a correct application
func AddCorrectApplication(accountID int64, fetchApplication bool) (app *entity.Application, err error) {
	applicationWithAccountID := correctApplication
	applicationWithAccountID.AccountID = accountID
	application, err := WriteApplication(applicationWithAccountID, fetchApplication)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// AddCorrectUser creates a correct user
func AddCorrectUser(accountID, applicationID int64, fetchUser bool) (usr *entity.ApplicationUser, err error) {
	userWithIDs := correctUser
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := WriteUser(userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectConnection creates a correct user connection
func AddCorrectConnection(accountID, applicationID, userFromID, userToID int64, fetchConnection bool) (con *entity.Connection, err error) {
	connectionWithIDs := correctConnection
	connectionWithIDs.AccountID = accountID
	connectionWithIDs.ApplicationID = applicationID
	connectionWithIDs.UserFromID = userFromID
	connectionWithIDs.UserToID = userToID
	connection, err := WriteConnection(connectionWithIDs, fetchConnection)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// AddCorrectEvent creates a correct event
func AddCorrectEvent(accountID, applicationID, userID int64, fetchEvent bool) (evn *entity.Event, err error) {
	eventWithIDs := correctEvent
	eventWithIDs.AccountID = accountID
	eventWithIDs.ApplicationID = applicationID
	eventWithIDs.UserID = userID
	event, err := WriteEvent(eventWithIDs, fetchEvent)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// CorrectAccount returns a correct account
func CorrectAccount() *entity.Account {
	return correctAccount
}

// CorrectAccountUser returns a correct account user
func CorrectAccountUser() *entity.AccountUser {
	return correctAccountUser
}

// CorrectApplication returns a correct application
func CorrectApplication() *entity.Application {
	return correctApplication
}

// CorrectUser returns a correct user
func CorrectUser() *entity.ApplicationUser {
	return correctUser
}

// CorrectEvent returns a correct event
func CorrectEvent() *entity.Event {
	return correctEvent
}
