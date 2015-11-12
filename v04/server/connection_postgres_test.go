// +build !bench

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/v04/entity"

	. "gopkg.in/check.v1"
)

func (s *ConnectionSuite) TestCreateConnectionAfterDisable(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": %q}`,
		userFrom.ID,
		userTo.ID,
		entity.ConnectionTypeFriend,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, entity.ConnectionTypeFriend)

	routeName = "deleteCurrentUserConnection"
	route = getComposedRoute(routeName, entity.ConnectionTypeFriend, userTo.ID)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	routeName = "createCurrentUserConnection"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection = &entity.Connection{}
	er = json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, entity.ConnectionTypeFriend)
}

func (s *ConnectionSuite) TestCreateFollowConnectionWithUserGeneratedEvent(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": %q}`,
		userFrom.ID,
		userTo.ID,
		entity.ConnectionTypeFollow,
	)

	routeName := "createCurrentUserConnection"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	connection := &entity.Connection{}
	er := json.Unmarshal([]byte(body), connection)
	c.Assert(er, IsNil)
	c.Assert(connection.UserFromID, Equals, userFrom.ID)
	c.Assert(connection.UserToID, Equals, userTo.ID)
	c.Assert(connection.Type, Equals, entity.ConnectionTypeFollow)

	payload = fmt.Sprintf(
		`{"type":%q, "target":{"id": "%d", "type": "tg_user"}, "visibility": %d}`,
		"tg_follow",
		userTo.ID,
		entity.EventConnections,
	)

	routeName = "createCurrentUserEvent"
	route = getComposedRoute(routeName)
	code, body, headerz, err := runRequestWithHeaders(routeName, route, payload, func(*http.Request) {}, signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(headerz.Get("Location"), Not(Equals), "")
	c.Assert(headerz.Get("Content-Type"), Equals, "application/json; charset=UTF-8")
	c.Assert(body, Not(Equals), "")

	receivedEvent := &entity.Event{}
	er = json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(er, IsNil)
	c.Assert(receivedEvent.ID, Not(Equals), "")
	c.Assert(receivedEvent.UserID, Equals, userFrom.ID)
	c.Assert(receivedEvent.Enabled, Equals, true)
	c.Assert(receivedEvent.Type, Equals, "tg_follow")
	c.Assert(receivedEvent.Target.ID.(string), Equals, strconv.FormatUint(userTo.ID, 10))
	c.Assert(receivedEvent.Target.Type, Equals, "tg_user")
	c.Assert(int(receivedEvent.Visibility), Equals, entity.EventConnections)

	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, userTo, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := struct {
		Count  int                               `json:"unread_events_count"`
		Events []entity.Event                    `json:"events"`
		Users  map[string]entity.ApplicationUser `json:"users"`
	}{}
	er = json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)

	c.Assert(response.Count, Equals, 1)
	c.Assert(len(response.Events), Equals, 1)
	c.Assert(len(response.Users), Equals, 1)
	c.Assert(response.Events[0].Type, Equals, "tg_follow")
}
