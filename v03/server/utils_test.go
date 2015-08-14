package server_test

import (
	"fmt"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/tgflake"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/fixtures"
)

type (
	AppUserByID []*entity.ApplicationUser
)

func (s AppUserByID) Len() int {
	return len(s)
}
func (s AppUserByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s AppUserByID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// AddCorrectAccount creates a correct account
func AddCorrectOrganization(fetchAccount bool) (*entity.Organization, errors.Error) {
	account, err := coreAcc.Create(&fixtures.CorrectAccount, fetchAccount)
	if err != nil {
		return nil, err[0]
	}

	return account, nil
}

// AddCorrectAccountUser creates a correct account user
func AddCorrectMember(accountID int64, fetchUser bool) (*entity.Member, errors.Error) {
	accountUserWithAccountID := fixtures.CorrectAccountUser
	accountUserWithAccountID.OrgID = accountID
	accountUser, err := coreAccUser.Create(&accountUserWithAccountID, fetchUser)
	if err != nil {
		return nil, err[0]
	}

	return accountUser, nil
}

// AddCorrectApplication creates a correct application
func AddCorrectApplication(accountID int64, fetchApplication bool) (*entity.Application, errors.Error) {
	applicationWithAccountID := fixtures.CorrectApplication
	applicationWithAccountID.OrgID = accountID
	application, err := coreApp.Create(&applicationWithAccountID, fetchApplication)
	if err != nil {
		return nil, err[0]
	}

	return application, nil
}

// AddCorrectUser creates a correct user
func AddCorrectUser(accountID, applicationID int64, fetchUser bool) (*entity.ApplicationUser, errors.Error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Password = "password"
	user, err := coreAppUser.Create(accountID, applicationID, &userWithIDs, fetchUser)
	if err != nil {
		return nil, err[0]
	}

	return user, nil
}

// AddCorrectUser2 creates a correct user
func AddCorrectUser2(accountID, applicationID int64, fetchUser bool) (*entity.ApplicationUser, errors.Error) {
	userWithIDs := fixtures.CorrectUser
	userWithIDs.Username = "demouser2"
	userWithIDs.Password = "password"
	userWithIDs.Email = "user2@tapglue.com"
	user, err := coreAppUser.Create(accountID, applicationID, &userWithIDs, fetchUser)
	if err != nil {
		return nil, err[0]
	}

	return user, nil
}

// AddCorrectConnection creates a correct user connection
func AddCorrectConnection(accountID, applicationID int64, userFromID, userToID uint64, fetchConnection bool) (*entity.Connection, errors.Error) {
	connectionWithIDs := fixtures.CorrectConnection
	connectionWithIDs.UserFromID = userFromID
	connectionWithIDs.UserToID = userToID
	connectionWithIDs.Type = "follow"
	connection, err := coreConn.Create(accountID, applicationID, &connectionWithIDs, fetchConnection)
	if err != nil {
		return nil, err[0]
	}

	return connection, nil
}

// AddCorrectEvent creates a correct event
func AddCorrectEvent(accountID, applicationID int64, userID uint64, fetchEvent bool) (*entity.Event, errors.Error) {
	eventWithIDs := fixtures.CorrectEvent
	eventWithIDs.UserID = userID
	event, err := coreEvt.Create(accountID, applicationID, userID, &eventWithIDs, fetchEvent)
	if err != nil {
		return nil, err[0]
	}

	return event, nil
}

func UpdateUser(accountID, applicationID int64, user entity.ApplicationUser) {
	_, err := coreAppUser.Update(accountID, applicationID, user, user, false)
	if err != nil {
		panic(err[0].InternalErrorWithLocation())
	}
}

// CorrectAccount returns a correct account
func CorrectOrganization() *entity.Organization {
	account := fixtures.CorrectAccount
	return &account
}

// CorrectAccountUser returns a correct account user
func CorrectMember() *entity.Member {
	accountUser := fixtures.CorrectAccountUser
	return &accountUser
}

// CorrectUserWithDefaults returns a new user entity with prepoulated defaults
func CorrectAccountUserWithDefaults(accountID, userNumber int64) *entity.Member {
	user := CorrectMember()
	user.OrgID = accountID
	user.Username = fmt.Sprintf("acc-%d-user-%d", user.OrgID, userNumber)
	user.Email = fmt.Sprintf("acc-%d-user-%d@tapglue-test.com", user.OrgID, userNumber)
	user.Password = fmt.Sprintf("password-acc-%d-user-%d", user.OrgID, userNumber)
	user.FirstName = fmt.Sprintf("acc-%d-user-%d-first-name", user.OrgID, userNumber)
	user.LastName = fmt.Sprintf("acc-%d-user-%d-last-name", user.OrgID, userNumber)

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
	userID, err := tgflake.FlakeNextID(applicationID, "users")
	if err != nil {
		panic(err)
	}
	user := CorrectUser()
	user.ID = userID
	user.Username = fmt.Sprintf("acc-%d-app-%d-user-%d", accountID, applicationID, userNumber)
	user.Email = fmt.Sprintf("acc-%d-app-%d-user-%d@tapglue-test.com", accountID, applicationID, userNumber)
	user.Password = fmt.Sprintf("password-acc-%d-app-%d-user-%d", accountID, applicationID, userNumber)
	user.FirstName = fmt.Sprintf("acc-%d-app-%d-user-%d-first-name", accountID, applicationID, userNumber)
	user.LastName = fmt.Sprintf("acc-%d-app-%d-user-%d-last-name", accountID, applicationID, userNumber)
	user.CustomID = fmt.Sprintf("acc-%d-app-%d-user-%d-custom-id", accountID, applicationID, userNumber)
	user.SocialIDs = map[string]string{
		"facebook": fmt.Sprintf("acc-%d-app-%d-user-%d-fb", accountID, applicationID, userNumber),
		"twitter":  fmt.Sprintf("acc-%d-app-%d-user-%d-tw", accountID, applicationID, userNumber),
		"gplus":    fmt.Sprintf("acc-%d-app-%d-user-%d-gpl", accountID, applicationID, userNumber),
		"abook":    fmt.Sprintf("acc-%d-app-%d-user-%d-abk", accountID, applicationID, userNumber),
	}

	return user
}

// CorrectEvent returns a correct event
func CorrectEvent(applicationID int64) *entity.Event {
	event := &fixtures.CorrectEvent
	event.ID, _ = tgflake.FlakeNextID(applicationID, "events")
	return event
}

func AddCorrectAccounts(numberOfAccounts int) []*entity.Organization {
	var err []errors.Error
	result := make([]*entity.Organization, numberOfAccounts)
	for i := 0; i < numberOfAccounts; i++ {
		account := CorrectOrganization()
		account.Name = fmt.Sprintf("acc-%d", i+1)
		account.Description = fmt.Sprintf("acc description %d", i+1)
		result[i], err = coreAcc.Create(account, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func AddCorrectAccountUsers(account *entity.Organization, numberOfAccountUsersPerAccount int) []*entity.Member {
	var err []errors.Error
	result := make([]*entity.Member, numberOfAccountUsersPerAccount)
	for i := 0; i < numberOfAccountUsersPerAccount; i++ {
		accountUser := CorrectAccountUserWithDefaults(account.ID, int64(i+1))
		accountUser.PublicAccountID = account.PublicID
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

func LoginMember(user *entity.Member) {
	sessionToken, err := coreAccUser.CreateSession(user)
	if err != nil {
		panic(err[0].InternalErrorWithLocation())
	}
	user.SessionToken = sessionToken
	timeNow := time.Now()
	user.LastLogin = &timeNow
	_, err = coreAccUser.Update(*user, *user, false)
	if err != nil {
		panic(err[0].InternalErrorWithLocation())
	}
}

func LoginMembers(users []*entity.Member) {
	for idx := range users {
		LoginMember(users[idx])
	}
}

func AddCorrectApplications(account *entity.Organization, numberOfApplicationsPerAccount int) []*entity.Application {
	var err []errors.Error
	result := make([]*entity.Application, numberOfApplicationsPerAccount)
	for i := 0; i < numberOfApplicationsPerAccount; i++ {
		application := CorrectApplication()
		application.OrgID = account.ID
		application.PublicOrgID = account.PublicID
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
func HookUp(accountID, applicationID int64, userFromID, userToID uint64) {
	connection := &entity.Connection{
		UserFromID: userFromID,
		UserToID:   userToID,
		Type:       "follow",
	}
	coreConn.Create(accountID, applicationID, connection, false)
	coreConn.Confirm(accountID, applicationID, connection, false)
}

// HookUpUsers creates connection between all users that you provide
// bidi:true will provide bidirectional connections
func HookUpUsers(accountID, applicationID int64, users []*entity.ApplicationUser, bidi bool) {
	if len(users) < 2 {
		return
	}

	for i := 1; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			HookUp(accountID, applicationID, users[i].ID, users[j].ID)
			if bidi {
				HookUp(accountID, applicationID, users[j].ID, users[i].ID)
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
func HookUpUsersCustom(accountID, applicationID int64, users []*entity.ApplicationUser) {
	if len(users) < 2 {
		return
	}

	if len(users) >= 2 {
		HookUp(accountID, applicationID, users[0].ID, users[1].ID)
		HookUp(accountID, applicationID, users[1].ID, users[0].ID)
	}

	if len(users) >= 3 {
		HookUp(accountID, applicationID, users[0].ID, users[2].ID)
		HookUp(accountID, applicationID, users[1].ID, users[2].ID)
	}

	if len(users) >= 4 {
		HookUp(accountID, applicationID, users[0].ID, users[3].ID)
		HookUp(accountID, applicationID, users[2].ID, users[3].ID)
	}

	if len(users) >= 5 {
		connection := &entity.Connection{
			UserFromID: users[0].ID,
			UserToID:   users[4].ID,
			Type:       "follow",
		}
		coreConn.Create(accountID, applicationID, connection, false)
	}

	if len(users) >= 6 {
		connection := &entity.Connection{
			UserFromID: users[4].ID,
			UserToID:   users[5].ID,
			Type:       "follow",
		}
		coreConn.Create(accountID, applicationID, connection, false)
	}

	if len(users) >= 8 {
		connection := &entity.Connection{
			UserFromID: users[6].ID,
			UserToID:   users[7].ID,
			Type:       "follow",
		}
		coreConn.Create(accountID, applicationID, connection, false)
	}
}

func LoginApplicationUser(accountID, applicationID int64, user *entity.ApplicationUser) {
	sessionToken, err := coreAppUser.CreateSession(accountID, applicationID, user)
	if err != nil {
		panic(err)
	}
	user.SessionToken = sessionToken
	timeNow := time.Now()
	user.LastLogin = &timeNow
	user, err = coreAppUser.Update(accountID, applicationID, *user, *user, true)
	if err != nil {
		panic(err)
	}
}

func LoginUsers(accountID, applicationID int64, users []*entity.ApplicationUser) {
	for idx := range users {
		LoginApplicationUser(accountID, applicationID, users[idx])
	}
}

func AddCorrectApplicationUsers(accountID int64, application *entity.Application, numberOfUsersPerApplication int, hookUpUsers bool) []*entity.ApplicationUser {
	var err []errors.Error
	result := make([]*entity.ApplicationUser, numberOfUsersPerApplication)
	for i := 0; i < numberOfUsersPerApplication; i++ {
		user := CorrectUserWithDefaults(accountID, application.ID, int64(i+1))
		password := user.Password
		user.Activated = true
		user.Deleted = entity.PFalse
		result[i], err = coreAppUser.Create(accountID, application.ID, user, true)
		if err != nil {
			panic(err[0].InternalErrorWithLocation())
		}
		result[i].OriginalPassword = password
	}

	if hookUpUsers {
		HookUpUsersCustom(accountID, application.ID, result)
	}

	return result
}

// AddCorrectUserEvents adds correct events to a user
// If numberOfEventsPerUser < 4 then events are common, else they are user specific (thus unique)
func AddCorrectUserEvents(accountID, applicationID int64, user *entity.ApplicationUser, numberOfEventsPerUser int) []*entity.Event {
	var err []errors.Error
	locations := []struct {
		Lat   float64
		Lon   float64
		Label string
	}{
		{Lat: 52.517220, Lon: 13.304817, Label: "dlsniper"},
		{Lat: 52.515931, Lon: 13.301226, Label: "gas"},
		{Lat: 52.520910, Lon: 13.295661, Label: "palace"},
		{Lat: 52.507686, Lon: 13.302404, Label: "ziko"},
		{Lat: 52.510017, Lon: 13.373451, Label: "cinestar"},
		{Lat: 52.517334, Lon: 13.389419, Label: "mercedes"},
		{Lat: 52.515585, Lon: 13.287820, Label: "onur"},
	}

	result := make([]*entity.Event, numberOfEventsPerUser)
	for i := 0; i < numberOfEventsPerUser; i++ {
		event := CorrectEvent(applicationID)
		event.Visibility = uint8(i%4*10 + 10)
		event.UserID = user.ID
		if event.Visibility == entity.EventPublic {
			event.Location = fmt.Sprintf("location-all-%d", i+1)
			event.Target = &entity.Object{
				ID:           fmt.Sprintf("target-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("target-%d-all", i+1)},
			}
			event.Object = &entity.Object{
				ID:           fmt.Sprintf("object-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("object-%d-all", i+1)},
			}
		} else if event.Visibility == entity.EventGlobal {
			event.Location = fmt.Sprintf("location-global-%d", i+1)
			event.Target = &entity.Object{
				ID:           fmt.Sprintf("target-global-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("target-global-%d-all", i+1)},
			}
			event.Object = &entity.Object{
				ID:           fmt.Sprintf("object-global-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("object-global-%d-all", i+1)},
			}
		} else {
			event.Location = fmt.Sprintf("location-%d", i+1)
			event.Target = &entity.Object{
				ID:           fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d", accountID, applicationID, user.ID, i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d-lall", accountID, applicationID, user.ID, i+1)},
			}
			event.Object = &entity.Object{
				ID:           fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d", accountID, applicationID, user.ID, i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d-lall", accountID, applicationID, user.ID, i+1)},
			}
		}

		// Some locations are more special then others (easier things to debug when it comes to the location search by geo coordinates)
		if i < len(locations) {
			event.Latitude = locations[i].Lat
			event.Longitude = locations[i].Lon
			event.Location = locations[i].Label
		}

		event.Images = map[string]*entity.Image{}
		event.Images["thumb_pic"] = &entity.Image{}
		event.Images["thumb_pic"].URL = "https://www.tapglue.com/img/box/newsfeed.jpg"
		event.Images["thumb_pic"].Width = 200
		event.Images["thumb_pic"].Heigth = 200
		event.Images["thumb_pic"].Type = "image/jpeg"

		event.Images["original_pic"] = &entity.Image{}
		event.Images["original_pic"].URL = "https://www.tapglue.com/img/original.jpg"
		event.Images["original_pic"].Width = 2048
		event.Images["original_pic"].Heigth = 2048
		event.Images["original_pic"].Type = "image/jpeg"

		result[i], err = coreEvt.Create(accountID, applicationID, user.ID, event, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func CorrectDeploy(
	numberOfOrganizations, numberOfMembersPerOrganization,
	numberOfApplicationsPerOrganization,
	numberOfUsersPerApplication, numberOfEventsPerUser int,
	hookUpUsers, loginUsers bool) []*entity.Organization {

	organizations := AddCorrectAccounts(numberOfOrganizations)

	for i := 0; i < numberOfOrganizations; i++ {
		organizations[i].Members = AddCorrectAccountUsers(organizations[i], numberOfMembersPerOrganization)
		if loginUsers {
			LoginMembers(organizations[i].Members)
		}

		organizations[i].Applications = AddCorrectApplications(organizations[i], numberOfApplicationsPerOrganization)

		for j := 0; j < numberOfApplicationsPerOrganization; j++ {
			organizations[i].Applications[j].Users = AddCorrectApplicationUsers(organizations[i].ID, organizations[i].Applications[j], numberOfUsersPerApplication, hookUpUsers)
			if loginUsers {
				LoginUsers(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users)
			}

			for k := 0; k < numberOfUsersPerApplication; k++ {
				organizations[i].Applications[j].Users[k].Events = AddCorrectUserEvents(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users[k], numberOfEventsPerUser)
			}
		}
	}

	return organizations
}

func testBootup(conf *config.Postgres) {
	db := v03PostgresClient.MainDatastore()

	existingSchemas, err := db.Query(`SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname ILIKE 'app_%_%'`)
	if err != nil {
		panic(err)
	}
	defer existingSchemas.Close()
	for existingSchemas.Next() {
		schemaName := ""
		err := existingSchemas.Scan(&schemaName)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(fmt.Sprintf(`DROP SCHEMA %s CASCADE`, schemaName))
		if err != nil {
			panic(err)
		}
	}

	queries := []string{
		`TRUNCATE TABLE tg.accounts RESTART IDENTITY`,
		`TRUNCATE TABLE tg.account_users RESTART IDENTITY`,
		`TRUNCATE TABLE tg.applications RESTART IDENTITY`,
		`TRUNCATE TABLE tg.account_user_sessions`,
	}

	for idx := range queries {
		_, err := db.Exec(queries[idx])
		if err != nil {
			panic(err)
		}
	}

	tgflake.RemoveAllFlakes()
}
