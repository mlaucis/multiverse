// +build !kinesis
// +build postgres

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/v02/entity"

	. "gopkg.in/check.v1"
)

func (s *ConnectionSuite) TestCreateConnectionAfterDisable(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	userFrom := application.Users[0]
	userTo := application.Users[1]

	LoginApplicationUser(accounts[0].ID, application.ID, userFrom)

	payload := fmt.Sprintf(
		`{"user_from_id":%d, "user_to_id":%d, "type": "friend"}`,
		userFrom.ID,
		userTo.ID,
	)

	routeName := "createConnection"
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
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)

	routeName = "deleteConnection"
	route = getComposedRoute(routeName, userTo.ID)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, userFrom, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	routeName = "createConnection"
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
	c.Assert(connection.Type, Equals, "friend")
	c.Assert(connection.Enabled, Equals, true)
}
