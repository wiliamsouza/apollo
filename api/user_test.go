package api

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/wiliamsouza/apollo/db"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
)

func (s *S) TestNewUser(c *gocheck.C) {
	result := `{"name":"Jhon Doe","email":"jhon@doe.com"}`
	defer db.Session.User().Remove(bson.M{"_id": "jhon@doe.com"})
	body := strings.NewReader(`{"name":"Jhon Doe","email":"jhon@doe.com","password":"12345"}`)
	request, _ := http.NewRequest("POST", "users", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewUser(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}
