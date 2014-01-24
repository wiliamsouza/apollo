package customer

import (
	"testing"

	"github.com/globocom/config"
	"github.com/wiliamsouza/apollo/db"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollo.conf")
	c.Check(err, gocheck.IsNil)
	config.Set("database:name", "apollo_user_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewUser(c *gocheck.C) {
	email := "jhon@doe.com"
	user, _ := NewUser("Jhon Doe", email, "12345")
	defer db.Session.User().Remove(bson.M{"_id": email})
	var userDb User
	_ = db.Session.User().Find(bson.M{"_id": email}).One(&userDb)
	c.Assert(userDb.Name, gocheck.Equals, user.Name)
	c.Assert(userDb.Email, gocheck.Equals, user.Email)
	c.Assert(userDb.Password, gocheck.Equals, user.Password)
	c.Assert(userDb.ApiKey, gocheck.Equals, user.ApiKey)
}

func (s *S) TestEncryptPassword(c *gocheck.C) {
	password := `12345`
	email := "jhon@doe.com"
	user := &User{Name: "Jhon Doe", Email: email, Password: password}
	defer db.Session.User().Remove(bson.M{"_id": email})
	user.EncryptPassword()
	c.Assert(password, gocheck.Not(gocheck.Equals), user.Password)
}

// TODO: How to test ApiKey token generation?

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
