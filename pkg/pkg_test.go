package pkg

import (
	"github.com/wiliamsouza/apollo/db"
	"github.com/globocom/config"
	"launchpad.net/gocheck"
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

//func (s *S) TestListPackagesShouldReturnAPackageList(c *gocheck.C) {
//	var pkg Package
//        pkg = &Package{Name: "bluetooth.tgz", Description: "Bluetooth test"}
//	var pkglist PackageList
//	pkglist = PackageList(pkg)
//	pkgList := ListPackages()
//	c.Assert(pkglist, gocheck.DeepEquals, pkgList)
//}
