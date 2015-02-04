package core

import (
	"time"

	"github.com/tapglue/backend/core/entity"
)

// Create a correct account
func AddCorrectAccount() *entity.Account {
	savedAccount, err := WriteAccount(correctAccount, true)
	if err != nil {
		panic(err)
	}

	return savedAccount
}

// Create a correct account user
func AddCorrectAccountUser() *entity.AccountUser {
	savedAccountUser, err := WriteAccountUser(correctAccountUser, false)
	if err != nil {
		panic(err)
	}

	return savedAccountUser
}

// Create correct account users
func AddCorrectAccountUsers() (accountUser1, accountUser2 *entity.AccountUser) {
	savedAccountUser1, err := WriteAccountUser(correctAccountUser, false)
	if err != nil {
		panic(err)
	}
	savedAccountUser2, err := WriteAccountUser(correctAccountUser, false)
	if err != nil {
		panic(err)
	}

	return savedAccountUser1, savedAccountUser2
}

// Create a correct applicaton
func AddCorrectApplication() *entity.Application {
	savedApplication, err := WriteApplication(correctApplication, false)
	if err != nil {
		panic(err)
	}

	return savedApplication
}

// Create correct applicatons
func AddCorrectApplications() (application1, application2 *entity.Application) {
	savedApplication1, err := WriteApplication(correctApplication, false)
	if err != nil {
		panic(err)
	}
	savedApplication2, err := WriteApplication(correctApplication, false)
	if err != nil {
		panic(err)
	}

	return savedApplication1, savedApplication2
}

// Create a correct user
func AddCorrectUser() *entity.User {
	correctUser.ID = time.Now().UTC().UnixNano()
	savedUser, err := WriteUser(correctUser, false)
	if err != nil {
		panic(err)
	}

	return savedUser
}

// Create correct users
func AddCorrectApplicationUsers() (user1, user2 *entity.User) {
	correctUser.ID = time.Now().UTC().UnixNano()
	savedUser1, err := WriteUser(correctUser, false)
	if err != nil {
		panic(err)
	}
	correctUser.ID = time.Now().UTC().UnixNano()
	savedUser2, err := WriteUser(correctUser, false)
	if err != nil {
		panic(err)
	}

	return savedUser1, savedUser2
}

// Create a correct event
func AddCorrectEvent() *entity.Event {
	savedEvent, err := WriteEvent(correctEvent, false)
	if err != nil {
		panic(err)
	}

	return savedEvent
}

// Create correct events
func AddCorrectEvents() (event1, event2 *entity.Event) {
	savedEvent1, err := WriteEvent(correctEvent, false)
	if err != nil {
		panic(err)
	}
	savedEvent2, err := WriteEvent(correctEvent, false)
	if err != nil {
		panic(err)
	}

	return savedEvent1, savedEvent2
}
