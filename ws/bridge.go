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
	webConn    []*wsConn
	runnerConn []*wsConn
}

type bridge struct {
	connections       map[string]*connection
	registerWeb       chan *wsConn
	unregisterWeb     chan *wsConn
	registerRunner    chan *wsConn
	unregisterRunner  chan *wsConn
	broadcastToWeb    chan *message
	broadcastToRunner chan *message
}

var Bridge = bridge{
	connections:       make(map[string]*connection),
	registerWeb:       make(chan *wsConn),
	unregisterWeb:     make(chan *wsConn),
	registerRunner:    make(chan *wsConn),
	unregisterRunner:  make(chan *wsConn),
	broadcastToWeb:    make(chan *message),
	broadcastToRunner: make(chan *message),
}

func (b *bridge) Run() {
	for {
		select {
		case conn := <-b.registerWeb:
			if c, ok := b.connections[conn.APIKey]; ok {
				c.webConn = append(c.webConn, conn)
			} else {
				c := &connection{webConn: []*wsConn{conn}, runnerConn: []*wsConn{}}
				b.connections[conn.APIKey] = c
			}
		case conn := <-b.unregisterWeb:
			_ = b.connections[conn.APIKey]
			//delete(c.webConn, conn)
			close(conn.send)
		case conn := <-b.registerRunner:
			if c, ok := b.connections[conn.APIKey]; ok {
				c.runnerConn = append(c.runnerConn, conn)
			} else {
				c := &connection{webConn: []*wsConn{}, runnerConn: []*wsConn{conn}}
				b.connections[conn.APIKey] = c
			}
		case conn := <-b.unregisterRunner:
			_ = b.connections[conn.APIKey]
			//delete(c.runnerConn, conn)
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
		case m := <-b.broadcastToRunner:
			if c, ok := b.connections[m.APIKey]; ok {
				for _, r := range c.runnerConn {
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
