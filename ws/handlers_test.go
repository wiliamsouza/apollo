package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

type testWebHandler struct{}

func (t testWebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Web(w, r)
}

type testRunnerHandler struct{}

func (t testRunnerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Runner(w, r)
}

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}

func (s *S) TestWebSocket(c *gocheck.C) {
	srv := httptest.NewServer(testWebHandler{})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestRunnerSocket(c *gocheck.C) {
	srv := httptest.NewServer(testRunnerHandler{})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}
