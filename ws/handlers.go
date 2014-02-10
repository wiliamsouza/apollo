package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/wiliamsouza/apollo/customer"
)

// Web handler websocket for web side
func Web(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	APIKey := vars["apikey"]
	u, err := customer.GetUserByAPIKey(APIKey)
	if err != nil {
		msg := "Invalid APIKey, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method != "GET" {
		msg := "Method not allowed"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// TODO: FIX: "http://" used here! Maybe set "origin" option in /etc/apollo.conf
	/**if origin := r.Header.Get("Origin"); origin != "http://"+r.Host {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}**/
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		msg := "Not a websocket handshake"
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//ws.Close()
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = ws.WriteMessage(messageType, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Runner handler websocket for runner side
func Runner(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	APIKey := vars["apikey"]
	u, err := customer.GetUserByAPIKey(APIKey)
	if err != nil {
		msg := "Invalid APIKey, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method != "GET" {
		msg := "Method not allowed"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		msg := "Not a websocket handshake"
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//ws.Close()
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = ws.WriteMessage(messageType, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
