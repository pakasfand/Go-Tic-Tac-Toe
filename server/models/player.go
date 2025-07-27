package models

import (
	. "shared_types"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn                 *websocket.Conn
	PlayerType           PlayerType
	Play_again_requested bool
	Id                   string
}
