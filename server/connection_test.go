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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createConnection request with an wrong name
func (s *ServerSuite) TestCreateConnection_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_from_id":"","user_to_id":""}`

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(response), connection)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(connection.AccountID, Equals, correctAccount.ID)
	c.Assert(connection.ApplicationID, Equals, correctApplication.ID)
	c.Assert(connection.UserFromID, Equals, correctUserFrom.ID)
	c.Assert(connection.UserToID, Equals, correctUserTo.ID)
	c.Assert(connection.Enabled, Equals, true)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	connection := &entity.Connection{}
	err = json.Unmarshal([]byte(response), connection)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test updateConnection request with an invalid name
func (s *ServerSuite) TestUpdateConnection_WrongValue(c *C) {
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteConnection request
func (s *ServerSuite) TestDeleteConnection_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test deleteConnection request with a wrong id
func (s *ServerSuite) TestDeleteConnection_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUserFrom, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctUserTo, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	AddCorrectConnection(correctAccount.ID, correctApplication.ID, correctUserFrom.ID, correctUserTo.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteConnection"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUserFrom.ID+1, correctUserTo.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}
