package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-martini/martini"
	"github.com/tsuru/config"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/pkg"
	"github.com/wiliamsouza/apollo/token"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Assert(err, gocheck.IsNil)
	config.Set("rsa:public", "../data/keys/rsa.pub")
	config.Set("rsa:private", "../data/keys/rsa")
	if os.Getenv("MONGODB_URL") != "" {
		config.Set("database:url", os.Getenv("MONGODB_URL"))
	} else {
		config.Set("database:url", "127.0.0.1:27017")
	}
	config.Set("database:name", "apollo_api_tests")
	token.LoadKeys()
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestListPackages(c *gocheck.C) {
	results := `[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}},{"filename":"package2.tgz","metadata":{"description":"Package2 ON/OFF test"}}]`
	filename1 := "package1.tgz"
	pkgFile, err := os.Open("../data/" + filename1)
	c.Assert(err, gocheck.IsNil)
	metaFile, err := os.Open("../data/metadata1.json")
	c.Assert(err, gocheck.IsNil)
	_, err = pkg.NewPackage(pkgFile, metaFile, filename1)
	c.Assert(err, gocheck.IsNil)
	filename2 := "package2.tgz"
	pkgFile2, err := os.Open("../data/" + filename2)
	c.Assert(err, gocheck.IsNil)
	metaFile2, err := os.Open("../data/metadata2.json")
	c.Assert(err, gocheck.IsNil)
	_, err = pkg.NewPackage(pkgFile2, metaFile2, filename2)
	c.Assert(err, gocheck.IsNil)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer pkgFile2.Close()
	defer metaFile2.Close()
	defer db.Session.Package().Remove(filename1)
	defer db.Session.Package().Remove(filename2)
	request, err := http.NewRequest("GET", "tests/packages", nil)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	response := httptest.NewRecorder()
	ListPackages(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestUploadPackage(c *gocheck.C) {
	results := `{"filename":"bluetooth.jar","metadata":{"description":"Bluetooth ON/OFF test"}}`
	filename := "bluetooth.jar"
	metadata := "bluetooth.json"
	pkgFile, err := os.Open("../data/" + filename)
	c.Assert(err, gocheck.IsNil)
	metaFile, err := os.Open("../data/" + metadata)
	c.Assert(err, gocheck.IsNil)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	pkgPart, err := writer.CreateFormFile("package", filename)
	c.Assert(err, gocheck.IsNil)
	_, err = io.Copy(pkgPart, pkgFile)
	c.Assert(err, gocheck.IsNil)
	metaPart, err := writer.CreateFormFile("metadata", metadata)
	c.Assert(err, gocheck.IsNil)
	_, err = io.Copy(metaPart, metaFile)
	c.Assert(err, gocheck.IsNil)
	writer.Close()
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	request, err := http.NewRequest("POST", "tests/packages", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	response := httptest.NewRecorder()
	UploadPackage(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDetailPackage(c *gocheck.C) {
	results := `{"filename":"package1.tgz","metadata":{"version":0.1,"description":"Package1 ON/OFF test","install":"adb push dist/package1.jar /data/local/tmp/","run":"adb shell uiautomator runtest package1.jar -c com.github.wiliamsouza.package1.Package1Test"}}`
	request, err := http.NewRequest("GET",
		"test/package/package1.tgz", nil)
	c.Assert(err, gocheck.IsNil)
	email := "jhon@doe.com"
	t, err := token.New(email)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	filename := "package1.tgz"
	metadata := "metadata1.json"
	pkgFile, err := os.Open("../data/" + filename)
	c.Assert(err, gocheck.IsNil)
	metaFile, err := os.Open("../data/" + metadata)
	c.Assert(err, gocheck.IsNil)
	_, err = pkg.NewPackage(pkgFile, metaFile, filename)
	c.Assert(err, gocheck.IsNil)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["filename"] = filename
	params := martini.Params(p)
	response := httptest.NewRecorder()
	DetailPackage(response, request, tk, params)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDownloadPackage(c *gocheck.C) {
	request, err := http.NewRequest("GET",
		"tests/packages/downloads/package1.tgz", nil)
	c.Assert(err, gocheck.IsNil)
	email := "jhon@doe.com"
	t, err := token.New(email)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	filename := "package1.tgz"
	metadata := "metadata1.json"
	pkgFile, err := os.Open("../data/" + filename)
	c.Assert(err, gocheck.IsNil)
	metaFile, err := os.Open("../data/" + metadata)
	c.Assert(err, gocheck.IsNil)
	_, err = pkg.NewPackage(pkgFile, metaFile, filename)
	c.Assert(err, gocheck.IsNil)
	defer pkgFile.Close()
	defer metaFile.Close()
	defer db.Session.Package().Remove(filename)
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["filename"] = filename
	params := martini.Params(p)
	response := httptest.NewRecorder()
	DownloadPackage(response, request, tk, params)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	md5 := response.HeaderMap["Etag"][0]
	c.Assert(ct, gocheck.Equals, "application/octet-stream")
	pkgDb, err := db.Session.Package().Open(filename)
	c.Assert(err, gocheck.IsNil)
	md5Db := pkgDb.MD5()
	c.Assert(md5Db, gocheck.Equals, md5)
}
