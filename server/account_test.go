/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"

	"fmt"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Test createAcccount request with a wrong key
func (s *ServerSuite) TestCreateAccount_WrongKey(c *C) {
	payload := "{namae:''}"

	req, err := http.NewRequest(
		"POST",
		getComposedRoute("createAccount"),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test createAcccount request with an invalid name
func (s *ServerSuite) TestCreateAccount_Invalid(c *C) {
	payload := `{"name":""}`

	req, err := http.NewRequest(
		"POST",
		getComposedRoute("createAccount"),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct createAccount request
func (s *ServerSuite) TestCreateAccount_OK(c *C) {
	correctAccount := utils.CorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, correctAccount.Name, correctAccount.Description)
	req, err := http.NewRequest(
		"POST",
		getComposedRoute("createAccount"),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusCreated)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
	c.Assert(account.Token, Not(Equals), "")
}

// Test a correct updateAccount request
func (s *ServerSuite) TestUpdateAccount_OK(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)
	req, err := http.NewRequest(
		"PUT",
		getComposedRoute("updateAccount", correctAccount.ID),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("updateAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	if account.ID < 1 {
		c.Fail()
	}
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Description, Equals, description)
	c.Assert(account.Enabled, Equals, true)
	//c.Assert(account.Token, Not(Equals), "")
}

// Test a correct updateAccount request with a wrong id
func (s *ServerSuite) TestUpdateAccount_WrongID(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, correctAccount.Name, description)
	req, err := http.NewRequest(
		"PUT",
		getComposedRoute("updateAccount", (correctAccount.ID+1)),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("updateAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct updateAccount request with an invalid description
func (s *ServerSuite) TestUpdateAccount_Invalid(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	payload := fmt.Sprintf(`{"name":"%s", "description":"","enabled":true}`, correctAccount.Name)
	req, err := http.NewRequest(
		"PUT",
		getComposedRoute("updateAccount", correctAccount.ID),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("updateAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("updateAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Not(Equals), "")
}

// Test a correct deleteAccount request
func (s *ServerSuite) TestDeleteAccount_OK(c *C) {
	account, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	req, err := http.NewRequest(
		"DELETE",
		getComposedRoute("deleteAccount", account.ID),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("deleteAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("deleteAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusNoContent)
}

// Test a correct deleteAccount request with a wrong id
func (s *ServerSuite) TestDeleteAccount_WrongID(c *C) {
	account, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	req, err := http.NewRequest(
		"DELETE",
		getComposedRoute("deleteAccount", (account.ID+1)),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("deleteAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("deleteAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}

// Test a correct getAccount request
func (s *ServerSuite) TestGetAccount_OK(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	req, err := http.NewRequest(
		"GET",
		getComposedRoute("getAccount", correctAccount.ID),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("getAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("getAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	account := &entity.Account{}
	err = json.Unmarshal([]byte(response), account)
	c.Assert(err, IsNil)
	c.Assert(account.ID, Equals, correctAccount.ID)
	c.Assert(account.Name, Equals, correctAccount.Name)
	c.Assert(account.Enabled, Equals, true)
	c.Assert(account.Token, Not(Equals), "")
}

// Test a correct getAccount request with a wrong id
func (s *ServerSuite) TestGetAccount_WrongID(c *C) {
	correctAccount, err := utils.AddCorrectAccount()
	c.Assert(err, IsNil)

	req, err := http.NewRequest(
		"GET",
		getComposedRoute("getAccount", (correctAccount.ID+1)),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("getAccount")

	m.HandleFunc(route.routePattern(apiVersion), customHandler("getAccount", route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusInternalServerError)
}
