//go:build !js
// +build !js

package main

import (
	_ "image/png"
	"log"

	. "client/game_state/states"
	. "client/models"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Tic Tac Toe - Multiplayer")

	conn, err := ConnectToServer("ws://127.0.0.1:8080/ws")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	game := NewGame(conn)
	defer game.Cleanup()

	go game.ReadMessages()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
