/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"
	. "gopkg.in/check.v1"
)

// Test create acccountUser request with a wrong key
func (s *ServerSuite) TestCreateAccountUser_WrongKey(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	payload := "{usrnamae:''}"

	req, err := http.NewRequest(
		"POST",
		getComposedRoute("createAccountUser", correctAccount.ID),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccountUser(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test create acccountUser request with an invalid name
func (s *ServerSuite) TestCreateAccountUser_Invalid(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	payload := `{"user_name":""}`

	req, err := http.NewRequest(
		"POST",
		getComposedRoute("createAccountUser", correctAccount.ID),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccountUser(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct createAccountUser request
func (s *ServerSuite) TestCreateAccountUser_OK(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	correctAccountUser := utils.CorrectAccountUser()

	currentRoute := "createAccountUser"

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"%s", "first_name": "%s", "last_name": "%s", "email": "%s"}`,
		correctAccountUser.Username,
		correctAccountUser.Password,
		correctAccountUser.FirstName,
		correctAccountUser.LastName,
		correctAccountUser.Email,
	)

	req, err := http.NewRequest(
		"POST",
		getComposedRoute(currentRoute, correctAccount.ID),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)
	token, err := storageClient.GenerateAccountToken(correctAccount)
	if err != nil {
		panic(err)
	}
	signRequest(token, req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute(currentRoute)

	m.HandleFunc(route.routePattern(apiVersion), customHandler(currentRoute, route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	accountUser := &entity.AccountUser{}
	err = json.Unmarshal([]byte(response), accountUser)
	c.Assert(err, IsNil)
	if accountUser.ID < 1 {
		c.Fail()
	}
	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
	c.Assert(accountUser.Enabled, Equals, true)
}

// // Test a correct updateAccountUser request
// func (s *ServerSuite) TestUpdateAccountUser_OK(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	password := "changed"
// 	payload := fmt.Sprintf(`{"user_name":"%s", "password":"%s","enabled":true}`, correctAccountUser.Username, password)
// 	req, err := http.NewRequest(
// 		"PUT",
// 		getComposedRoute("updateAccountUser", correctAccountUser.AccountID, correctAccountUser.ID),
// 		strings.NewReader(payload),
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader(payload, req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("updateAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusOK)
// 	response := w.Body.String()
// 	c.Assert(response, Not(Equals), "")

// 	accountUser := &entity.AccountUser{}
// 	err = json.Unmarshal([]byte(response), accountUser)
// 	c.Assert(err, IsNil)
// 	if accountUser.ID < 1 {
// 		c.Fail()
// 	}
// 	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
// 	c.Assert(accountUser.Password, Equals, password)
// 	c.Assert(accountUser.Enabled, Equals, true)
// }

// // Test a correct updateAccountUser request with a wrong id
// func (s *ServerSuite) TestUpdateAccountUser_WrongID(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	password := "changed"
// 	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccountUser.Username, password)
// 	req, err := http.NewRequest(
// 		"PUT",
// 		getComposedRoute("updateAccountUser", correctAccountUser.AccountID, (correctAccountUser.ID+1)),
// 		strings.NewReader(payload),
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader(payload, req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("updateAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusInternalServerError)
// }

// // Test a correct updateAccountUser request with an invalid description
// func (s *ServerSuite) TestUpdateAccountUser_Invalid(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	payload := fmt.Sprintf(`{"user_name":"%s", "password":"","enabled":true}`, correctAccountUser.Username)
// 	req, err := http.NewRequest(
// 		"PUT",
// 		getComposedRoute("updateAccountUser", correctAccountUser.AccountID, correctAccountUser.ID),
// 		strings.NewReader(payload),
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader(payload, req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("updateAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusBadRequest)
// 	c.Assert(w.Body.String(), Not(Equals), "")
// }

// // Test a correct deleteAccountUser request
// func (s *ServerSuite) TestDeleteAccountUser_OK(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	c.Assert(err, IsNil)

// 	req, err := http.NewRequest(
// 		"DELETE",
// 		getComposedRoute("deleteAccountUser", correctAccountUser.AccountID, correctAccountUser.ID),
// 		nil,
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader("", req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("deleteAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("deleteAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusNoContent)
// }

// // Test a correct deleteAccountUser request with a wrong id
// func (s *ServerSuite) TestDeleteAccountUser_WrongID(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	c.Assert(err, IsNil)

// 	req, err := http.NewRequest(
// 		"DELETE",
// 		getComposedRoute("deleteAccountUser", correctAccountUser.AccountID, (correctAccountUser.ID+1)),
// 		nil,
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader("", req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("deleteAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("deleteAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusInternalServerError)
// }

// // Test a correct getAccountUser request
// func (s *ServerSuite) TestGetAccountUser_OK(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	c.Assert(err, IsNil)

// 	req, err := http.NewRequest(
// 		"GET",
// 		getComposedRoute("getAccountUser", correctAccountUser.AccountID, correctAccountUser.ID),
// 		nil,
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader("", req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("getAccountUser")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("getAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusOK)
// 	response := w.Body.String()
// 	c.Assert(response, Not(Equals), "")

// 	accountUser := &entity.AccountUser{}
// 	err = json.Unmarshal([]byte(response), accountUser)
// 	c.Assert(err, IsNil)
// 	c.Assert(accountUser.ID, Equals, correctAccountUser.ID)
// 	c.Assert(accountUser.Username, Equals, correctAccountUser.Username)
// 	c.Assert(accountUser.Enabled, Equals, true)
// }

// // Test a correct getAccountUser request with a wrong id
// func (s *ServerSuite) TestGetAccountUser_WrongID(c *C) {
// 	correctAccountUser, err := utils.AddCorrectAccountUser()
// 	c.Assert(err, IsNil)

// 	req, err := http.NewRequest(
// 		"GET",
// 		getComposedRoute("getAccountUser", correctAccountUser.AccountID, (correctAccountUser.ID+1)),
// 		nil,
// 	)
// 	c.Assert(err, IsNil)

// 	clHeader("", req)

// 	w := httptest.NewRecorder()
// 	m := mux.NewRouter()
// 	route := getRoute("getAccount")

// 	m.HandleFunc(route.routePattern(apiVersion), customHandler("getAccountUser", route, nil, logChan)).Methods(route.method)
// 	m.ServeHTTP(w, req)

// 	c.Assert(w.Code, Equals, http.StatusInternalServerError)
// }
