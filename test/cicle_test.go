package test

import (
	"os"
	"testing"

	"github.com/tsuru/config"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Check(err, gocheck.IsNil)
	if os.Getenv("MONGODB_URL") != "" {
		config.Set("database:url", os.Getenv("MONGODB_URL"))
	} else {
		config.Set("database:url", "127.0.0.1:27017")
	}
	config.Set("database:name", "apollo_test_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewCicle(c *gocheck.C) {
	cicle := Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test, err := NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(test.Id)
	var cicleDb Cicle
	err = db.Session.Cicle().FindId(test.Id).One(&cicleDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(cicleDb, gocheck.DeepEquals, test)
}

func (s *S) TestListCicles(c *gocheck.C) {
	cicle1 := Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	cicle2 := Cicle{Name: "Test yakju", Device: "yakju",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test1, err := NewCicle(cicle1)
	c.Assert(err, gocheck.IsNil)
	test2, err := NewCicle(cicle2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(test1.Id)
	defer db.Session.Cicle().RemoveId(test2.Id)
	cicleList := CicleList{test1, test2}
	cicleListDb, err := ListCicles()
	c.Assert(err, gocheck.IsNil)
	c.Assert(cicleListDb, gocheck.DeepEquals, cicleList)
}

func (s *S) TestDetailCicle(c *gocheck.C) {
	cicle := Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test, err := NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(test.Id)
	cicleDb, err := DetailCicle(test.Id.Hex())
	c.Assert(err, gocheck.IsNil)
	c.Assert(cicleDb, gocheck.DeepEquals, test)
}

func (s *S) TestmodifyCicle(c *gocheck.C) {
	cicle1 := Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test, err := NewCicle(cicle1)
	cicle2 := Cicle{Id: test.Id, Name: "Test yakju", Device: "yakju",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(test.Id)
	err = ModifyCicle(test.Id.Hex(), cicle2)
	c.Assert(err, gocheck.IsNil)
	var cicleDb Cicle
	err = db.Session.Cicle().FindId(test.Id).One(&cicleDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(cicleDb, gocheck.DeepEquals, cicle2)
}

func (s *S) TestRemoveCicle(c *gocheck.C) {
	cicle := Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test, err := NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	err = RemoveCicle(test.Id.Hex())
	c.Assert(err, gocheck.IsNil)
	lenght, err := db.Session.Cicle().FindId(test.Id).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}
