package pkg

import (
	"os"
	"testing"

	"github.com/tsuru/config"
	"github.com/wiliamsouza/apollo/db"
	"gopkg.in/mgo.v2/bson"
	"launchpad.net/gocheck"
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
	defer db.Session.Package().Remove(filename)
	var pkgDb Package
	_ = db.Session.Package().Files.Find(bson.M{"filename": filename}).Select(bson.M{"filename": 1, "metadata.description": 1}).One(&pkgDb)
	c.Assert(pkgDb, gocheck.DeepEquals, pkg)
}

func (s *S) TestListPackages(c *gocheck.C) {
	filename2 := "package2.tgz"
	pkgFile2, _ := os.Open("../data/" + filename2)
	metaFile2, _ := os.Open("../data/metadata2.json")
	pkg2, _ := NewPackage(pkgFile2, metaFile2, filename2)
	filename3 := "package3.tgz"
	pkgFile3, _ := os.Open("../data/" + filename3)
	metaFile3, _ := os.Open("../data/metadata3.json")
	pkg3, _ := NewPackage(pkgFile3, metaFile3, filename3)
	defer db.Session.Package().Remove(filename2)
	defer db.Session.Package().Remove(filename3)
	pkgList := PackageList{pkg2, pkg3}
	pkgListDb, _ := ListPackages()
	c.Assert(pkgList, gocheck.DeepEquals, pkgListDb)
}

func (s *S) TestDetailPackage(c *gocheck.C) {
	filename := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/metadata1.json")
	_, _ = NewPackage(pkgFile, metaFile, filename)
	pkg, _ := DetailPackage(filename)
	var pkgDb Package
	_ = db.Session.Package().Files.Find(bson.M{"filename": filename}).One(&pkgDb)
	defer db.Session.Package().Remove(filename)
	c.Assert(pkgDb, gocheck.DeepEquals, pkg)
}

func (s *S) TestGetPackage(c *gocheck.C) {
	filename := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/metadata1.json")
	_, _ = NewPackage(pkgFile, metaFile, filename)
	pkg, _ := GetPackage(filename)
	md5 := pkg.MD5()
	pkgDb, _ := db.Session.Package().Open(filename)
	md5Db := pkgDb.MD5()
	defer db.Session.Package().Remove(filename)
	c.Assert(md5Db, gocheck.Equals, md5)
}
