package game_state

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState interface {
	Draw(screen *ebiten.Image)
	Update() error
}
