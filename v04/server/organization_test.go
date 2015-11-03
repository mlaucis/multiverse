package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/multiverse/v04/entity"

	. "gopkg.in/check.v1"
)

func (s *OrganizationSuite) TestCreateOrganization_WrongKey(c *C) {
	payload := "{namae:''}"

	routeName := "createOrganization"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(nil, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *OrganizationSuite) TestCreateOrganization_WrongValue(c *C) {
	payload := `{"name":""}`

	routeName := "createOrganization"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(nil, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *OrganizationSuite) TestCreateOrganization_OK(c *C) {
	organization := CorrectOrganization()
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s"}`, organization.Name, organization.Description)

	routeName := "createOrganization"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(nil, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedOrganization := &entity.Organization{}
	er := json.Unmarshal([]byte(body), receivedOrganization)
	c.Assert(er, IsNil)
	if receivedOrganization.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedOrganization.ID, Not(Equals), "")
	c.Assert(receivedOrganization.Name, Equals, organization.Name)
	c.Assert(receivedOrganization.Enabled, Equals, true)
}

func (s *OrganizationSuite) TestUpdateOrganization_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, organization.Name, description)

	routeName := "updateOrganization"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")

	receivedOrganization := &entity.Organization{}
	er := json.Unmarshal([]byte(body), receivedOrganization)
	c.Assert(er, IsNil)
	if receivedOrganization.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedOrganization.Name, Equals, organization.Name)
	c.Assert(receivedOrganization.Description, Equals, description)
	c.Assert(receivedOrganization.Enabled, Equals, true)
}

func (s *OrganizationSuite) TestUpdateOrganization_WrongID(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true}`, organization.Name, description)

	routeName := "updateOrganization"
	route := getComposedRoute(routeName, organization.PublicID+"1")
	code, _, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
}

func (s *OrganizationSuite) TestUpdateOrganizationMalformedPayload(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	description := "changed"
	payload := fmt.Sprintf(`{"name":"%s", "description":"%s","enabled":true`, organization.Name, description)

	routeName := "updateOrganization"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}`+"\n")
}

func (s *OrganizationSuite) TestDeleteOrganization_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "deleteOrganization"
	route := getComposedRoute(routeName, organization.PublicID)
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *OrganizationSuite) TestDeleteOrganization_WrongID(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "deleteOrganization"
	route := getComposedRoute(routeName, organization.PublicID+"1")
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Equals, `{"errors":[{"code":6003,"message":"organization mismatch"}]}`+"\n")
}

func (s *OrganizationSuite) TestGetOrganization_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	routeName := "getOrganization"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)

	c.Assert(body, Not(Equals), "")

	receivedOrganization := &entity.Organization{}
	er := json.Unmarshal([]byte(body), receivedOrganization)
	c.Assert(er, IsNil)
	c.Assert(receivedOrganization.PublicID, Equals, organization.PublicID)
	c.Assert(receivedOrganization.Name, Equals, organization.Name)
	c.Assert(receivedOrganization.Enabled, Equals, true)
}

func (s *OrganizationSuite) TestGetOrganization_WrongID(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	routeName := "getOrganization"
	route := getComposedRoute(routeName, organization.PublicID+"a")
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusBadRequest)
}
