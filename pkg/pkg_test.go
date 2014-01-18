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

/*
func (s *S) TestListPackages(c *gocheck.C) {
	var pkg Package
        pkg = &Package{Filename: "bluetooth.tgz", Description: "Bluetooth test"}
	var pkglist PackageList
	pkglist = PackageList(pkg)
	pkgList := ListPackages()
	c.Assert(pkglist, gocheck.DeepEquals, pkgList)
}
*/
