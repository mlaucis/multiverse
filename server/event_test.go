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

// Test createEvent request with a wrong key
func (s *ServerSuite) TestCreateEvent_WrongKey(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := "{verbea:''}"

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createEvent request with an wrong name
func (s *ServerSuite) TestCreateEvent_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	payload := `{"verb":"","language":""}`

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct createEvent request
func (s *ServerSuite) TestCreateEvent_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent := CorrectEvent()
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s"}`,
		correctEvent.Verb,
		correctEvent.Language,
	)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "createEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	event := &entity.Event{}
	err = json.Unmarshal([]byte(response), event)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(event.AccountID, Equals, correctAccount.ID)
	c.Assert(event.ApplicationID, Equals, correctApplication.ID)
	c.Assert(event.UserID, Equals, correctUser.ID)
	c.Assert(event.Enabled, Equals, true)
}

// Test a correct updateEvent request
func (s *ServerSuite) TestUpdateEvent_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false}`,
		correctEvent.Verb,
		correctEvent.Language,
	)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	event := &entity.Event{}
	err = json.Unmarshal([]byte(response), event)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(event.AccountID, Equals, correctAccount.ID)
	c.Assert(event.ApplicationID, Equals, correctApplication.ID)
	c.Assert(event.UserID, Equals, correctUser.ID)
	c.Assert(event.Enabled, Equals, false)
}

// Test updateEvent request with a wrong id
func (s *ServerSuite) TestUpdateEvent_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false}`,
		correctEvent.Verb,
		correctEvent.Language,
	)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID+1)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test updateEvent request with a wrong value
func (s *ServerSuite) TestUpdateEvent_WrongValue(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"", "language":"%s", "enabled":false}`,
		correctEvent.Language,
	)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID)
	w, err := runRequest(routeName, route, payload, token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteEvent request
func (s *ServerSuite) TestDeleteEvent_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test deleteEvent request with a wrong id
func (s *ServerSuite) TestDeleteEvent_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "deleteEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID+1)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct getEvent request
func (s *ServerSuite) TestGetEvent_OK(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "getEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	event := &entity.Event{}
	err = json.Unmarshal([]byte(response), event)

	c.Assert(err, IsNil)
	c.Assert(event.AccountID, Equals, correctAccount.ID)
	c.Assert(event.ApplicationID, Equals, correctApplication.ID)
	c.Assert(event.UserID, Equals, correctUser.ID)
	c.Assert(event.Enabled, Equals, true)
}

// Test getEvent request with a wrong id
func (s *ServerSuite) TestGetEvent_WrongID(c *C) {
	correctAccount, err := AddCorrectAccount(true)
	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	correctEvent, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	token, err := storageClient.GenerateApplicationToken(correctApplication)
	c.Assert(err, IsNil)

	routeName := "getEvent"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID, correctEvent.ID+1)
	w, err := runRequest(routeName, route, "", token)
	c.Assert(err, IsNil)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}
