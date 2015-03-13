/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/core/entity"

	. "gopkg.in/check.v1"
)

// Test createConnection request with a wrong key
func (s *ServerSuite) TestCreateConnection_WrongKey(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := "{usrfromidea:''}"

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createConnection request with an wrong name
func (s *ServerSuite) TestCreateConnection_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_from_id":"","user_to_id":""}`

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createConnection request
func (s *ServerSuite) TestCreateConnection_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d}`,
		correctUserFrom.ID,
		correctUserTo.ID,
	)

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(body), connection)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(connection.AccountID, Equals, correctAccount.ID)
	c.Assert(connection.ApplicationID, Equals, correctApplication.ID)
	c.Assert(connection.UserFromID, Equals, correctUserFrom.ID)
	c.Assert(connection.UserToID, Equals, correctUserTo.ID)
	c.Assert(connection.Enabled, Equals, false)
}

// Test a correct updateConnection request
func (s *ServerSuite) TestUpdateConnection_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctConnection, err := AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserFromID,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(body), connection)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(connection.AccountID, Equals, correctAccount.ID)
	c.Assert(connection.ApplicationID, Equals, correctApplication.ID)
	c.Assert(connection.UserFromID, Equals, correctUserFrom.ID)
	c.Assert(connection.UserToID, Equals, correctUserTo.ID)
	c.Assert(connection.Enabled, Equals, false)
}

// Test updateConnection request with a wrong id
func (s *ServerSuite) TestUpdateConnection_WrongID(c *C) {
	c.Skip("forced the correct user id using the contexts")
	return
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctConnection, err := AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserFromID+1,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	code, _, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test updateConnection request with an invalid name
func (s *ServerSuite) TestUpdateConnection_WrongValue(c *C) {
	c.Skip("skip because we now force things to be correct in the contexts")
	return
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctConnection, err := AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_from_id":10, "user_to_id":%d, "enabled":false}`,
		correctConnection.UserToID,
	)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteConnection request
func (s *ServerSuite) TestDeleteConnection_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	code, _, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test deleteConnection request with a wrong id
func (s *ServerSuite) TestDeleteConnection_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID+1, correctUserTo.ID)
	code, _, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusInternalServerError)
}
