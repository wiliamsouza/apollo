package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"

	"github.com/wiliamsouza/apollo/test"
)

// NewCicle create new cicle
func NewCicle(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o test.Cicle
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	newCicle, err := test.NewCicle(o)
	if err != nil {
		msg := "Error creating new cicle, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(newCicle)
	if err != nil {
		msg := "Error generating json result, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// ListCicles list cicles
func ListCicles(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	cicles, err := test.ListCicles()
	if err != nil {
		msg := "Error getting cicle list: "
		http.Error(w, msg+err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&cicles)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// DetailCicle detail cicle
func DetailCicle(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := filepath.Base(r.URL.Path)
	cicle, err := test.DetailCicle(id)
	if err != nil {
		msg := "Error getting cicle detail: "
		http.Error(w, msg+err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&cicle)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// ModifyCicle modify cicle
func ModifyCicle(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o test.Cicle
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	id := filepath.Base(r.URL.Path)
	err = test.ModifyCicle(id, o)
	if err != nil {
		msg := "Error updating cicle, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteCicle delete cicle
func DeleteCicle(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	id := filepath.Base(r.URL.Path)
	err := test.RemoveCicle(id)
	if err != nil {
		http.Error(w, "Error deleting cicle, "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
