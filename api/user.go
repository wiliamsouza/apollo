package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/wiliamsouza/apollo/customer"
)

type requestUser struct {
	Name     string
	Email    string
	Password string
}

type responseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parssing request body: "+err.Error(), http.StatusInternalServerError)
	}
	var u requestUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, "Error parssing json request: "+err.Error(), http.StatusInternalServerError)
	}
	newUser, err := customer.NewUser(u.Name, u.Email, u.Password)
	if err != nil {
		http.Error(w, "Error creating new user: "+err.Error(), http.StatusInternalServerError)
	}
	user := responseUser{Name: newUser.Name, Email: newUser.Email}
	result, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error generating json result: "+err.Error(), http.StatusInternalServerError)
	}
	w.Write(result)
}
