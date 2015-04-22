/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"fmt"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/fixtures"
)

// AddCorrectAccount creates a correct account
func AddCorrectAccount(fetchAccount bool) (acc *entity.Account, err errors.Error) {
	account, err := coreAcc.Create(&fixtures.CorrectAccount, fetchAccount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectAccountUser(accountID int64, fetchUser bool) (accUsr *entity.AccountUser, err errors.Error) {
	accountUserWithAccountID := fixtures.CorrectAccountUser
	accountUserWithAccountID.AccountID = accountID
	accountUser, err := coreAccUser.Create(&accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err
	}

	return accountUser, nil
}

// AddCorrectApplication creates a correct application
func AddCorrectApplication(accountID int64, fetchApplication bool) (*entity.Application, errors.Error) {
	applicationWithAccountID := fixtures.CorrectApplication
	applicationWithAccountID.AccountID = accountID
	application, err := coreApp.Create(&applicationWithAccountID, fetchApplication)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// AddCorrectUser creates a correct user
func AddCorrectUser(accountID, applicationID int64, fetchUser bool) (usr *entity.ApplicationUser, err errors.Error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Password = "password"
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := coreAppUser.Create(&userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectUser2 creates a correct user
func AddCorrectUser2(accountID, applicationID int64, fetchUser bool) (usr *entity.ApplicationUser, err errors.Error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Username = "demouser2"
	userWithIDs.Password = "password"
	userWithIDs.Email = "user2@tapglue.com"
	userWithIDs.AccountID = accountID
	userWithIDs.ApplicationID = applicationID
	user, err := coreAppUser.Create(&userWithIDs, fetchUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddCorrectConnection creates a correct user connection
func AddCorrectConnection(accountID, applicationID, userFromID, userToID int64, fetchConnection bool) (con *entity.Connection, err errors.Error) {
	connectionWithIDs := fixtures.CorrectConnection
	connectionWithIDs.AccountID = accountID
	connectionWithIDs.ApplicationID = applicationID
	connectionWithIDs.UserFromID = userFromID
	connectionWithIDs.UserToID = userToID
	connection, err := coreConn.Create(&connectionWithIDs, fetchConnection)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// AddCorrectEvent creates a correct event
func AddCorrectEvent(accountID, applicationID, userID int64, fetchEvent bool) (evn *entity.Event, err errors.Error) {
	eventWithIDs := fixtures.CorrectEvent
	eventWithIDs.AccountID = accountID
	eventWithIDs.ApplicationID = applicationID
	eventWithIDs.UserID = userID
	event, err := coreEvt.Create(&eventWithIDs, fetchEvent)
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

// CorrectUserWithDefaults returns a new user entity with prepoulated defaults
func CorrectAccountUserWithDefaults(accountID, userNumber int64) *entity.AccountUser {
	user := CorrectAccountUser()
	user.AccountID = accountID
	user.Username = fmt.Sprintf("acc-%d-user-%d", user.AccountID, userNumber)
	user.Email = fmt.Sprintf("acc-%d-user-%d@tapglue-test.com", user.AccountID, userNumber)
	user.Password = fmt.Sprintf("acc-%d-user-%d", user.AccountID, userNumber)
	user.FirstName = fmt.Sprintf("acc-%d-user-%d-first-name", user.AccountID, userNumber)
	user.LastName = fmt.Sprintf("acc-%d-user-%d-last-name", user.AccountID, userNumber)

	return user
}

// CorrectApplication returns a correct application
func CorrectApplication() *entity.Application {
	application := fixtures.CorrectApplication
	return &application
}

// CorrectUser returns a correct user
func CorrectUser() *entity.ApplicationUser {
	applicationUser := fixtures.CorrectUser
	return &applicationUser
}

// CorrectUserWithDefaults returns a new user entity with prepoulated defaults
func CorrectUserWithDefaults(accountID, applicationID, userNumber int64) *entity.ApplicationUser {
	user := CorrectUser()
	user.AccountID = accountID
	user.ApplicationID = applicationID
	user.Username = fmt.Sprintf("acc-%d-app-%d-user-%d", user.AccountID, user.ApplicationID, userNumber)
	user.Email = fmt.Sprintf("acc-%d-app-%d-user-%d@tapglue-test.com", user.AccountID, user.ApplicationID, userNumber)
	user.Password = fmt.Sprintf("acc-%d-app-%d-user-%d", user.AccountID, user.ApplicationID, userNumber)
	user.FirstName = fmt.Sprintf("acc-%d-app-%d-user-%d-first-name", user.AccountID, user.ApplicationID, userNumber)
	user.LastName = fmt.Sprintf("acc-%d-app-%d-user-%d-last-name", user.AccountID, user.ApplicationID, userNumber)
	user.SocialIDs = map[string]string{
		"facebook": fmt.Sprintf("acc-%d-app-%d-user-%d-fb", user.AccountID, user.ApplicationID, userNumber),
		"twitter":  fmt.Sprintf("acc-%d-app-%d-user-%d-tw", user.AccountID, user.ApplicationID, userNumber),
		"gplus":    fmt.Sprintf("acc-%d-app-%d-user-%d-gpl", user.AccountID, user.ApplicationID, userNumber),
		"abook":    fmt.Sprintf("acc-%d-app-%d-user-%d-abk", user.AccountID, user.ApplicationID, userNumber),
	}

	return user
}

// CorrectEvent returns a correct event
func CorrectEvent() *entity.Event {
	event := &fixtures.CorrectEvent
	return event
}

func AddCorrectAccounts(numberOfAccounts int) []*entity.Account {
	var err errors.Error
	result := make([]*entity.Account, numberOfAccounts)
	for i := 0; i < numberOfAccounts; i++ {
		account := CorrectAccount()
		account.Name = fmt.Sprintf("acc-%d", i+1)
		account.Description = fmt.Sprintf("acc description %d", i+1)
		result[i], err = coreAcc.Create(account, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func AddCorrectAccountUsers(account *entity.Account, numberOfAccountUsersPerAccount int) []*entity.AccountUser {
	var err errors.Error
	result := make([]*entity.AccountUser, numberOfAccountUsersPerAccount)
	for i := 0; i < numberOfAccountUsersPerAccount; i++ {
		accountUser := CorrectAccountUserWithDefaults(account.ID, int64(i+1))
		password := accountUser.Password
		accountUser.Activated = true
		result[i], err = coreAccUser.Create(accountUser, true)
		result[i].OriginalPassword = password
		if err != nil {
			panic(err)
		}
	}

	return result
}

func LoginAccountUsers(users []*entity.AccountUser) {
	for idx := range users {
		sessionToken, err := coreAccUser.CreateSession(users[idx])
		if err != nil {
			panic(err)
		}
		users[idx].SessionToken = sessionToken
		users[idx].LastLogin = time.Now()
		_, err = coreAccUser.Update(*users[idx], *users[idx], false)
		if err != nil {
			panic(err)
		}
	}
}

func AddCorrectApplications(account *entity.Account, numberOfApplicationsPerAccount int) []*entity.Application {
	var err errors.Error
	result := make([]*entity.Application, numberOfApplicationsPerAccount)
	for i := 0; i < numberOfApplicationsPerAccount; i++ {
		application := CorrectApplication()
		application.AccountID = account.ID
		application.Name = fmt.Sprintf("acc-%d-app-%d", account.ID, i+1)
		application.Description = fmt.Sprintf("acc %d app %d", account.ID, i+1)
		result[i], err = coreApp.Create(application, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

// HookUp create a connection between two users provided
func HookUp(accountID, applicationID, userFromID, userToID int64) {
	connection := &entity.Connection{
		AccountID:     accountID,
		ApplicationID: applicationID,
		UserFromID:    userFromID,
		UserToID:      userToID,
	}
	coreConn.Create(connection, false)
	coreConn.Confirm(connection, false)
}

// HookUpUsers creates connection between all users that you provide
// bidi:true will provide bidirectional connections
func HookUpUsers(users []*entity.ApplicationUser, bidi bool) {
	if len(users) < 2 {
		return
	}

	for i := 1; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			HookUp(users[0].AccountID, users[0].ApplicationID, users[i].ID, users[j].ID)
			if bidi {
				HookUp(users[0].AccountID, users[0].ApplicationID, users[j].ID, users[i].ID)
			}
		}
	}
}

// HookUpUsersCustom creates a custom connection web between supplied users based on the number of
// users supplied. If the number is greater than 8, all the users > 8 are not connected in any way
// The connection table is defined below:
// 1->2 1->3 1->4 1->5
// 2->1 2->3
// 3->4
// 5->6
// 7->8
func HookUpUsersCustom(users []*entity.ApplicationUser) {
	if len(users) < 2 {
		return
	}

	if len(users) >= 2 {
		HookUp(users[0].AccountID, users[0].ApplicationID, users[0].ID, users[1].ID)
		HookUp(users[0].AccountID, users[0].ApplicationID, users[1].ID, users[0].ID)
	}

	if len(users) >= 3 {
		HookUp(users[0].AccountID, users[0].ApplicationID, users[0].ID, users[2].ID)
		HookUp(users[0].AccountID, users[0].ApplicationID, users[1].ID, users[2].ID)
	}

	if len(users) >= 4 {
		HookUp(users[0].AccountID, users[0].ApplicationID, users[0].ID, users[3].ID)
		HookUp(users[0].AccountID, users[0].ApplicationID, users[2].ID, users[3].ID)
	}

	if len(users) >= 5 {
		connection := &entity.Connection{
			AccountID:     users[0].AccountID,
			ApplicationID: users[0].ApplicationID,
			UserFromID:    users[0].ID,
			UserToID:      users[4].ID,
		}
		coreConn.Create(connection, false)
	}

	if len(users) >= 6 {
		connection := &entity.Connection{
			AccountID:     users[0].AccountID,
			ApplicationID: users[0].ApplicationID,
			UserFromID:    users[4].ID,
			UserToID:      users[5].ID,
		}
		coreConn.Create(connection, false)
	}

	if len(users) >= 8 {
		connection := &entity.Connection{
			AccountID:     users[0].AccountID,
			ApplicationID: users[0].ApplicationID,
			UserFromID:    users[6].ID,
			UserToID:      users[7].ID,
		}
		coreConn.Create(connection, false)
	}
}

func LoginUsers(users []*entity.ApplicationUser) {
	for idx := range users {
		sessionToken, err := coreAppUser.CreateSession(users[idx])
		if err != nil {
			panic(err)
		}
		users[idx].SessionToken = sessionToken
		users[idx].LastLogin = time.Now()
		_, err = coreAppUser.Update(*users[idx], *users[idx], false)
		if err != nil {
			panic(err)
		}
	}
}

func AddCorrectApplicationUsers(application *entity.Application, numberOfUsersPerApplication int, hookUpUsers bool) []*entity.ApplicationUser {
	var err errors.Error
	result := make([]*entity.ApplicationUser, numberOfUsersPerApplication)
	for i := 0; i < numberOfUsersPerApplication; i++ {
		user := CorrectUserWithDefaults(application.AccountID, application.ID, int64(i+1))
		password := user.Password
		user.Activated = true
		result[i], err = coreAppUser.Create(user, true)
		result[i].OriginalPassword = password
		if err != nil {
			panic(err)
		}
	}

	if hookUpUsers {
		HookUpUsersCustom(result)
	}

	return result
}

// AddCorrectUserEvents adds correct events to a user
// If numberOfEventsPerUser < 4 then events are common, else they are user specific (thus unique)
func AddCorrectUserEvents(user *entity.ApplicationUser, numberOfEventsPerUser int) []*entity.Event {
	var err errors.Error
	locations := []struct {
		Lat   float64
		Lon   float64
		Label string
	}{
		{Lat: 52.5169257, Lon: 13.3065105, Label: "dlsniper"},
		{Lat: 52.5148045, Lon: 13.3000390, Label: "gas"},
		{Lat: 52.5177294, Lon: 13.2938378, Label: "palace"},
		{Lat: 52.5104167, Lon: 13.3003824, Label: "ziko"},
		{Lat: 52.5120818, Lon: 13.3762879, Label: "cinestar"},
		{Lat: 52.5157576, Lon: 13.3873319, Label: "mercedes"},
	}

	result := make([]*entity.Event, numberOfEventsPerUser)
	for i := 0; i < numberOfEventsPerUser; i++ {
		event := CorrectEvent()
		event.AccountID = user.AccountID
		event.ApplicationID = user.ApplicationID
		event.UserID = user.ID
		if i < 4 {
			event.Location = fmt.Sprintf("location-all-%d", i+1)
			event.Target = &entity.Object{
				ID:          fmt.Sprintf("target-%d", i+1),
				DisplayName: map[string]string{"all": fmt.Sprintf("target-%d-all", i+1)},
			}
			event.Object = &entity.Object{
				ID:          fmt.Sprintf("object-%d", i+1),
				DisplayName: map[string]string{"all": fmt.Sprintf("object-%d-all", i+1)},
			}
		} else {
			event.Location = fmt.Sprintf("location-%d", i+1)
			event.Target = &entity.Object{
				ID:          fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d", user.AccountID, user.ApplicationID, user.ID, i+1),
				DisplayName: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d-lall", user.AccountID, user.ApplicationID, user.ID, i+1)},
			}
			event.Object = &entity.Object{
				ID:          fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d", user.AccountID, user.ApplicationID, user.ID, i+1),
				DisplayName: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d-lall", user.AccountID, user.ApplicationID, user.ID, i+1)},
			}
		}
		if i < 6 {
			event.Latitude = locations[i].Lat
			event.Longitude = locations[i].Lon
			event.Location = locations[i].Label
		}
		result[i], err = coreEvt.Create(event, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func CorrectDeploy(numberOfAccounts, numberOfAccountUsersPerAccount, numberOfApplicationsPerAccount, numberOfUsersPerApplication, numberOfEventsPerUser int, hookUpUsers, loginUsers bool) []*entity.Account {
	accounts := AddCorrectAccounts(numberOfAccounts)

	for i := 0; i < numberOfAccounts; i++ {
		accounts[i].Users = AddCorrectAccountUsers(accounts[i], numberOfAccountUsersPerAccount)
		if loginUsers {
			LoginAccountUsers(accounts[i].Users)
		}

		accounts[i].Applications = AddCorrectApplications(accounts[i], numberOfApplicationsPerAccount)

		for j := 0; j < numberOfApplicationsPerAccount; j++ {
			accounts[i].Applications[j].Users = AddCorrectApplicationUsers(accounts[i].Applications[j], numberOfUsersPerApplication, hookUpUsers)
			if loginUsers {
				LoginUsers(accounts[i].Applications[j].Users)
			}

			for k := 0; k < numberOfUsersPerApplication; k++ {
				accounts[i].Applications[j].Users[k].Events = AddCorrectUserEvents(accounts[i].Applications[j].Users[k], numberOfEventsPerUser)
			}
		}
	}

	return accounts
}
