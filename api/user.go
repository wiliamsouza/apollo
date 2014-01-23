package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wiliamsouza/apollo/customer"
)

type requestUser struct {
	Name     string
	Email    string
	Password string
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var u requestUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, "Error parssing json request: "+err.Error(), http.StatusInternalServerError)
	}
	user, err := customer.NewUser(u.Name, u.Email, u.Password)
	fmt.Println(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	result, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(result)
}
