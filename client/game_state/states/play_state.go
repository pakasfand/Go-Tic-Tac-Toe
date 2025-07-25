package states

import (
	"image/color"

	. "shared_types"
	. "client/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var gridLinesImage *ebiten.Image

type PlayState struct {
	Game *Game
}

func (m *PlayState) Draw(screen *ebiten.Image) {
	g := m.Game
	
	screen.Fill(color.RGBA{255, 255, 255, 255})

	g.drawTiles(screen)
	drawGridLines(screen)
}

func init() {
	createGridLines()
}

func createGridLines() {
	gridLinesImage = ebiten.NewImage(ScreenWidth, ScreenHeight)
	
	// Draw horizontal lines
	for i := 1; i < 3; i++ {
		gridLinePosition := float32(i * TileHeight)
		vector.StrokeLine(gridLinesImage, 0, gridLinePosition, float32(ScreenWidth), gridLinePosition, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	}

	// Draw vertical lines
	for i := 1; i < 3; i++ {
		gridLinePosition := float32(i * TileWidth)
		vector.StrokeLine(gridLinesImage, gridLinePosition, 0, gridLinePosition, float32(ScreenHeight), float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	}
}

func drawGridLines(screen *ebiten.Image) {
	if gridLinesImage != nil {
		screen.DrawImage(gridLinesImage, nil)
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