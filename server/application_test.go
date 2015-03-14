/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/core/entity"

	"github.com/tapglue/backend/utils"
	. "gopkg.in/check.v1"
)

// Test createApplication request with a wrong key
func (s *ServerSuite) TestCreateApplication_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := "{namae:''}"

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createApplication request with an wrong name
func (s *ServerSuite) TestCreateApplication_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := `{"name":""}`

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createApplication request
func (s *ServerSuite) TestCreateApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application := CorrectApplication()

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"%s", "url": "%s"}`,
		application.Name,
		application.Description,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	err = json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(err, IsNil)
	if receivedApplication.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.Description, Equals, application.Description)
	c.Assert(receivedApplication.URL, Equals, application.URL)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct updateApplication request
func (s *ServerSuite) TestUpdateApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	err = json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(err, IsNil)
	if receivedApplication.ID < 1 {
		c.Fail()
	}

	c.Assert(err, IsNil)
	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.URL, Equals, application.URL)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct updateApplication request with a wrong id
func (s *ServerSuite) TestUpdateApplication_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID+1)
	code, _, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateApplication request with an invalid description
func (s *ServerSuite) TestUpdateApplication_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct updateApplication request with a wrong token
func (s *ServerSuite) TestUpdateApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, sessionToken, 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusUnauthorized)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteApplication request
func (s *ServerSuite) TestDeleteApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteApplication request with a wrong id
func (s *ServerSuite) TestDeleteApplication_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	correctAccountUser, err := AddCorrectAccountUser(correctAccount.ID, true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID+1)
	code, _, err := runRequest(routeName, route, "", correctAccount.AuthToken, getAccountUserSessionToken(correctAccountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct deleteApplication request with a wrong token
func (s *ServerSuite) TestDeleteApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusUnauthorized)
}

// Test a correct getApplication request
func (s *ServerSuite) TestGetApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	err = json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(err, IsNil)

	c.Assert(receivedApplication.ID, Equals, application.ID)
	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.Description, Equals, application.Description)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct getApplication request with a wrong id
func (s *ServerSuite) TestGetApplication_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID+1)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getApplication request with a wrong token
func (s *ServerSuite) TestGetApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "getApplication"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusUnauthorized)
}
