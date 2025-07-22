package types

type GameData struct {
	Board  [3][3]TileState `json:"board"`
	Turn   PlayerType      `json:"turn"`
	Winner string          `json:"winner,omitempty"`
}

func (gd *GameData) IsGameOver() bool {
	// Check for winner first
	if gd.GetWinner() != "" {
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

func (gd *GameData) GetWinner() string {
	// Check rows
	for row := 0; row < 3; row++ {
		if gd.Board[row][0] != TileStateEmpty &&
			gd.Board[row][0] == gd.Board[row][1] &&
			gd.Board[row][1] == gd.Board[row][2] {
			if gd.Board[row][0] == TileStateCross {
				return "Cross"
			} else {
				return "Circle"
			}
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if gd.Board[0][col] != TileStateEmpty &&
			gd.Board[0][col] == gd.Board[1][col] &&
			gd.Board[1][col] == gd.Board[2][col] {
			if gd.Board[0][col] == TileStateCross {
				return "Cross"
			} else {
				return "Circle"
			}
		}
	}

	// Check diagonals
	if gd.Board[0][0] != TileStateEmpty &&
		gd.Board[0][0] == gd.Board[1][1] &&
		gd.Board[1][1] == gd.Board[2][2] {
		if gd.Board[0][0] == TileStateCross {
			return "Cross"
		} else {
			return "Circle"
		}
	}

	if gd.Board[0][2] != TileStateEmpty &&
		gd.Board[0][2] == gd.Board[1][1] &&
		gd.Board[1][1] == gd.Board[2][0] {
		if gd.Board[0][2] == TileStateCross {
			return "Cross"
		} else {
			return "Circle"
		}
	}

	return ""
}
