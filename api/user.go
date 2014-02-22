package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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

type detailUser struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	APIKey    string    `json:"apikey"`
	Created   time.Time `json:"created"`
	LastLogin time.Time `json:"lastlogin"`
}

// NewUser create new user
func NewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parssing request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var u requestUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, "Error parssing json request: "+err.Error(), http.StatusBadRequest)
		return
	}
	newUser, err := customer.NewUser(u.Name, u.Email, u.Password)
	if err != nil {
		http.Error(w, "Error creating new user: "+err.Error(), http.StatusBadRequest)
		return
	}
	user := responseUser{Name: newUser.Name, Email: newUser.Email}
	result, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, "Error generating json result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// DetailUser detail user
func DetailUser(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, err := customer.DetailUser(vars["email"])
	if err != nil {
		http.Error(w, "Error getting user detail: "+err.Error(), http.StatusNotFound)
		return
	}
	response := detailUser{Name: user.Name, Email: user.Email, APIKey: user.APIKey, Created: user.Created, LastLogin: user.LastLogin}
	result, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "Error generating json result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
