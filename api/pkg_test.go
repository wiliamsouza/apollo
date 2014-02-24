package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/globocom/config"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/pkg"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Assert(err, gocheck.IsNil)
	config.Set("rsa:public", "../data/keys/rsa.pub")
	config.Set("rsa:private", "../data/keys/rsa")
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "apollo_api_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestListPackages(c *gocheck.C) {
	results := `[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}},{"filename":"package2.tgz","metadata":{"description":"Package2 ON/OFF test"}}]`
	filename1 := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename1)
	metaFile, _ := os.Open("../data/metadata1.json")
	_, _ = pkg.NewPackage(pkgFile, metaFile, filename1)
	filename2 := "package2.tgz"
	pkgFile2, _ := os.Open("../data/" + filename2)
	metaFile2, _ := os.Open("../data/metadata2.json")
	_, _ = pkg.NewPackage(pkgFile2, metaFile2, filename2)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer pkgFile2.Close()
	defer metaFile2.Close()
	defer db.Session.Package().Remove(filename1)
	defer db.Session.Package().Remove(filename2)
	request, _ := http.NewRequest("GET", "tests/packages", nil)
	response := httptest.NewRecorder()
	ListPackages(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestUploadPackage(c *gocheck.C) {
	results := `{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}`
	filename := "package1.tgz"
	metadata := "metadata1.json"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/" + metadata)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	pkgPart, _ := writer.CreateFormFile("package", filename)
	metaPart, _ := writer.CreateFormFile("metadata", metadata)
	_, _ = io.Copy(pkgPart, pkgFile)
	_, _ = io.Copy(metaPart, metaFile)
	writer.Close()
	request, _ := http.NewRequest("POST", "tests/packages", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	UploadPackage(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDetailPackage(c *gocheck.C) {
	results := `{"filename":"package1.tgz","metadata":{"version":0.1,"description":"Package1 ON/OFF test","install":"adb push dist/package1.jar /data/local/tmp/","run":"adb shell uiautomator runtest package1.jar -c com.github.wiliamsouza.package1.Package1Test"}}`
	request, _ := http.NewRequest("GET", "test/package/package1.tgz", nil)
	filename := "package1.tgz"
	metadata := "metadata1.json"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/" + metadata)
	_, _ = pkg.NewPackage(pkgFile, metaFile, filename)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	response := httptest.NewRecorder()
	DetailPackage(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)

}

func (s *S) TestDownloadPackage(c *gocheck.C) {
	request, _ := http.NewRequest("GET", "tests/packages/downloads/package1.tgz", nil)
	filename := "package1.tgz"
	metadata := "metadata1.json"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/" + metadata)
	_, _ = pkg.NewPackage(pkgFile, metaFile, filename)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	response := httptest.NewRecorder()
	DownloadPackage(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	md5 := response.HeaderMap["Etag"][0]
	c.Assert(ct, gocheck.Equals, "application/octet-stream")
	pkgDb, _ := db.Session.Package().Open(filename)
	md5Db := pkgDb.MD5()
	c.Assert(md5Db, gocheck.Equals, md5)
}
