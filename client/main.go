package main

import (
	_ "image/png"
	"log"
	"syscall/js"

	. "client/game_state/states"
	. "client/models"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Tic Tac Toe - Multiplayer")

	// Connect to server
	conn := connectToServer()
	if conn.IsUndefined() || conn.IsNull() {
		log.Fatal("Failed to connect to server.")
	}

	game := NewGame(conn)
	defer game.Cleanup()

	go game.ReadMessages()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func connectToServer() js.Value {
	url := "ws://localhost:8080/ws"
	log.Printf("connecting to %s", url)

	ws := js.Global().Get("WebSocket").New(url)

	log.Println("WebSocket connection initiated")
	return ws
}
