/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/v02/entity"

	. "gopkg.in/check.v1"
)

func (s *ServerSuite) TestCreateUser_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := "{usernamae:''}"

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ServerSuite) TestCreateUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := `{"user_name":""}`

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ServerSuite) TestCreateUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

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

	routeName := "createUser"
	route := getComposedRoute(routeName, account.ID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	if receivedUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

func (s *ServerSuite) TestUpdateUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s", "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	if receivedUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

func (s *ServerSuite) TestUpdateUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"changed", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, _, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

func (s *ServerSuite) TestUpdateUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "first_name":"", "last_name": "%s",  "email": "%s",  "url": "%s",  "password": "%s",  "enabled": true}`,
		user.Username,
		user.LastName,
		user.Email,
		user.URL,
		user.Password,
	)

	routeName := "updateUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *ServerSuite) TestUpdateUserMalformedPayloadFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	payload := fmt.Sprintf(`{"user_name":"%s"`, user.Username)

	routeName := "updateUser"
	route := getComposedRoute(routeName, user.AccountID, user.ApplicationID, user.ID)
	code, body, err := runRequest(routeName, route, payload, accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to update the user (1)\nunexpected end of JSON input")
}

func (s *ServerSuite) TestDeleteUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *ServerSuite) TestDeleteUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (9)")
}

func (s *ServerSuite) TestDeleteUserInvalidID(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, true, true)
	user := accounts[0].Applications[0].Users[0]

	routeName := "deleteUser"
	route := getComposedRouteString(routeName, fmt.Sprintf("%d", user.AccountID), fmt.Sprintf("%d", user.ApplicationID), "90876543211234567890")
	code, body, err := runRequest(routeName, route, "", accounts[0].Applications[0].AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to parse application user id\nstrconv.ParseInt: parsing \"90876543211234567890\": value out of range")
}

func (s *ServerSuite) TestGetUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	er := json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.AccountID, Equals, account.ID)
	c.Assert(receivedUser.ApplicationID, Equals, application.ID)
	c.Assert(receivedUser.Username, Equals, user.Username)
	c.Assert(receivedUser.Enabled, Equals, true)
}

func (s *ServerSuite) TestGetUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	user, err := AddCorrectUser(account.ID, application.ID, true)
	c.Assert(err, IsNil)

	routeName := "getUser"
	route := getComposedRoute(routeName, account.ID, application.ID, user.ID+1)
	code, _, err := runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

func (s *ServerSuite) TestLoginUserWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
}

func (s *ServerSuite) TestRefreshSessionOnOriginalTokenFailsAfterDoubleUserLogin(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.Token, Not(Equals), user.SessionToken)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, user.SessionToken)

	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (12)\nsession mismatch")
}

func (s *ServerSuite) TestLoginUserAfterLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	initialToken := sessionToken.Token

	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken = struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.Token, Not(Equals), initialToken)
}

func (s *ServerSuite) TestLoginAndRefreshSessionWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"session_token": "%s"}`, user.SessionToken)

	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
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

