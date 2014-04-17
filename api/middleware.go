package api

import (
	"net/http"

	"github.com/go-martini/martini"

	"github.com/wiliamsouza/apollo/token"
)

// AuthN handle json jwt token authentication and inject a token intance on the request.
func AuthN() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, context martini.Context) {
		t, err := token.Validate(req)
		tk := &token.Token{Email: t.Claims["email"].(string), Exp: t.Claims["exp"].(float64)}
		context.Map(tk)
		if err != nil {
			msg := "Error not authorized: "
			http.Error(res, msg+err.Error(), http.StatusUnauthorized)
		}
	}
}
