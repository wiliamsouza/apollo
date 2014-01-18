package pkg

import (
	"encoding/json"
	"github.com/wiliamsouza/apollo/db"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"mime/multipart"
	"strings"
)

type Package struct {
	Filename    string `json:"filename"`
	Description string `json:"description"`
}

func (pkg *Package) String() string {
	parts := make([]string, 1, 2)
	if pkg.Description != "" {
		parts = append(parts, pkg.Description)
	}
	return strings.Join(parts, " ")
}

type PackageList []Package

func (packages PackageList) MarshalJSON() ([]byte, error) {
	m := make(map[string]string, len(packages))
	for _, pkg := range packages {
		m[pkg.Filename] = pkg.String()
	}
	return json.Marshal(m)
}

type Metadata struct {
	Version     float32 `json:"version"`
	Description string  `json:"description"`
	Install     string  `json:"install"`
	Run         string  `json:"run"`
}

func NewPackage(file, meta multipart.File, filename string) (Package, error) {
	pkg := Package{}
	gfs, err := db.Session.Package().Create(filename)
	_, err = io.Copy(gfs, file)
	data, err := ioutil.ReadAll(meta)
	metadata := Metadata{}
	json.Unmarshal(data, &metadata)
	gfs.SetMeta(metadata)
	pkg.Filename = filename
	pkg.Description = metadata.Description
	err = gfs.Close()
	return pkg, err
}

func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Files.Find(nil).Select(bson.M{"filename": 1, "metadata.description": 1}).Sort("filename").All(&packages)
	return PackageList(packages), err
}