func (s *ServerSuite) TestLoginChangePasswordLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"password": "%s"}`, "newPass")

	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	newSessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)
	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginRefreshSessionLogoutWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, updatedToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, updatedToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "\"logged out\"\n")
}

func (s *ServerSuite) TestLogoutUserWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
}

func (s *ServerSuite) TestLoginLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)
	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	newSession := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSession)
	c.Assert(er, IsNil)
	c.Assert(newSession.UserID, Equals, user.ID)
	c.Assert(newSession.Token, Not(Equals), "")
	c.Assert(newSession.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginChangeUsernameLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_name": "%s"}`, "newUserName")
	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Username, Equals, "newUserName")
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Username = "newUserName"
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	newSessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)
	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginChangeEmailLogoutLoginWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"email": "%s"}`, "newUserEmail@tapglue.com")
	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Email, Equals, "newUserEmail@tapglue.com")
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Email = "newUserEmail@tapglue.com"
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	newSessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er = json.Unmarshal([]byte(body), &newSessionToken)
	c.Assert(er, IsNil)

	c.Assert(newSessionToken.UserID, Equals, user.ID)
	c.Assert(newSessionToken.Token, Not(Equals), "")
	c.Assert(newSessionToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginDisableLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = `{"enabled": false}`

	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Enabled, Equals, false)
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Enabled = false
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, "404 failed to login the user (3)\nuser is disabled")
}

func (s *ServerSuite) TestLoginDeleteLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = `{"enabled": false}`

	routeName = "deleteUser"
	route = getComposedRoute(routeName, application.ID, application.ID, user.ID)
	code, _, err = runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName = "loginUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, "404 failed to login the user (3)\nuser is disabled")
}

func (s *ServerSuite) TestRefreshSessionWithoutLoginFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, "random session token stuff")
	routeName := "refreshUserSession"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "random session stuff", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (2)")

}

func (s *ServerSuite) TestLoginLogoutRefreshSessionFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "logged out")

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (10)")
}

func (s *ServerSuite) TestLoginChangePasswordRefreshWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"password": "%s"}`, "newPass")

	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginChangeUsernameRefreshWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_name": "%s"}`, "newUserName")
	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Username, Equals, "newUserName")
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Username = "newUserName"
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginChangeEmailRefreshWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"email": "%s"}`, "newUserEmail@tapglue.com")
	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Email, Equals, "newUserEmail@tapglue.com")
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	user.Password = ""
	user.Events = nil
	user.Image = nil
	user.Email = "newUserEmail@tapglue.com"
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	payload = fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		"newPass",
	)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "refreshUserSession"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := sessionToken
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.UserID, Equals, sessionToken.UserID)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
}

func (s *ServerSuite) TestLoginRefreshDifferentUserFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	// REFRESH USER SESSION
	payload := fmt.Sprintf(`{"session_token": "%s"}`, user1.SessionToken)
	routeName := "refreshUserSession"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user2.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (9)")

}

func (s *ServerSuite) TestLoginLogoutLogoutFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "\"logged out\"\n")

	// LOGOUT USER
	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "logged out")
}

func (s *ServerSuite) TestLoginLogoutDifferentUserFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, false)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	// LOGOUT USER
	payload := fmt.Sprintf(`{"session_token": "%s"}`, user1.SessionToken)
	routeName := "logoutUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user2.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "\"logged out\"\n")
}

func (s *ServerSuite) TestLoginChangeUsernameGetEventWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 1, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	payload = fmt.Sprintf(`{"user_name": "%s"}`, "newUserName")
	routeName = "updateUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, user.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	updatedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), updatedUser)
	c.Assert(er, IsNil)
	c.Assert(updatedUser.Username, Equals, "newUserName")
	// WE need these to make DeepEquals work
	updatedUser.SessionToken = user.SessionToken
	updatedUser.OriginalPassword = user.OriginalPassword
	updatedUser.Image = nil
	updatedUser.LastLogin = user.LastLogin
	updatedUser.Events = user.Events
	user.Password = ""
	user.Image = nil
	user.Username = "newUserName"
	user.Activated = true
	c.Assert(updatedUser, DeepEquals, user)

	// GET EVENT
	routeName = "getEvent"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID, user.Events[0].ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	event := &entity.Event{}
	er = json.Unmarshal([]byte(body), &event)
	c.Assert(er, IsNil)
	c.Assert(event, DeepEquals, user.Events[0])
}

func (s *ServerSuite) TestLoginChangeUsernameExistingUsernameFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"user_name": "%s"}`, user2.Username)
	routeName := "updateUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 username already in use")
}

func (s *ServerSuite) TestLoginChangeUsernameSameUsernameFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"user_name": "%s"}`, user2.Username)
	routeName := "updateUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 username already in use")
}

func (s *ServerSuite) TestLoginChangeEmailExistingEmailFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"email": "%s"}`, user2.Email)
	routeName := "updateUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 email address already in use")
}

func (s *ServerSuite) TestLoginChangeEmailSameEmailFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	payload := fmt.Sprintf(`{"email": "%s"}`, user2.Email)
	routeName := "updateUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 email address already in use")
}

