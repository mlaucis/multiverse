// +build !bench

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/tapglue/multiverse/v04/entity"

	. "gopkg.in/check.v1"
)

/****************************************************************/
/******************** CREATECONNECTION TESTS ********************/
/****************************************************************/

// Test createCurrentUserConnection request with a wrong key
func (s *ConnectionSuite) TestCreateConnection_WrongKey(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := "{usrfromidea:''}"

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createCurrentUserConnection request with an wrong name
func (s *ConnectionSuite) TestCreateConnection_WrongValue(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := `{"user_from_id":"","user_to_id":""}`

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Not(Equals), http.StatusCreated)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createCurrentUserConnection request
func (s *ConnectionSuite) TestCreateConnection_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": "friend"}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)
}

func (s *ConnectionSuite) TestCreateConnectionTwice(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": "friend"}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(body, Equals, "\"\"\n")
}

func (s *ConnectionSuite) TestCreateFriendConnection(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_to_id":%d}`,
		userTo.ID,
	)

	routeName := "createCurrentUserFriendConnectionAlias"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)
}

func (s *ConnectionSuite) TestCreateFollowConnection(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_to_id":%d}`,
		userTo.ID,
	)

	routeName := "createCurrentUserFollowConnectionAlias"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "follow")
	c.Assert(connection.Enabled, Equals, true)
}

func (s *ConnectionSuite) TestCreateConnectionWithCustomIDs_OK(c *C) {
	c.Skip("we should implement a different logic here, as it's hard to reuse the same field for int and string and then int again in Go")

	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_to_id":%d, "type": "friend"}`,
		userTo.CustomID,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in
func (s *ConnectionSuite) TestCreateConnectionAfterLogin(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "follow"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in and refreshes session with the new token
func (s *ConnectionSuite) TestCreateConnectionAfterLoginRefreshNewToken(c *C) {
	c.Skip("Skip this for now as we don't have the endpoint in the docs yet")
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), userFrom.SessionToken)

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in and refreshes session with the old token
func (s *ConnectionSuite) TestCreateConnectionAfterLoginRefreshOldToken_Works(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
}

// Test to create connections after a user logs in and logs out
func (s *ConnectionSuite) TestCreateConnectionAfterLoginLogout(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test to create connections after a user logs in and logs out and logs in again
func (s *ConnectionSuite) TestCreateConnectionAfterLoginLogoutLogin(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payloadLogout := fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payloadLogout, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
}

// Test to create connections after a user logs in and refreshes session and logs out
func (s *ConnectionSuite) TestCreateConnectionAfterLoginRefreshLogout(c *C) {
	c.Skip("Skip this for now as we don't have the endpoint in the docs yet")
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		userFrom.Email,
		userFrom.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token
	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (10)")
}

// Test to create connections and check the follower, followedby and connectionsevents lists
func (s *ConnectionSuite) TestCreateFollowConnectionAndCheckLists(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(account.ID, application.ID, userFrom)
	LoginApplicationUser(account.ID, application.ID, userTo)

	payload := fmt.Sprintf(`{"user_to_id":%d,  "type": "follow"}`, userTo.ID)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
	c.Assert(connection.Type, Equals, "follow")

	// Check connetions list
	routeName = "getUserFollows"
	route = getComposedRoute(routeName, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections := struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userTo.ID)

	// Check followedBy list
	routeName = "getUserFollowers"
	route = getComposedRoute(routeName, userTo.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userFrom.ID)

	// Check activity feed events
	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int            `json:"unread_events_count"`
		Events []entity.Event `json:"events"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(response.Events[0].ID, Equals, userTo.Events[len(userTo.Events)-1].ID)
}

// Test to create connections and check the friend lists
func (s *ConnectionSuite) TestCreateFriendConnectionAndCheckLists(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(`{"user_to_id":%d,  "type": "friend"}`, userTo.ID)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
	c.Assert(connection.Type, Equals, "friend")

	// Check connetions list
	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections := struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userTo.ID)

	// Check followedBy list
	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userTo.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userFrom.ID)

	// Check activity feed events
	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int            `json:"unread_events_count"`
		Events []entity.Event `json:"events"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(response.Events[0].ID, Equals, userTo.Events[len(userTo.Events)-1].ID)
}

// Test to create connections and check the friend lists
func (s *ConnectionSuite) TestCreateFriendConnectionAfterDeletingTheSameFriendConnection(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(`{"user_to_id":%d,  "type": "friend"}`, userTo.ID)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
	c.Assert(connection.Type, Equals, "friend")

	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections := struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userTo.ID)

	// Check followedBy list
	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userTo.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userFrom.ID)

	routeName = "deleteCurrentUserConnection"
	route = getComposedRoute(routeName, userTo.ID)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(`{"user_to_id":%d,  "type": "friend"}`, userTo.ID)
	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection = &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
	c.Assert(connection.Type, Equals, "friend")

	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userTo.ID)

	// Check followedBy list
	routeName = "getUserFriends"
	route = getComposedRoute(routeName, userTo.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = struct {
		Users      []entity.ApplicationUser `json:"users"`
		UsersCount int                      `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections.Users), Equals, 1)
	c.Assert(userConnections.UsersCount, Equals, 1)
	c.Assert(userConnections.Users[0].ID, Equals, userFrom.ID)

	// Check activity feed events
	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int            `json:"unread_events_count"`
		Events []entity.Event `json:"events"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(response.Events[0].ID, Equals, userTo.Events[len(userTo.Events)-1].ID)
}

// Test to create connections if users are already connected
func (s *ConnectionSuite) TestCreateConnectionUsersAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, userTo.ID)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, _, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
}

