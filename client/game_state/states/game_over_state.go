package states

import (
	"image/color"

	. "shared_types"
	. "client/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type GameOverState struct {
	Game *Game
}

func (m *GameOverState) Draw(screen *ebiten.Image) {
	g := m.Game
	
	var msg string
	switch g.ServerGameData.GetWinner() {
	case "Cross":
		msg = "Cross wins! Press R to restart."
	case "Circle":
		msg = "Circle wins! Press R to restart."
	default:
		msg = "It's a draw! Press R to restart."
	}

	bounds, _ := font.BoundString(basicfont.Face7x13, msg)
	textW := (bounds.Max.X - bounds.Min.X).Ceil()
	textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

	x := (ScreenWidth - textW) / 2
	y := (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil()

	text.Draw(screen, msg, basicfont.Face7x13, x, y, color.White)

	instruction := "Press Esc to navigate back to main menu."
	bounds, _ = font.BoundString(basicfont.Face7x13, instruction)
	textW = (bounds.Max.X - bounds.Min.X).Ceil()
	x = (ScreenWidth - textW) / 2
	y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() + 30

	text.Draw(screen, instruction, basicfont.Face7x13, x, y, color.White)
}

func (m *GameOverState) Update() error {
	g := m.Game

	if g.isKeyJustReleased(ebiten.KeyR) {
		g.sendMessage(OutboundMessage{Type: MessageTypeRequestPlayAgain, GameID: g.gameID})
		// g.gameState = ClientStateMenu
	} else if g.isKeyJustReleased(ebiten.KeyEscape) {
		// Reset state
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
	return nil
}