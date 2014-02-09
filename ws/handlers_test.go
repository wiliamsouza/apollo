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

type testWebHandler struct {
	Key string
}

func (t testWebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Web(w, r, map[string]string{"apikey": t.Key})
}

type testRunnerHandler struct {
	Key string
}

func (t testRunnerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Runner(w, r, map[string]string{"apikey": t.Key})
}

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}

func (s *S) TestWebSocket(c *gocheck.C) {
	apiKey := "secret-key"
	srv := httptest.NewServer(testWebHandler{apiKey})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestRunnerSocket(c *gocheck.C) {
	apiKey := "secret-key"
	srv := httptest.NewServer(testRunnerHandler{apiKey})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}
