package models

import (
	. "shared_types"
)

type OutboundMessage struct {
	Type   MessageType `json:"type"`
	GameID string      `json:"game_id,omitempty"`
	Move   *GameMove   `json:"move,omitempty"`
}