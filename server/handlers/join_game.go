package handlers

import (
	"encoding/json"
	"log"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func HandleJoinGame(c *websocket.Conn, gameID string) {
	game := FindGame(gameID)
	if game == nil {
		// Game not found
		response := OutboundMessage{
			Type: "error",
			Data: "Game not found",
		}
		responseBytes, _ := json.Marshal(response)
		c.WriteMessage(websocket.TextMessage, responseBytes)
		return
	}

	if game.IsFull() {
		// Game is full
		response := OutboundMessage{
			Type: "error",
			Data: "Game is full",
		}
		responseBytes, _ := json.Marshal(response)
		c.WriteMessage(websocket.TextMessage, responseBytes)
		return
	}

	// Add player
	if game.Players[0].PlayerType == PlayerTypeCircle {
		game.AddPlayer(c, PlayerTypeCross)
	} else {
		game.AddPlayer(c, PlayerTypeCircle)
	}

	// Send success response
	response := OutboundMessage{
		Type:   "game_joined",
		GameID: game.ID,
	}
	responseBytes, _ := json.Marshal(response)
	c.WriteMessage(websocket.TextMessage, responseBytes)

	log.Printf("Player joined game %s as Circle", game.ID)

	// If game is now full, notify both players that game can start
	if game.IsFull() {
		startGameMsg := OutboundMessage{
			Type:     "game_start",
			Data:     "Game is ready to start!",
			GameData: &game.State,
		}
		startGameBytes, _ := json.Marshal(startGameMsg)

		for _, player := range game.Players {
			player.Conn.WriteMessage(websocket.TextMessage, startGameBytes)
		}
	}
}