/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"

	. "gopkg.in/check.v1"
)

/****************************************************************/
/******************** CREATECONNECTION TESTS ********************/
/****************************************************************/

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
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createConnection"
	route := getComposedRoute(routeName, application.AccountID, application.ID, userFrom.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.AccountID, Equals, application.AccountID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in
func (s *ServerSuite) TestCreateConnectionAfterLogin(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.AccountID, Equals, account.ID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in and refreshes session with the new token
func (s *ServerSuite) TestCreateConnectionAfterLoginRefreshNewToken(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(body, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.AccountID, Equals, account.ID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)
}

// Test to create connections after a user logs in and refreshes session with the old token
func (s *ServerSuite) TestCreateConnectionAfterLoginRefreshOldToken(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (12)\nsession mismatch")
}

// Test to create connections after a user logs in and logs out
func (s *ServerSuite) TestCreateConnectionAfterLoginLogout(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutUser"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "\"logged out\"\n")

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (10)")
}

// Test to create connections after a user logs in and logs out and logs in again
func (s *ServerSuite) TestCreateConnectionAfterLoginLogoutLogin(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payloadLogout := fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutUser"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payloadLogout, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)

	routeName = "loginUser"
	route = getComposedRoute(routeName, account.ID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)

	connection := &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
}

// Test to create connections after a user logs in and refreshes session and logs out
func (s *ServerSuite) TestCreateConnectionAfterLoginRefreshLogout(c *C) {
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

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, userFrom.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	userFrom.SessionToken = sessionToken.Token
	payload = fmt.Sprintf(`{"session_token": "%s"}`, userFrom.SessionToken)

	routeName = "logoutUser"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "\"logged out\"\n")

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (10)")
}

// Test to create connections and check the follower, followedby and connectionsevents lists
func (s *ServerSuite) TestCreateConnectionAndCheckLists(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)

	c.Assert(connection.AccountID, Equals, account.ID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Enabled, Equals, true)

	// Check connetions list
	routeName = "getConnectionList"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections := []entity.User{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections), Equals, 1)
	c.Assert(userConnections[0].ID, Equals, userTo.ID)

	// Check followedBy list
	routeName = "getFollowerList"
	route = getComposedRoute(routeName, account.ID, application.ID, userTo.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, userTo.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userConnections = []entity.User{}
	er = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(er, IsNil)

	c.Assert(len(userConnections), Equals, 1)
	c.Assert(userConnections[0].ID, Equals, userFrom.ID)

	// Check activity feed events
	routeName = "getConnectionEventList"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	userToEvents := []entity.Event{}
	er = json.Unmarshal([]byte(body), &userToEvents)
	c.Assert(er, IsNil)

	c.Assert(len(userToEvents), Equals, 2)
	c.Assert(userToEvents[0].ID, Equals, userTo.Events[len(userTo.Events)-1].ID)
	c.Assert(userToEvents[1].ID, Equals, userTo.Events[len(userTo.Events)-2].ID)
}

// Test to create connections if users are already connected
func (s *ServerSuite) TestCreateConnectionUsersAlreadyConnected(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
	c.Assert(body, Equals, "500 failed to write the user connection (3)")
}

