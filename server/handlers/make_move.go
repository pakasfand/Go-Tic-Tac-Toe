package handlers

import (
	"encoding/json"
	"log"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func HandleMakeMove(c *websocket.Conn, gameID string, playerId string, move *GameMove) {
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

	player := game.FindPlayer(playerId)
	if !game.MakeMove(move, player) {
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
		Type:     MessageTypeMakeMove,
		Move:     move,
		GameData: &game.State,
	}
	game.Broadcast(moveMsg)

	// Check for game end
	if game.State.IsGameOver() {

		var winnerType, hasWinner = game.State.GetWinnerType()
		if hasWinner {
			game.State.Winner = game.FindPlayerFromType(PlayerType(winnerType)).Id
			log.Printf("Game %s ended. Winner: %s", game.ID, game.State.Winner)
		} else {
			game.State.Winner = ""
		}

		gameOverMsg := OutboundMessage{
			Type:     "game_over",
			GameData: &game.State,
		}
		game.Broadcast(gameOverMsg)
	}
}
