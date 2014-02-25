package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/customer"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/token"
)

func (s *S) TestNewOrganization(c *gocheck.C) {
	result := `{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}`
	defer db.Session.Organization().RemoveId("doegroup")
	body := strings.NewReader(result)
	request, err := http.NewRequest("POST", "organizations", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestNewOrganizationInvalidJson(c *gocheck.C) {
	result := "Error parssing json request, unexpected end of JSON input\n"
	invalid := `{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]`
	body := strings.NewReader(invalid)
	request, err := http.NewRequest("POST", "organizations", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestListOrganizations(c *gocheck.C) {
	results := `[{"name":"janecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jane@doe.com"]},{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}]`
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name1 := "jhoncorp"
	name2 := "janecorp"
	defer db.Session.Organization().RemoveId(name1)
	defer db.Session.Organization().RemoveId(name2)
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o1 := customer.Organization{Name: name1, Teams: []customer.Team{t}, Admins: []string{jhon}}
	o2 := customer.Organization{Name: name2, Teams: []customer.Team{t}, Admins: []string{jane}}
	_, err := customer.NewOrganization(o1)
	c.Assert(err, gocheck.IsNil)
	_, err = customer.NewOrganization(o2)
	c.Assert(err, gocheck.IsNil)
	request, err := http.NewRequest("GET", "organizations", nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	ListOrganizations(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDetailOrganization(c *gocheck.C) {
	result := `{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}`
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "doecorp"
	defer db.Session.Organization().RemoveId(name)
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DetailOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestModifyOrganization(c *gocheck.C) {
	b := `{"name":"jhoncorp","teams":[{"name":"Test","users":["jane@doe.com"]}],"admins":["jane@doe.com"]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb customer.Organization
	err = db.Session.Organization().FindId(name).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganizationAddNewTeamAndAdmins(c *gocheck.C) {
	b := `{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]},{"name":"Dev","users":["jhon@doe.com"]}],"admins":["jhon@doe.com","jane@doe.com"]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb customer.Organization
	err = db.Session.Organization().FindId(name).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganizationRemoveAllTeam(c *gocheck.C) {
	b := `{"name":"jhoncorp","teams":[],"admins":["jhon@doe.com","jane@doe.com"]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t1 := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	t2 := customer.Team{Name: "Dev", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t1, t2}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb customer.Organization
	err = db.Session.Organization().FindId(name).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganizationRemoveOneTeamOneTeamUser(c *gocheck.C) {
	b := `{"name":"jhoncorp","teams":[{"name":"Dev","users":["jhon@doe.com"]}],"admins":["jhon@doe.com","jane@doe.com"]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t1 := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	t2 := customer.Team{Name: "Dev", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t1, t2}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb customer.Organization
	err = db.Session.Organization().FindId(name).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganizationRemoveOneAdmin(c *gocheck.C) {
	b := `{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]},{"name":"Dev","users":["jane@doe.com","jhon@doe.com"]}],"admins":["jane@doe.com"]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t1 := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	t2 := customer.Team{Name: "Dev", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t1, t2}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb customer.Organization
	err = db.Session.Organization().FindId(name).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganizationRemoveAllAdminShouldReturnError(c *gocheck.C) {
	result := "Error updating organization, can not remove all organization admins\n"
	b := `{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]},{"name":"Dev","users":["jane@doe.com","jhon@doe.com"]}],"admins":[]}`
	body := strings.NewReader(b)
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "jhoncorp"
	t1 := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	t2 := customer.Team{Name: "Dev", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t1, t2}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(name)
	var org customer.Organization
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestDeleteOrganization(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "doecorp"
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", name)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DeleteOrganization(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	lenght, err := db.Session.Organization().FindId(name).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}

func (s *S) TestDeleteOrganizationNoExist(c *gocheck.C) {
	result := "Error deleting organization, error removing organization: not found\n"
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	name := "exist"
	defer db.Session.Organization().RemoveId(name)
	t := customer.Team{Name: "Test", Users: []string{jhon, jane}}
	o := customer.Organization{Name: name, Teams: []customer.Team{t}, Admins: []string{jhon}}
	_, err := customer.NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("organizations/%s", "noexist")
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DeleteOrganization(response, request, tk)
	lenght, err := db.Session.Organization().FindId(name).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 1)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}
