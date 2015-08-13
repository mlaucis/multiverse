package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"

	"strings"

	. "gopkg.in/check.v1"
)

// Test create acccountUser request with a wrong key
func (s *AccountUserSuite) TestCreateAccountUser_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	payload := "{usrnamae:''}"

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test create acccountUser request with an wrong name
func (s *AccountUserSuite) TestCreateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	payload := `{"user_name":""}`

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createAccountUser request
func (s *AccountUserSuite) TestCreateAccountUser_OK(c *C) {
	account := CorrectDeploy(1, 0, 0, 0, 0, false, true)[0]
	accountUser := CorrectAccountUser()
	accountUser.Username += "-asdafasdasda"
	accountUser.Email = account.PublicID + "." + accountUser.Email

	payload := fmt.Sprintf(
		`{"user_name":%q, "password":%q, "first_name": %q, "last_name": %q, "email": %q}`,
		accountUser.Username,
		accountUser.Password,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedAccountUser := &entity.AccountUser{}
	er := json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(er, IsNil)
	if receivedAccountUser.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedAccountUser.ID, Not(Equals), "")
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Email, Equals, accountUser.Email)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
	c.Assert(receivedAccountUser.Password, Equals, "")
}

// Test a correct updateAccountUser request
func (s *AccountUserSuite) TestUpdateAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedAccountUser := &entity.AccountUser{}
	er := json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(er, IsNil)
	if receivedAccountUser.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Email, Equals, accountUser.Email)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
	c.Assert(receivedAccountUser.Password, Equals, "")
}

// Test a correct updateAccountUser request with a wrong id
func (s *AccountUserSuite) TestUpdateAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID+"1")
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusConflict)
}

// Test a correct updateAccountUser request with an invalid description
func (s *AccountUserSuite) TestUpdateAccountUser_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "email"}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct updateAccountUser request with a wrong token
func (s *AccountUserSuite) TestUpdateAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		accountUser.Username,
		accountUser.FirstName,
		accountUser.LastName,
		accountUser.Email,
	)

	routeName := "updateAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct deleteAccountUser request
func (s *AccountUserSuite) TestDeleteAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccountUser request with a wrong id
func (s *AccountUserSuite) TestDeleteAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct deleteAccountUser request with a wrong token
func (s *AccountUserSuite) TestDeleteAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID+"a")
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, `{"errors":[{"code":7004,"message":"account user not found"}]}`+"\n")
}

// Test a correct getAccountUser request
func (s *AccountUserSuite) TestGetAccountUser_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedAccountUser := &entity.AccountUser{}
	er := json.Unmarshal([]byte(body), receivedAccountUser)
	c.Assert(er, IsNil)
	c.Assert(receivedAccountUser.PublicID, Equals, accountUser.PublicID)
	c.Assert(receivedAccountUser.Username, Equals, accountUser.Username)
	c.Assert(receivedAccountUser.Email, Equals, accountUser.Email)
	c.Assert(receivedAccountUser.Enabled, Equals, true)
	c.Assert(receivedAccountUser.Password, Equals, "")
}

func (s *AccountUserSuite) TestGetAccountUserListWorks(c *C) {
	numAccountUsers := 3
	accounts := CorrectDeploy(2, numAccountUsers, 0, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]

	routeName := "getAccountUserList"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := &struct {
		AccountUsers []*entity.AccountUser `json:"accountUsers"`
	}{}
	er := json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)
	accountUsers := response.AccountUsers
	c.Assert(numAccountUsers, Equals, 3)
	for idx := range accountUsers {
		c.Logf("pass %d", idx)
		c.Assert(accountUsers[idx].Password, Equals, "")

		accountUsers[idx].ID = account.Users[numAccountUsers-1-idx].ID
		accountUsers[idx].AccountID = account.Users[numAccountUsers-1-idx].AccountID
		accountUsers[idx].SessionToken = account.Users[numAccountUsers-1-idx].SessionToken
		accountUsers[idx].UpdatedAt = account.Users[numAccountUsers-1-idx].UpdatedAt
		accountUsers[idx].Password = account.Users[numAccountUsers-1-idx].Password
		accountUsers[idx].LastLogin = account.Users[numAccountUsers-1-idx].LastLogin
		accountUsers[idx].OriginalPassword = account.Users[numAccountUsers-1-idx].OriginalPassword
		c.Assert(accountUsers[idx], DeepEquals, account.Users[numAccountUsers-1-idx])
	}
}

