// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"

	"github.com/tsuru/config"
	"gopkg.in/mgo.v2"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	if os.Getenv("MONGODB_URL") != "" {
		config.Set("database:url", os.Getenv("MONGODB_URL"))
	} else {
		config.Set("database:url", "127.0.0.1:27017")
	}
	config.Set("database:name", "apollo_db_tests")
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

func (s *S) TestSessionOrganizationShouldReturnAMongoCollection(c *gocheck.C) {
	var org *mgo.Collection
	org = Session.Organization()
	cOrg := Session.DB.C("organization")
	c.Assert(org, gocheck.DeepEquals, cOrg)
}

func (s *S) TestSessionUserShouldReturnAMongoCollection(c *gocheck.C) {
	var user *mgo.Collection
	user = Session.User()
	cUser := Session.DB.C("user")
	c.Assert(user, gocheck.DeepEquals, cUser)
}

func (s *S) TestSessionDeviceShouldReturnAMongoCollection(c *gocheck.C) {
	var device *mgo.Collection
	device = Session.Device()
	cDevice := Session.DB.C("device")
	c.Assert(device, gocheck.DeepEquals, cDevice)
}

func (s *S) TestConnect(c *gocheck.C) {
	Connect()
	c.Assert(Session.DB.Name, gocheck.Equals, "apollo_db_tests")
	err := Session.DB.Session.Ping()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestConnectDefaultSettings(c *gocheck.C) {
	if os.Getenv("MONGODB_URL") != "" {
		c.Skip("Running inside a CI")
	}
	oldURL, _ := config.Get("database:url")
	defer config.Set("database:url", oldURL)
	oldName, _ := config.Get("database:name")
	defer config.Set("database:name", oldName)
	config.Unset("database:url")
	config.Unset("database:name")
	Connect()
	c.Assert(Session.DB.Name, gocheck.Equals, "apollo")
	c.Assert(Session.DB.Session.LiveServers(), gocheck.DeepEquals, []string{"127.0.0.1:27017"})
}
