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

// Test createUser request with a wrong key
func (s *ServerSuite) TestCreateUser_WrongKey(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	payload := "{usernamae:''}"

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createUser request with an wrong name
func (s *ServerSuite) TestCreateUser_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_name":""}`

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createUser request
func (s *ServerSuite) TestCreateUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	correctUser := CorrectUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"%s", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "auth_token": "%s"}`,
		correctUser.Username,
		correctUser.FirstName,
		correctUser.LastName,
		correctUser.Email,
		correctUser.URL,
		correctUser.Password,
		correctUser.AuthToken,
	)

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(body), user)
	c.Assert(err, IsNil)
	if user.ID < 1 {
		c.Fail()
	}

	c.Assert(err, IsNil)
	c.Assert(user.AccountID, Equals, correctAccount.ID)
	c.Assert(user.ApplicationID, Equals, correctApplication.ID)
	c.Assert(user.Username, Equals, correctUser.Username)
	c.Assert(user.Enabled, Equals, true)
}

// Test a correct updateUser request
func (s *ServerSuite) TestUpdateUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "auth_token": "%s", "enabled": true}`,
		correctUser.Username,
		correctUser.LastName,
		correctUser.Email,
		correctUser.URL,
		correctUser.Password,
		correctUser.AuthToken,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(body), user)
	c.Assert(err, IsNil)
	if user.ID < 1 {
		c.Fail()
	}
	c.Assert(err, IsNil)
	c.Assert(user.AccountID, Equals, correctAccount.ID)
	c.Assert(user.ApplicationID, Equals, correctApplication.ID)
	c.Assert(user.Username, Equals, correctUser.Username)
	c.Assert(user.Enabled, Equals, true)
}

// Test a correct updateUser request with a wrong id
func (s *ServerSuite) TestUpdateUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "auth_token": "%s", "enabled": true}`,
		correctUser.Username,
		correctUser.LastName,
		correctUser.Email,
		correctUser.URL,
		correctUser.Password,
		correctUser.AuthToken,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	code, _, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateUser request with an invalid name
func (s *ServerSuite) TestUpdateUser_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "auth_token": "%s", "enabled": true}`,
		correctUser.Username,
		correctUser.LastName,
		correctUser.Email,
		correctUser.URL,
		correctUser.Password,
		correctUser.AuthToken,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteUser request
func (s *ServerSuite) TestDeleteUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, _, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteUser request with a wrong id
func (s *ServerSuite) TestDeleteUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	code, _, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getUser request
func (s *ServerSuite) TestGetUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(body), user)
	c.Assert(err, IsNil)

	c.Assert(user.AccountID, Equals, correctAccount.ID)
	c.Assert(user.ApplicationID, Equals, correctApplication.ID)
	c.Assert(user.Username, Equals, correctUser.Username)
	c.Assert(user.Enabled, Equals, true)
}

// Test a correct getUser request with a wrong id
func (s *ServerSuite) TestGetUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	code, _, err := runRequest(routeName, route, "", correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct loginUser request
func (s *ServerSuite) TestLoginUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		correctUser.Email,
		"password",
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
}

// Test a correct logoutUser request
func (s *ServerSuite) TestLogoutUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		correctUser.Email,
		"password",
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err := runRequest(routeName, route, payload, correctApplication.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	routeName = "logoutUser"
	route = getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	code, body, err = runRequest(routeName, route, payload, correctApplication.AuthToken, getApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
}
