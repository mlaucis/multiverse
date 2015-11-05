// +build !bench

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/server"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

func (s *ApplicationUserSuite) TestCreateUser_WrongKey(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	payload := "{usernamae:''}"

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestCreateUser_WrongValue(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	payload := `{"user_name":""}`

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestCreateUser_NoPassword(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	payload := `{"user_name": "dlsniper"}`

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestCreateUser_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	user := CorrectUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"%s", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s"}`,
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

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.SessionToken, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestCreateUserBareDetailsOK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	user := CorrectUser()

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)

	payload := fmt.Sprintf(
		`{"user_name":%q, "password": %q}`,
		user.Username,
		user.Password,
	)

	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Username, Equals, user.Username)

	payload = fmt.Sprintf(
		`{"email": %q,  "password": %q}`,
		user.Email,
		user.Password,
	)

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedUser = &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Email, Equals, user.Email)
}

func (s *ApplicationUserSuite) TestCreateAndLoginUser_OK(c *C) {
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
}

func (s *ApplicationUserSuite) TestCreateAndLoginExistingUser_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name": %q, "first_name": %q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.URL,
		user.OriginalPassword,
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
	c.Assert(receivedUser.ID, Equals, user.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.SessionToken, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestUpdateUser_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s", "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	c.Assert(receivedUser.Username, Equals, user.Username)
}

func (s *ApplicationUserSuite) TestUpdateUser_WrongID(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, _, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, false))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestUpdateUser_WrongValue(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(
		`{"user_name":%q, "first_name":%q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q,  "enabled": %t}`,
		user.Username,
		user.FirstName,
		user.LastName,
		user2.Email,
		user.URL,
		user.Password,
		user.Enabled,
	)

	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestUpdateUserMalformedPayloadFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	payload := fmt.Sprintf(`{"user_name":"%s"`, user.Username)

	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(accounts[0].Applications[0], user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestDeleteUser_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 3, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	user2 := application.Users[1]
	user3 := application.Users[2]

	routeName := "deleteCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, _, err := runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	routeName = "deleteApplicationUser"
	route = getComposedRoute(routeName, user3.ID)
	code, _, err = runRequest(routeName, route, "", signApplicationRequest(application, user2, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *ApplicationUserSuite) TestGetUser_OK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 10, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	user2 := application.Users[1]
	user10 := application.Users[9]

	routeName := "getCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.IsFriend, IsNil)
	c.Assert(receivedUser.IsFollower, IsNil)
	c.Assert(receivedUser.IsFollowed, IsNil)
	c.Assert(receivedUser.CreatedAt, Not(IsNil))
	c.Assert(receivedUser.UpdatedAt, Not(IsNil))

	routeName = "getApplicationUser"
	route = getComposedRoute(routeName, user2.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	receivedUser = &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.Username, Equals, user2.Username)
	c.Assert(*receivedUser.IsFriend, Equals, false)
	c.Assert(*receivedUser.IsFollower, Equals, true)
	c.Assert(*receivedUser.IsFollowed, Equals, true)
	c.Assert(receivedUser.CreatedAt, IsNil)
	c.Assert(receivedUser.UpdatedAt, IsNil)
	c.Assert(strings.Contains(body, `created_at":null`), Equals, false)
	c.Assert(strings.Contains(body, `updated_at":null`), Equals, false)

	routeName = "getApplicationUser"
	route = getComposedRoute(routeName, user10.ID)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	receivedUser = &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.Username, Equals, user10.Username)
	c.Assert(*receivedUser.IsFriend, Equals, false)
	c.Assert(*receivedUser.IsFollower, Equals, false)
	c.Assert(*receivedUser.IsFollowed, Equals, false)
	c.Assert(receivedUser.CreatedAt, IsNil)
	c.Assert(receivedUser.UpdatedAt, IsNil)
	c.Assert(strings.Contains(body, `created_at":null`), Equals, false)
	c.Assert(strings.Contains(body, `updated_at":null`), Equals, false)
}

func (s *ApplicationUserSuite) TestGetUser_WrongID(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	routeName := "getApplicationUser"
	route := getComposedRoute(routeName, user.ID+1)
	code, _, err := runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestLoginUserWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.ID, Equals, user.ID)
	c.Assert(receivedUser.SessionToken, Not(Equals), "")
	c.Assert(receivedUser.Email, Equals, user.Email)
}

