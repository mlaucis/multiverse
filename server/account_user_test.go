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

// Test create acccountUser request with a wrong key
func (s *ServerSuite) TestCreateAccountUser_WrongKey(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	payload := "{usrnamae:''}"

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test create acccountUser request with an wrong name
func (s *ServerSuite) TestCreateAccountUser_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	payload := `{"user_name":""}`

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccountUser request
func (s *ServerSuite) TestCreateAccountUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser := CorrectAccountUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"%s", "first_name": "%s", "last_name": "%s", "email": "%s"}`,
		correctAccountUser.Username,
		correctAccountUser.Password,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), accountUser)
	c.Assert(err, IsNil)
	if accountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
	c.Assert(accountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request
func (s *ServerSuite) TestUpdateAccountUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		correctAccountUser.Username,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), accountUser)
	c.Assert(err, IsNil)
	if accountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
	c.Assert(accountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request with a wrong id
func (s *ServerSuite) TestUpdateAccountUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		correctAccountUser.Username,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID+1)
	code, _, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccountUser request with an invalid description
func (s *ServerSuite) TestUpdateAccountUser_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		correctAccountUser.Username,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct updateAccountUser request with a wrong token
func (s *ServerSuite) TestUpdateAccountUser_WrongToken(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		correctAccountUser.Username,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)
	c.Assert(err, IsNil)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, _, err := runRequest(routeName, route, payload, "")

	c.Assert(code, Equals, http.StatusBadRequest)
}

// Test a correct deleteAccountUser request
func (s *ServerSuite) TestDeleteAccountUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, _, err := runRequest(routeName, route, "", "")

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccountUser request with a wrong id
func (s *ServerSuite) TestDeleteAccountUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", "")

	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct deleteAccountUser request with a wrong token
func (s *ServerSuite) TestDeleteAccountUser_WrongToken(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, _, err := runRequest(routeName, route, "", "")

	c.Assert(code, Equals, http.StatusBadRequest)
}

// Test a correct getAccountUser request
func (s *ServerSuite) TestGetAccountUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, body, err := runRequest(routeName, route, "", "")

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), accountUser)
	c.Assert(err, IsNil)
	c.Assert(accountUser.ID, Equals, correctAccountUser.ID)
	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
	c.Assert(accountUser.Enabled, Equals, true)
}

// Test a correct getAccountUser request with a wrong id
func (s *ServerSuite) TestGetAccountUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", "")

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccountUser request with a wrong token
func (s *ServerSuite) TestGetAccountUser_WrongToken(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, correctAccountUser.AccountID, correctAccountUser.ID)
	code, _, err := runRequest(routeName, route, "", "")

	c.Assert(code, Equals, http.StatusBadRequest)
}
