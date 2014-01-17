package pkg

import (
	"github.com/wiliamsouza/apollo/db"
	"encoding/json"
	"mime/multipart"
	"strings"
	"io"
	"io/ioutil"
)

type Package struct {
	Name string
	Description string
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
                m[pkg.Name] = pkg.String()
        }
        return json.Marshal(m)
}

type Metadata struct {
	Name string `json:"name"`
	Version float32 `json:"version"`
	Description string `json:"description"`
	Install string `json:"install"`
	Run string `json:"run"`
}

func NewPackage(file, meta, multipart.File, filename string) (Package, error) {
	gfs, err : db.Session.Package().Create(filename)

	err = io.Copy(gfs, file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(meta)
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{}
	json.Unmarshal(data, &metadata)
	gfs.SetMeta(metadata)

	pkg := Package{}
	pkg.Name = filename
	pkg.Description = metadata.Description

	err = gfs.Close()
	if err != nil {
		return nil, err
	}

	return pkg, err
}

func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Find(nil).Sort("filename").All(&packages)
	return PackageList(packages), err
}