// Test a correct getAccountUser request with a wrong id
func (s *AccountUserSuite) TestGetAccountUser_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct getAccountUser request with a wrong token
func (s *AccountUserSuite) TestGetAccountUser_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, er := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(er, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "getAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *AccountUserSuite) TestAccountUserMalformedPaylodsFail(c *C) {
	accounts := CorrectDeploy(1, 1, 1, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]

	scenarios := []struct {
		Payload      string
		RouteName    string
		Route        string
		StatusCode   int
		ResponseBody string
	}{
		// 0
		{
			Payload:      "{",
			RouteName:    "updateAccountUser",
			Route:        getComposedRoute("updateAccountUser", account.PublicID, accountUser.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 1
		{
			Payload:      "{",
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 2
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4003,"message":"both username and email are empty"}]}` + "\n",
		},
		// 3
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "tap@glue.com", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"errors":[{"code":7004,"message":"account user not found"}]}` + "\n",
		},
		// 4
		{
			Payload:      fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, "tap@glue.com", accountUser.OriginalPassword),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"errors":[{"code":7004,"message":"account user not found"}]}` + "\n",
		},
		// 5
		{
			Payload:      fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, accountUser.Username, "fake"),
			RouteName:    "loginAccountUser",
			Route:        getComposedRoute("loginAccountUser"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4011,"message":"different passwords"}]}` + "\n",
		},
		// 6
		{
			Payload:      "{",
			RouteName:    "refreshAccountUserSession",
			Route:        getComposedRoute("refreshAccountUserSession", account.PublicID, accountUser.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 7
		{
			Payload:      fmt.Sprintf(`{"session": "%s"}`, "fake"),
			RouteName:    "refreshAccountUserSession",
			Route:        getComposedRoute("refreshAccountUserSession", account.PublicID, accountUser.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4012,"message":"session token mismatch"}]}` + "\n",
		},
	}

	for idx := range scenarios {
		code, body, err := runRequest(scenarios[idx].RouteName, scenarios[idx].Route, scenarios[idx].Payload, signAccountRequest(account, accountUser, true, true))
		c.Logf("pass: %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, scenarios[idx].StatusCode)
		c.Assert(body, Equals, scenarios[idx].ResponseBody)
	}
}

func (s *AccountUserSuite) TestLoginRefreshSessionLogoutAccountUserWorks(c *C) {
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
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(nil, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusCreated)

	sessionToken := struct {
		UserID       string `json:"id"`
		AccountToken string `json:"account_token"`
		Token        string `json:"token"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	}{}
	er := json.Unmarshal([]byte(body), &sessionToken)
	c.Assert(er, IsNil)
	c.Assert(sessionToken.UserID, Equals, user.PublicID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.AccountToken, Equals, account.AuthToken)
	c.Assert(sessionToken.FirstName, Equals, user.FirstName)
	c.Assert(sessionToken.LastName, Equals, user.LastName)

	user.SessionToken = sessionToken.Token

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"token": "%s"}`, sessionToken.Token)
	routeName = "refreshAccountUserSession"
	route = getComposedRoute(routeName, account.PublicID, user.PublicID)
	code, body, err = runRequest(routeName, route, payload, signAccountRequest(account, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := struct {
		Token string `json:"token"`
	}{}
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
	user.SessionToken = sessionToken.Token

	// LOGOUT USER
	payload = fmt.Sprintf(`{"token": "%s"}`, updatedToken.Token)
	routeName = "logoutAccountUser"
	route = getComposedRoute(routeName, account.PublicID, user.PublicID)
	code, body, err = runRequest(routeName, route, payload, signAccountRequest(account, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *AccountUserSuite) TestCreateAccountUserDoubleEmailCheckMessage(c *C) {
	account := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	accountUser := account.Users[0]

	payload := fmt.Sprintf(
		`{"user_name":%q, "password":%q, "first_name": %q, "last_name": %q, "email": %q}`,
		accountUser.Username,
		accountUser.OriginalPassword,
		accountUser.FirstName,
		accountUser.LastName,
		"new+"+accountUser.Email,
	)

	routeName := "createAccountUser"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")

	receivedResponse := &entity.ErrorsResponse{}
	er := json.Unmarshal([]byte(body), receivedResponse)
	c.Assert(er, IsNil)
	c.Assert(len(receivedResponse.Errors), Equals, 1)
	c.Assert(receivedResponse.Errors[0].Code, Equals, errmsg.ErrApplicationUserUsernameInUse.Code())
	c.Assert(receivedResponse.Errors[0].Message, Equals, errmsg.ErrApplicationUserUsernameInUse.Error())
}

func (s *AccountUserSuite) TestDeleteAccountUserNewRequestFail(c *C) {
	account := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	accountUser := account.Users[0]

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	routeName = "getAccountUserList"
	route = getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Not(Equals), "")
}

func (s *AccountUserSuite) TestLoginAfterDeleteAccount(c *C) {
	account := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	accountUser := account.Users[0]

	routeName := "deleteAccountUser"
	route := getComposedRoute(routeName, account.PublicID, accountUser.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		accountUser.Email,
		accountUser.OriginalPassword,
	)

	routeName = "loginAccountUser"
	route = getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(nil, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusNotFound)
}
