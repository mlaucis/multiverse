package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/multiverse/utils"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"

	"strings"

	. "gopkg.in/check.v1"
)

func (s *MemberSuite) TestCreateMember_WrongKey(c *C) {
	organization, err := AddCorrectOrganization(true)
	payload := "{usrnamae:''}"

	routeName := "createMember"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *MemberSuite) TestCreateMember_WrongValue(c *C) {
	organization, err := AddCorrectOrganization(true)
	payload := `{"user_name":""}`

	routeName := "createMember"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, nil, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *MemberSuite) TestCreateMember_OK(c *C) {
	organization := CorrectDeploy(1, 0, 0, 0, 0, false, true)[0]
	member := CorrectMember()
	member.Username += "-asdafasdasda"
	member.Email = organization.PublicID + "." + member.Email

	payload := fmt.Sprintf(
		`{"user_name":%q, "password":%q, "first_name": %q, "last_name": %q, "email": %q}`,
		member.Username,
		member.Password,
		member.FirstName,
		member.LastName,
		member.Email,
	)

	routeName := "createMember"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedMember := &entity.Member{}
	er := json.Unmarshal([]byte(body), receivedMember)
	c.Assert(er, IsNil)
	if receivedMember.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedMember.ID, Not(Equals), "")
	c.Assert(receivedMember.Username, Equals, member.Username)
	c.Assert(receivedMember.Email, Equals, member.Email)
	c.Assert(receivedMember.Enabled, Equals, true)
	c.Assert(receivedMember.Password, Equals, "")
}

func (s *MemberSuite) TestUpdateMember_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		member.Username,
		member.FirstName,
		member.LastName,
		member.Email,
	)

	routeName := "updateMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedMember := &entity.Member{}
	er := json.Unmarshal([]byte(body), receivedMember)
	c.Assert(er, IsNil)
	if receivedMember.PublicID == "" {
		c.Fail()
	}
	c.Assert(receivedMember.Username, Equals, member.Username)
	c.Assert(receivedMember.Email, Equals, member.Email)
	c.Assert(receivedMember.Enabled, Equals, true)
	c.Assert(receivedMember.Password, Equals, "")
}

func (s *MemberSuite) TestUpdateMember_WrongID(c *C) {
	organinzation, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organinzation.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"changed", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		member.Username,
		member.FirstName,
		member.LastName,
		member.Email,
	)

	routeName := "updateMember"
	route := getComposedRoute(routeName, organinzation.PublicID, member.PublicID+"1")
	code, _, err := runRequest(routeName, route, payload, signOrganizationRequest(organinzation, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusConflict)
}

