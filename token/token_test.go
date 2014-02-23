package token

import (
	"fmt"
	"net/http"

	"launchpad.net/gocheck"
)

func (s *S) TestNewToken(c *gocheck.C) {
	LoadKeys()
	_, err := New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestVaidateToken(c *gocheck.C) {
	email := "jhon@doe.com"
	LoadKeys()
	t, err := New(email)
	c.Assert(err, gocheck.IsNil)
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	token, err := Validate(request)
	c.Assert(err, gocheck.IsNil)
	c.Assert(token.Claims["email"], gocheck.Equals, email)
}
