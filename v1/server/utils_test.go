/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"github.com/tapglue/backend/v1/core"
	"github.com/tapglue/backend/v1/entity"
	"github.com/tapglue/backend/v1/fixtures"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount(fetchAccount bool) (acc *entity.Account, err error) {
	account, err := core.WriteAccount(&fixtures.CorrectAccount, fetchAccount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser(accountID int64, fetchUser bool) (accUsr *entity.AccountUser, err error) {
	accountUserWithAccountID := fixtures.CorrectAccountUser
	accountUserWithAccountID.AccountID = accountID
	accountUser, err := core.WriteAccountUser(&accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// AddCorrectApplication creates a correct application
func AddCorrectApplication(accountID int64, fetchApplication bool) (app *entity.Application, err error) {
	applicationWithAccountID := fixtures.CorrectApplication
	applicationWithAccountID.AccountID = accountID
	application, err := core.WriteApplication(&applicationWithAccountID, fetchApplication)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// AddCorrectUser creates a correct user
func AddCorrectUser(accountID, applicationID int64, fetchUser bool) (usr *entity.User, err error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Password = "password"
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := core.WriteUser(&userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectUser2 creates a correct user
func AddCorrectUser2(accountID, applicationID int64, fetchUser bool) (usr *entity.User, err error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Username = "demouser2"
	userWithIDs.Password = "password"
	userWithIDs.Email = "user2@tapglue.com"
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := core.WriteUser(&userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectConnection creates a correct user connection
func AddCorrectConnection(accountID, applicationID, userFromID, userToID int64, fetchConnection bool) (con *entity.Connection, err error) {
	connectionWithIDs := fixtures.CorrectConnection
	connectionWithIDs.AccountID = accountID
	connectionWithIDs.ApplicationID = applicationID
	connectionWithIDs.UserFromID = userFromID
	connectionWithIDs.UserToID = userToID
	connection, err := core.WriteConnection(&connectionWithIDs, fetchConnection)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// AddCorrectEvent creates a correct event
func AddCorrectEvent(accountID, applicationID, userID int64, fetchEvent bool) (evn *entity.Event, err error) {
	eventWithIDs := fixtures.CorrectEvent
	eventWithIDs.AccountID = accountID
	eventWithIDs.ApplicationID = applicationID
	eventWithIDs.UserID = userID
	event, err := core.WriteEvent(&eventWithIDs, fetchEvent)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// CorrectAccount returns a correct account
func CorrectAccount() *entity.Account {
	account := fixtures.CorrectAccount
	return &account
}

// CorrectAccountUser returns a correct account user
func CorrectAccountUser() *entity.AccountUser {
	accountUser := fixtures.CorrectAccountUser
	return &accountUser
}

// CorrectApplication returns a correct application
func CorrectApplication() *entity.Application {
	application := fixtures.CorrectApplication
	return &application
}

// CorrectUser returns a correct user
func CorrectUser() *entity.User {
	applicationUser := fixtures.CorrectUser
	return &applicationUser
}

// CorrectEvent returns a correct event
func CorrectEvent() *entity.Event {
	event := &fixtures.CorrectEvent
	return event
}
