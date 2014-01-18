package pkg

import (
	"encoding/json"
	"github.com/wiliamsouza/apollo/db"
	"io"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"mime/multipart"
)

type Package struct {
	Filename string `json:"filename"`
	Metadata Meta   `json:"metadata"`
}

type Meta struct {
	Description string `json:"description"`
}

type PackageList []Package

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
	err = gfs.Close()
	_ = db.Session.Package().Files.Find(bson.M{"filename": filename}).Select(bson.M{"filename": 1, "metadata.description": 1, "_id": 0}).One(&pkg)
	return pkg, err
}

func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Files.Find(nil).Select(bson.M{"filename": 1, "metadata.description": 1, "_id": 0}).Sort("filename").All(&packages)
	return PackageList(packages), err
}
