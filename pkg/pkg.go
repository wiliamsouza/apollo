package pkg

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/wiliamsouza/apollo/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Package struct {
	Filename string   `json:"filename"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Version     float32 `json:"version,omitempty"`
	Description string  `json:"description"`
	Install     string  `json:"install,omitempty"`
	Run         string  `json:"run,omitempty"`
}

type PackageList []Package

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

func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Files.Find(nil).Select(bson.M{"filename": 1, "metadata.description": 1}).Sort("filename").All(&packages)
	return PackageList(packages), err
}

func DetailPackage(filename string) (Package, error) {
	var pkg Package
	err := db.Session.Package().Files.Find(bson.M{"filename": filename}).One(&pkg)
	return pkg, err
}

func GetPackage(filename string) (*mgo.GridFile, error) {
	file, err := db.Session.Package().Open(filename)
	return file, err
}
