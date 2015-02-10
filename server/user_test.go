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

// Test a correct createUser request
func (s *ServerSuite) TestCreateUser_OK(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
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
}