// Test to create connections if users are from different appIDs
// TODO Remove this for good. Our IDs are completly separated now, individual per app but due to their generating scheme
// they'll most likely overlap at points in time (whic is ok).
/*func (s *ConnectionSuite) TestCreateConnectionUsersFromDifferentApps(c *C) {
	accounts := CorrectDeploy(1, 0, 2, 2, 0, false, true)
	account := accounts[0]
	application1 := account.Applications[0]
	application2 := account.Applications[1]
	app1UserFrom := application1.Users[0]
	app2UserTo := application2.Users[0]

	payload := fmt.Sprintf(`{"user_to_id":%d, "type": "friend"}`, app2UserTo.ID)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application1, app1UserFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, `{"errors":[{"code":1001,"message":"application user not found"},{"code":1001,"message":"application user not found"},{"code":1000,"message":"user not activated"}]}`+"\n")
}*/

// Test to create connections if users are not activated
func (s *ConnectionSuite) TestCreateConnectionUsersNotActivated(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]
	LoginApplicationUser(account.ID, application.ID, userFrom)

	payload := `{"activated": false}`

	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test to create connections if users are not enabled
func (s *ConnectionSuite) TestCreateConnectionUsersNotEnabled(c *C) {
	c.Skip("not impletented")
}

// Test to create connections if one user are not activated
func (s *ConnectionSuite) TestCreateConnectionOneUserNotActivated(c *C) {
	c.Skip("not impletented")
}

// Test to create connections if one user are not enabled
func (s *ConnectionSuite) TestCreateConnectionOneUserNotEnabled(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************** UPDATECONNECTION TESTS ********************/
/****************************************************************/

// Test a correct updateConnection request
func (s *ConnectionSuite) TestUpdateConnection_OK(c *C) {
	c.Skip("not available in 0.3")
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, userTo.ID)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, false)
}

// Test to update connections after a user logs in
func (s *ConnectionSuite) TestUpdateConnectionAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and refreshes session
func (s *ConnectionSuite) TestUpdateConnectionAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and logs out
func (s *ConnectionSuite) TestUpdateConnectionAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and logs out and logs in again
func (s *ConnectionSuite) TestUpdateConnectionAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and refreshes session and logs out
func (s *ConnectionSuite) TestUpdateConnectionAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to update connections and check the follower, followedby and connectionsevents lists
func (s *ConnectionSuite) TestUpdateConnectionAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

// Test to update connections to enable it and check the follower, followedby and connectionsevents lists
func (s *ConnectionSuite) TestUpdateConnectionEnableAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

// Test to update connections to disable it and check the follower, followedby and connectionsevents lists
func (s *ConnectionSuite) TestUpdateConnectionDisableAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

/****************************************************************/
/******************** DELETECONNECTION TESTS ********************/
/****************************************************************/

// Test a correct deleteCurrentUserConnection request
func (s *ConnectionSuite) TestDeleteConnection_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(account.ID, application.ID, userFrom)

	routeName := "deleteCurrentUserConnection"
	route := getComposedRoute(routeName, userTo.ID)
	code, _, err := runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test deleteCurrentUserConnection request with a wrong id
