package pkg

import (
	"github.com/wiliamsouza/apollo/db"
	"github.com/globocom/config"
	"launchpad.net/gocheck"
	"testing"
)

setup := `{
	"name": "bluetooth"
	"version": 0.1
	"description": "Bluetooth ON/Off test"
	"install": "adb push dist/bluetooth.jar /data/local/tmp/"
	"run": "adb shell uiautomator runtest bluetooth.jar -c com.github.wiliamsouza.bluetooth.BluetoothTest"}`


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

/*
func (s *S) TestListPackagesShouldReturnAPackageList(c *gocheck.C) {
	var pkg Package
        pkg = &Package{Name: "bluetooth.tgz", Description: "Bluetooth test"}
	var pkglist PackageList
	pkglist = PackageList(pkg)
	pkgList := ListPackages()
	c.Assert(pkglist, gocheck.DeepEquals, pkgList)
}
*/
