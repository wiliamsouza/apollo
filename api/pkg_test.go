package api

import (
	"github.com/globocom/config"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/pkg"
	"launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollo.conf")
	c.Assert(err, gocheck.IsNil)
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "apollo_api_tests")
	db.Connect()
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestPackageList(c *gocheck.C) {
	filename := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/metadata1.json")
	_, _ = pkg.NewPackage(pkgFile, metaFile, filename)

	filename2 := "package2.tgz"
	pkgFile2, _ := os.Open("../data/" + filename2)
	metaFile2, _ := os.Open("../data/metadata2.json")
	_, _ = pkg.NewPackage(pkgFile2, metaFile2, filename2)

	results := `[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}},{"filename":"package2.tgz","metadata":{"description":"Package2 ON/OFF test"}}]`
	request, _ := http.NewRequest("GET", "tests/packages", nil)
	response := httptest.NewRecorder()
	packageList(response, request)

	c.Assert(response.Code, gocheck.Equals, http.StatusOK)

	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")

	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestPackageUpload(c *gocheck.C) {
	request, _ := http.NewRequest("POST", "tests/packages", nil)
	response := httptest.NewRecorder()
	packageUpload(response, request)

	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)

	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
}

func (s *S) TestPackageDownload(c *gocheck.C) {
	request, _ := http.NewRequest("GET", "tests/packages/pkg.tgz", nil)
	response := httptest.NewRecorder()
	packageDownload(response, request)

	c.Assert(response.Code, gocheck.Equals, http.StatusOK)

	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
}
