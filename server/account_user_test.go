/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"

	. "gopkg.in/check.v1"
)

// Test create acccountUser request with a wrong key
func (s *ServerSuite) TestCreateAccountUser_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	payload := "{usrnamae:''}"

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test create acccountUser request with an wrong name
func (s *ServerSuite) TestCreateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	payload := `{"user_name":""}`

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccountUser request
func (s *ServerSuite) TestCreateAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser := CorrectAccountUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"%s", "first_name": "%s", "last_name": "%s", "email": "%s"}`,
		accountUser.Username,
		accountUser.Password,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	if receivedAccountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request
func (s *ServerSuite) TestUpdateAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	if receivedAccountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request with a wrong id
func (s *ServerSuite) TestUpdateAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccountUser request with an invalid description
func (s *ServerSuite) TestUpdateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct updateAccountUser request with a wrong token
func (s *ServerSuite) TestUpdateAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)
	c.Assert(err, IsNil)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(code, Equals, http.StatusBadRequest)
}

// Test a correct deleteAccountUser request
func (s *ServerSuite) TestDeleteAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccountUser request with a wrong id
func (s *ServerSuite) TestDeleteAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct deleteAccountUser request with a wrong token
func (s *ServerSuite) TestDeleteAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)

	c.Assert(code, Equals, http.StatusUnauthorized)
}

// Test a correct getAccountUser request
func (s *ServerSuite) TestGetAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	c.Assert(receivedAccountUser.ID, Equals, accountUser.ID)
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

// Test a correct getAccountUser request with a wrong id
func (s *ServerSuite) TestGetAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccountUser request with a wrong token
func (s *ServerSuite) TestGetAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)

	c.Assert(code, Equals, http.StatusUnauthorized)
}
