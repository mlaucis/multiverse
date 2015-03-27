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

	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator/keys"

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
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
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
	err = json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
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
	err = json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
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
	err = json.Unmarshal([]byte(body), receivedEvent)

	c.Assert(err, IsNil)
	c.Assert(receivedEvent.AccountID, Equals, account.ID)
	c.Assert(receivedEvent.ApplicationID, Equals, application.ID)
	c.Assert(receivedEvent.UserID, Equals, user.ID)
	c.Assert(receivedEvent.Enabled, Equals, true)
}

// Test a correct getEventList request
func (s *ServerSuite) TestGetEventList_OK(c *C) {
	c.Skip("this should be implemented properly")
	return

	correctAccount, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(correctAccount.ID, true)
	c.Assert(err, IsNil)

	correctUser, err := AddCorrectUser(correctAccount.ID, correctApplication.ID, true)
	c.Assert(err, IsNil)

	event1, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	event2, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	event3, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	event4, err := AddCorrectEvent(correctAccount.ID, correctApplication.ID, correctUser.ID, true)
	c.Assert(err, IsNil)

	routeName := "getEventList"
	route := getComposedRoute(routeName, correctAccount.ID, correctApplication.ID, correctUser.ID)
	code, body, err := runRequest(routeName, route, "", correctApplication.AuthToken, createApplicationUserSessionToken(correctUser), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	_, _, _, _ = event1, event2, event3, event4

	// TODO Check EventList body

	// event := &entity.Event{}
	// err = json.Unmarshal([]byte(body), event)

	// c.Assert(err, IsNil)
	// c.Assert(event.AccountID, Equals, correctAccount.ID)
	// c.Assert(event.ApplicationID, Equals, correctApplication.ID)
	// c.Assert(event.UserID, Equals, correctUser.ID)
	// c.Assert(event.Enabled, Equals, true)
}

// Test getEvent request with a wrong id
func (s *ServerSuite) TestGetEvent_WrongID(c *C) {
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

func (s *ServerSuite) TestGeoLocationSearch(c *C) {
	accounts := CorrectDeploy(1, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	events := []*entity.Event{
		CorrectEvent(),
		CorrectEvent(),
		CorrectEvent(),
		CorrectEvent(),
		CorrectEvent(),
		CorrectEvent(),
	}

	locations := []struct {
		Lat   float64
		Lon   float64
		Label string
	}{
		{Lat: 52.5169257, Lon: 13.3065105, Label: "dlsniper"},
		{Lat: 52.5148045, Lon: 13.3000390, Label: "gas"},
		{Lat: 52.5177294, Lon: 13.2938378, Label: "palace"},
		{Lat: 52.5104167, Lon: 13.3003824, Label: "ziko"},
		{Lat: 52.5120818, Lon: 13.3762879, Label: "cinestar"},
		{Lat: 52.5157576, Lon: 13.3873319, Label: "mercedes"},
	}

	routeName := "createEvent"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user.ID)

	for idx := range events {
		events[idx].Latitude = locations[idx].Lat
		events[idx].Longitude = locations[idx].Lon
		events[idx].Location = locations[idx].Label

		payload := fmt.Sprintf(
			`{"verb":"%s", "language":"%s", "location": "%s", "latitude": %.7f, "longitude": %.7f}`,
			events[idx].Verb,
			events[idx].Language,
			events[idx].Location,
			events[idx].Latitude,
			events[idx].Longitude,
		)

		code, body, err := runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
		c.Assert(err, IsNil)

		c.Assert(code, Equals, http.StatusCreated)

		c.Assert(body, Not(Equals), "")

		receivedEvent := &entity.Event{}
		err = json.Unmarshal([]byte(body), receivedEvent)
		c.Assert(err, IsNil)

		c.Assert(err, IsNil)
		c.Assert(receivedEvent.AccountID, Equals, application.AccountID)
		c.Assert(receivedEvent.ApplicationID, Equals, application.ID)
		c.Assert(receivedEvent.UserID, Equals, user.ID)
		c.Assert(receivedEvent.Enabled, Equals, true)
		c.Assert(receivedEvent.Verb, Equals, events[idx].Verb)
		c.Assert(receivedEvent.Language, Equals, events[idx].Language)
		c.Assert(receivedEvent.Latitude, Equals, events[idx].Latitude)
		c.Assert(receivedEvent.Longitude, Equals, events[idx].Longitude)
		c.Assert(receivedEvent.Location, Equals, events[idx].Location)
	}

	routeName = "getGeoEventList"
	route = getComposedRoute(routeName, application.AccountID, application.ID, locations[0].Lat, locations[0].Lon, 25000.0)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	var receivedEvents []*entity.Event
	err = json.Unmarshal([]byte(body), &receivedEvents)
	c.Assert(err, IsNil)

	expectedOrder := []string{"dlsniper", "gas", "ziko", "palace", "cinestar", "mercedes"}

	for idx := range receivedEvents {
		c.Assert(receivedEvents[idx].Location, Equals, expectedOrder[idx])
	}
}

func (s *ServerSuite) TestGetLocation(c *C) {
	accounts := CorrectDeploy(1, 1, 2, 7, true, true)
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
	err = json.Unmarshal([]byte(body), &events)
	c.Assert(err, IsNil)
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

	requestRoute := server.GetRoute(routeName, apiVersion)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	if err != nil {
		panic(err)
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

	requestRoute := server.GetRoute(routeName, apiVersion)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		nil,
	)
	if err != nil {
		panic(err)
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
