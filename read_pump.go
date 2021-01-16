package wshub

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.
		SetReadLimit(maxMessageSize)
	c.conn.
		SetReadDeadline(time.Now().Add(pongWait))
	c.conn.
		SetPongHandler(
			func(string) error {
				c.conn.SetReadDeadline(time.Now().Add(pongWait))
				return nil
			},
		)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(
			bytes.Replace(message, newline, space, -1),
		)
		for _, HFunc := range c.hub.handlers {
			go HFunc(&message)
		}
		c.hub.Broadcast <- message
	}
}
