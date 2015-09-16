// +build !bench

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/multiverse/v02/entity"

	. "gopkg.in/check.v1"
)

// Test createAcccount request with a wrong key
func (s *AccountSuite) TestCreateAccount_WrongKey(c *C) {
	payload := "{namae:''}"

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(nil, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createAcccount request with an wrong name
func (s *AccountSuite) TestCreateAccount_WrongValue(c *C) {
	payload := `{"name":""}`

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(nil, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccount request
func (s *AccountSuite) TestCreateAccount_OK(c *C) {
	account := CorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, account.Name, account.Description)

	routeName := "createAccount"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(nil, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	if receivedAccount.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedAccount.ID, Not(Equals), "")
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct updateAccount request
func (s *AccountSuite) TestUpdateAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	if receivedAccount.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Description, Equals, description)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct updateAccount request with a wrong id
func (s *AccountSuite) TestUpdateAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.PublicID+"1")
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
}

// Test a correct updateAccount request with a wrong id
func (s *AccountSuite) TestUpdateAccountMalformedPayload(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true`, account.Name, description)

	routeName := "updateAccount"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}`+"\n")
}

// Test a correct deleteAccount request
func (s *AccountSuite) TestDeleteAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccount request with a wrong id
func (s *AccountSuite) TestDeleteAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "deleteAccount"
	route := getComposedRoute(routeName, account.PublicID+"1")
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":6003,"message":"account mismatch"}]}`+"\n")
}

// Test a correct getAccount request
func (s *AccountSuite) TestGetAccount_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedAccount := &entity.Account{}
	er := json.Unmarshal([]byte(body), receivedAccount)
	c.Assert(er, IsNil)
	c.Assert(receivedAccount.PublicID, Equals, account.PublicID)
	c.Assert(receivedAccount.Name, Equals, account.Name)
	c.Assert(receivedAccount.Enabled, Equals, true)
}

// Test a correct getAccount request with a wrong id
func (s *AccountSuite) TestGetAccount_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccount"
	route := getComposedRoute(routeName, account.PublicID+"a")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
}
