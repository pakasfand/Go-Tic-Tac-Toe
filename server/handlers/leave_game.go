package handlers

import (
	. "server/models"
	types "shared_types"

	"github.com/gorilla/websocket"
)

func HandleLeaveGame(c *websocket.Conn, gameID string) {
	game := FindGame(gameID)
	if game != nil && game.DisconnectGame() {
		game.Broadcast(OutboundMessage{Type: types.MessageOpponentDisconnect})
	}
}
