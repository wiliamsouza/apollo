package pkg

import (
	"github.com/globocom/config"
	"github.com/wiliamsouza/apollo/db"
	"labix.org/v2/mgo/bson"
	"launchpad.net/gocheck"
	"os"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollo.conf")
	c.Check(err, gocheck.IsNil)
	config.Set("database:name", "apollo_pkg_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewPackage(c *gocheck.C) {
	filename := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/metadata1.json")
	pkg, _ := NewPackage(pkgFile, metaFile, filename)
	var pkgDb Package
	_ = db.Session.Package().Files.Find(bson.M{"filename": filename}).Select(bson.M{"filename": 1, "metadata.description": 1}).One(&pkgDb)
	c.Assert(pkgDb, gocheck.DeepEquals, pkg)
}

func (s *S) TestListPackages(c *gocheck.C) {
	filename1 := "package1.tgz"
	pkgFile1, _ := os.Open("../data/" + filename1)
	metaFile1, _ := os.Open("../data/metadata1.json")
	pkg1, _ := NewPackage(pkgFile1, metaFile1, filename1)

	filename2 := "package2.tgz"
	pkgFile2, _ := os.Open("../data/" + filename2)
	metaFile2, _ := os.Open("../data/metadata2.json")
	pkg2, _ := NewPackage(pkgFile2, metaFile2, filename2)

	// TODO: Description is not store in DB.
	//pkg1 := Package{Filename: "package1.tgz", Description: "Package1 ON/Off test"}

	pkgList := PackageList{pkg1, pkg2}
	pkgListDb, _ := ListPackages()
	c.Assert(pkgList, gocheck.DeepEquals, pkgListDb)
}