func (s *MemberSuite) TestUpdateMember_WrongValue(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "email"}`,
		member.Username,
		member.FirstName,
		member.LastName,
	)

	routeName := "updateMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusBadRequest)
	c.Assert(body, Not(Equals), "")
}

func (s *MemberSuite) TestUpdateMember_WrongToken(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	payload := fmt.Sprintf(
		`{"user_name":"%s", "password":"", "first_name": "%s", "last_name": "%s", "email": "%s", "enabled": true}`,
		member.Username,
		member.FirstName,
		member.LastName,
		member.Email,
	)

	routeName := "updateMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, _, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *MemberSuite) TestDeleteMember_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "deleteMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *MemberSuite) TestDeleteMember_WrongID(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "deleteMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))

	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *MemberSuite) TestDeleteMember_WrongToken(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "deleteMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID+"a")
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Equals, `{"errors":[{"code":7004,"message":"member not found"}]}`+"\n")
}

func (s *MemberSuite) TestGetMember_OK(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	LoginMember(member)

	routeName := "getMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	c.Assert(strings.Contains(body, "created_at"), Equals, true)

	receivedMember := &entity.Member{}
	er := json.Unmarshal([]byte(body), receivedMember)
	c.Assert(er, IsNil)
	c.Assert(receivedMember.PublicID, Equals, member.PublicID)
	c.Assert(receivedMember.Username, Equals, member.Username)
	c.Assert(receivedMember.Email, Equals, member.Email)
	c.Assert(receivedMember.Enabled, Equals, true)
	c.Assert(receivedMember.Password, Equals, "")
}

func (s *MemberSuite) TestGetMemberListWorks(c *C) {
	numOrganizationUsers := 3
	organizations := CorrectDeploy(2, numOrganizationUsers, 0, 0, 0, false, true)
	organization := organizations[0]
	member := organization.Members[0]

	routeName := "getMemberList"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")

	response := &struct {
		Members []*entity.Member `json:"members"`
	}{}
	er := json.Unmarshal([]byte(body), &response)
	c.Assert(er, IsNil)
	memberstUsers := response.Members
	c.Assert(numOrganizationUsers, Equals, 3)
	for idx := range memberstUsers {
		c.Logf("pass %d", idx)
		c.Assert(memberstUsers[idx].Password, Equals, "")

		memberstUsers[idx].ID = organization.Members[numOrganizationUsers-1-idx].ID
		memberstUsers[idx].OrgID = organization.Members[numOrganizationUsers-1-idx].OrgID
		memberstUsers[idx].SessionToken = organization.Members[numOrganizationUsers-1-idx].SessionToken
		memberstUsers[idx].UpdatedAt = organization.Members[numOrganizationUsers-1-idx].UpdatedAt
		memberstUsers[idx].Password = organization.Members[numOrganizationUsers-1-idx].Password
		memberstUsers[idx].LastLogin = organization.Members[numOrganizationUsers-1-idx].LastLogin
		memberstUsers[idx].OriginalPassword = organization.Members[numOrganizationUsers-1-idx].OriginalPassword
		c.Assert(memberstUsers[idx], DeepEquals, organization.Members[numOrganizationUsers-1-idx])
	}
}

func (s *MemberSuite) TestGetMember_WrongID(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	routeName := "getMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID+"1")
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))

	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *MemberSuite) TestGetMember_WrongToken(c *C) {
	organization, err := AddCorrectOrganization(true)
	c.Assert(err, IsNil)

	member, err := AddCorrectMember(organization.ID, true)
	c.Assert(err, IsNil)

	sessionToken, er := utils.Base64Decode(getMemberSessionToken(member))
	c.Assert(er, IsNil)

	sessionToken = utils.Base64Encode(sessionToken + "a")

	routeName := "getMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, false))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
}

func (s *MemberSuite) TestMemberMalformedPaylodsFail(c *C) {
	organizations := CorrectDeploy(1, 1, 1, 0, 0, false, true)
	organization := organizations[0]
	members := organization.Members[0]

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
			RouteName:    "updateMember",
			Route:        getComposedRoute("updateMember", organization.PublicID, members.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 1
		{
			Payload:      "{",
			RouteName:    "loginMember",
			Route:        getComposedRoute("loginMember"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 2
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "", members.OriginalPassword),
			RouteName:    "loginMember",
			Route:        getComposedRoute("loginMember"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4003,"message":"both username and email are empty"}]}` + "\n",
		},
		// 3
		{
			Payload:      fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "tap@glue.com", members.OriginalPassword),
			RouteName:    "loginMember",
			Route:        getComposedRoute("loginMember"),
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"errors":[{"code":7004,"message":"member not found"}]}` + "\n",
		},
		// 4
		{
			Payload:      fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, "tap@glue.com", members.OriginalPassword),
			RouteName:    "loginMember",
			Route:        getComposedRoute("loginMember"),
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"errors":[{"code":7004,"message":"member not found"}]}` + "\n",
		},
		// 5
		{
			Payload:      fmt.Sprintf(`{"user_name": "%s", "password": "%s"}`, members.Username, "fake"),
			RouteName:    "loginMember",
			Route:        getComposedRoute("loginMember"),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4011,"message":"different passwords"}]}` + "\n",
		},
		// 6
		{
			Payload:      "{",
			RouteName:    "refreshMemberSession",
			Route:        getComposedRoute("refreshMemberSession", organization.PublicID, members.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":5001,"message":"unexpected end of JSON input"}]}` + "\n",
		},
		// 7
		{
			Payload:      fmt.Sprintf(`{"session": "%s"}`, "fake"),
			RouteName:    "refreshMemberSession",
			Route:        getComposedRoute("refreshMemberSession", organization.PublicID, members.PublicID),
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"errors":[{"code":4012,"message":"session token mismatch"}]}` + "\n",
		},
	}

	for idx := range scenarios {
		code, body, err := runRequest(scenarios[idx].RouteName, scenarios[idx].Route, scenarios[idx].Payload, signOrganizationRequest(organization, members, true, true))
		c.Logf("pass: %d", idx)
		c.Assert(err, IsNil)
		c.Assert(code, Equals, scenarios[idx].StatusCode)
		c.Assert(body, Equals, scenarios[idx].ResponseBody)
	}
}

func (s *MemberSuite) TestLoginRefreshSessionLogoutMemberWorks(c *C) {
	organizations := CorrectDeploy(1, 1, 0, 0, 0, false, false)
	organization := organizations[0]
	member := organization.Members[0]

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		member.Email,
		member.OriginalPassword,
	)

	routeName := "loginMember"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(nil, nil, true, true))
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
	c.Assert(sessionToken.UserID, Equals, member.PublicID)
	c.Assert(sessionToken.Token, Not(Equals), "")
	c.Assert(sessionToken.AccountToken, Equals, organization.AuthToken)
	c.Assert(sessionToken.FirstName, Equals, member.FirstName)
	c.Assert(sessionToken.LastName, Equals, member.LastName)

	member.SessionToken = sessionToken.Token

	// REFRESH USER SESSION
	payload = fmt.Sprintf(`{"token": "%s"}`, sessionToken.Token)
	routeName = "refreshMemberSession"
	route = getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, body, err = runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusCreated)
	c.Assert(body, Not(Equals), "")
	updatedToken := struct {
		Token string `json:"token"`
	}{}
	er = json.Unmarshal([]byte(body), &updatedToken)
	c.Assert(er, IsNil)
	c.Assert(updatedToken.Token, Not(Equals), sessionToken.Token)
	member.SessionToken = sessionToken.Token

	// LOGOUT USER
	payload = fmt.Sprintf(`{"token": "%s"}`, updatedToken.Token)
	routeName = "logoutMember"
	route = getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, body, err = runRequest(routeName, route, payload, signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)
}

func (s *MemberSuite) TestCreateMemberDoubleEmailCheckMessage(c *C) {
	organization := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	member := organization.Members[0]

	payload := fmt.Sprintf(
		`{"user_name":%q, "password":%q, "first_name": %q, "last_name": %q, "email": %q}`,
		member.Username,
		member.OriginalPassword,
		member.FirstName,
		member.LastName,
		"new+"+member.Email,
	)

	routeName := "createMember"
	route := getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(organization, nil, true, true))
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

func (s *MemberSuite) TestDeleteMemberNewRequestFail(c *C) {
	organization := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	member := organization.Members[0]

	routeName := "deleteMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	routeName = "getMemberList"
	route = getComposedRoute(routeName, organization.PublicID)
	code, body, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNotFound)
	c.Assert(body, Not(Equals), "")
}

func (s *MemberSuite) TestLoginAfterDeleteOrganization(c *C) {
	organization := CorrectDeploy(1, 1, 0, 0, 0, false, true)[0]
	member := organization.Members[0]

	routeName := "deleteMember"
	route := getComposedRoute(routeName, organization.PublicID, member.PublicID)
	code, _, err := runRequest(routeName, route, "", signOrganizationRequest(organization, member, true, true))

	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusNoContent)

	payload := fmt.Sprintf(
		`{"email": "%s", "password": "%s"}`,
		member.Email,
		member.OriginalPassword,
	)

	routeName = "loginMember"
	route = getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, payload, signOrganizationRequest(nil, nil, true, true))
	c.Assert(err, IsNil)
	c.Assert(body, Not(Equals), "")
	c.Assert(code, Equals, http.StatusNotFound)
}
