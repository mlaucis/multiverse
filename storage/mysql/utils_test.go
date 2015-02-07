package mysql

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

// Create a correct account
func AddCorrectAccount() *entity.Account {
	savedAccount, err := AddAccount(correctAccount)
	if err != nil {
		panic(err)
	}

	return savedAccount
}

/*
// Create a correct account user
func AddCorrectAccountUser() *entity.AccountUser {
	savedAccountUser, err := AddAccountUser(AddCorrectAccount().ID, correctAccountUser)
	if err != nil {
		panic(err)
	}

	return savedAccountUser
}

// Create correct account users
func AddCorrectAccountUsers() (accountUser1, accountUser2 *entity.AccountUser) {
	savedAccount := AddCorrectAccount()
	savedAccountUser1, err := AddAccountUser(savedAccount.ID, correctAccountUser)
	if err != nil {
		panic(err)
	}
	savedAccountUser2, err := AddAccountUser(savedAccount.ID, correctAccountUser)
	if err != nil {
		panic(err)
	}

	return savedAccountUser1, savedAccountUser2
}

// Create a correct applicaton
func AddCorrectAccountApplication() *entity.Application {
	savedApplication, err := AddAccountApplication(AddCorrectAccount().ID, correctApplication)
	if err != nil {
		panic(err)
	}

	return savedApplication
}

// Create correct applicatons
func AddCorrectAccountApplications() (application1, application2 *entity.Application) {
	savedAccount := AddCorrectAccount()
	savedApplication1, err := AddAccountApplication(savedAccount.ID, correctApplication)
	if err != nil {
		panic(err)
	}
	savedApplication2, err := AddAccountApplication(savedAccount.ID, correctApplication)
	if err != nil {
		panic(err)
	}

	return savedApplication1, savedApplication2
}

// Create a correct user
func AddCorrectApplicationUser() *entity.User {
	correctUser.Token = RandomToken()
	savedUser, err := AddApplicationUser(AddCorrectAccountApplication().ID, correctUser)
	if err != nil {
		panic(err)
	}

	return savedUser
}

// Create correct users
func AddCorrectApplicationUsers() (user1, user2 *entity.User) {
	savedApplication := AddCorrectAccountApplication()
	correctUser.Token = RandomToken()
	savedUser1, err := AddApplicationUser(savedApplication.ID, correctUser)
	if err != nil {
		panic(err)
	}
	correctUser.Token = RandomToken()
	savedUser2, err := AddApplicationUser(savedApplication.ID, correctUser)
	if err != nil {
		panic(err)
	}

	return savedUser1, savedUser2
}


// Create a correct event
func AddCorrectEvent() *entity.Event {
	savedSession := AddCorrectUserSession()
	UpdateEvent(savedSession.applicationId, savedSession.ID, savedSession.UserToken)
	savedEvent, err := AddSessionEvent(correctEvent)
	if err != nil {
		panic(err)
	}

	return savedEvent
}

// Create correct events
func AddCorrectEvents() (event1, event2 *entity.Event) {
	savedSession := AddCorrectUserSession()
	UpdateEvent(savedSession.applicationId, savedSession.ID, savedSession.UserToken)
	savedEvent1, err := AddSessionEvent(correctEvent)
	if err != nil {
		panic(err)
	}
	savedEvent2, err := AddSessionEvent(correctEvent)
	if err != nil {
		panic(err)
	}

	return savedEvent1, savedEvent2
}
*/
// RandomToken returns a random Token
func RandomToken() string {
	// String length
	size := 32

	// Create random string
	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	// Encode to base64 string
	rs := base64.URLEncoding.EncodeToString(rb)

	// Return string
	return rs
}

/*
// UpdateEvent updates correctEvent struct
func UpdateEvent(applicationId, sessionID uint64, token string) {
	correctEvent.applicationId = applicationId
	correctEvent.SessionID = sessionID
	correctEvent.UserToken = token
}
*/
