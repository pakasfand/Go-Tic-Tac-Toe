package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func HandleCreateGame(c *websocket.Conn) {
	game := CreateNewGame()
	playerType := PlayerType(rand.Intn(2))
	game.AddPlayer(c, playerType)

	// Send game ID back to client
	response := OutboundMessage{
		Type:       "game_created",
		GameID:     game.ID,
		PlayerType: playerType,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, responseBytes)
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}

	log.Printf("Player joined game %s as Cross", game.ID)
}