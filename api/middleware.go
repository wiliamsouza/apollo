package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/wiliamsouza/apollo/token"
)

// AuthNHandleFunc handle jwt token authentication.
type AuthNHandleFunc func(http.ResponseWriter, *http.Request, *jwt.Token)

func (h AuthNHandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := token.Validate(r)
	if err != nil {
		msg := "Error not authorized: "
		http.Error(w, msg+err.Error(), http.StatusUnauthorized)
		return
	}
	h(w, r, token)
}

// CORS Cross-Origin Resource Sharing.
type CORS struct {
	handler http.Handler
}

// CORSHandle handle Cross-Origin Resource Sharing.
func CORSHandle(handler http.Handler) *CORS {
	return &CORS{handler: handler}
}

func (s *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Get allow origin list from apollod.conf
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s.handler.ServeHTTP(w, r)
}

// PreFlightHandleFunc handle browser preflight OPTIONS requests.
type PreFlightHandleFunc func(http.ResponseWriter, *http.Request)

func (h PreFlightHandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "OPTIONS" {
		h(w, r)
		return
	}
	// TODO: Get allow headers list from apollod.conf
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	w.WriteHeader(http.StatusOK)
	return
}
