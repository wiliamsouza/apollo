package api

import (
	"encoding/json"
	"github.com/wiliamsouza/apollo/pkg"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

func listPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	packages, err := pkg.ListPackages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&packages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func uploadPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	pkgFile, pkgHeader, err := r.FormFile("package")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metaFile, _, err := r.FormFile("metadata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newPkg, err := pkg.NewPackage(pkgFile, metaFile, pkgHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(&newPkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func detailPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	filename := filepath.Base(r.URL.Path)
	pkg, err := pkg.DetailPackage(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&pkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func downloadPackage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	filename := filepath.Base(r.URL.Path)
	pkg, err := pkg.GetPackage(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pkg.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(pkg.Size(), 10))
	w.Header().Set("ETag", pkg.MD5())
	io.Copy(w, pkg)

}