func (s *ApplicationUserSuite) TestLoginUserWorksWithUsernameOrEmail(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	routeName := "loginCurrentUserApplicationUser"
	route := getComposedRoute(routeName)

	payload := fmt.Sprintf(
		`{"username": %q, "password": %q}`,
		user.Email,
		user.OriginalPassword,
	)

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

	payload = fmt.Sprintf(
		`{"username": %q, "password": %q}`,
		user.Username,
		user.OriginalPassword,
	)

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken = struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
}

func (s *ApplicationUserSuite) TestLoginUserWithDetails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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
		entity.ApplicationUser
		Token string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.ID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.Email, Equals, user.Email)
	c.Assert(sessionToken.Password, Equals, "")
}

func (s *ApplicationUserSuite) TestRefreshSessionOnOriginalTokenFailsAfterDoubleUserLogin(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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
	c.Assert(sessionToken.Token, Not(Equals), user.SessionToken)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, user.SessionToken)

	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
}

func (s *ApplicationUserSuite) TestLoginUserAfterLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	initialToken := sessionToken.Token

	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken = struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.Token, Not(Equals), initialToken)
}

func (s *ApplicationUserSuite) TestLoginAndRefreshSessionWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	payload = fmt.Sprintf(`{"session_token": "%s"}`, user.SessionToken)

	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	refreshSessionToken := struct {
		Token string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &refreshSessionToken)
	c.Assert(er, IsNil)
	c.Assert(refreshSessionToken.Token, Not(Equals), "")
	c.Assert(refreshSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginRefreshSessionLogoutWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)

	user.SessionToken = sessionToken.Token

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, updatedToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *ApplicationUserSuite) TestLogoutUserWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, "", signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *ApplicationUserSuite) TestLoginLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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
	c.Assert(err, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)
	routeName = "loginCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	newSession := struct {
		UserID uint64 `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSession)
	c.Assert(er, IsNil)
	c.Assert(newSession.UserID, Equals, user.ID)
	c.Assert(newSession.Token, Not(Equals), "")
	c.Assert(newSession.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestRefreshSessionWithoutLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]

	// REFRESH USER SESSION
	payload := fmt.Sprintf(`{"session_token": "%s"}`, "random session token stuff")
	routeName := "refreshCurrentUserApplicationUserSession"
	route := getComposedRoute(routeName)
	code, _, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestLoginLogoutRefreshSessionFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestLoginChangePasswordRefreshWorks(c *C) {
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

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginChangeUsernameRefreshWorks(c *C) {
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
	c.Assert(err, IsNil)

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

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginChangeEmailRefreshWorks(c *C) {
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

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshCurrentUserApplicationUserSession"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ApplicationUserSuite) TestLoginLogoutLogoutFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

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

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutCurrentUserApplicationUser"
	route = getComposedRoute(routeName)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationUserSuite) TestLoginChangeUsernameGetEventWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	event := user.Events[0]

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

	// GET EVENT
	routeName = "getEvent"
	route = getComposedRoute(routeName, user.ID, event.ID)
	code, body, err = runRequest(routeName, route, payload, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	receivedEvent := &entity.Event{}
	er = json.Unmarshal([]byte(body), receivedEvent)
	c.Assert(er, IsNil)
	compareEvents(c, event, receivedEvent)
}

func (s *ApplicationUserSuite) TestLoginChangeUsernameExistingUsernameFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"user_name": "%s"}`, user2.Username)
	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":1008,"message":"username already in use"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestLoginChangeUsernameSameUsernameFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"user_name": "%s"}`, user2.Username)
	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":1008,"message":"username already in use"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestLoginChangeEmailExistingEmailFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"email": "%s"}`, user2.Email)
	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":1002,"message":"email address already in use"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestLoginChangeEmailSameEmailFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"email": "%s"}`, user2.Email)
	routeName := "updateCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, user1, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":1002,"message":"email address already in use"}]}`+"\n")
}

func (s *ApplicationUserSuite) TestCreateUserAutoBindSocialAccounts(c *C) {
	c.Skip("We've decided not to automatically create connections on user creation, for now at least")
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]

	user2 := CorrectUserWithDefaults(application.OrgID, application.ID, 2)
	user2.Enabled = true
	user2.Activated = true
	user2.SocialConnectionsIDs = map[string][]string{
		"facebook": []string{user1.SocialIDs["facebook"]},
	}

	payloadByte, err := json.Marshal(user2)
	c.Assert(err, IsNil)
	payload := string(payloadByte)

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, er := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(er, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(err, IsNil)
	c.Assert(receivedUser.ID, Not(Equals), 0)
	user2.OriginalPassword, receivedUser.OriginalPassword = user2.Password, user2.Password
	user2.Password = ""
	user2.CreatedAt = receivedUser.CreatedAt
	user2.UpdatedAt = receivedUser.UpdatedAt
	user2.LastLogin = receivedUser.LastLogin
	user2.ID = receivedUser.ID
	receivedUser.Images, user2.Images = nil, nil
	c.Assert(receivedUser, DeepEquals, user2)

	// Check connetions list
	routeName = "getUserFollowed"
	route = getComposedRoute(routeName, user1.ID)
	code, body, er = runRequest(routeName, route, "", signApplicationRequest(application, user1, true, true))
	c.Assert(er, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "[]\n")

	userConnections := []entity.ApplicationUser{}
	err = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(err, IsNil)

	c.Assert(len(userConnections), Equals, 1)
	c.Assert(userConnections[0].ID, Equals, receivedUser.ID)
}

func (s *ApplicationUserSuite) TestLoginRefreshLogoutMalformedPayloadFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	iterations := []struct {
		Payload   string
		RouteName string
		Route     string
		Code      int
		Body      string
	}{
		// 0
		{
			Payload:   fmt.Sprintf(`{"email": "%s", "password": "%s"`, user.Email, user.OriginalPassword),
			RouteName: "loginCurrentUserApplicationUser",
			Route:     getComposedRoute("loginCurrentUserApplicationUser"),
			Code:      http.StatusBadRequest,
			Body:      `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 1
		{
			Payload:   fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "tap@glue", user.OriginalPassword),
			RouteName: "loginCurrentUserApplicationUser",
			Route:     getComposedRoute("loginCurrentUserApplicationUser"),
			Code:      http.StatusNotFound,
			Body:      `{"errors":[{"code":1001,"message":"application user not found"}]}` + "\n",
		},
		// 2
		{
			Payload:   fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, "", user.OriginalPassword),
			RouteName: "loginCurrentUserApplicationUser",
			Route:     getComposedRoute("loginCurrentUserApplicationUser"),
			Code:      http.StatusBadRequest,
			Body:      `{"errors":[{"code":4003,"message":"both username and email are empty"}]}` + "\n",
		},
		// 3
		{
			Payload:   fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, "tapg", user.OriginalPassword),
			RouteName: "loginCurrentUserApplicationUser",
			Route:     getComposedRoute("loginCurrentUserApplicationUser"),
			Code:      http.StatusNotFound,
			Body:      `{"errors":[{"code":1001,"message":"application user not found"}]}` + "\n",
		},
		// 4
		{
			Payload:   fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, user.Username, "nothing"),
			RouteName: "loginCurrentUserApplicationUser",
			Route:     getComposedRoute("loginCurrentUserApplicationUser"),
			Code:      http.StatusBadRequest,
			Body:      `{"errors":[{"code":4001,"message":"authentication error"}]}` + "\n",
		},
		// 5
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"`, user.SessionToken),
			RouteName: "refreshCurrentUserApplicationUserSession",
			Route:     getComposedRoute("refreshCurrentUserApplicationUserSession"),
			Code:      http.StatusBadRequest,
			Body:      `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 6
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"}`, "nothing"),
			RouteName: "refreshCurrentUserApplicationUserSession",
			Route:     getComposedRoute("refreshCurrentUserApplicationUserSession"),
			Code:      http.StatusBadRequest,
			Body:      `{"errors":[{"code":4012,"message":"session token mismatch"}]}` + "\n",
		},
	}

	for idx := range iterations {
		code, body, err := runRequest(
			iterations[idx].RouteName,
			iterations[idx].Route,
			iterations[idx].Payload,
			signApplicationRequest(application, user, true, true))
		c.Logf("pass %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, iterations[idx].Code)
		c.Assert(body, Equals, iterations[idx].Body)
	}
}

