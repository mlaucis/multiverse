package server_test

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/tgflake"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/fixtures"

	. "gopkg.in/check.v1"
)

type AppUserByID []*entity.ApplicationUser

func (s AppUserByID) Len() int {
	return len(s)
}
func (s AppUserByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s AppUserByID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// AddCorrectOrganization creates a correct organization
func AddCorrectOrganization(fetchOrganization bool) (*entity.Organization, errors.Error) {
	organization, err := coreAcc.Create(&fixtures.CorrectOrganization, fetchOrganization)
	if err != nil {
		return nil, err[0]
	}

	return organization, nil
}

// AddCorrectMember creates a correct member
func AddCorrectMember(orgID int64, fetchUser bool) (*entity.Member, errors.Error) {
	memberWithOrgID := fixtures.CorrectMember
	memberWithOrgID.OrgID = orgID
	member, err := coreAccUser.Create(&memberWithOrgID, fetchUser)
	if err != nil {
		return nil, err[0]
	}

	return member, nil
}

func UpdateUser(orgID, applicationID int64, user entity.ApplicationUser) {
	_, err := coreAppUser.Update(orgID, applicationID, user, user, false)
	if err != nil {
		panic(err[0].InternalErrorWithLocation())
	}
}

// CorrectOrganization returns a correct organization
func CorrectOrganization() *entity.Organization {
	organization := fixtures.CorrectOrganization
	return &organization
}

// CorrectMember returns a correct member
func CorrectMember() *entity.Member {
	member := fixtures.CorrectMember
	return &member
}

// CorrectMemberWithDefaults returns a new user entity with prepoulated defaults
func CorrectMemberWithDefaults(orgID, userNumber int64) *entity.Member {
	user := CorrectMember()
	user.OrgID = orgID
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
func CorrectUserWithDefaults(orgID, applicationID, userNumber int64) *entity.ApplicationUser {
	userID, err := tgflake.FlakeNextID(applicationID, "users")
	if err != nil {
		panic(err)
	}
	user := CorrectUser()
	user.ID = userID
	user.Username = fmt.Sprintf("acc-%d-app-%d-user-%d", orgID, applicationID, userNumber)
	user.Email = fmt.Sprintf("acc-%d-app-%d-user-%d@tapglue-test.com", orgID, applicationID, userNumber)
	user.Password = fmt.Sprintf("password-acc-%d-app-%d-user-%d", orgID, applicationID, userNumber)
	user.FirstName = fmt.Sprintf("acc-%d-app-%d-user-%d-first-name", orgID, applicationID, userNumber)
	user.LastName = fmt.Sprintf("acc-%d-app-%d-user-%d-last-name", orgID, applicationID, userNumber)
	user.CustomID = fmt.Sprintf("acc-%d-app-%d-user-%d-custom-id", orgID, applicationID, userNumber)
	user.SocialIDs = map[string]string{
		"facebook": fmt.Sprintf("acc-%d-app-%d-user-%d-fb", orgID, applicationID, userNumber),
		"twitter":  fmt.Sprintf("acc-%d-app-%d-user-%d-tw", orgID, applicationID, userNumber),
		"gplus":    fmt.Sprintf("acc-%d-app-%d-user-%d-gpl", orgID, applicationID, userNumber),
		"abook":    fmt.Sprintf("acc-%d-app-%d-user-%d-abk", orgID, applicationID, userNumber),
	}

	return user
}

// CorrectEvent returns a correct event
func CorrectEvent(applicationID int64) *entity.Event {
	event := fixtures.CorrectEvent
	event.ID, _ = tgflake.FlakeNextID(applicationID, "events")
	return &event
}

func AddCorrectOrganizations(numberOfOrganizations int) []*entity.Organization {
	var err []errors.Error
	result := make([]*entity.Organization, numberOfOrganizations)
	for i := 0; i < numberOfOrganizations; i++ {
		org := CorrectOrganization()
		org.Name = fmt.Sprintf("acc-%d", i+1)
		org.Description = fmt.Sprintf("acc description %d", i+1)
		result[i], err = coreAcc.Create(org, true)
		if err != nil {
			panic(err)
		}
	}

	return result
}

func AddCorrectMembers(organization *entity.Organization, numberOfMembersPerOrg int) []*entity.Member {
	var err []errors.Error
	result := make([]*entity.Member, numberOfMembersPerOrg)
	for i := 0; i < numberOfMembersPerOrg; i++ {
		member := CorrectMemberWithDefaults(organization.ID, int64(i+1))
		member.PublicAccountID = organization.PublicID
		password := member.Password
		result[i], err = coreAccUser.Create(member, true)
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

func AddCorrectApplications(org *entity.Organization, numberOfApplicationsPerOrg int) []*entity.Application {
	var err []errors.Error
	result := make([]*entity.Application, numberOfApplicationsPerOrg)
	for i := 0; i < numberOfApplicationsPerOrg; i++ {
		application := CorrectApplication()
		application.OrgID = org.ID
		application.PublicOrgID = org.PublicID
		application.Name = fmt.Sprintf("acc-%d-app-%d", org.ID, i+1)
		application.Description = fmt.Sprintf("acc %d app %d", org.ID, i+1)
		result[i], err = coreApp.Create(application, true)
		if err != nil {
			panic(err)
		}
		_, err := coreAppRedis.Create(application, false)
		if err != nil {
			panic(err)
		}
	}

	return result
}

// HookUp create a connection between two users provided
func HookUp(orgID, applicationID int64, userFromID, userToID uint64) {
	connection := entity.Connection{
		UserFromID: userFromID,
		UserToID:   userToID,
		Type:       entity.ConnectionTypeFollow,
	}
	coreConn.Create(orgID, applicationID, &connection)
	coreConn.Confirm(orgID, applicationID, &connection, false)
}

// HookUpUsers creates connection between all users that you provide
// bidi:true will provide bidirectional connections
func HookUpUsers(orgID, applicationID int64, users []*entity.ApplicationUser, bidi bool) {
	if len(users) < 2 {
		return
	}

	for i := 1; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			HookUp(orgID, applicationID, users[i].ID, users[j].ID)
			if bidi {
				HookUp(orgID, applicationID, users[j].ID, users[i].ID)
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
func HookUpUsersCustom(orgID, applicationID int64, users []*entity.ApplicationUser) {
	if len(users) < 2 {
		return
	}

	if len(users) >= 2 {
		HookUp(orgID, applicationID, users[0].ID, users[1].ID)
		HookUp(orgID, applicationID, users[1].ID, users[0].ID)
	}

	if len(users) >= 3 {
		HookUp(orgID, applicationID, users[0].ID, users[2].ID)
		HookUp(orgID, applicationID, users[1].ID, users[2].ID)
	}

	if len(users) >= 4 {
		HookUp(orgID, applicationID, users[0].ID, users[3].ID)
		HookUp(orgID, applicationID, users[2].ID, users[3].ID)
	}

	if len(users) >= 5 {
		connection := entity.Connection{
			UserFromID: users[0].ID,
			UserToID:   users[4].ID,
			Type:       entity.ConnectionTypeFollow,
		}
		coreConn.Create(orgID, applicationID, &connection)
	}

	if len(users) >= 6 {
		connection := entity.Connection{
			UserFromID: users[4].ID,
			UserToID:   users[5].ID,
			Type:       entity.ConnectionTypeFollow,
		}
		coreConn.Create(orgID, applicationID, &connection)
	}

	if len(users) >= 8 {
		connection := entity.Connection{
			UserFromID: users[6].ID,
			UserToID:   users[7].ID,
			Type:       entity.ConnectionTypeFollow,
		}
		coreConn.Create(orgID, applicationID, &connection)
	}
}

func BenchHookUpUsersCustom(orgID, applicationID int64, users []*entity.ApplicationUser, numberOfConnectionsPerUser int) {
	if len(users) < numberOfConnectionsPerUser {
		return
	}

	connection := &entity.Connection{
		Type: entity.ConnectionTypeFriend,
	}

	for i := 0; i < len(users)-numberOfConnectionsPerUser; i++ {
		connection.UserFromID = users[i].ID
		for j := 1; j <= numberOfConnectionsPerUser; j++ {
			connection.UserToID = users[i+j].ID
			coreConn.Create(orgID, applicationID, connection)
		}
	}
}

func LoginApplicationUser(orgID, applicationID int64, user *entity.ApplicationUser) {
	sessionToken, err := coreAppUser.CreateSession(orgID, applicationID, user)
	if err != nil {
		panic(err)
	}
	user.SessionToken = sessionToken
	timeNow := time.Now()
	user.LastLogin = &timeNow
	user, err = coreAppUser.Update(orgID, applicationID, *user, *user, true)
	if err != nil {
		panic(err)
	}
}

func LoginUsers(orgID, applicationID int64, users []*entity.ApplicationUser) {
	for idx := range users {
		LoginApplicationUser(orgID, applicationID, users[idx])
	}
}

func AddCorrectApplicationUsers(orgID int64, application *entity.Application, numberOfUsersPerApplication int, hookUpUsers bool) []*entity.ApplicationUser {
	var err []errors.Error
	result := make([]*entity.ApplicationUser, numberOfUsersPerApplication)
	for i := 0; i < numberOfUsersPerApplication; i++ {
		result[i] = CorrectUserWithDefaults(orgID, application.ID, int64(i+1))
		result[i].OriginalPassword = result[i].Password
		result[i].Deleted = entity.PFalse
		err = coreAppUser.Create(orgID, application.ID, result[i])
		if err != nil {
			panic(err[0].InternalErrorWithLocation())
		}
	}

	if hookUpUsers {
		HookUpUsersCustom(orgID, application.ID, result)
	}

	return result
}

func BenchAddCorrectApplicationUsers(orgID int64, application *entity.Application, numberOfUsersPerApplication, numberOfConnectionsPerUser int, hookUpUsers bool) []*entity.ApplicationUser {
	var err []errors.Error
	result := make([]*entity.ApplicationUser, numberOfUsersPerApplication)
	for i := 0; i < numberOfUsersPerApplication; i++ {
		result[i] = CorrectUserWithDefaults(orgID, application.ID, int64(i+1))
		result[i].OriginalPassword = result[i].Password
		result[i].Deleted = entity.PFalse
		result[i].Metadata = *result[i]
		err = coreAppUser.Create(orgID, application.ID, result[i])
		if err != nil {
			panic(err[0].InternalErrorWithLocation())
		}
	}

	return result
}

// AddCorrectUserEvents adds correct events to a user
// If numberOfEventsPerUser < 4 then events are common, else they are user specific (thus unique)
func AddCorrectUserEvents(orgID, applicationID int64, user *entity.ApplicationUser, numberOfEventsPerUser int, fetchEvents bool) []*entity.Event {
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
		result[i] = CorrectEvent(applicationID)
		result[i].Visibility = uint8(i%4*10 + 10)
		result[i].UserID = user.ID
		if result[i].Visibility == entity.EventPublic {
			result[i].Location = fmt.Sprintf("location-all-%d", i+1)
			result[i].Target = &entity.Object{
				ID:           fmt.Sprintf("target-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("target-%d-all", i+1)},
			}
			result[i].Object = &entity.Object{
				ID:           fmt.Sprintf("object-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("object-%d-all", i+1)},
			}
		} else if result[i].Visibility == entity.EventGlobal {
			result[i].Location = fmt.Sprintf("location-global-%d", i+1)
			result[i].Target = &entity.Object{
				ID:           fmt.Sprintf("target-global-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("target-global-%d-all", i+1)},
			}
			result[i].Object = &entity.Object{
				ID:           fmt.Sprintf("object-global-%d", i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("object-global-%d-all", i+1)},
			}
		} else {
			result[i].Location = fmt.Sprintf("location-%d", i+1)
			result[i].Target = &entity.Object{
				ID:           fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d", orgID, applicationID, user.ID, i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-target-%d-lall", orgID, applicationID, user.ID, i+1)},
			}
			result[i].Object = &entity.Object{
				ID:           fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d", orgID, applicationID, user.ID, i+1),
				DisplayNames: map[string]string{"all": fmt.Sprintf("acc-%d-app-%d-usr-%d-object-%d-lall", orgID, applicationID, user.ID, i+1)},
			}
		}

		// Some locations are more special then others (easier things to debug when it comes to the location search by geo coordinates)
		if i < len(locations) {
			result[i].Latitude = locations[i].Lat
			result[i].Longitude = locations[i].Lon
			result[i].Location = locations[i].Label
		}

		result[i].Images = map[string]entity.Image{}
		result[i].Images["thumb_pic"] = entity.Image{
			URL:    "https://www.tapglue.com/img/box/newsfeed.jpg",
			Width:  200,
			Heigth: 200,
			Type:   "image/jpeg",
		}

		result[i].Images["original_pic"] = entity.Image{
			URL:    "https://www.tapglue.com/img/original.jpg",
			Width:  2048,
			Heigth: 2048,
			Type:   "image/jpeg",
		}

		err = coreEvt.Create(orgID, applicationID, user.ID, result[i])

		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	return result
}

func CorrectDeploy(
	numberOfOrganizations, numberOfMembersPerOrganization,
	numberOfApplicationsPerOrganization,
	numberOfUsersPerApplication, numberOfEventsPerUser int,
	hookUpUsers, loginUsers bool) []*entity.Organization {

	organizations := AddCorrectOrganizations(numberOfOrganizations)

	for i := 0; i < numberOfOrganizations; i++ {
		organizations[i].Members = AddCorrectMembers(organizations[i], numberOfMembersPerOrganization)
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
				organizations[i].Applications[j].Users[k].Events = AddCorrectUserEvents(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users[k], numberOfEventsPerUser, true)
			}
		}
	}

	return organizations
}

func CorrectDeployBench(
	numberOfOrganizations, numberOfMembersPerOrganization, numberOfApplicationsPerOrganization,
	numberOfUsersPerApplication, numberOfConnectionsPerUser, numberOfEventsPerUser int, hookUpUsers, loginUsers bool) []*entity.Organization {
	organizations := AddCorrectOrganizations(numberOfOrganizations)

	for i := 0; i < numberOfOrganizations; i++ {
		log.Printf("Generating organization %d members...\n", i+1)
		organizations[i].Members = AddCorrectMembers(organizations[i], numberOfMembersPerOrganization)
		if loginUsers {
			log.Printf("Logging in organization %d members...\n", i+1)
			LoginMembers(organizations[i].Members)
		}
		organizations[i].Members = []*entity.Member{}
		runtime.GC()

		log.Printf("Generating application for organization %d...\n", i+1)
		organizations[i].Applications = AddCorrectApplications(organizations[i], numberOfApplicationsPerOrganization)

		for j := 0; j < numberOfApplicationsPerOrganization; j++ {
			log.Printf("Generating application users for application %d ...\n", organizations[i].Applications[j].ID)
			organizations[i].Applications[j].Users = BenchAddCorrectApplicationUsers(organizations[i].ID, organizations[i].Applications[j], numberOfUsersPerApplication, numberOfConnectionsPerUser, hookUpUsers)
			wg := new(sync.WaitGroup)
			if loginUsers {
				wg.Add(1)
				go func(w *sync.WaitGroup) {
					defer w.Done()
					log.Printf("Logging in application users for application %d...\n", organizations[i].Applications[j].ID)
					LoginUsers(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users)
					runtime.GC()
				}(wg)
			}

			if hookUpUsers {
				wg.Add(1)
				go func(w *sync.WaitGroup) {
					defer w.Done()
					log.Printf("Hooking up users for organization %d, application %d", organizations[i].ID, organizations[i].Applications[j].ID)
					BenchHookUpUsersCustom(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users, numberOfConnectionsPerUser)
					runtime.GC()
				}(wg)
			}

			wg.Wait()

			for k := 0; k < numberOfUsersPerApplication; k++ {
				log.Printf("Generating event for organization %d, application %d user %d...\n", organizations[i].ID, organizations[i].Applications[j].ID, k+1)
				AddCorrectUserEvents(organizations[i].ID, organizations[i].Applications[j].ID, organizations[i].Applications[j].Users[k], numberOfEventsPerUser, false)
				runtime.GC()
			}

			organizations[i].Applications[j] = &entity.Application{}
			runtime.GC()
		}

		organizations[i] = &entity.Organization{}
		runtime.GC()
	}

	organizations = []*entity.Organization{}
	runtime.GC()

	return organizations
}

func populateEventsForUsers(orgID, appID int64, eventsPerUser int) {
	existingUsersRows, err := v03PostgresClient.SlaveDatastore(-1).
		Query(fmt.Sprintf("SELECT json_data->>'id' FROM app_%d_%d.users ORDER BY json_data->>'created_at' ASC", orgID, appID))
	if err != nil {
		panic(err)
	}
	defer existingUsersRows.Close()
	var userID uint64
	for existingUsersRows.Next() {
		err := existingUsersRows.Scan(&userID)
		if err != nil {
			panic(err)
		}
		log.Printf("Generating events for user %d\n", userID)
		AddCorrectUserEvents(orgID, appID, &entity.ApplicationUser{ID: userID}, eventsPerUser, false)
	}
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

	conn := rateLimitPool.Get()
	defer conn.Close()
	conn.Do("EVAL", "return redis.call('del', unpack(redis.call('keys', ARGV[1])))", 0, "test:*")
	conn.Do("EVAL", "return redis.call('del', unpack(redis.call('keys', ARGV[1])))", 0, "applications:*")
}

func compareUsers(c *C, expectedUser, obtainedUser *entity.ApplicationUser) {
	if obtainedUser.CreatedAt != nil && expectedUser.CreatedAt != nil {
		c.Assert(obtainedUser.CreatedAt.Sub(*expectedUser.CreatedAt), Equals, time.Duration(0))
	} else if obtainedUser.CreatedAt != expectedUser.CreatedAt {
		c.Assert(obtainedUser.CreatedAt, Equals, expectedUser.CreatedAt)
	}

	if obtainedUser.UpdatedAt != nil && expectedUser.UpdatedAt != nil {
		c.Assert(obtainedUser.UpdatedAt.Sub(*expectedUser.UpdatedAt) > 0, Equals, true)
	} else if obtainedUser.UpdatedAt != expectedUser.UpdatedAt {
		c.Assert(obtainedUser.UpdatedAt, Equals, expectedUser.UpdatedAt)
	}

	if obtainedUser.Deleted != nil && expectedUser.Deleted != nil {
		c.Assert(*obtainedUser.Deleted, Equals, *expectedUser.Deleted)
	} else if obtainedUser.Deleted != expectedUser.Deleted {
		c.Assert(obtainedUser.Deleted, Equals, expectedUser.Deleted)
	}

	if obtainedUser.FriendCount != nil && expectedUser.FriendCount != nil {
		c.Assert(*obtainedUser.FriendCount, Equals, *expectedUser.FriendCount)
	} else if obtainedUser.FriendCount != expectedUser.FriendCount {
		c.Assert(obtainedUser.FriendCount, Equals, expectedUser.FriendCount)
	}

	if obtainedUser.FollowerCount != nil && expectedUser.FollowerCount != nil {
		c.Assert(*obtainedUser.FollowerCount, Equals, *expectedUser.FollowerCount)
	} else if obtainedUser.FollowerCount != expectedUser.FollowerCount {
		c.Assert(obtainedUser.FollowerCount, Equals, expectedUser.FollowerCount)
	}

	if obtainedUser.FollowedCount != nil && expectedUser.FollowedCount != nil {
		c.Assert(*obtainedUser.FollowedCount, Equals, *expectedUser.FollowedCount)
	} else if obtainedUser.FollowedCount != expectedUser.FollowedCount {
		c.Assert(obtainedUser.FollowedCount, Equals, expectedUser.FollowedCount)
	}

	for idx := range expectedUser.SocialIDs {
		if _, ok := obtainedUser.SocialIDs[idx]; !ok {
			c.Fatalf("social key %q not found in expected user %#v", idx, obtainedUser.SocialIDs)
		}
		c.Assert(obtainedUser.SocialIDs[idx], Equals, expectedUser.SocialIDs[idx])
	}

	// WE need these to make DeepEquals work
	obtainedUser.SessionToken = expectedUser.SessionToken
	obtainedUser.OriginalPassword = expectedUser.OriginalPassword
	obtainedUser.Images = nil
	obtainedUser.LastLogin = expectedUser.LastLogin
	obtainedUser.LastRead = expectedUser.LastRead
	obtainedUser.Enabled = expectedUser.Enabled
	expectedUser.Password = ""
	expectedUser.Events = nil
	expectedUser.Images = nil
	obtainedUser.Deleted = expectedUser.Deleted
	obtainedUser.FriendCount = expectedUser.FriendCount
	obtainedUser.FollowerCount = expectedUser.FollowerCount
	obtainedUser.FollowedCount = expectedUser.FollowedCount
	obtainedUser.CreatedAt = expectedUser.CreatedAt
	obtainedUser.UpdatedAt = expectedUser.UpdatedAt

	// TODO: Better inspection for metadata is required
	obtainedUser.Metadata, expectedUser.Metadata = nil, nil

	c.Assert(obtainedUser, DeepEquals, expectedUser)
}

func compareEvents(c *C, expectedEvent, obtainedEvent *entity.Event) {
	c.Assert(obtainedEvent.Object.ID, Equals, expectedEvent.Object.ID)
	c.Assert(obtainedEvent.Object.Type, Equals, expectedEvent.Object.Type)
	c.Assert(obtainedEvent.CreatedAt.Sub(*expectedEvent.CreatedAt), Equals, time.Duration(0))
	c.Assert(obtainedEvent.UpdatedAt.Sub(*expectedEvent.UpdatedAt), Equals, time.Duration(0))

	obtainedEvent.Object, expectedEvent.Object = nil, nil
	obtainedEvent.CreatedAt, expectedEvent.CreatedAt = nil, nil
	obtainedEvent.UpdatedAt, expectedEvent.UpdatedAt = nil, nil

	// TODO: Better inspection for metadata is required
	obtainedEvent.Metadata, expectedEvent.Metadata = nil, nil

	c.Assert(obtainedEvent, DeepEquals, expectedEvent)
}
