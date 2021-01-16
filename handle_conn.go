package wshub

import (
	"log"
	"net/http"
)

func (h *WSHub) HandleConnection(w http.ResponseWriter, req *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Client{hub: h, conn: conn, Send: make(chan []byte, 256)}

	if h.ResolveScope != nil {
		c.Scope = h.ResolveScope(req)
	}

	c.hub.Register <- c

	go c.writePump()
	go c.readPump()
}
