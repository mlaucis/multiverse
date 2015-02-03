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

	"fmt"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Test create acccount request with a wrong key
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

// Test a correct create account request
func (s *ServerSuite) TestCreateAccount_Correct(c *C) {
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

// Test getAccount
func (s *ServerSuite) TestGetAccount_OK(c *C) {
	account := AddCorrectAccount()

	req, err := http.NewRequest(
		"GET",
		getComposedRoute("getAccount", account.ID),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	route := getRoute("getAccount")

	m.HandleFunc(route.routePattern(apiVersion), route.handlers).Methods(route.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusOK)
	response := w.Body.String()
	c.Assert(response, Not(Equals), "")

	accountGet := &entity.Account{}
	err = json.Unmarshal([]byte(response), accountGet)
	c.Assert(err, IsNil)
	c.Assert(accountGet.ID, Equals, account.ID)
	c.Assert(accountGet.Name, Equals, account.Name)
	c.Assert(accountGet.Enabled, Equals, true)
	c.Assert(accountGet.Token, Not(Equals), "")
}
