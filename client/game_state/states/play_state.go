package states

import (
	"image/color"

	. "shared_types"
	. "client/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayState struct {
	Game *Game
}

func (m *PlayState) Draw(screen *ebiten.Image) {
	g := m.Game
	
	screen.Fill(color.RGBA{255, 255, 255, 255})

	g.drawTiles(screen)
	drawGridLines(screen)
}

func drawGridLines(screen *ebiten.Image) {
	for i := 0; i < 3; i++ {
		gridLine := ebiten.NewImage(ScreenWidth, ScreenHeight)
		gridLinePosition := float32(i * TileHeight)
		vector.StrokeLine(gridLine, 0, gridLinePosition, float32(ScreenWidth), gridLinePosition, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
		screen.DrawImage(gridLine, nil)
	}

	for i := 0; i < 3; i++ {
		gridLine := ebiten.NewImage(ScreenWidth, ScreenHeight)
		gridLinePosition := float32(i * TileWidth)
		vector.StrokeLine(gridLine, gridLinePosition, 0, gridLinePosition, float32(ScreenHeight), float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
		screen.DrawImage(gridLine, nil)
	}
}

func (m *PlayState) Update() error {
	g := m.Game

	if g.isKeyJustReleased(ebiten.KeyEscape) {
		// Reset to menu
		g.ServerGameData = GameData{
			Board: [3][3]TileState{
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
			},
			Turn: 0,
		}
		g.gameData = ClientStateMenu
		g.StateMachine.SetState(&MenuState{Game: g})
		g.gameID = ""
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		g.TryToSelectTileAtPosition(x, y)
	}

	return nil
}