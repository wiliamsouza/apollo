package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/dgrijalva/jwt-go"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/token"
)

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func testHandlerTokenFunc(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (s *S) TestAuthNHandleFunc(c *gocheck.C) {
	request, err := http.NewRequest("POST", "/protected", nil)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	_, err = token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	AuthNHandleFunc(testHandlerTokenFunc).ServeHTTP(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	c.Assert(response.Body.String(), gocheck.Equals, "success")
}

func (s *S) TestCORSHandle(c *gocheck.C) {
	request, err := http.NewRequest("POST", "/protected", nil)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	_, err = token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	CORSHandle(AuthNHandleFunc(testHandlerTokenFunc)).ServeHTTP(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ac := response.HeaderMap["Access-Control-Allow-Origin"][0]
	c.Assert(ac, gocheck.Equals, "*")
	c.Assert(response.Body.String(), gocheck.Equals, "success")
}

func (s *S) TestPreFlightHandleFunc(c *gocheck.C) {
	request, err := http.NewRequest("OPTIONS", "/protected", nil)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	_, err = token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	PreFlightHandleFunc(testHandlerFunc).ServeHTTP(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ac := response.HeaderMap["Access-Control-Allow-Headers"][0]
	c.Assert(ac, gocheck.Equals, "Content-Type, Accept")
}
