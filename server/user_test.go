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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createUser request with an wrong name
func (s *ServerSuite) TestCreateUser_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_name":""}`

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(response), user)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(response), user)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteUser request
func (s *ServerSuite) TestDeleteUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test a correct deleteUser request with a wrong id
func (s *ServerSuite) TestDeleteUser_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct getUser request
func (s *ServerSuite) TestGetUser_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	user := &entity.User{}
	err = json.Unmarshal([]byte(response), user)
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

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID+1)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}
