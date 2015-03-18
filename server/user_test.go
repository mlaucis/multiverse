/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/core/entity"

	. "gopkg.in/check.v1"
)

// Test createUser request with a wrong key
func (s *ServerSuite) TestCreateUser_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := "{usernamae:''}"

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createUser request with an wrong name
func (s *ServerSuite) TestCreateUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_name":""}`

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createUser request
func (s *ServerSuite) TestCreateUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user := CorrectUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"%s", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s"}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.User{}
	err = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(err, IsNil)
	if receivedUser.ID < 1 {
		c.Fail()
	}

	c.Assert(err, IsNil)
	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

// Test a correct updateUser request
func (s *ServerSuite) TestUpdateUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s", "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.User{}
	err = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(err, IsNil)
	if receivedUser.ID < 1 {
		c.Fail()
	}
	c.Assert(err, IsNil)
	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

// Test a correct updateUser request with a wrong id
func (s *ServerSuite) TestUpdateUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, _, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateUser request with an invalid name
func (s *ServerSuite) TestUpdateUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteUser request
func (s *ServerSuite) TestDeleteUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteUser request with a wrong id
func (s *ServerSuite) TestDeleteUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusUnauthorized)
}

// Test a correct getUser request
func (s *ServerSuite) TestGetUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.User{}
	err = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(err, IsNil)

	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

// Test a correct getUser request with a wrong id
func (s *ServerSuite) TestGetUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct loginUser request
func (s *ServerSuite) TestLoginUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"password",
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	err = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken = struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	err = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
}

// Test a correct logoutUser request
func (s *ServerSuite) TestLogoutUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"password",
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	err = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	routeName = "logoutUser"
	route = getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
}
