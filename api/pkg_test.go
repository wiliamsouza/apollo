package api

import (
	"bytes"
	"fmt"
	"github.com/globocom/config"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/pkg"
	"io"
	"launchpad.net/gocheck"
	"mime/multipart"
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

func (s *S) TestListPackages(c *gocheck.C) {
	results := `[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}},{"filename":"package2.tgz","metadata":{"description":"Package2 ON/OFF test"}}]`
	filename := "package1.tgz"
	pkgFile, _ := os.Open("../data/" + filename)
	metaFile, _ := os.Open("../data/metadata1.json")
	_, _ = pkg.NewPackage(pkgFile, metaFile, filename)
	filename2 := "package2.tgz"
	pkgFile2, _ := os.Open("../data/" + filename2)
	metaFile2, _ := os.Open("../data/metadata2.json")
	_, _ = pkg.NewPackage(pkgFile2, metaFile2, filename2)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer pkgFile2.Close()
	defer metaFile2.Close()
	request, _ := http.NewRequest("GET", "tests/packages", nil)
	response := httptest.NewRecorder()
	listPackages(response, request)
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
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	pkgPart, err := writer.CreateFormFile("package", filename)
	if err != nil {
		fmt.Println("CreateFormFile")
	}
	metaPart, _ := writer.CreateFormFile("metadata", metadata)
	_, _ = io.Copy(pkgPart, pkgFile)
	_, _ = io.Copy(metaPart, metaFile)
	writer.Close()
	request, _ := http.NewRequest("POST", "tests/packages", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	uploadPackage(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDownloadPackage(c *gocheck.C) {
	request, _ := http.NewRequest("GET", "tests/packages/pkg.tgz", nil)
	response := httptest.NewRecorder()
	downloadPackage(response, request)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
}
