package models

import (
	. "shared_types"
)

type InboundMessage struct {
	Type   MessageType `json:"type"`
	GameID string      `json:"game_id,omitempty"`
	Move   *GameMove   `json:"move,omitempty"`
	PlayerId string      `json:"player_id,omitempty"`
}