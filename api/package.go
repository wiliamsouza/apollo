package api

import (
	"net/http"
)

func packageList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func packageUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func packageDownload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
