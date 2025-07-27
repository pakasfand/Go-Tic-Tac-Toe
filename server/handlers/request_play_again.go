package handlers

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func HandleRequestPlayAgain(c *websocket.Conn, gameID string) {
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

	for i, player := range game.Players {
		if player.Conn == c {
			game.Players[i].Play_again_requested = true
			break
		}
	}

	// Check if any player hasn't requested to play again yet
	for _, player := range game.Players {
		if !player.Play_again_requested {
			return
		}
	}

	// Reset the game state
	game.State = GameData{
		Board: [3][3]TileState{
			{TileStateEmpty, TileStateEmpty, TileStateEmpty},
			{TileStateEmpty, TileStateEmpty, TileStateEmpty},
			{TileStateEmpty, TileStateEmpty, TileStateEmpty},
		},
		CurrentPlayerId: game.Players[rand.IntN(len(game.Players))].Id,
	}

	resetMsg := OutboundMessage{
		Type:     "game_reset",
		GameData: &game.State,
	}
	game.Broadcast(resetMsg)

	log.Printf("Game %s has been reset for a new round", game.ID)

	// Reset play again flag
	for i := range game.Players {
		game.Players[i].Play_again_requested = false
	}
}