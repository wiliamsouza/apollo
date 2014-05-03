package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"

	"github.com/wiliamsouza/apollo/pkg"
	"github.com/wiliamsouza/apollo/token"
)

// ListPackages list packages
func ListPackages(w http.ResponseWriter, r *http.Request, token *token.Token) {
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

// UploadPackage upload package
func UploadPackage(w http.ResponseWriter, r *http.Request, token *token.Token) {
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

// DetailPackage detail package
func DetailPackage(w http.ResponseWriter, r *http.Request, token *token.Token, p martini.Params) {
	w.Header().Set("Content-Type", "application/json")
	filename := p["filename"]
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

// DownloadPackage download package
func DownloadPackage(w http.ResponseWriter, r *http.Request, token *token.Token, p martini.Params) {
	filename := p["filename"]
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
