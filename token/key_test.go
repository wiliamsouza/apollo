package token

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/globocom/config"

	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Assert(err, gocheck.IsNil)
	config.Set("rsa:public", "../data/keys/rsa.pub")
	config.Set("rsa:private", "../data/keys/rsa")
}

func (s *S) TestLoadKeys(c *gocheck.C) {
	public := bytes.Equal(PublicKey, []byte(""))
	c.Assert(public, gocheck.Equals, true)

	private := bytes.Equal(PrivateKey, []byte(""))
	c.Assert(private, gocheck.Equals, true)

	LoadKeys()

	publicFile, err := config.GetString("rsa:public")
	c.Assert(err, gocheck.IsNil)
	publicKey, err := ioutil.ReadFile(publicFile)
	c.Assert(err, gocheck.IsNil)
	loadedPublic := bytes.Equal(PublicKey, publicKey)
	c.Assert(loadedPublic, gocheck.Equals, true)

	privateFile, err := config.GetString("rsa:private")
	c.Assert(err, gocheck.IsNil)
	privateKey, err := ioutil.ReadFile(privateFile)
	c.Assert(err, gocheck.IsNil)
	loadedPrivate := bytes.Equal(PrivateKey, privateKey)
	c.Assert(loadedPrivate, gocheck.Equals, true)
}
