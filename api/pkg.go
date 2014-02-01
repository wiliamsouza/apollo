package api

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/wiliamsouza/apollo/pkg"
)

func ListPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UploadPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func DetailPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func DownloadPackage(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	io.Copy(w, pkg)

}
