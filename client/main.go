package main

import (
	_ "image/png"
	"log"
	"net/url"

	. "client/models"
	. "client/game_state/states"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Tic Tac Toe - Multiplayer")

	// Connect to server
	conn := connectToServer()
	if conn == nil {
		log.Fatal("Failed to connect to server.")
	}

	game := NewGame(conn)
	defer game.Cleanup()

	go game.ReadMessages()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func connectToServer() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial: %v", err)
		return nil
	}

	log.Println("Connected to server successfully")
	return c
}