package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/customer"
	"github.com/wiliamsouza/apollo/db"
)

func (s *S) TestNewUser(c *gocheck.C) {
	result := `{"name":"Jhon Doe","email":"jhon@doe.com"}`
	defer db.Session.User().Remove(bson.M{"_id": "jhon@doe.com"})
	body := strings.NewReader(`{"name":"Jhon Doe","email":"jhon@doe.com","password":"12345"}`)
	request, err := http.NewRequest("POST", "users", body)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewUser(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestDetailUser(c *gocheck.C) {
	email := "jhon@doe.com"
	user, err := customer.NewUser("Jhon Doe", email, "12345")
	defer db.Session.User().RemoveId(email)
	c.Assert(err, gocheck.IsNil)
	detail := detailUser{Name: user.Name, Email: user.Email,
		APIKey: user.APIKey, Created: user.Created,
		LastLogin: user.LastLogin}
	result, err := json.Marshal(&detail)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("users/%s", email)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DetailUser(response, request, map[string]string{"email": email})
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, string(result))
}

func (s *S) TestAuthenticate(c *gocheck.C) {
	email := "jhon@doe.com"
	_, err := customer.NewUser("Jhon Doe", email, "12345")
	defer db.Session.User().RemoveId(email)
	c.Assert(err, gocheck.IsNil)
	body := strings.NewReader(`{"email":"jhon@doe.com","password":"12345"}`)
	request, err := http.NewRequest("POST", "users/authenticate", body)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	Authenticate(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
}
