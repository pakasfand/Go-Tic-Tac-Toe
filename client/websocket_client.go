//go:build !js

package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn *websocket.Conn
}

func NewWebSocketClient() *WebSocketClient {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial: %v", err)
		return nil
	}

	log.Println("Connected to server successfully")
	return &WebSocketClient{Conn: c}
}
