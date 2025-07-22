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

type MenuState struct {
	Game *Game
}

func (m *MenuState) Draw(screen *ebiten.Image) {
	g := m.Game
	screen.Fill(color.RGBA{0x2c, 0x3e, 0x50, 0xff}) // Dark blue background

	title := "Tic Tac Toe"
	bounds, _ := font.BoundString(basicfont.Face7x13, title)
	textW := (bounds.Max.X - bounds.Min.X).Ceil()
	textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

	x := (ScreenWidth - textW) / 2
	y := (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() - 100

	text.Draw(screen, title, basicfont.Face7x13, x, y, color.White)

	if g.gameID != "" && g.gameData == ClientStateMenu {
		// Show game ID for sharing
		gameIDText := "Game ID: " + g.gameID
		bounds, _ = font.BoundString(basicfont.Face7x13, gameIDText)
		textW = (bounds.Max.X - bounds.Min.X).Ceil()
		x = (ScreenWidth - textW) / 2
		y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() - 50

		text.Draw(screen, gameIDText, basicfont.Face7x13, x, y, color.RGBA{0x00, 0xff, 0x00, 0xff})

		waitingText := "Waiting for opponent to join..."
		bounds, _ = font.BoundString(basicfont.Face7x13, waitingText)
		textW = (bounds.Max.X - bounds.Min.X).Ceil()
		x = (ScreenWidth - textW) / 2
		y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() - 20

		text.Draw(screen, waitingText, basicfont.Face7x13, x, y, color.Gray{Y: 200})
	} else {
		// Show menu options
		option1 := "Press 1 to Create New Game"
		bounds, _ = font.BoundString(basicfont.Face7x13, option1)
		textW = (bounds.Max.X - bounds.Min.X).Ceil()
		x = (ScreenWidth - textW) / 2
		y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil()

		text.Draw(screen, option1, basicfont.Face7x13, x, y, color.White)

		option2 := "Press 2 to Join Existing Game"
		bounds, _ = font.BoundString(basicfont.Face7x13, option2)
		textW = (bounds.Max.X - bounds.Min.X).Ceil()
		x = (ScreenWidth - textW) / 2
		y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() + 30

		text.Draw(screen, option2, basicfont.Face7x13, x, y, color.White)
	}
}

func (m *MenuState) Update() error {
	g := m.Game

	// Handle menu navigation
	if g.isKeyJustReleased(ebiten.Key1) {
		// Create new game
		g.sendMessage(OutboundMessage{Type: MessageTypeCreateGame})
	} else if g.isKeyJustReleased(ebiten.Key2) {
		// Join existing game
		m.Game.gameData = ClientStateEnteringGameID
		m.Game.StateMachine.SetState(&EnteringGameIdState{Game: g})
	}
	return nil
}
