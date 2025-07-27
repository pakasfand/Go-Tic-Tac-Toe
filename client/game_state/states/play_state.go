package states

import (
	"image/color"

	. "client/models"
	. "shared_types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var gridLinesImage *ebiten.Image
var headerOutlineImage *ebiten.Image

type PlayState struct {
	Game *Game
}

func (m *PlayState) Draw(screen *ebiten.Image) {
	g := m.Game

	screen.Fill(color.RGBA{255, 255, 255, 255})

	g.drawTiles(screen)
	drawGridLines(screen)
	g.drawWinningLine(screen)
	drawHeaderOutline(screen)
	drawHeader(g, screen)
}

func drawHeader(g *Game, screen *ebiten.Image) {
	var msg string
	var textColor color.Color
	if g.GameData.IsGameOver() {
		if g.GameData.Winner == "" {
			msg = "Draw!"
			textColor = color.RGBA{255, 255, 255, 255}
		} else if g.GameData.Winner == g.PlayerId {
			msg = "You Win!"
			textColor = color.RGBA{0, 128, 0, 255}
		} else {
			msg = "You Lose!"
			textColor = color.RGBA{128, 0, 0, 255}
		}
	} else {
		if g.GameData.CurrentPlayerId == g.PlayerId {
			msg = "Your Turn"
			textColor = color.RGBA{0, 128, 0, 255}
		} else {
			msg = "Opponent's Turn"
			textColor = color.RGBA{128, 0, 0, 255}
		}
	}
	if g.opponentDisconnected {
		msg = "Opponent Disconnected!"
		textColor = color.RGBA{255, 255, 255, 255}

		// Draw instructions
		instruction := "Press Esc to navigate back the main menu."

		bounds, _ := font.BoundString(basicfont.Face7x13, instruction)
		textW := (bounds.Max.X - bounds.Min.X).Ceil()
		textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

		x := (HeaderWidth - textW) / 2
		y := (HeaderHeight+textH+50)/2 - bounds.Max.Y.Ceil()

		text.Draw(screen, instruction, basicfont.Face7x13, x, y, textColor)
	}
	bounds, _ := font.BoundString(basicfont.Face7x13, msg)
	textW := (bounds.Max.X - bounds.Min.X).Ceil()
	textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

	x := (HeaderWidth - textW) / 2
	y := (HeaderHeight+textH)/2 - bounds.Max.Y.Ceil()

	text.Draw(screen, msg, basicfont.Face7x13, x, y, textColor)
}

func init() {
	createGridLines()
	createHeaderOutline()
}

func createGridLines() {
	gridLinesImage = ebiten.NewImage(ScreenWidth, ScreenHeight)

	// Draw horizontal lines
	for i := 1; i < 3; i++ {
		gridLinePosition := HeaderHeight + float32(i*TileHeight)
		vector.StrokeLine(gridLinesImage, 0, gridLinePosition, float32(ScreenWidth), gridLinePosition, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	}

	// Draw vertical lines
	for i := 1; i < 3; i++ {
		gridLinePosition := float32(i * TileWidth)
		vector.StrokeLine(gridLinesImage, gridLinePosition, HeaderHeight, gridLinePosition, float32(ScreenHeight), float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	}

}

func createHeaderOutline() {
	// Draw Header outline
	vector.StrokeLine(gridLinesImage, TileWidth, HeaderHeight-50, TileWidth*2, HeaderHeight-50, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	vector.StrokeLine(gridLinesImage, TileWidth, HeaderHeight-75, TileWidth*2, HeaderHeight-75, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)

	vector.StrokeLine(gridLinesImage, TileWidth, HeaderHeight-50, TileWidth, HeaderHeight-75, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
	vector.StrokeLine(gridLinesImage, TileWidth*2, HeaderHeight-50, TileWidth*2, HeaderHeight-75, float32(GridLineWidth), color.RGBA{0x00, 0x00, 0x00, 0xff}, false)
}

func drawGridLines(screen *ebiten.Image) {
	if gridLinesImage != nil {
		screen.DrawImage(gridLinesImage, nil)
	}
}

func drawHeaderOutline(screen *ebiten.Image) {
	if headerOutlineImage != nil {
		screen.DrawImage(headerOutlineImage, nil)
	}
}

func (m *PlayState) Update() error {
	g := m.Game

	if g.isKeyJustReleased(ebiten.KeyEscape) {
		// Reset to menu
		g.GameData = GameData{
			Board: [3][3]TileState{
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
				{TileStateEmpty, TileStateEmpty, TileStateEmpty},
			},
		}
		g.SendMessage(OutboundMessage{Type: MessageLeaveGame, GameID: g.GameID})
		g.clientState = ClientStateMenu
		g.StateMachine.SetState(&MenuState{Game: g})
		g.GameID = ""
	}

	if g.isMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		g.TryToSelectTileAtPosition(x, y)
	}

	return nil
}
