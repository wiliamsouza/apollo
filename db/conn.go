// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package db provides util functions to deal with Gandalf's database.
package db

import (
	"github.com/globocom/config"
	"labix.org/v2/mgo"
)

type session struct {
	DB *mgo.Database
}

// The global Session that must be used by users.
var Session session

// Connect uses database:url and database:name settings in config file and
// connects to the database. If it cannot connect or these settings are not
// defined, it will panic.
func Connect() {
	url, _ := config.GetString("database:url")
	if url == "" {
		url = "localhost:27017"
	}
	name, _ := config.GetString("database:name")
	if name == "" {
		name = "apollo"
	}
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	Session.DB = s.DB(name)
}

// Package returns a reference to the "package" GridFS in MongoDB.
func (s *session) Package() *mgo.GridFS {
	return s.DB.GridFS("fs")
}

// Cicle returns a reference to the "cicle" collection in MongoDB.
func (s *session) Cicle() *mgo.Collection {
	return s.DB.C("cicle")
}

// Plan returns a reference to the "plan" collection in MongoDB.
func (s *session) Plan() *mgo.Collection {
	return s.DB.C("plan")
}

// Case returns a reference to the "case" collection in MongoDB.
func (s *session) Case() *mgo.Collection {
	return s.DB.C("case")
}