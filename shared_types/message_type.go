package types

type MessageType string

const (
	MessageTypeInit             MessageType = "init"
	MessageTypeCreateGame       MessageType = "create_game"
	MessageTypeJoinGame         MessageType = "join_game"
	MessageTypeMakeMove         MessageType = "make_move"
	MessageTypeRequestPlayAgain MessageType = "request_play_again"
	MessageLeaveGame            MessageType = "leave_game"
	MessageOpponentDisconnect   MessageType = "opponent_disconnect"
)
