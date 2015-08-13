// +build !kinesis
// +build postgres

package server_test

import (
	"net/http"

	. "gopkg.in/check.v1"
)

// Test deleteEvent request with a wrong id
func (s *EventSuite) TestDeleteEvent_WrongID(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	event := user.Events[0]

	routeName := "deleteCurrentUserEvent"
	route := getComposedRoute(routeName, event.ID+1)
	code, _, err := runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}
