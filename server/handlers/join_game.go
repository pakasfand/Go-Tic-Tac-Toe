package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
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
	var playerId string
	if game.Players[0].PlayerType == PlayerTypeCircle {
		playerId = game.AddPlayer(c, PlayerTypeCross)
	} else {
		playerId = game.AddPlayer(c, PlayerTypeCircle)
	}

	// Send success response
	response := OutboundMessage{
		Type:       "game_joined",
		GameID:     game.ID,
		PlayerType: game.Players[1].PlayerType,
		GameData:   &game.State,
		PlayerId:  playerId,
	}
	responseBytes, _ := json.Marshal(response)
	c.WriteMessage(websocket.TextMessage, responseBytes)

	player := game.FindPlayer(playerId)
	if player.PlayerType == PlayerTypeCircle {
		log.Printf("Player joined game %s as Circle", game.ID)
	} else {
		log.Printf("Player joined game %s as Cross", game.ID)
	}

	// If game is now full, notify both players that game can start
	if game.IsFull() {
		// Randomly select which player starts
		game.State.CurrentPlayerId = game.Players[rand.Intn(len(game.Players))].Id
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
