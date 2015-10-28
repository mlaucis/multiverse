// +build !bench

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/multiverse/v04/entity"

	. "gopkg.in/check.v1"
)

func (s *ApplicationUserSuite) TestLoginDisableLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	var zero int64 = 0
	user.FriendCount = &zero
	user.FollowerCount = &zero
	user.FollowedCount = &zero

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = `{"enabled": false}`

	user.Enabled = false
	routeName = "updateCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	compareUsers(c, user, updatedUser)

	payload = fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, `{"errors":[{"code":1001,"message":"application user not found"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestLoginDeleteLogoutFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	routeName = "deleteCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestLoginDeleteLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = `{"enabled": false}`

	routeName = "deleteCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, `{"errors":[{"code":1001,"message":"application user not found"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestLoginChangeUsernameLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	var zero int64 = 0
	user.FriendCount = &zero
	user.FollowerCount = &zero
	user.FollowedCount = &zero

	payload := fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	user.Username = "newUserName"
	payload = fmt.Sprintf(`{"user_name": "%s"}`, user.Username)
	routeName = "updateCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	compareUsers(c, user, updatedUser)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(
		`{"user_name": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	newSessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)
	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginChangeEmailLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	var zero int64 = 0
	user.FriendCount = &zero
	user.FollowerCount = &zero
	user.FollowedCount = &zero

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	user.Email = "newUserEmail@tapglue.com"
	payload = fmt.Sprintf(`{"email": "%s"}`, user.Email)
	routeName = "updateCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	compareUsers(c, user, updatedUser)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	newSessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)

	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginChangePasswordLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	var zero int64 = 0
	user.FriendCount = &zero
	user.FollowerCount = &zero
	user.FollowedCount = &zero

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"password": "%s"}`, "newPass")

	routeName = "updateCurrentApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	compareUsers(c, user, updatedUser)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	newSessionToken := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)
	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestDeleteOnEventsOnUserDeleteWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	// GET EVENT
	routeName := "deleteCurrentUserConnection"
	route := getComposedRoute(routeName, entity.ConnectionTypeFollow, user2.ID)
	code, body, err := runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(body, Equals, "\"\"\n")

	// GET EVENTS LIST
	routeName = "getEventList"
	route = getComposedRoute(routeName, user1.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	response := struct {
		Events      []entity.Event `json:"events"`
		EventsCount int            `json:"events_count"`
	}{}
	er := json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)
	c.Assert(response.EventsCount, Equals, 2)
	compareEvents(c, user1.Events[1], &response.Events[0])
	compareEvents(c, user1.Events[0], &response.Events[1])

	// Check connetions list
	routeName = "getUserFollowers"
	route = getComposedRoute(routeName, user1.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	followers := struct {
		Users      []*entity.ApplicationUser `json:"users"`
		UsersCount int                       `json:"users_count"`
	}{}
	er = json.Unmarshal([]byte(body), &followers)
	c.Assert(er, IsNil)
	c.Assert(followers.UsersCount, Equals, 1)
	c.Assert(followers.Users[0].Username, DeepEquals, user2.Username)

	routeName = "getUserFollows"
	route = getComposedRoute(routeName, user1.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	// GET EVENTS LIST
	routeName = "getCurrentUserFeed"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *ApplicationUserSuite) TestCreateAndLoginExistingUserTwice_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, false, false)
	application := accounts[0].Applications[0]

	user := CorrectUser()

	payload := fmt.Sprintf(
		`{"user_name": %q, "first_name": %q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &struct {
		entity.ApplicationUser
		SessionToken string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.SessionToken, Not(Equals), "")

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedUser1 := &struct {
		entity.ApplicationUser
		SessionToken string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), receivedUser1)
	c.Assert(er, IsNil)
	c.Assert(receivedUser1.ID, Equals, receivedUser.ID)
	c.Assert(receivedUser1.Username, Equals, receivedUser.Username)
	c.Assert(receivedUser1.SessionToken, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestCreateAndLoginExistingUserTwiceDifferentPasswordFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, false, false)
	application := accounts[0].Applications[0]

	user := CorrectUser()

	payload := fmt.Sprintf(
		`{"user_name": %q, "first_name": %q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &struct {
		entity.ApplicationUser
		SessionToken string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.SessionToken, Not(Equals), "")

	payload = fmt.Sprintf(
		`{"user_name": %q, "first_name": %q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.URL,
		user.Password+"as",
	)

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
}