func (s *ConnectionSuite) TestDeleteConnection_WrongID(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	routeName := "deleteCurrentUserConnection"
	route := getComposedRoute(routeName, userTo.ID+1)
	code, _, err := runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

// Test to delete connections after a user logs in
func (s *ConnectionSuite) TestDeleteConnectionAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and refreshes session
func (s *ConnectionSuite) TestDeleteConnectionAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and logs out
func (s *ConnectionSuite) TestDeleteConnectionAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and logs out and logs in again
func (s *ConnectionSuite) TestDeleteConnectionAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and refreshes session and logs out
func (s *ConnectionSuite) TestDeleteConnectionAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections and check the follower, followedby and connectionsevents lists
func (s *ConnectionSuite) TestDeleteConnectionAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

/****************************************************************/
/******************** GETCONNECTIONLIST TESTS *******************/
/****************************************************************/

// Test to get the list of connections of the user (followsUsers)
func (s *ConnectionSuite) TestGetConnectionList(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in
func (s *ConnectionSuite) TestGetConnectionListAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and refreshes session
func (s *ConnectionSuite) TestGetConnectionListAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and logs out
func (s *ConnectionSuite) TestGetConnectionListAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and logs out and logs in again
func (s *ConnectionSuite) TestGetConnectionListAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and refreshes session and logs out
func (s *ConnectionSuite) TestGetConnectionListAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a connected user
func (s *ConnectionSuite) TestGetConnectionListOfConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a non-connected user
func (s *ConnectionSuite) TestGetConnectionListOfNonConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a user from different app
func (s *ConnectionSuite) TestGetConnectionListOfUserFromDifferentApp(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************* GETFOLLOWEDBYUSERS TESTS *******************/
/****************************************************************/

// Test to get the list of connections of the user (followedByUsers)
func (s *ConnectionSuite) TestGetFollowedByUsersList(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a connected user
func (s *ConnectionSuite) TestGetFollowedByUsersListOfConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a non-connected user
func (s *ConnectionSuite) TestGetFollowedByUsersListOfNonConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a user from different app
func (s *ConnectionSuite) TestUsersListOfUserFromDifferentApp(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************** CONFIRMCONNECTION TESTS *******************/
/****************************************************************/

// Test if the lists are created after confirming a connection
func (s *ConnectionSuite) TestConfirmConnectionLists(c *C) {
	c.Skip("not impletented")
}

func (s *ConnectionSuite) TestConfirmConnection(c *C) {
	c.Skip("We don't support confirming connections for now so we can disable this")
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, user1)

	payload := fmt.Sprintf(`{"user_from_id":%q, "user_to_id":%d, "type": "friend", "enabled": false}`, user1.ID, user2.ID)
	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, user1.ID)
	c.Assert(connection.UserToID, Equals, user2.ID)
	c.Assert(connection.Enabled, Equals, false)

	payload = fmt.Sprintf(`{"user_from_id":%q, "user_to_id":%d, "type":"friend", "enabled": true}`, user1.ID, user2.ID)
	routeName = "confirmConnection"
	route = getComposedRoute(routeName, user2.ID)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection = &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, user1.ID)
	c.Assert(connection.UserToID, Equals, user2.ID)
	c.Assert(connection.Enabled, Equals, true)
}

/****************************************************************/
/***************** CREATESOCIALCONNECTIONS TESTS ****************/
/****************************************************************/

// Test to create connections from the social accounts
func (s *ConnectionSuite) TestCreateSocialConnection(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)
}

// Test to create a social connection from users of differnt apps
func (s *ConnectionSuite) TestCreateSocialConnectionDifferentApp(c *C) {
	c.Skip("not impletented")
}

// Test to create a social connection from users of differnt network
func (s *ConnectionSuite) TestCreateSocialConnectionDifferentNetwork(c *C) {
	c.Skip("not impletented")
}

// Test to create a social connection from users who previously disabled the connection
func (s *ConnectionSuite) TestCreateSocialConnectionWhenConnectionDisabled(c *C) {
	c.Skip("not impletented")
}

func (s *ConnectionSuite) TestConnectionMalformedPayloadFails(c *C) {
	c.Skip("Skip this for now")
	accounts := CorrectDeploy(1, 1, 1, 12, 0, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user12 := application.Users[11]

	iterations := []struct {
		Payload   string
		RouteName string
		Route     string
		Code      int
		Body      string
	}{
		/*
			We don't have the update endpoint anymore so we disable this for now

			{
				Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`, user1.ID, user2.ID),
				RouteName: "updateConnection",
				Route:     getComposedRouteString("updateConnection", fmt.Sprintf("%d", application.AccountID), fmt.Sprintf("%d", application.ID), fmt.Sprintf("%d", user1.ID), "90876543211234567890"),
				Code:      http.StatusBadRequest,
				Body:      "400 failed to update the connection (1)\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range",
			},
			{
				Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false`, user1.ID, user2.ID),
				RouteName: "updateConnection",
				Route:     getComposedRoute("updateConnection", application.AccountID, application.ID, user1.ID, user2.ID),
				Code:      http.StatusBadRequest,
				Body:      "400 failed to update the connection (4)\nunexpected end of JSON input",
			},
			{
				Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`, user1.ID, 0),
				RouteName: "updateConnection",
				Route:     getComposedRoute("updateConnection", application.AccountID, application.ID, user1.ID, user2.ID),
				Code:      http.StatusBadRequest,
				Body:      "400 failed to update the connection (6)\nuser_to mismatch",
			},
			{
				Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`, user1.ID, user12.ID),
				RouteName: "updateConnection",
				Route:     getComposedRoute("updateConnection", application.AccountID, application.ID, user1.ID, user12.ID),
				Code:      http.StatusNotFound,
				Body:      "404 failed to update the connection (3)\nusers are not connected",
			},
		*/
		// 0
		{
			Payload:   "",
			RouteName: "deleteCurrentUserConnection",
			Route:     getComposedRouteString("deleteCurrentUserConnection", user12.ID),
			Code:      http.StatusNotFound,
			Body:      `{"message":"connection not found"}` + "\n",
		},
		// 1
		{
			Payload:   "",
			RouteName: "deleteCurrentUserConnection",
			Route:     getComposedRoute("deleteCurrentUserConnection", user12.ID),
			Code:      http.StatusNotFound,
			Body:      `{"message":"connection not found"}` + "\n",
		},
		// 2
		{
			Payload:   fmt.Sprintf(`{"user_from_id":%q, "user_to_id":%d, "enabled":false}`, user1.ID, user1.ID),
			RouteName: "createCurrentUserConnection",
			Route:     getComposedRoute("createCurrentUserConnection", user12.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to create connection (2)\nuser is connecting with itself",
		},
		// 3
		{
			Payload:   "{",
			RouteName: "confirmConnection",
			Route:     getComposedRoute("confirmConnection"),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to confirm the connection (1)\nunexpected end of JSON input",
		},
		// 4
		{
			Payload:   fmt.Sprintf(`{"user_from_id":%q, "user_to_id":%d, "enabled":false}`, user1.ID, 13),
			RouteName: "confirmConnection",
			Route:     getComposedRoute("confirmConnection"),
			Code:      http.StatusBadRequest,
			Body:      "400 user does not exists",
		},
		// 5
		{
			Payload:   "",
			RouteName: "createCurrentUserSocialConnections",
			Route:     getComposedRoute("createCurrentUserSocialConnections"),
			Code:      http.StatusNotFound,
			Body:      "404 social connecting failed (1)\nunexpected social platform",
		},
		// 6
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %q}`, "13"),
			RouteName: "createCurrentUserSocialConnections",
			Route:     getComposedRoute("createCurrentUserSocialConnections"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (3)\nuser mismatch",
		},
		// 7
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %q, "social_platform": "%s"}`, user1.ID, "fake"),
			RouteName: "createCurrentUserSocialConnections",
			Route:     getComposedRoute("createCurrentUserSocialConnections"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (3)\nplatform mismatch",
		},
		// 8
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %q, "social_platform": "%s"`, user1.ID, "fake"),
			RouteName: "createCurrentUserSocialConnections",
			Route:     getComposedRoute("createCurrentUserSocialConnections"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (2)\nunexpected end of JSON input",
		},
	}

	for idx := range iterations {
		c.Logf("pass %d", idx)
		code, body, err := runRequest(iterations[idx].RouteName, iterations[idx].Route, iterations[idx].Payload, signApplicationRequest(application, user1, true, true))
		c.Assert(err, IsNil)
		c.Assert(code, Equals, iterations[idx].Code)
		c.Assert(body, Equals, iterations[idx].Body)
	}
}

func (s *ConnectionSuite) TestCreateSocialConnectionFriendsAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)

	payload, er = json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     user2.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			userFrom.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName = "createCurrentUserSocialConnections"
	code, body, err = runRequest(routeName, route, string(payload), signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers = struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	sort.Sort(AppUserByID(connectedUsers.Users))
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, userFrom.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)
}

func (s *ConnectionSuite) TestCreateSocialConnectionFollowsAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "follow",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)

	payload, er = json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     user2.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			userFrom.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "follow",
	})
	c.Assert(er, IsNil)

	routeName = "createCurrentUserSocialConnections"
	code, body, err = runRequest(routeName, route, string(payload), signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers = struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	sort.Sort(AppUserByID(connectedUsers.Users))
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, userFrom.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)
}

func (s *ConnectionSuite) TestCreateSocialConnectionFollowsFriendAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "follow",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)

	payload, er = json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     user2.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			userFrom.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName = "createCurrentUserSocialConnections"
	code, body, err = runRequest(routeName, route, string(payload), signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers = struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	sort.Sort(AppUserByID(connectedUsers.Users))
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, userFrom.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)
}

func (s *ConnectionSuite) TestCreateSocialConnectionFriendFollowsAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)

	payload, er = json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     user2.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			userFrom.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "follow",
	})
	c.Assert(er, IsNil)

	routeName = "createCurrentUserSocialConnections"
	code, body, err = runRequest(routeName, route, string(payload), signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers = struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	sort.Sort(AppUserByID(connectedUsers.Users))
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, userFrom.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)
}

func (s *ConnectionSuite) TestGetConnectionsCount(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 10, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	userFriend := application.Users[9]

	payload := fmt.Sprintf(
		`{"user_to_id":%d}`,
		userFriend.ID,
	)

	routeName := "createCurrentUserFriendConnectionAlias"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, user.ID)
	c.Assert(connection.UserToID, Equals, userFriend.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)

	routeName = "getCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(*receivedUser.FriendCount, Equals, int64(1))
	c.Assert(*receivedUser.FollowerCount, Equals, int64(1))
	c.Assert(*receivedUser.FollowedCount, Equals, int64(3))
}

func (s *ConnectionSuite) TestCreateFriendsConnectionWithEvent(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": "friend"}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName) + "?with_event=true"
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)

	// We need to wait a bit for the background event to be created
	time.Sleep(10 * time.Millisecond)

	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int                               `json:"unread_events_count"`
		Events []entity.Event                    `json:"events"`
		Users  map[string]entity.ApplicationUser `json:"users"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(len(response.Users), Equals, 1)
	c.Assert(response.Events[0].Type, Equals, "tg_friend")
}

func (s *ConnectionSuite) TestCreateFollowConnectionWithEvent(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": "follow"}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName) + "?with_event=true"
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, "follow")
	c.Assert(connection.Enabled, Equals, true)

	// We need to wait a bit for the background event to be created
	time.Sleep(10 * time.Millisecond)

	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int                               `json:"unread_events_count"`
		Events []entity.Event                    `json:"events"`
		Users  map[string]entity.ApplicationUser `json:"users"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(len(response.Users), Equals, 1)
	c.Assert(response.Events[0].Type, Equals, "tg_follow")
}

func (s *ConnectionSuite) TestCreateSocialConnectionWithEvent(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 5, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]

	userFrom := application.Users[0]
	user2 := application.Users[1]
	user4 := application.Users[3]

	payload, er := json.Marshal(struct {
		UserFromID     string   `json:"platform_user_id"`
		SocialPlatform string   `json:"platform"`
		ConnectionsIDs []string `json:"connection_ids"`
		Type           string   `json:"type"`
	}{
		UserFromID:     userFrom.SocialIDs["facebook"],
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
		Type: "friend",
	})
	c.Assert(er, IsNil)

	routeName := "createCurrentUserSocialConnections"
	route := getComposedRoute(routeName) + "?with_event=true"
	code, body, err := runRequest(routeName, route, string(payload), signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(connectedUsers.UsersCount, Equals, 2)
	c.Assert(connectedUsers.Users[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers.Users[1].ID, Equals, user4.ID)

	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int                               `json:"unread_events_count"`
		Events []entity.Event                    `json:"events"`
		Users  map[string]entity.ApplicationUser `json:"users"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(len(response.Users), Equals, 1)
	c.Assert(response.Events[0].Type, Equals, "tg_friend")
}
