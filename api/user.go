package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-martini/martini"

	"github.com/wiliamsouza/apollo/customer"
	"github.com/wiliamsouza/apollo/token"
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

type responseToken struct {
	Token string `json:"token"`
}

type authenticateUser struct {
	Email    string
	Password string
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var u requestUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		msg := "Error parssing json request: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	newUser, err := customer.NewUser(u.Name, u.Email, u.Password)
	if err != nil {
		msg := "Error creating new user: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	user := responseUser{Name: newUser.Name, Email: newUser.Email}
	result, err := json.Marshal(&user)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// DetailUser detail user
func DetailUser(w http.ResponseWriter, r *http.Request, t *token.Token, p martini.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	email := p["email"]
	if email != t.Email {
		msg := "Error getting user detail other user: "
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	user, err := customer.DetailUser(email)
	if err != nil {
		msg := "Error getting user detail: "
		http.Error(w, msg+err.Error(), http.StatusNotFound)
		return
	}
	response := detailUser{Name: user.Name, Email: user.Email,
		APIKey: user.APIKey, Created: user.Created,
		LastLogin: user.LastLogin}
	result, err := json.Marshal(&response)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Authenticate user using email and password and issue a JSON Web Token
func Authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var u authenticateUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		msg := "Error parssing json request: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	user, err := customer.GetUserByEmail(u.Email)
	if err != nil {
		msg := "Error user not found: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	err = user.ValidatePassword(u.Password)
	if err != nil {
		msg := "Error invalid password: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	newToken, err := token.New(user.Email)
	if err != nil {
		msg := "Error generating token: "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	t := responseToken{Token: newToken}
	result, err := json.Marshal(&t)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
