package handlers

import (
	"encoding/json"
	"log"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func HandleMakeMove(c *websocket.Conn, gameID string, move GameMove) {
	game := FindGame(gameID)
	if game == nil {
		response := OutboundMessage{
			Type: "error",
			Data: "Game not found",
		}
		responseBytes, _ := json.Marshal(response)
		c.WriteMessage(websocket.TextMessage, responseBytes)
		return
	}

	var playerType PlayerType
	for _, player := range game.Players {
		if player.Conn == c {
			playerType = player.PlayerType
			break
		}
	}

	if !game.MakeMove(move, playerType) {
		response := OutboundMessage{
			Type: "error",
			Data: "Invalid move or not your turn",
		}
		responseBytes, _ := json.Marshal(response)
		c.WriteMessage(websocket.TextMessage, responseBytes)
		return
	}

	// Broadcast the move to all players
	moveMsg := OutboundMessage{
		Type: MessageTypeMakeMove,
		Move: &move,
	}
	game.Broadcast(moveMsg)

	// Check for game end
	if game.State.IsGameOver() {
		game.State.Winner = game.State.GetWinner()
		gameOverMsg := OutboundMessage{
			Type:     "game_over",
			GameData: &game.State,
		}
		game.Broadcast(gameOverMsg)
		log.Printf("Game %s ended. Winner: %s", game.ID, game.State.Winner)
	}
}