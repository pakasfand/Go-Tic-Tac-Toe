package models

import (
	. "shared_types"
)

type OutboundMessage struct {
	Type       MessageType `json:"type"`
	GameID     string      `json:"game_id,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Move       *GameMove   `json:"move,omitempty"`
	GameData   *GameData   `json:"game_data,omitempty"`
	PlayerType PlayerType  `json:"player_type,omitempty"`
	PlayerId   string      `json:"player_id,omitempty"`
}