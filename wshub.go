package wshub

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// WSHub ...
type WSHub struct {
	// Status of WShub
	run bool
	// Map of clients
	clients map[*Client]bool
	// List of handlers
	handlers []func(msg *[]byte)
	// Broadcast chanell
	Broadcast chan []byte
	// ??
	ScopeBroadcast chan ScopedMsg
	// Register client chanel
	Register chan *Client
	// Unregister client chanel
	Unregister chan *Client
	// Scope resolver
	ResolveScope func(req *http.Request) string
}

// NewWSHub ...
// create new websocket pool instance
func NewWSHub() *WSHub {
	return &WSHub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		run:        false,
	}
}

// Run ...
func (h *WSHub) Run() {
	if h.run {
		return
	}
	h.run = true
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case scopedmsg := <-h.ScopeBroadcast:
			for client := range h.clients {
				if client.Scope == scopedmsg.Scope {
					select {
					case client.Send <- scopedmsg.Msg:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
