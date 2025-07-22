package types

type TileState uint8

const (
	TileStateEmpty TileState = iota
	TileStateCross
	TileStateCircle
)