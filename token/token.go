package token

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token is thin layer to jwt.Token and is used to be inject
type Token struct {
	Email string
	Exp   float64
}

// New generate a JWT token in string format
func New(email string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims["email"] = email
	t.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token, err := t.SignedString(PublicKey)
	return token, err
}

// Validate a token try parser from Authorization header
// or access_token parameter
func Validate(r *http.Request) (*jwt.Token, error) {
	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) ([]byte, error) {
		return PublicKey, nil
	})
	return token, err
}
