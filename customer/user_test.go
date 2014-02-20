package customer

import (
	"testing"

	"github.com/globocom/config"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Check(err, gocheck.IsNil)
	config.Set("database:name", "apollo_customer_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewUser(c *gocheck.C) {
	email := "jhon@doe.com"
	user, err := NewUser("Jhon Doe", email, "12345")
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(email)
	var userDb User
	err = db.Session.User().FindId(email).One(&userDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(userDb.Name, gocheck.Equals, user.Name)
	c.Assert(userDb.Email, gocheck.Equals, user.Email)
	c.Assert(userDb.Password, gocheck.Equals, user.Password)
	c.Assert(userDb.APIKey, gocheck.Equals, user.APIKey)
}

func (s *S) TestEncryptPassword(c *gocheck.C) {
	password := `12345`
	email := "jhon@doe.com"
	user := &User{Name: "Jhon Doe", Email: email, Password: password}
	defer db.Session.User().RemoveId(email)
	user.EncryptPassword()
	c.Assert(password, gocheck.Not(gocheck.Equals), user.Password)
}

// TODO: How to test APIKey token generation?

func (s *S) TestValidateEmail(c *gocheck.C) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"jhon@gmail.com", true},
		{"doe@apollolab.com.br", true},
		{"jane+doe@gmail.com", true},
		{"janie2", false},
		{"g4oph4er", false},
		{"g0o-ph3er", false},
	}
	for _, t := range tests {
		u := User{Email: t.input}
		v, _ := u.ValidateEmail()
		if v != t.expected {
			c.Errorf("Is %q valid? Want %v. Got %v.", t.input, t.expected, v)
		}
	}
}

func (s *S) TestGetUserByAPIKey(c *gocheck.C) {
	email := "jhon@doe.com"
	user, err := NewUser("Jhon Doe", email, "12345")
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(email)
	userK, err := GetUserByAPIKey(user.APIKey)
	c.Assert(err, gocheck.IsNil)
	c.Assert(userK.APIKey, gocheck.DeepEquals, user.APIKey)
}

func (s *S) TestDetailUser(c *gocheck.C) {
	email := "jhon@doe.com"
	user, err := NewUser("Jhon Doe", email, "12345")
	c.Assert(err, gocheck.IsNil)
	defer db.Session.User().RemoveId(email)
	userDb, err := DetailUser(email)
	c.Assert(err, gocheck.IsNil)
	c.Assert(userDb.Name, gocheck.Equals, user.Name)
	c.Assert(userDb.Email, gocheck.Equals, user.Email)
	c.Assert(userDb.Password, gocheck.Equals, user.Password)
	c.Assert(userDb.APIKey, gocheck.Equals, user.APIKey)
}
