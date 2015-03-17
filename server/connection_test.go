/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"

	. "gopkg.in/check.v1"
)

// Test createConnection request with a wrong key
func (s *ServerSuite) TestCreateConnection_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	correctUser, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := "{usrfromidea:''}"

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createConnection request with an wrong name
func (s *ServerSuite) TestCreateConnection_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	correctUser, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_from_id":"","user_to_id":""}`

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createConnection request
func (s *ServerSuite) TestCreateConnection_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser2(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(body), connection)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(connection.AccountID, Equals, account.ID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test a correct updateConnection request
func (s *ServerSuite) TestUpdateConnection_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser2(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	correctConnection, err := AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserFromID,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(body), connection)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(connection.AccountID, Equals, account.ID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, false)
}

// Test a correct updateConnection request
func (s *ServerSuite) TestUpdateConnection_NotCrossUpdate(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser2(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	correctConnection, err := AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserFromID,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userTo), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusUnauthorized)

	c.Assert(body, Equals, "401 failed to check session token (9)")
}

// Test updateConnection request with a wrong id
func (s *ServerSuite) TestUpdateConnection_WrongID(c *C) {
	c.Skip("forced the correct user id using the contexts")
	return
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	correctConnection, err := AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserFromID+1,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, _, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test updateConnection request with an invalid name
func (s *ServerSuite) TestUpdateConnection_WrongValue(c *C) {
	c.Skip("skip because we now force things to be correct in the contexts")
	return
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	connection, err := AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":10, "user_to_id":%d, "enabled":false}`,
		connection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteConnection request
func (s *ServerSuite) TestDeleteConnection_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser2(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	_, err = AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test deleteConnection request with a wrong id
func (s *ServerSuite) TestDeleteConnection_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	userTo, err := AddCorrectUser2(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	_, err = AddCorrectConnection(account.ID, application.ID, userFrom.ID, userTo.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID+1, userTo.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusUnauthorized)
}

func (s *ServerSuite) TestAddSocialConnection(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	userFrom := CorrectUser()
	userFrom.AccountID = account.ID
	userFrom.ApplicationID = application.ID
	userFrom.Username = "social1"
	userFrom.Email = "social-connection@user1.com"
	userFrom.SocialIDs = map[string]string{"facebook": "fb-id-1"}
	userFrom, err = core.WriteUser(userFrom, true)
	c.Assert(err, IsNil)

	user2 := CorrectUser()
	user2.AccountID = account.ID
	user2.ApplicationID = application.ID
	user2.Username = "social2"
	user2.Email = "social-connection@user2.com"
	user2.SocialIDs = map[string]string{"facebook": "fb-id-2"}
	user2, err = core.WriteUser(user2, true)
	c.Assert(err, IsNil)

	user3 := CorrectUser()
	user3.AccountID = account.ID
	user3.ApplicationID = application.ID
	user3.Username = "social3"
	user3.Email = "social-connection@user3.com"
	user3.SocialIDs = map[string]string{"facebook": "fb-id-3"}
	user3, err = core.WriteUser(user3, true)
	c.Assert(err, IsNil)

	user4 := CorrectUser()
	user4.AccountID = account.ID
	user4.ApplicationID = application.ID
	user4.Username = "social4"
	user4.Email = "social-connection@user4.com"
	user4.SocialIDs = map[string]string{"facebook": "fb-id-4"}
	user4, err = core.WriteUser(user4, true)
	c.Assert(err, IsNil)

	user5 := CorrectUser()
	user5.AccountID = account.ID
	user5.ApplicationID = application.ID
	user5.Username = "social5"
	user5.Email = "social-connection@user5.com"
	user5.SocialIDs = map[string]string{"facebook": "fb-id-5"}
	user5, err = core.WriteUser(user5, true)
	c.Assert(err, IsNil)

	payload, err := json.Marshal(struct {
		UserFromID int64 `json:"user_from_id"`
		SocialPlatform string `json:"social_platform"`
		ConnectionsIDs []string `json:"connection_ids"`
	}{
		UserFromID: userFrom.ID,
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
	})

	routeName := "createSocialConnections"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, "facebook")
	code, body, err := runRequest(routeName, route, string(payload), application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")
}
