package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/entity"

	"math/rand"

	. "gopkg.in/check.v1"
)

// Test createApplication request with a wrong key
func (s *ApplicationSuite) TestCreateApplication_WrongKey(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := "{namae:''}"

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test createApplication request with an wrong name
func (s *ApplicationSuite) TestCreateApplication_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	payload := `{"name":""}`

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

// Test a correct createApplication request
func (s *ApplicationSuite) TestCreateApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	application := CorrectApplication()

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"%s", "url": "%s"}`,
		application.Name,
		application.Description,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "createApplication"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)

	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	er := json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(er, IsNil)
	if receivedApplication.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedApplication.ID, Not(Equals), "")
	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.Description, Equals, application.Description)
	c.Assert(receivedApplication.URL, Equals, application.URL)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct updateApplication request
func (s *ApplicationSuite) TestUpdateApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID)
	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	er := json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(er, IsNil)
	if receivedApplication.PublicID == "" {
		c.Fail()
	}

	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.URL, Equals, application.URL)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct updateApplication request with a wrong id
func (s *ApplicationSuite) TestUpdateApplication_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	application.PublicAccountID = account.PublicID

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, application.PublicAccountID, application.PublicID+"a")
	code, _, er := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(er, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct updateApplication request with an invalid description
func (s *ApplicationSuite) TestUpdateApplication_WrongValue(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	application.PublicAccountID = account.PublicID

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"", "url": "%s", "enabled": true}`,
		application.Name,
		application.URL,
	)
	c.Assert(err, IsNil)

	routeName := "updateApplication"
	route := getComposedRoute(routeName, application.PublicAccountID, application.PublicID)
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct updateApplication request with a wrong token
func (s *ApplicationSuite) TestUpdateApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	correctApplication, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	correctApplication.PublicAccountID = account.PublicID

	payload := fmt.Sprintf(
		`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true}`,
		correctApplication.Name,
		correctApplication.URL,
	)
	c.Assert(err, IsNil)

	sessionToken, er := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(er, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "updateApplication"
	route := getComposedRoute(routeName, correctApplication.PublicAccountID, correctApplication.PublicID)
	code, _, err := runRequest(routeName, route, payload, signAccountRequest(account, accountUser, false, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct deleteApplication request
func (s *ApplicationSuite) TestDeleteApplication_OK(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

// Test a correct deleteApplication request with a wrong id
func (s *ApplicationSuite) TestDeleteApplication_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct deleteApplication request with a wrong token
func (s *ApplicationSuite) TestDeleteApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	sessionToken, er := utils.Base64Decode(getAccountUserSessionToken(accountUser))
	c.Assert(er, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "deleteApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct getApplication request
func (s *ApplicationSuite) TestGetApplication_OK(c *C) {
	accounts := CorrectDeploy(1, 1, 1, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]
	application := account.Applications[rand.Intn(1)]

	routeName := "getApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedApplication := &entity.Application{}
	er := json.Unmarshal([]byte(body), receivedApplication)
	c.Assert(er, IsNil)
	c.Assert(receivedApplication.PublicID, Equals, application.PublicID)
	c.Assert(receivedApplication.Name, Equals, application.Name)
	c.Assert(receivedApplication.Description, Equals, application.Description)
	c.Assert(receivedApplication.Enabled, Equals, true)
}

// Test a correct getApplication request with a wrong id
func (s *ApplicationSuite) TestGetApplication_WrongID(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	routeName := "getApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID+"a")
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

// Test a correct getApplication request with a wrong token
func (s *ApplicationSuite) TestGetApplication_WrongToken(c *C) {
	account, err := AddCorrectAccount(true)
	c.Assert(err, IsNil)

	accountUser, err := AddCorrectAccountUser(account.ID, true)
	c.Assert(err, IsNil)

	application, err := AddCorrectApplication(account.ID, true)
	c.Assert(err, IsNil)

	LoginAccountUser(accountUser)

	routeName := "getApplication"
	route := getComposedRoute(routeName, account.PublicID, application.PublicID)
	code, _, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *ApplicationSuite) TestGetApplicationListWorks(c *C) {
	accounts := CorrectDeploy(2, 1, 1, 0, 0, false, true)
	account := accounts[0]
	accountUser := account.Users[0]
	application := account.Applications[0]

	routeName := "getApplications"
	route := getComposedRoute(routeName, account.PublicID)
	code, body, err := runRequest(routeName, route, "", signAccountRequest(account, accountUser, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := &struct {
		Applications []*entity.Application `json:"applications"`
	}{}

	er := json.Unmarshal([]byte(body), response)
	c.Assert(er, IsNil)
	c.Assert(len(response.Applications), Equals, 1)
	application.Users = nil
	application.ID = 0
	c.Assert(response.Applications[0], DeepEquals, application)
}

func (s *ApplicationSuite) TestApplicationMalformedPayloadsFails(c *C) {
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
			RouteName:    "updateApplication",
			Route:        getComposedRoute("updateApplication", account.PublicID, application.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
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
