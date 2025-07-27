package states

import (
	"image/color"

	. "client/models"
	. "shared_types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type GameOverState struct {
	Game             *Game
	rematchRequested bool
}

func (m *GameOverState) Draw(screen *ebiten.Image) {
	var msg string
	if m.rematchRequested {
		msg = "Rematch requested. Waiting for opponent to accept..."
	} else {
		msg = "Press R to request a rematch."
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
		m.rematchRequested = true
		g.SendMessage(OutboundMessage{Type: MessageTypeRequestPlayAgain, GameID: g.GameID})
	} else if g.isKeyJustReleased(ebiten.KeyEscape) {
		// Reset state
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
	return nil
}
