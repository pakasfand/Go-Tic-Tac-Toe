package types

type ClientState int

const (
	ClientStateMenu ClientState = iota
	ClientStateEnteringGameID
	ClientStatePlaying
	ClientStateGameOver
)