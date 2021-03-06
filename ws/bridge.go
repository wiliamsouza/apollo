package ws

import (
	"github.com/gorilla/websocket"
)

type message struct {
	msg    []byte
	APIKey string
}

type wsConn struct {
	ws     *websocket.Conn
	send   chan []byte
	APIKey string
}

func (conn *wsConn) reader(broadcast chan *message) {
	for {
		_, msg, err := conn.ws.ReadMessage()
		if err != nil {
			break
		}
		broadcast <- &message{msg: msg, APIKey: conn.APIKey}
	}
	conn.ws.Close()
}

func (conn *wsConn) writer() {
	for msg := range conn.send {
		err := conn.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	conn.ws.Close()
}

type connection struct {
	webConn   []*wsConn
	agentConn []*wsConn
}

type bridge struct {
	connections      map[string]*connection
	registerWeb      chan *wsConn
	unregisterWeb    chan *wsConn
	registerAgent    chan *wsConn
	unregisterAgent  chan *wsConn
	broadcastToWeb   chan *message
	broadcastToAgent chan *message
}

// Bridge store connections on both sides
var Bridge = bridge{
	connections:      make(map[string]*connection),
	registerWeb:      make(chan *wsConn),
	unregisterWeb:    make(chan *wsConn),
	registerAgent:    make(chan *wsConn),
	unregisterAgent:  make(chan *wsConn),
	broadcastToWeb:   make(chan *message),
	broadcastToAgent: make(chan *message),
}

// Run is used to coordinate all channel message exchange.
func (b *bridge) Run() {
	for {
		select {
		case conn := <-b.registerWeb:
			if c, ok := b.connections[conn.APIKey]; ok {
				c.webConn = append(c.webConn, conn)
			} else {
				c := &connection{webConn: []*wsConn{conn}, agentConn: []*wsConn{}}
				b.connections[conn.APIKey] = c
			}
		case conn := <-b.unregisterWeb:
			_ = b.connections[conn.APIKey]
			//delete(c.webConn, conn)
			close(conn.send)
		case conn := <-b.registerAgent:
			if c, ok := b.connections[conn.APIKey]; ok {
				c.agentConn = append(c.agentConn, conn)
			} else {
				c := &connection{webConn: []*wsConn{}, agentConn: []*wsConn{conn}}
				b.connections[conn.APIKey] = c
			}
		case conn := <-b.unregisterAgent:
			_ = b.connections[conn.APIKey]
			//delete(c.agentConn, conn)
			close(conn.send)
		case m := <-b.broadcastToWeb:
			if c, ok := b.connections[m.APIKey]; ok {
				for _, w := range c.webConn {
					select {
					case w.send <- m.msg:
					default:
						//delete(h.connections, c)
						close(w.send)
						go w.ws.Close()
					}
				}
			}
		case m := <-b.broadcastToAgent:
			if c, ok := b.connections[m.APIKey]; ok {
				for _, r := range c.agentConn {
					select {
					case r.send <- m.msg:
					default:
						//delete(h.connections, c)
						close(r.send)
						go r.ws.Close()
					}
				}
			}
		}
	}
}
