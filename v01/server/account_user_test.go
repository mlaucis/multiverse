/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/entity"

	. "gopkg.in/check.v1"
)

// Test create acccountUser request with a wrong key
func (s *ServerSuite) TestCreateAccountUser_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	payload := "{usrnamae:''}"

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test create acccountUser request with an wrong name
func (s *ServerSuite) TestCreateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	payload := `{"user_name":""}`

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccountUser request
func (s *ServerSuite) TestCreateAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser := CorrectAccountUser()

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"%s", "first_name": "%s", "last_name": "%s", "email": "%s"}`,
		accountUser.Username,
		accountUser.Password,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, "", 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	if receivedAccountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request
func (s *ServerSuite) TestUpdateAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	if receivedAccountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

// Test a correct updateAccountUser request with a wrong id
func (s *ServerSuite) TestUpdateAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccountUser request with an invalid description
func (s *ServerSuite) TestUpdateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "email"}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, payload, account.AuthToken, getAccountUserSessionToken(accountUser), 2)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct updateAccountUser request with a wrong token
func (s *ServerSuite) TestUpdateAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)
	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, payload, account.AuthToken, sessionToken, 2)

	c.Assert(code, Equals, http.StatusUnauthorized)
}

// Test a correct deleteAccountUser request
func (s *ServerSuite) TestDeleteAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccountUser request with a wrong id
func (s *ServerSuite) TestDeleteAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct deleteAccountUser request with a wrong token
func (s *ServerSuite) TestDeleteAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)

	c.Assert(code, Equals, http.StatusUnauthorized)
}

// Test a correct getAccountUser request
func (s *ServerSuite) TestGetAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, body, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedAccountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(err, IsNil)
	c.Assert(receivedAccountUser.ID, Equals, accountUser.ID)
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
}

func (s *ServerSuite) TestGetAccountUserListWorks(c *C) {
	accounts := CorrectDeploy(2, 3, 0, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]

	routeName := "getAccountUserList"
	route := getComposedRoute(routeName, accountUser.AccountID)
	code, body, err := runRequest(routeName, route, "", account.AuthToken, accountUser.SessionToken, 2)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := &struct {
		AccountUsers []*entity.AccountUser `json:"accountUsers"`
	}{}
	err = json.Unmarshal([]byte(body), &response)
	c.Assert(err, IsNil)
	accountUsers := response.AccountUsers
	c.Assert(len(accountUsers), Equals, 3)
	for idx := range accountUsers {
		accountUsers[idx].SessionToken = account.Users[len(accountUsers)-1-idx].SessionToken
		accountUsers[idx].UpdatedAt = account.Users[len(accountUsers)-1-idx].UpdatedAt
		accountUsers[idx].Password = account.Users[len(accountUsers)-1-idx].Password
		accountUsers[idx].LastLogin = account.Users[len(accountUsers)-1-idx].LastLogin
		accountUsers[idx].OriginalPassword = account.Users[len(accountUsers)-1-idx].OriginalPassword
		c.Assert(accountUsers[idx], DeepEquals, account.Users[len(accountUsers)-1-idx])
	}
}

// Test a correct getAccountUser request with a wrong id
func (s *ServerSuite) TestGetAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID+1)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, getAccountUserSessionToken(accountUser), 2)

	c.Assert(code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccountUser request with a wrong token
func (s *ServerSuite) TestGetAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, err := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(err, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, accountUser.AccountID, accountUser.ID)
	code, _, err := runRequest(routeName, route, "", account.AuthToken, sessionToken, 2)

	c.Assert(code, Equals, http.StatusUnauthorized)
}

func (s *ServerSuite) TestAccountUserMalformedPaylodsFail(c *C) {
	accounts := CorrectDeploy(1, 1, 1, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]
	application := account.Applications[0]

	scenarios := []struct {
		Payload      string
		RouteName    string
		Route        string
		StatusCode   int
		ResponseBody string
	}{
		{
			Payload:      "{",
			RouteName:    "updateAccountUser",
			Route:        getComposedRoute("updateAccountUser", application.AccountID, accountUser.ID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to update the user (1)\nunexpected EOF",
		},
		{
			Payload:      "{",
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to login the user (1)\nunexpected EOF",
		},
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to login the user (2)\nboth username and email are empty. please use one of them",
		},
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "tap@glue.com", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to login the user (3)\naccount user not found",
		},
		{
			Payload:      fmt.Sprintf(`{"username": "%s", "password": "%s"}`, "tap@glue.com", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to login the user (4)\naccount user not found",
		},
		{
			Payload:      fmt.Sprintf(`{"username": "%s", "password": "%s"}`, accountUser.Username, "fake"),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusUnauthorized,
			ResponseBody: "401 failed to login the user (6)",
		},
		{
			Payload:      "{",
			RouteName:    "refreshAccountUserSession",
			Route:        getComposedRoute("refreshAccountUserSession", account.ID, accountUser.ID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to refresh session token (1)\nunexpected EOF",
		},
		{
			Payload:      fmt.Sprintf(`{"session": "%s"}`, "fake"),
			RouteName:    "refreshAccountUserSession",
			Route:        getComposedRoute("refreshAccountUserSession", account.ID, accountUser.ID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to refresh session token (2) \nsession token mismatch",
		},
		{
			Payload:      "{",
			RouteName:    "logoutAccountUser",
			Route:        getComposedRoute("logoutAccountUser", account.ID, accountUser.ID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to logout the user (1)\nunexpected EOF",
		},
		{
			Payload:      fmt.Sprintf(`{"session": "%s"}`, "fake"),
			RouteName:    "logoutAccountUser",
			Route:        getComposedRoute("logoutAccountUser", account.ID, accountUser.ID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: "400 failed to logout the user (2) \nsession token mismatch",
		},
	}

	for idx := range scenarios {
		code, body, err := runRequest(scenarios[idx].RouteName, scenarios[idx].Route, scenarios[idx].Payload, account.AuthToken, accountUser.SessionToken, 2)
		c.Logf("pass: %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, scenarios[idx].StatusCode)
		c.Assert(body, Equals, scenarios[idx].ResponseBody)
	}
}

func (s *ServerSuite) TestLoginRefreshSessionLogoutAccountUserWorks(c *C) {
	accounts := CorrectDeploy(1, 1, 0, 0, 0, false, false)
	account := accounts[0]
	user := account.Users[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		user.Email,
		user.OriginalPassword,
	)

	routeName := "loginAccountUser"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, "", "", 2)
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	sessionToken := struct {
		UserID       int64  `json:"id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{}
	err = json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(err, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.ID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.AccountToken, Equals, account.AuthToken)
	c.Assert(sessionToken.FirstName, Equals, user.FirstName)
	c.Assert(sessionToken.LastName, Equals, user.LastName)

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"token": "%s"}`, sessionToken.Token)
	routeName = "refreshAccountUserSession"
	route = getComposedRoute(routeName, account.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, account.AuthToken, sessionToken.Token, 2)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(err, IsNil)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)

	// LOGOUT USER
	payload = fmt.Sprintf(`{"token": "%s"}`, updatedToken.Token)
	routeName = "logoutAccountUser"
	route = getComposedRoute(routeName, account.ID, user.ID)
	code, body, err = runRequest(routeName, route, payload, account.AuthToken, updatedToken.Token, 2)
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Equals, "\"logged out\"\n")
}