func (s *ApplicationUserSuite) TestSearch(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 10, 0, true, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]
	user2 := application.Users[1]
	user3 := application.Users[2]
	user3.Enabled = false
	UpdateUser(accounts[0].ID, application.ID, *user3)

	user.Deleted, user2.Deleted, user3.Deleted = nil, nil, nil
	user.CreatedAt, user2.CreatedAt, user3.CreatedAt = nil, nil, nil
	user.UpdatedAt, user2.UpdatedAt, user3.UpdatedAt = nil, nil, nil
	user2.IsFriend, user2.IsFollower, user2.IsFollowed = entity.PFalse, entity.PTrue, entity.PTrue

	iterations := []struct {
		Payload   string
		RouteName string
		Route     string
		Code      int
		Response  []*entity.ApplicationUser
	}{
		//0
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=dlsniper",
			Code:      http.StatusNoContent,
			Response:  []*entity.ApplicationUser{},
		},
		// 1
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=florin@tapglue.com",
			Code:      http.StatusNoContent,
			Response:  []*entity.ApplicationUser{},
		},
		// 2
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=" + user2.Username,
			Code:      http.StatusOK,
			Response:  []*entity.ApplicationUser{user2},
		},
		// 3
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=" + user2.Email,
			Code:      http.StatusOK,
			Response:  []*entity.ApplicationUser{user2},
		},
		// 4
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=" + user3.Email,
			Code:      http.StatusNoContent,
			Response:  []*entity.ApplicationUser{},
		},
		// 5
		{
			Payload:   "",
			RouteName: "searchApplicationUser",
			Route:     getComposedRoute("searchApplicationUser") + "?q=" + user.Email,
			Code:      http.StatusNoContent,
			Response:  []*entity.ApplicationUser{},
		},
	}

	for idx := range iterations {
		code, body, err := runRequest(
			iterations[idx].RouteName,
			iterations[idx].Route,
			iterations[idx].Payload,
			signApplicationRequest(application, user, true, true))
		c.Logf("pass %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, iterations[idx].Code)
		response := struct {
			Users      []*entity.ApplicationUser `json:"users"`
			UsersCount int                       `json:"users_count"`
		}{}
		c.Logf("response body: %s", body)
		er := json.Unmarshal([]byte(body), &response)
		c.Assert(er, IsNil)
		if response.UsersCount > 0 {
			compareUsers(c, iterations[idx].Response[0], response.Users[0])
		} else {
			c.Assert(response.Users, DeepEquals, iterations[idx].Response)
		}
	}
}

