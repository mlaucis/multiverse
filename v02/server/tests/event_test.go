/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tapglue/backend/v01/validator/keys"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Test createEvent request with a wrong key
func (s *ServerSuite) TestCreateEvent_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := "{verbea:''}"

	routeName := "createEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createEvent request with an wrong name
func (s *ServerSuite) TestCreateEvent_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := `{"verb":"","language":""}`

	routeName := "createEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createEvent request
func (s *ServerSuite) TestCreateEvent_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event := CorrectEvent()
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s"}`,
		event.Verb,
		event.Language,
	)

	routeName := "createEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedEvent := &entity.Event{}
	er := json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(er, IsNil)
	c.Assert(receivedEvent.AccountID, Equals, account.ID)
	c.Assert(receivedEvent.ApplicationID, Equals, application.ID)
	c.Assert(receivedEvent.UserID, Equals, user.ID)
	c.Assert(receivedEvent.Enabled, Equals, true)
	c.Assert(receivedEvent.Verb, Equals, event.Verb)
	c.Assert(receivedEvent.Language, Equals, event.Language)
}

// Test a correct updateEvent request
func (s *ServerSuite) TestUpdateEvent_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false}`,
		event.Verb,
		event.Language,
	)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedEvent := &entity.Event{}
	er := json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(er, IsNil)
	c.Assert(receivedEvent.AccountID, Equals, account.ID)
	c.Assert(receivedEvent.ApplicationID, Equals, application.ID)
	c.Assert(receivedEvent.UserID, Equals, user.ID)
	c.Assert(receivedEvent.Enabled, Equals, false)
}

// Test updateEvent request with a wrong id
func (s *ServerSuite) TestUpdateEvent_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	correctEvent, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false}`,
		correctEvent.Verb,
		correctEvent.Language,
	)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, correctEvent.ID+1)
	code, _, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

func (s *ServerSuite) TestUpdateEventMalformedIDFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false}`,
		user.Events[0].Verb,
		user.Events[0].Language,
	)

	routeName := "updateEvent"
	route := getComposedRouteString(routeName, fmt.Sprintf("%d", user.AccountID), fmt.Sprintf("%d", user.ApplicationID), fmt.Sprintf("%d", user.ID), "90876543211234567890")
	code, body, err := runRequest(routeName, route, payload, accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to update the event (1)\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range")
}

func (s *ServerSuite) TestUpdateEventMalformedPayloadFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s", "enabled":false`,
		user.Events[0].Verb,
		user.Events[0].Language,
	)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, user.AccountID, user.ApplicationID, user.ID, user.Events[0].ID)
	code, body, err := runRequest(routeName, route, payload, accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to update the event (2)\nunexpected end of JSON input")
}

