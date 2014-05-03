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
	_, err = token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	m := martini.Classic()
	m.Post("/protected", AuthN(), func() int {
		return http.StatusOK
	})
	m.ServeHTTP(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
}
