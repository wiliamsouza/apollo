package ws

import (
	//"net/http"
	//"net/http/httptest"
	"testing"

	//"github.com/gorilla/websocket"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}

// TODO: Learn how to test websocket
/**func (s *S) TestWebSocket(c *gocheck.C) {
	apiKey := "secret-key"
	srv := httptest.NewServer(testWebHandler{apiKey})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestAgentSocket(c *gocheck.C) {
	apiKey := "secret-key"
	srv := httptest.NewServer(testAgentHandler{apiKey})
	defer srv.Close()
	header := http.Header{"Origin": {srv.URL}}
	_, _, err := websocket.DefaultDialer.Dial(httpToWs(srv.URL), header)
	c.Assert(err, gocheck.IsNil)
}**/
