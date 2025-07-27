package types

type GameMove struct {
	Row int `json:"row"`
	Col int `json:"col"`
	PlayerType PlayerType `json:"player_type"`
}
