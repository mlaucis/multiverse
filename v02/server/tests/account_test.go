/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/v02/entity"

	. "gopkg.in/check.v1"
)

// Test createAcccount request with a wrong key
func (s *ServerSuite) TestCreateAccount_WrongKey(c *C) {
	payload := "{namae:''}"

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, nilSigner)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createAcccount request with an wrong name
func (s *ServerSuite) TestCreateAccount_WrongValue(c *C) {
	payload := `{"name":""}`

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, nilSigner)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccount request
func (s *ServerSuite) TestCreateAccount_OK(c *C) {
	account := CorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, account.Name, account.Description)

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, nilSigner)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	if receivedAccount.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct updateAccount request
func (s *ServerSuite) TestUpdateAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	if receivedAccount.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Description, Equals, description)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct updateAccount request with a wrong id
func (s *ServerSuite) TestUpdateAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.ID+1)
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccount request with a wrong id
func (s *ServerSuite) TestUpdateAccountMalformedPayload(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to update the account (1)\nunexpected end of JSON input")
}

// Test a correct updateAccount request with an invalid description
func (s *ServerSuite) TestUpdateAccount_Invalid(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(`{"name":"%s", "description":"","enabled":true}`, account.Name)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteAccount request
func (s *ServerSuite) TestDeleteAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccount request with a wrong id
func (s *ServerSuite) TestDeleteAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID+1)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, "404 The resource for the provided id doesn't exist")
}

// Test a correct getAccount request
func (s *ServerSuite) TestGetAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	c.Assert(receivedAccount.ID, Equals, account.ID)
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct getAccount request with a wrong id
func (s *ServerSuite) TestGetAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, account.ID+1)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}