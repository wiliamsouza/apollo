package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wiliamsouza/apollo/customer"
)

// NewOrganization create new organization
func NewOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o customer.Organization
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	newOrganization, err := customer.NewOrganization(o)
	if err != nil {
		msg := "Error creating new machine, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(newOrganization)
	if err != nil {
		msg := "Error generating json result, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// ListOrganizations list organizations
func ListOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	organizations, err := customer.ListOrganizations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&organizations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// DetailOrganization detail organization
func DetailOrganization(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	organization, err := customer.DetailOrganization(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&organization)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// ModifyOrganization modify organization
func ModifyOrganization(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o customer.Organization
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	name := vars["name"]
	err = customer.ModifyOrganization(name, o)
	if err != nil {
		msg := "Error updating machine, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteOrganization delete organization
func DeleteOrganization(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	name := vars["name"]
	err := customer.RemoveOrganization(name)
	if err != nil {
		http.Error(w, "Error deleting organization, "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
