/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/core/entity"

	"fmt"

	. "gopkg.in/check.v1"
)

// Test createAcccount request with a wrong key
func (s *ServerSuite) TestCreateAccount_WrongKey(c *C) {
	payload := "{namae:''}"

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createAcccount request with an wrong name
func (s *ServerSuite) TestCreateAccount_WrongValue(c *C) {
	payload := `{"name":""}`

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccount request
func (s *ServerSuite) TestCreateAccount_OK(c *C) {
	correctAccount := CorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, correctAccount.Name, correctAccount.Description)

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(body), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
}

// Test a correct updateAccount request
func (s *ServerSuite) TestUpdateAccount_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(body), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Description, Equals, description)
	c.Assert(account.Enabled, Equals, true)
}

// Test a correct updateAccount request with a wrong id
func (s *ServerSuite) TestUpdateAccount_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID+1)
	code, _, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccount request with an invalid description
func (s *ServerSuite) TestUpdateAccount_Invalid(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	payload := fmt.Sprintf(`{"name":"%s", "description":"","enabled":true}`, correctAccount.Name)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteAccount request
func (s *ServerSuite) TestDeleteAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID)
	code, _, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccount request with a wrong id
func (s *ServerSuite) TestDeleteAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID+1)
	code, _, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccount request
func (s *ServerSuite) TestGetAccount_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	code, body, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(body), account)
	c.Assert(err, IsNil)
	c.Assert(account.ID, Equals, correctAccount.ID)
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
}

// Test a correct getAccount request with a wrong id
func (s *ServerSuite) TestGetAccount_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, correctAccount.ID+1)
	code, _, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}
