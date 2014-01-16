// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/globocom/config"
	"labix.org/v2/mgo"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "apollo_tests")
	Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	Session.DB.DropDatabase()
}

func (s *S) TestSessionPackageShouldReturnAMongoGridFS(c *gocheck.C) {
	var rep *mgo.GridFS
	rep = Session.Package()
	fsRep := Session.DB.GridFS("fs")
	c.Assert(rep, gocheck.DeepEquals, fsRep)
}

func (s *S) TestSessionCicleShouldReturnAMongoCollection(c *gocheck.C) {
	var cicle *mgo.Collection
	cicle = Session.Cicle()
	cCicle := Session.DB.C("cicle")
	c.Assert(cicle, gocheck.DeepEquals, cCicle)
}

func (s *S) TestSessionPlanShouldReturnAMongoCollection(c *gocheck.C) {
	var plan *mgo.Collection
	plan = Session.Plan()
	cPlan := Session.DB.C("plan")
	c.Assert(plan, gocheck.DeepEquals, cPlan)
}

func (s *S) TestSessionCaseShouldReturnAMongoCollection(c *gocheck.C) {
	var ccase *mgo.Collection
	ccase = Session.Case()
	cCase := Session.DB.C("case")
	c.Assert(ccase, gocheck.DeepEquals, cCase)
}

func (s *S) TestConnect(c *gocheck.C) {
	Connect()
	c.Assert(Session.DB.Name, gocheck.Equals, "apollo_tests")
	err := Session.DB.Session.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestConnectDefaultSettings(c *gocheck.C) {
	oldURL, _ := config.Get("database:url")
	defer config.Set("database:url", oldURL)
	oldName, _ := config.Get("database:name")
	defer config.Set("database:name", oldName)
	config.Unset("database:url")
	config.Unset("database:name")
	Connect()
	c.Assert(Session.DB.Name, gocheck.Equals, "apollo")
	c.Assert(Session.DB.Session.LiveServers(), gocheck.DeepEquals, []string{"localhost:27017"})
}
