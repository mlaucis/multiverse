/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount(fetchAccount bool) (acc *entity.Account, err error) {
	account, err := core.WriteAccount(correctAccount, fetchAccount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser(accountID int64, fetchUser bool) (accUsr *entity.AccountUser, err error) {
	accountUserWithAccountID := correctAccountUser
	accountUserWithAccountID.AccountID = accountID
	accountUser, err := core.WriteAccountUser(accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// AddCorrectApplication creates a correct application
func AddCorrectApplication(accountID int64, fetchApplication bool) (app *entity.Application, err error) {
	applicationWithAccountID := correctApplication
	applicationWithAccountID.AccountID = accountID
	application, err := core.WriteApplication(applicationWithAccountID, fetchApplication)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// AddCorrectUser creates a correct user
func AddCorrectUser(accountID, applicationID int64, fetchUser bool) (usr *entity.User, err error) {
	correctUser.Password = "password"
	userWithIDs := correctUser
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := core.WriteUser(userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectUser2 creates a correct user
func AddCorrectUser2(accountID, applicationID int64, fetchUser bool) (usr *entity.User, err error) {
	correctUser.Password = "password"
	correctUser.Email = "user2@tapglue.com"
	userWithIDs := correctUser
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := core.WriteUser(userWithIDs, fetchUser)
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
	connection, err := core.WriteConnection(connectionWithIDs, fetchConnection)
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
	event, err := core.WriteEvent(eventWithIDs, fetchEvent)
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
func CorrectUser() *entity.User {
	return correctUser
}

// CorrectEvent returns a correct event
func CorrectEvent() *entity.Event {
	return correctEvent
}