func (s *ServerSuite) TestLoginDeleteLogoutFails(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, false)
	application := accounts[0].Applications[0]
	user := application.Users[0]

	payload := fmt.Sprintf(
		`{"username": "%s", "password": "%s"}`,
		user.Username,
		user.OriginalPassword,
	)

	routeName := "loginUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, err := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	sessionToken := struct {
		UserID int64  `json:"id"`
		Token  string `json:"session_token"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)

	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	user.SessionToken = sessionToken.Token

	routeName = "deleteUser"
	route = getComposedRoute(routeName, application.ID, application.ID, user.ID)
	code, _, err = runRequest(routeName, route, "", application.AuthToken, createApplicationUserSessionToken(user), 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload = fmt.Sprintf(`{"session_token": "%s"}`, sessionToken.Token)
	routeName = "logoutUser"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, application.AuthToken, sessionToken.Token, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, "400 failed to check session token (12)\nsession mismatch")
}

func (s *ServerSuite) TestCreateUserAutoBindSocialAccounts(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]

	user2 := CorrectUserWithDefaults(application.AccountID, application.ID, 2)
	user2.Enabled = true
	user2.Activated = true
	user2.SocialConnectionsIDs = map[string][]string{
		"facebook": []string{user1.SocialIDs["facebook"]},
	}

	payloadByte, err := json.Marshal(user2)
	c.Assert(err, IsNil)
	payload := string(payloadByte)

	routeName := "createUser"
	route := getComposedRoute(routeName, application.AccountID, application.ID)
	code, body, er := runRequest(routeName, route, payload, application.AuthToken, "", 3)
	c.Assert(er, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedUser := &entity.ApplicationUser{}
	err = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(err, IsNil)
	if receivedUser.ID < 1 {
		c.Fail()
	}
	user2.OriginalPassword, receivedUser.OriginalPassword = user2.Password, user2.Password
	user2.Password = ""
	user2.CreatedAt = receivedUser.CreatedAt
	user2.UpdatedAt = receivedUser.UpdatedAt
	user2.LastLogin = receivedUser.LastLogin
	user2.ID = receivedUser.ID
	receivedUser.Image, user2.Image = nil, nil
	c.Assert(receivedUser, DeepEquals, user2)

	// Check connetions list
	routeName = "getConnectionList"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, er = runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(er, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "[]\n")

	userConnections := []entity.ApplicationUser{}
	err = json.Unmarshal([]byte(body), &userConnections)
	c.Assert(err, IsNil)

	c.Assert(len(userConnections), Equals, 1)
	c.Assert(userConnections[0].ID, Equals, receivedUser.ID)
}

func (s *ServerSuite) TestDeleteOnEventsOnUserDeleteWorks(c *C) {
	accounts := CorrectDeploy(1, 0, 1, 2, 2, true, true)
	application := accounts[0].Applications[0]
	user1 := application.Users[0]
	user2 := application.Users[1]

	// GET EVENT
	routeName := "deleteConnection"
	route := getComposedRoute(routeName, application.AccountID, application.ID, user1.ID, user2.ID)
	code, body, err := runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
	c.Assert(body, Equals, "\"\"\n")

	// GET EVENTS LIST
	routeName = "getEventList"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	events := []*entity.Event{}
	er := json.Unmarshal([]byte(body), &events)
	c.Assert(er, IsNil)
	c.Assert(len(events), Equals, 2)
	c.Assert(events[0], DeepEquals, user1.Events[1])
	c.Assert(events[1], DeepEquals, user1.Events[0])

	// Check connetions list
	routeName = "getConnectionList"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "[]\n")

	// GET EVENTS LIST
	routeName = "getConnectionEventList"
	route = getComposedRoute(routeName, application.AccountID, application.ID, user1.ID)
	code, body, err = runRequest(routeName, route, "", application.AuthToken, user1.SessionToken, 3)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "[]\n")
}

func (s *ServerSuite) TestLoginRefreshLogoutMalformedPayloadFails(c *C) {
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
		{
			Payload:   fmt.Sprintf(`{"email": "%s", "password": "%s"`, user.Email, user.OriginalPassword),
			RouteName: "loginUser",
			Route:     getComposedRoute("loginUser", application.AccountID, application.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to login the user (1)\nunexpected end of JSON input",
		},
		{
			Payload:   fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "tap@glue", user.OriginalPassword),
			RouteName: "loginUser",
			Route:     getComposedRoute("loginUser", application.AccountID, application.ID),
			Code:      http.StatusInternalServerError,
			Body:      "500 failed to retrieve the application user (1)",
		},
		{
			Payload:   fmt.Sprintf(`{"username": "%s", "password": "%s"}`, "", user.OriginalPassword),
			RouteName: "loginUser",
			Route:     getComposedRoute("loginUser", application.AccountID, application.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 both username and email are empty",
		},
		{
			Payload:   fmt.Sprintf(`{"username": "%s", "password": "%s"}`, "tapg", user.OriginalPassword),
			RouteName: "loginUser",
			Route:     getComposedRoute("loginUser", application.AccountID, application.ID),
			Code:      http.StatusInternalServerError,
			Body:      "500 failed to retrieve the application user (1)",
		},
		{
			Payload:   fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user.Username, "nothing"),
			RouteName: "loginUser",
			Route:     getComposedRoute("loginUser", application.AccountID, application.ID),
			Code:      http.StatusInternalServerError,
			Body:      "500 failed to check the account user credentials (5)\ninvalid user credentials",
		},
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"`, user.SessionToken),
			RouteName: "refreshUserSession",
			Route:     getComposedRoute("refreshUserSession", application.AccountID, application.ID, user.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to refresh the session token (1)\nunexpected end of JSON input",
		},
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"}`, "nothing"),
			RouteName: "refreshUserSession",
			Route:     getComposedRoute("refreshUserSession", application.AccountID, application.ID, user.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to refresh the session token (2)\nsession token mismatch",
		},
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"`, user.SessionToken),
			RouteName: "logoutUser",
			Route:     getComposedRoute("logoutUser", application.AccountID, application.ID, user.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to logout the user (1)\nunexpected end of JSON input",
		},
		{
			Payload:   fmt.Sprintf(`{"session_token": "%s"}`, "nothing"),
			RouteName: "logoutUser",
			Route:     getComposedRoute("logoutUser", application.AccountID, application.ID, user.ID),
			Code:      http.StatusBadRequest,
			Body:      "400 failed to logout the user (2)\nsession token mismatch",
		},
	}

	for idx := range iterations {
		code, body, err := runRequest(iterations[idx].RouteName, iterations[idx].Route, iterations[idx].Payload, application.AuthToken, user.SessionToken, 3)
		c.Logf("pass %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, iterations[idx].Code)
		c.Assert(body, Equals, iterations[idx].Body)
	}
}
