/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"

	"fmt"

	. "gopkg.in/check.v1"
)

// Test createAcccount request with a wrong key
func (s *ServerSuite) TestCreateAccount_WrongKey(c *C) {
	payload := "{namae:''}"

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createAcccount request with an wrong name
func (s *ServerSuite) TestCreateAccount_WrongValue(c *C) {
	payload := `{"name":""}`

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct createAccount request
func (s *ServerSuite) TestCreateAccount_OK(c *C) {
	correctAccount := utils.CorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, correctAccount.Name, correctAccount.Description)

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
	c.Assert(account.Token, Not(Equals), "")
}

// Test a correct updateAccount request
func (s *ServerSuite) TestUpdateAccount_OK(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Description, Equals, description)
	c.Assert(account.Enabled, Equals, true)
	//c.Assert(account.Token, Not(Equals), "")
}

// Test a correct updateAccount request with a wrong id
func (s *ServerSuite) TestUpdateAccount_WrongID(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID+1)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccount request with an invalid description
func (s *ServerSuite) TestUpdateAccount_Invalid(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"","enabled":true}`, correctAccount.Name)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteAccount request
func (s *ServerSuite) TestDeleteAccount_OK(c *C) {
	account, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID)
	w, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccount request with a wrong id
func (s *ServerSuite) TestDeleteAccount_WrongID(c *C) {
	account, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.ID+1)
	w, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccount request
func (s *ServerSuite) TestGetAccount_OK(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	c.Assert(account.ID, Equals, correctAccount.ID)
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
	c.Assert(account.Token, Not(Equals), "")
}

// Test a correct getAccount request with a wrong id
func (s *ServerSuite) TestGetAccount_WrongID(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, correctAccount.ID+1)
	w, err := runRequest(routeName, route, "", "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}
