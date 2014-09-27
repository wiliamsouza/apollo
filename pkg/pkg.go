package pkg

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/wiliamsouza/apollo/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Package represents GridFS file
type Package struct {
	Filename string   `json:"filename"`
	Metadata Metadata `json:"metadata"`
}

// Metadata represents GridFS file metadata
type Metadata struct {
	Version     float32 `json:"version,omitempty"`
	Description string  `json:"description"`
	Install     string  `json:"install,omitempty"`
	Run         string  `json:"run,omitempty"`
}

// PackageList a list of packages
type PackageList []Package

// NewPackage create a new package
func NewPackage(file, meta multipart.File, filename string) (Package, error) {
	pkg := Package{}
	gfs, err := db.Session.Package().Create(filename)
	_, err = io.Copy(gfs, file)
	data, err := ioutil.ReadAll(meta)
	metadata := Metadata{}
	json.Unmarshal(data, &metadata)
	gfs.SetMeta(metadata)
	err = gfs.Close()
	_ = db.Session.Package().Files.Find(bson.M{"filename": filename}).Select(bson.M{"filename": 1, "metadata.description": 1}).One(&pkg)
	return pkg, err
}

// ListPackages list packages
func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Files.Find(nil).Select(bson.M{"filename": 1, "metadata.description": 1}).Sort("filename").All(&packages)
	return PackageList(packages), err
}

// DetailPackage detail package
func DetailPackage(filename string) (Package, error) {
	var pkg Package
	err := db.Session.Package().Files.Find(bson.M{"filename": filename}).One(&pkg)
	return pkg, err
}

// GetPackage get a package
func GetPackage(filename string) (*mgo.GridFile, error) {
	file, err := db.Session.Package().Open(filename)
	return file, err
}
