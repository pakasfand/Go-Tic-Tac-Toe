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

type EnteringGameIdState struct {
	Game *Game
}

func (m *EnteringGameIdState) Draw(screen *ebiten.Image) {
	g := m.Game

	screen.Fill(color.RGBA{0x2c, 0x3e, 0x50, 0xff}) // Dark blue background

	title := "Tic Tac Toe"
	bounds, _ := font.BoundString(basicfont.Face7x13, title)
	textW := (bounds.Max.X - bounds.Min.X).Ceil()
	textH := (bounds.Max.Y - bounds.Min.Y).Ceil()

	x := (ScreenWidth - textW) / 2
	y := (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() - 100

	text.Draw(screen, title, basicfont.Face7x13, x, y, color.White)

	// Show input prompt with 6-digit format
	displayBuffer := g.inputBuffer
	for len(displayBuffer) < 6 {
		displayBuffer += "_"
	}
	prompt := "Enter Game ID: " + displayBuffer
	bounds, _ = font.BoundString(basicfont.Face7x13, prompt)
	textW = (bounds.Max.X - bounds.Min.X).Ceil()
	x = (ScreenWidth - textW) / 2
	y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil()

	text.Draw(screen, prompt, basicfont.Face7x13, x, y, color.White)

	// Show instructions based on current input length
	var instructions string
	if len(g.inputBuffer) == 6 {
		instructions = "Press Enter to join, Escape to cancel"
	} else {
		instructions = "Enter 6-digit Game ID, Escape to cancel"
	}
	bounds, _ = font.BoundString(basicfont.Face7x13, instructions)
	textW = (bounds.Max.X - bounds.Min.X).Ceil()
	x = (ScreenWidth - textW) / 2
	y = (ScreenHeight+textH)/2 - bounds.Max.Y.Ceil() + 30

	text.Draw(screen, instructions, basicfont.Face7x13, x, y, color.Gray{Y: 150})
}

func (m *EnteringGameIdState) Update() error {
	g := m.Game

	// Only allow input if we haven't reached 6 digits yet
	if len(g.inputBuffer) < 6 {
		for key := ebiten.Key0; key <= ebiten.Key9; key++ {
			if g.isKeyJustReleased(key) {
				g.handleTextInput(rune('0' + int(key-ebiten.Key0)))
			}
		}
	}

	if g.isKeyJustReleased(ebiten.KeyBackspace) {
		g.handleTextInput('\b')
	}
	if g.isKeyJustReleased(ebiten.KeyEnter) {
		if len(g.inputBuffer) == 6 {
			g.SendMessage(OutboundMessage{Type: MessageTypeJoinGame, GameID: g.inputBuffer})
			g.inputBuffer = ""
		}
	}
	if g.isKeyJustReleased(ebiten.KeyEscape) {
		g.clientState = ClientStateMenu
		g.StateMachine.SetState(&MenuState{Game: g})
		g.inputBuffer = ""
		g.GameID = ""
	}

	return nil
}