// Test to create connections if users are from different appIDs
func (s *ServerSuite) TestCreateConnectionUsersFromDifferentApps(c *C) {
	accounts := CorrectDeploy(1, 0, 2, 2, 0, false, true)
	account := accounts[0]
	application1 := account.Applications[0]
	application2 := account.Applications[1]
	app1UserFrom := application1.Users[0]
	app2UserTo := application2.Users[0]

	payload := fmt.Sprintf(`{"user_to_id":%d}`, app2UserTo.ID)

	routeName := "createConnection"
	route := getComposedRoute(routeName, account.ID, application2.ID, app1UserFrom.ID)
	code, body, err := runRequest(routeName, route, payload, application2.AuthToken, app1UserFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (7)")
}

// Test to create connections if users are not activated
func (s *ServerSuite) TestCreateConnectionUsersNotActivated(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	payload := `{"activated": false}`

	routeName := "updateUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, userTo.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, userTo.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(`{"user_to_id":%d}`, userTo.ID)

	routeName = "createConnection"
	route = getComposedRoute(routeName, account.ID, application.ID, userFrom.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, userFrom.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test to create connections if users are not enabled
func (s *ServerSuite) TestCreateConnectionUsersNotEnabled(c *C) {
	c.Skip("not impletented")
}

// Test to create connections if one user are not activated
func (s *ServerSuite) TestCreateConnectionOneUserNotActivated(c *C) {
	c.Skip("not impletented")
}

// Test to create connections if one user are not enabled
func (s *ServerSuite) TestCreateConnectionOneUserNotEnabled(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************** UPDATECONNECTION TESTS ********************/
/****************************************************************/

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
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
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

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (9)")
}

// Test updateConnection request with a wrong id
func (s *ServerSuite) TestUpdateConnection_WrongID(c *C) {
	c.Skip("forced the correct user id using the contexts")
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

// Test to update connections after a user logs in
func (s *ServerSuite) TestUpdateConnectionAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and refreshes session
func (s *ServerSuite) TestUpdateConnectionAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and logs out
func (s *ServerSuite) TestUpdateConnectionAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and logs out and logs in again
func (s *ServerSuite) TestUpdateConnectionAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to update connections after a user logs in and refreshes session and logs out
func (s *ServerSuite) TestUpdateConnectionAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to update connections and check the follower, followedby and connectionsevents lists
func (s *ServerSuite) TestUpdateConnectionAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

// Test to update connections to enable it and check the follower, followedby and connectionsevents lists
func (s *ServerSuite) TestUpdateConnectionEnableAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

// Test to update connections to disable it and check the follower, followedby and connectionsevents lists
func (s *ServerSuite) TestUpdateConnectionDisableAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

/****************************************************************/
/******************** DELETECONNECTION TESTS ********************/
/****************************************************************/

// Test a correct deleteConnection request
func (s *ServerSuite) TestDeleteConnection_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, true, true)
	account := accounts[0]
	application := account.Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, userTo.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, userFrom.SessionToken, 3)
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
	code, body, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (9)")
}

// Test to delete connections after a user logs in
func (s *ServerSuite) TestDeleteConnectionAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and refreshes session
func (s *ServerSuite) TestDeleteConnectionAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and logs out
func (s *ServerSuite) TestDeleteConnectionAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and logs out and logs in again
func (s *ServerSuite) TestDeleteConnectionAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections after a user logs in and refreshes session and logs out
func (s *ServerSuite) TestDeleteConnectionAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to delete connections and check the follower, followedby and connectionsevents lists
func (s *ServerSuite) TestDeleteConnectionAndCheckLists(c *C) {
	c.Skip("not impletented")
	//followerList
	//followedByList
	//connectionsEventsList
}

/****************************************************************/
/******************** GETCONNECTIONLIST TESTS *******************/
/****************************************************************/

// Test to get the list of connections of the user (followsUsers)
func (s *ServerSuite) TestGetConnectionList(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in
func (s *ServerSuite) TestGetConnectionListAfterLogin(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and refreshes session
func (s *ServerSuite) TestGetConnectionListAfterLoginRefresh(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and logs out
func (s *ServerSuite) TestGetConnectionListAfterLoginLogout(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and logs out and logs in again
func (s *ServerSuite) TestGetConnectionListAfterLoginLogoutLogin(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of the user after a user logs in and refreshes session and logs out
func (s *ServerSuite) TestGetConnectionListAfterLoginRefreshLogout(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a connected user
func (s *ServerSuite) TestGetConnectionListOfConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a non-connected user
func (s *ServerSuite) TestGetConnectionListOfNonConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a user from different app
func (s *ServerSuite) TestGetConnectionListOfUserFromDifferentApp(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************* GETFOLLOWEDBYUSERS TESTS *******************/
/****************************************************************/

// Test to get the list of connections of the user (followedByUsers)
func (s *ServerSuite) TestGetFollowedByUsersList(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a connected user
func (s *ServerSuite) TestGetFollowedByUsersListOfConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a non-connected user
func (s *ServerSuite) TestGetFollowedByUsersListOfNonConnection(c *C) {
	c.Skip("not impletented")
}

// Test to get the list of connections of a user from different app
func (s *ServerSuite) TestUsersListOfUserFromDifferentApp(c *C) {
	c.Skip("not impletented")
}

/****************************************************************/
/******************** CONFIRMCONNECTION TESTS *******************/
/****************************************************************/

// Test if the lists are created after confirming a connection
func (s *ServerSuite) TestConfirmConnectionLists(c *C) {
	c.Skip("not impletented")
}

func (s *ServerSuite) TestConfirmConnection(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled": false}`, user1.ID, user2.ID)
	routeName := "createConnection"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.AccountID, Equals, application.AccountID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, user1.ID)
	c.Assert(connection.UserToID, Equals, user2.ID)
	c.Assert(connection.Enabled, Equals, false)

	payload = fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled": true}`, user1.ID, user2.ID)
	routeName = "confirmConnection"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection = &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.AccountID, Equals, application.AccountID)
	c.Assert(connection.ApplicationID, Equals, application.ID)
	c.Assert(connection.UserFromID, Equals, user1.ID)
	c.Assert(connection.UserToID, Equals, user2.ID)
	c.Assert(connection.Enabled, Equals, true)
}

/****************************************************************/
/***************** CREATESOCIALCONNECTIONS TESTS ****************/
/****************************************************************/

// Test to create connections from the social accounts
func (s *ServerSuite) TestCreateSocialConnection(c *C) {
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

	payload, er := json.Marshal(struct {
		UserFromID     int64    `json:"user_from_id"`
		SocialPlatform string   `json:"social_platform"`
		ConnectionsIDs []string `json:"connection_ids"`
	}{
		UserFromID:     userFrom.ID,
		SocialPlatform: "facebook",
		ConnectionsIDs: []string{
			user2.SocialIDs["facebook"],
			user4.SocialIDs["facebook"],
		},
	})
	c.Assert(er, IsNil)

	routeName := "createSocialConnections"
	route := getComposedRoute(routeName, account.ID, application.ID, userFrom.ID, "facebook")
	code, body, err := runRequest(routeName, route, string(payload), application.AuthToken, createApplicationUserSessionToken(userFrom), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(body, Not(Equals), "[]\n")

	connectedUsers := []*entity.User{}
	er = json.Unmarshal([]byte(body), &connectedUsers)
	c.Assert(er, IsNil)
	c.Assert(len(connectedUsers), Equals, 2)
	c.Assert(connectedUsers[0].ID, Equals, user2.ID)
	c.Assert(connectedUsers[1].ID, Equals, user4.ID)
}

// Test to create a social connection from users of differnt apps
func (s *ServerSuite) TestCreateSocialConnectionDifferentApp(c *C) {
	c.Skip("not impletented")
}

// Test to create a social connection from users of differnt network
func (s *ServerSuite) TestCreateSocialConnectionDifferentNetwork(c *C) {
	c.Skip("not impletented")
}

// Test to create a social connection from users who previously disabled the connection
func (s *ServerSuite) TestCreateSocialConnectionWhenConnectionDisabled(c *C) {
	c.Skip("not impletented")
}

func (s *ServerSuite) TestConnectionMalformedPayloadFails(c *C) {
	accounts := CorrectDeploy(1, 1, 1, 12, 0, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]
	user12 := application.Users[11]

	iterations := []struct {
		Payload   string
		RouteName string
		Route     string
		Code      int
		Body      string
	}{
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
		{
			Payload:   "",
			RouteName: "deleteConnection",
			Route:     getComposedRouteString("deleteConnection", fmt.Sprintf("%d", application.AccountID), fmt.Sprintf("%d", application.ID), fmt.Sprintf("%d", user1.ID), "90876543211234567890"),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to delete the connection(1)\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range",
		},
		{
			Payload:   "",
			RouteName: "deleteConnection",
			Route:     getComposedRoute("deleteConnection", application.AccountID, application.ID, user1.ID, user12.ID),
			Code:      http.StatusNotFound,
			Body:      "404 failed to delete the connection (2)",
		},
		{
			Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`, user1.ID, user1.ID),
			RouteName: "createConnection",
			Route:     getComposedRoute("createConnection", application.AccountID, application.ID, user1.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to create connection (2)\nuser is connecting with itself",
		},
		{
			Payload:   "{",
			RouteName: "confirmConnection",
			Route:     getComposedRoute("confirmConnection", application.AccountID, application.ID, user1.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to confirm the connection (1)\nunexpected end of JSON input",
		},
		{
			Payload:   fmt.Sprintf(`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`, user1.ID, 13),
			RouteName: "confirmConnection",
			Route:     getComposedRoute("confirmConnection", application.AccountID, application.ID, user1.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 user does not exists",
		},
		{
			Payload:   "",
			RouteName: "createSocialConnections",
			Route:     getComposedRoute("createSocialConnections", application.AccountID, application.ID, user1.ID, "fake"),
			Code:      http.StatusNotFound,
			Body:      "404 social connecting failed (1)\nunexpected social platform",
		},
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %d}`, 13),
			RouteName: "createSocialConnections",
			Route:     getComposedRoute("createSocialConnections", application.AccountID, application.ID, user1.ID, "facebook"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (3)\nuser mismatch",
		},
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %d, "social_platform": "%s"}`, user1.ID, "fake"),
			RouteName: "createSocialConnections",
			Route:     getComposedRoute("createSocialConnections", application.AccountID, application.ID, user1.ID, "facebook"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (3)\nplatform mismatch",
		},
		{
			Payload:   fmt.Sprintf(`{"user_from_id": %d, "social_platform": "%s"`, user1.ID, "fake"),
			RouteName: "createSocialConnections",
			Route:     getComposedRoute("createSocialConnections", application.AccountID, application.ID, user1.ID, "facebook"),
			Code:      http.StatusBadRequest,
			Body:      "400 social connecting failed (2)\nunexpected end of JSON input",
		},
	}

	for idx := range iterations {
		code, body, err := runRequest(iterations[idx].RouteName, iterations[idx].Route, iterations[idx].Payload, application.AuthToken, user1.SessionToken, 3)
		c.Logf("pass %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, iterations[idx].Code)
		c.Assert(body, Equals, iterations[idx].Body)
	}
}
