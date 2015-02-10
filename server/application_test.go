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

// Test createApplication request with a wrong key
func (s *ServerSuite) TestCreateApplication_WrongKey(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	payload := "{namae:''}"

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "createApplication"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createApplication request with an wrong name
func (s *ServerSuite) TestCreateApplication_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	payload := `{"name":""}`

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "createApplication"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct createApplication request
func (s *ServerSuite) TestCreateApplication_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication := CorrectApplication()

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"%s", "url": "%s"}`,
		correctApplication.Name,
		correctApplication.Description,
		correctApplication.URL,
	)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "createApplication"
	route := getComposedRoute(routeName, correctAccount.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	application := &entity.Application{}
	err = json.Unmarshal([]byte(response), application)
	c.Assert(err, IsNil)
	if application.ID < 1 {
		c.Fail()
	}
	c.Assert(application.Name, Equals, correctApplication.Name)
	c.Assert(application.Description, Equals, correctApplication.Description)
	c.Assert(application.URL, Equals, correctApplication.URL)
	c.Assert(application.Enabled, Equals, true)
}

// Test a correct updateApplication request
func (s *ServerSuite) TestUpdateApplication_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	application := &entity.Application{}
	err = json.Unmarshal([]byte(response), application)
	c.Assert(err, IsNil)
	if application.ID < 1 {
		c.Fail()
	}
	c.Assert(application.Name, Equals, correctApplication.Name)
	c.Assert(application.URL, Equals, correctApplication.URL)
	c.Assert(application.Enabled, Equals, true)
}

// Test a correct updateApplication request with a wrong id
func (s *ServerSuite) TestUpdateApplication_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID+1)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct updateApplication request with an invalid description
func (s *ServerSuite) TestUpdateApplication_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"name":"%s", "description":"", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct updateApplication request with a wrong token
func (s *ServerSuite) TestUpdateApplication_WrongToken(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, payload, "")
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteApplication request
func (s *ServerSuite) TestDeleteApplication_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, "", token)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test a correct deleteApplication request with a wrong id
func (s *ServerSuite) TestDeleteApplication_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID+1)
	w, err := runRequest(routeName, route, "", token)

	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct deleteApplication request with a wrong token
func (s *ServerSuite) TestDeleteApplication_WrongToken(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, "", "")

	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
}

// Test a correct getApplication request
func (s *ServerSuite) TestGetApplication_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, "", token)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	application := &entity.Application{}
	err = json.Unmarshal([]byte(response), application)
	c.Assert(err, IsNil)
	c.Assert(application.ID, Equals, correctApplication.ID)
	c.Assert(application.Name, Equals, correctApplication.Name)
	c.Assert(application.Description, Equals, correctApplication.Description)
	c.Assert(application.Enabled, Equals, true)
}

// Test a correct getApplication request with a wrong id
func (s *ServerSuite) TestGetApplication_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateAccountToken(correctAccount)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID+1)
	w, err := runRequest(routeName, route, "", token)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct getApplication request with a wrong token
func (s *ServerSuite) TestGetApplication_WrongToken(c *C) {
	c.Skip("To be refactored to use sessions")
	return

	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, correctApplication.AccountID, correctApplication.ID)
	w, err := runRequest(routeName, route, "", "")

	c.Assert(w.Code, Equals, http.StatusBadRequest)
}
