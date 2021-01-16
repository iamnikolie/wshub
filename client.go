package wshub

import "github.com/gorilla/websocket"

type Client struct {
	hub   *WSHub
	conn  *websocket.Conn
	Scope string
	Send  chan []byte
}