// Test updateEvent request with a wrong value
func (s *ServerSuite) TestUpdateEvent_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"verb":"", "language":"%s", "enabled":false}`,
		event.Language,
	)

	routeName := "updateEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct deleteEvent request
func (s *ServerSuite) TestDeleteEvent_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test deleteEvent request with a wrong id
func (s *ServerSuite) TestDeleteEvent_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID+1)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusInternalServerError)
}

func (s *ServerSuite) TestDeleteEventMalformedIDFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	routeName := "deleteEvent"
	route := getComposedRouteString(routeName, fmt.Sprintf("%d", user.AccountID), fmt.Sprintf("%d", user.ApplicationID), fmt.Sprintf("%d", user.ID), "90876543211234567890")
	code, body, err := runRequest(routeName, route, "", accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to delete the event (1)\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range")
}

// Test a correct getEvent request
func (s *ServerSuite) TestGetEvent_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	routeName := "getEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedEvent := &entity.Event{}
	er := json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(er, IsNil)
	c.Assert(receivedEvent.AccountID, Equals, account.ID)
	c.Assert(receivedEvent.ApplicationID, Equals, application.ID)
	c.Assert(receivedEvent.UserID, Equals, user.ID)
	c.Assert(receivedEvent.Enabled, Equals, true)
}

// Test a correct getEventList request
func (s *ServerSuite) TestGetEventList_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 5, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	routeName := "getEventList"
	route := getComposedRoute(routeName, user.AccountID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	events := []*entity.Event{}
	er := json.Unmarshal([]byte(body), &events)
	c.Assert(er, IsNil)
	c.Assert(len(events), Equals, len(user.Events))
	for idx := range events {
		c.Assert(events[idx], DeepEquals, user.Events[4-idx])
	}
}

// Test getEvent request with a wrong id
func (s *ServerSuite) TestGetEventWrongIDFails(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	c.Assert(err, IsNil)

	routeName := "getEvent"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID+1)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusInternalServerError)
}

func (s *ServerSuite) TestGetEventMalformedIDFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	routeName := "getEvent"
	route := getComposedRouteString(routeName, fmt.Sprintf("%d", user.AccountID), fmt.Sprintf("%d", user.ApplicationID), fmt.Sprintf("%d", user.ID), "90876543211234567890")
	code, body, err := runRequest(routeName, route, "", accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 read event failed (1)\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range")
}

func (s *ServerSuite) TestGeoLocationSearch(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 6, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	routeName := "getGeoEventList"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user.Events[0].Latitude, user.Events[0].Longitude, 25000.0)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	var receivedEvents []*entity.Event
	er := json.Unmarshal([]byte(body), &receivedEvents)
	c.Assert(er, IsNil)

	expectedOrder := []string{"dlsniper", "gas", "ziko", "palace", "cinestar", "mercedes"}

	for idx := range receivedEvents {
		c.Assert(receivedEvents[idx].Location, Equals, expectedOrder[idx])
	}
}

func (s *ServerSuite) TestGeoLocationInvalidSearchDataFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 6, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	routeName := "getGeoEventList"

	scenarios := []struct {
		Latitude     string
		Longitude    string
		Radius       string
		StatusCode   int
		ResponseBody string
	}{
		{
			Latitude:     fmt.Sprintf("%.7f", user.Events[0].Latitude),
			Longitude:    fmt.Sprintf("%.7f", user.Events[0].Longitude),
			Radius:       "-25000.0",
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to read the event by geo (4)\nLocation radius can't be smaller than 2 meters",
		},
		{
			Latitude:     "0.0.0",
			Longitude:    fmt.Sprintf("%.7f", user.Events[0].Longitude),
			Radius:       "25000",
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to read the event by geo (1)\nstrconv.ParseFloat: parsing \"0.0.0\": invalid syntax",
		},
		{
			Latitude:     fmt.Sprintf("%.7f", user.Events[0].Latitude),
			Longitude:    "0.0.0",
			Radius:       "25000",
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to read the event by geo (2)\nstrconv.ParseFloat: parsing \"0.0.0\": invalid syntax",
		},
		{
			Latitude:     fmt.Sprintf("%.7f", user.Events[0].Latitude),
			Longitude:    fmt.Sprintf("%.7f", user.Events[0].Longitude),
			Radius:       "0.0.0",
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to read the event by geo (3)\nstrconv.ParseFloat: parsing \"0.0.0\": invalid syntax",
		},
	}

	for idx := range scenarios {
		route := getComposedRouteString(routeName, fmt.Sprintf("%d", application.AccountID), fmt.Sprintf("%d", application.ID), scenarios[idx].Latitude, scenarios[idx].Longitude, scenarios[idx].Radius)
		code, body, err := runRequest(routeName, route, "", application.AuthToken, user.SessionToken, 3)
		c.Logf("pass: %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, scenarios[idx].StatusCode)
		c.Assert(body, Equals, scenarios[idx].ResponseBody)
	}
}

func (s *ServerSuite) TestGetLocation(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 7, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	routeName := "getLocationEventList"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.Events[0].Location)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	events := []*entity.Event{}
	er := json.Unmarshal([]byte(body), &events)
	c.Assert(er, IsNil)
	c.Assert(len(events), Equals, 2)
	c.Assert(events[0], DeepEquals, user2.Events[0])
	c.Assert(events[1], DeepEquals, user1.Events[0])
}

func (s *ServerSuite) TestGetObjectEvents(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 7, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	routeName := "getObjectEventList"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.Events[0].Object.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	events := []*entity.Event{}
	er := json.Unmarshal([]byte(body), &events)
	c.Assert(er, IsNil)
	c.Assert(len(events), Equals, 2)
	c.Assert(events[0], DeepEquals, user2.Events[0])
	c.Assert(events[1], DeepEquals, user1.Events[0])
}

func BenchmarkCreateEvent1_Write(b *testing.B) {
	account, err := AddCorrectAccount(true)
	if err != nil {
		panic(err)
	}
	application, err := AddCorrectApplication(account.ID, true)
	if err != nil {
		panic(err)
	}
	user, err := AddCorrectUser(account.ID, application.ID, true)
	if err != nil {
		panic(err)
	}
	event := CorrectEvent()

	payload := fmt.Sprintf(
		`{"verb":"%s", "language":"%s"}`,
		event.Verb,
		event.Language,
	)

	routeName := "createEvent"
	routePath := getComposedRoute(routeName, account.ID, application.ID, user.ID)

	requestRoute := server.GetRoute(routeName)

	req, er := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	if er != nil {
		panic(er)
	}

	createCommonRequestHeaders(req)
	if application.AuthToken != "" {
		err := keys.SignRequest(application.AuthToken, requestRoute.Scope, apiVersion, 2, req)
		if err != nil {
			panic(err)
		}
	}
	req.Header.Set("x-tapglue-session", createApplicationUserSessionToken(user))

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(requestRoute.RoutePattern(apiVersion), server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)

	for i := 1; i <= b.N; i++ {
		m.ServeHTTP(w, req)
	}
}

func BenchmarkCreateEvent2_Read(b *testing.B) {
	account, err := AddCorrectAccount(true)
	if err != nil {
		panic(err)
	}
	application, err := AddCorrectApplication(account.ID, true)
	if err != nil {
		panic(err)
	}
	user, err := AddCorrectUser(account.ID, application.ID, true)
	if err != nil {
		panic(err)
	}
	event, err := AddCorrectEvent(account.ID, application.ID, user.ID, true)
	if err != nil {
		panic(err)
	}

	routeName := "getEvent"
	routePath := getComposedRoute(routeName, account.ID, application.ID, user.ID, event.ID)

	requestRoute := server.GetRoute(routeName)

	req, er := http.NewRequest(
		requestRoute.Method,
		routePath,
		nil,
	)
	if er != nil {
		panic(er)
	}

	createCommonRequestHeaders(req)
	if application.AuthToken != "" {
		err := keys.SignRequest(application.AuthToken, requestRoute.Scope, apiVersion, 2, req)
		if err != nil {
			panic(err)
		}
	}
	req.Header.Set("x-tapglue-session", createApplicationUserSessionToken(user))

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(requestRoute.RoutePattern(apiVersion), server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, true)).
		Methods(requestRoute.Method)

	for i := 1; i <= b.N; i++ {
		m.ServeHTTP(w, req)
	}
}
