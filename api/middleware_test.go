package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/token"
)

func (s *S) TestAuthN(c *gocheck.C) {
	request, err := http.NewRequest("POST", "/protected", nil)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	context := martini.New().createContext(response, request)
	AuthN()(response, request, context)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	c.Assert(response.Body.String(), gocheck.Equals, "success")
}
