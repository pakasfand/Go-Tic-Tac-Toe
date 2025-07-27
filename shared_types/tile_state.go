package types

type TileState uint8

const (
	TileStateCross TileState = iota
	TileStateCircle
	TileStateEmpty
)