package game_state

import "github.com/hajimehoshi/ebiten/v2"

type GameStateMachine struct {
	State GameState
}

func (stateMachine *GameStateMachine) SetState(newState GameState) {
	stateMachine.State = newState
}

func (stateMachine *GameStateMachine) Update() error {
	return stateMachine.State.Update()
}

func (stateMachine *GameStateMachine) Draw(screen *ebiten.Image) {
	stateMachine.State.Draw(screen)
}
