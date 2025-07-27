package types

type GameData struct {
	Board           [3][3]TileState `json:"board"`
	CurrentPlayerId string          `json:"current_player_id"`
	Winner          string          `json:"winner,omitempty"`
}

func (gd *GameData) IsGameOver() bool {
	// Check for winner first
	_, hasWinner := gd.GetWinnerType()

	if hasWinner {
		return true
	}

	// Check if board is full (draw)
	for _, row := range gd.Board {
		for _, tile := range row {
			if tile == TileStateEmpty {
				return false
			}
		}
	}
	return true
}

func (gd *GameData) GetWinnerType() (TileState, bool) {
	// Check rows
	for row := 0; row < 3; row++ {
		tileType := gd.Board[row][0]
		if tileType != TileStateEmpty &&
			tileType == gd.Board[row][1] &&
			gd.Board[row][1] == gd.Board[row][2] {
			return tileType, true
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		tileType := gd.Board[0][col]
		if tileType != TileStateEmpty &&
			tileType == gd.Board[1][col] &&
			gd.Board[1][col] == gd.Board[2][col] {
			return tileType, true
		}
	}

	// Check diagonals
	tileType := gd.Board[0][0]
	if tileType != TileStateEmpty &&
		tileType == gd.Board[1][1] &&
		gd.Board[1][1] == gd.Board[2][2] {
		return tileType, true
	}

	tileType = gd.Board[0][2]
	if tileType != TileStateEmpty &&
		tileType == gd.Board[1][1] &&
		gd.Board[1][1] == gd.Board[2][0] {
		return tileType, true
	}

	return TileStateEmpty, false
}