func (s *ApplicationUserSuite) TestGetUserWithoutSessionFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]

	routeName := "getCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, "", signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "{\"errors\":[{\"code\":4013,\"message\":\"session token missing from request\"}]}\n")
}

func (s *ApplicationUserSuite) TestCreateUser_NullOK(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 0, 0, true, true)
	application := accounts[0].Applications[0]

	payload := `{"custom_id":"demo","first_name":"firstName1","last_name":"lastName2","images":{"profile": null}, "metadata":{"public":false}}`

	routeName := "createApplicationUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signApplicationRequest(application, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)

	c.Assert(body, Not(Equals), "")
}

func BenchmarkCreateAndLogin(b *testing.B) {
	var (
		accounts     = CorrectDeploy(1, 0, 1, 0, 0, false, false)
		application  = accounts[0].Applications[0]
		routeName    = "createApplicationUser"
		requestRoute = getRoute(routeName)
		routePath    = getComposedRoute(routeName)
		user         = CorrectUser()
		signer       = signApplicationRequest(application, nil, true, true)
		m            = mux.NewRouter()
		w            = httptest.NewRecorder()
		requests     = []*http.Request{}
	)

	for i := 0; i < b.N; i++ {
		payload := fmt.Sprintf(
			`{"user_name": %q, "first_name": %q, "last_name": %q,  "email": %q,  "url": %q,  "password": %q}`,
			fmt.Sprintf("%d-%s", i, user.Username),
			user.FirstName,
			user.LastName,
			fmt.Sprintf("%d-%s", i, user.Email),
			user.URL,
			user.Password,
		)

		req, err := http.NewRequest(
			requestRoute.Method,
			routePath,
			strings.NewReader(payload),
		)
		if err != nil {
			b.Fatal(err)
		}

		createCommonRequestHeaders(req)
		signer(req)

		requests = append(requests, req)
	}

	m.HandleFunc(
		requestRoute.RoutePattern(),
		server.CustomHandler(
			requestRoute,
			mainLogChan,
			errorLogChan,
			"test",
			false,
			true,
		),
	).Methods(requestRoute.Method)

	b.ResetTimer()

	for _, req := range requests {
		m.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("wrong response %d with body %s", w.Code, w.Body.String())
		}
	}
}
