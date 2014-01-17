package pkg

import (
	"github.com/wiliamsouza/apollo/db"
	"encoding/json"
	"strings"
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

func ListPackages() (PackageList, error) {
	var packages []Package
	err := db.Session.Package().Find(nil).Sort("_id").All(&packages)
	return PackageList(packages), err
}
