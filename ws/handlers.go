package ws

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"

	"github.com/wiliamsouza/apollo/customer"
)

// Web handler websocket for web side
func Web(w http.ResponseWriter, r *http.Request, p martini.Params) {
	APIKey := p["apikey"]
	_, err := customer.GetUserByAPIKey(APIKey)
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
	// TODO: FIX: "http://" used here! Maybe set "origin" option in /etc/apollos.conf
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
	conn := &wsConn{send: make(chan []byte, 256), ws: ws, APIKey: APIKey}
	Bridge.registerWeb <- conn
	defer func() { Bridge.unregisterWeb <- conn }()
	go conn.writer()
	conn.reader(Bridge.broadcastToAgent)
}

// Agent handler websocket for agent side
func Agent(w http.ResponseWriter, r *http.Request, p martini.Params) {
	APIKey := p["apikey"]
	_, err := customer.GetUserByAPIKey(APIKey)
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
	conn := &wsConn{send: make(chan []byte, 256), ws: ws, APIKey: APIKey}
	Bridge.registerAgent <- conn
	defer func() { Bridge.unregisterAgent <- conn }()
	go conn.writer()
	conn.reader(Bridge.broadcastToWeb)
}
