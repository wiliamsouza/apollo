package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Web handler websocket from web page side
func Web(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}
	// TODO: FIX: "http://" used here! Maybe set "origin" option in /etc/apollo.conf
	if origin := r.Header.Get("Origin"); origin != "http://"+r.Host {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ws.Close()
	/**for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = ws.WriteMessage(messageType, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}**/
}

// Runner handler websocket from runner side
func Runner(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ws.Close()
	/**for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = ws.WriteMessage(messageType, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}**/
}
