package states

import (
	"encoding/json"
	"image/color"
	"log"
	"syscall/js"
	"time"

	gs "client/game_state"
	. "client/models"
	. "shared_types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var crossImage *ebiten.Image
var circleImage *ebiten.Image

type Game struct {
	Conn                  js.Value
	incomingMessages      chan []byte
	clientState           ClientState
	GameID                string
	inputBuffer           string
	prevKeyStates         map[ebiten.Key]bool
	prevMouseButtonStates map[ebiten.MouseButton]bool
	PlayerType            PlayerType
	PlayerId              string
	GameData              GameData
	StateMachine          gs.GameStateMachine
	opponentDisconnected  bool
}

func NewGame(conn js.Value) *Game {
	game := &Game{
		Conn:                  conn,
		incomingMessages:      make(chan []byte, 10),
		clientState:           ClientStateMenu,
		prevKeyStates:         make(map[ebiten.Key]bool),
		prevMouseButtonStates: make(map[ebiten.MouseButton]bool),
		StateMachine:          gs.GameStateMachine{},
	}
	game.StateMachine.SetState(&MenuState{Game: game})
	return game
}

func init() {
	crossImage = importImage("assets/cross.png")
	circleImage = importImage("assets/circle.png")
}

func importImage(filePath string) *ebiten.Image {
	var image *ebiten.Image
	var err error

	image, _, err = ebitenutil.NewImageFromFile(filePath)
	if err != nil {
		log.Fatalf("failed to load image from path: %v. %v", filePath, err)
	}

	return image
}

func (g *Game) isKeyJustReleased(key ebiten.Key) bool {
	currentPressed := ebiten.IsKeyPressed(key)
	wasPressed := g.prevKeyStates[key]
	g.prevKeyStates[key] = currentPressed
	return wasPressed && !currentPressed
}

func (g *Game) isMouseButtonJustReleased(mouseButton ebiten.MouseButton) bool {
	currentPressed := ebiten.IsMouseButtonPressed(mouseButton)
	wasPressed := g.prevMouseButtonStates[mouseButton]
	g.prevMouseButtonStates[mouseButton] = currentPressed
	return wasPressed && !currentPressed
}

func (g *Game) Cleanup() {
	if g.Conn.Truthy() {
		g.Conn.Call("close")
	}
}

func (g *Game) ReadMessages() {
	onMessage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		data := event.Get("data").String()
		g.incomingMessages <- []byte(data)
		return nil
	})
	g.Conn.Set("onmessage", onMessage)
}

func (g *Game) Update() error {
	// Handle WebSocket messages (non-blocking)
	select {
	case message := <-g.incomingMessages:
		log.Printf("received: %s", message)
		g.handleServerMessage(message)
	default:
		// No message, continue
	}

	return g.StateMachine.Update()
}

func (g *Game) handleServerMessage(message []byte) {
	var response struct {
		Type            string      `json:"type"`
		GameID          string      `json:"game_id,omitempty"`
		Data            interface{} `json:"data,omitempty"`
		Move            *GameMove   `json:"move,omitempty"`
		GameData        *GameData   `json:"game_data,omitempty"`
		PlayerType      PlayerType  `json:"player_type,omitempty"`
		PlayerId        string      `json:"player_id,omitempty"`
		CurrentPlayerId string      `json:"current_player_id,omitempty"`
	}

	if err := json.Unmarshal(message, &response); err != nil {
		log.Printf("Error unmarshaling server message: %v", err)
		return
	}

	switch response.Type {
	case "game_created":
		g.GameID = response.GameID
		g.PlayerType = response.PlayerType
		g.PlayerId = response.PlayerId
		log.Printf("Game created with ID: %s", g.GameID)
		if g.PlayerType == PlayerTypeCircle {
			log.Printf("Player joined game as Circle")
		} else {
			log.Printf("Player joined game as Cross")
		}
	case "game_joined":
		g.GameID = response.GameID
		g.PlayerType = response.PlayerType
		g.PlayerId = response.PlayerId
		log.Printf("Joined game with ID: %s", g.GameID)
		if g.PlayerType == PlayerTypeCircle {
			log.Printf("Player joined game as Circle")
		} else {
			log.Printf("Player joined game as Cross")
		}
	case "game_start":
		log.Printf("Game is ready to start!")
		g.clientState = ClientStatePlaying
		g.StateMachine.SetState(&PlayState{Game: g})
		g.GameData.Board = response.GameData.Board
		g.GameData.CurrentPlayerId = response.GameData.CurrentPlayerId
		g.opponentDisconnected = false
	case "make_move":
		if response.Move.PlayerType == PlayerTypeCross {
			g.GameData.Board[response.Move.Row][response.Move.Col] = TileStateCross
		} else {
			g.GameData.Board[response.Move.Row][response.Move.Col] = TileStateCircle
		}

		g.GameData.CurrentPlayerId = response.GameData.CurrentPlayerId
	case "game_over":
		g.GameData = *response.GameData
		time.AfterFunc(5*time.Second, func() {
			if g.clientState == ClientStatePlaying {
				g.clientState = ClientStateGameOver
				g.StateMachine.SetState(&GameOverState{Game: g})
			}
		})
	case "game_reset":
		g.clientState = ClientStatePlaying
		g.StateMachine.SetState(&PlayState{Game: g})
		g.GameData.Board = response.GameData.Board
		g.GameData.CurrentPlayerId = response.GameData.CurrentPlayerId
	case "opponent_disconnect":
		g.opponentDisconnected = true
	case "error":
		log.Printf("Server error: %v", response.Data)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(screen)
}

func (g *Game) TryToSelectTileAtPosition(x, y int) bool {
	tileX := x / TileWidth
	tileY := (y - HeaderHeight) / TileHeight

	if tileX < 0 || tileX >= 3 || tileY < 0 || tileY >= 3 {
		return false
	}

	if g.GameData.Board[tileY][tileX] != TileStateEmpty {
		return false
	}

	// Send move to server
	move := GameMove{Row: tileY, Col: tileX}
	g.SendMessage(OutboundMessage{Type: MessageTypeMakeMove, GameID: g.GameID, Move: &move, PlayerId: g.PlayerId})

	return true
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			tileState := g.GameData.Board[row][col]
			if tileState == TileStateEmpty {
				continue
			}

			tileX := col * TileWidth
			tileY := HeaderHeight + row*TileHeight

			var tileImage *ebiten.Image
			switch tileState {
			case TileStateCross:
				tileImage = crossImage
			case TileStateCircle:
				tileImage = circleImage
			}

			// Reset the transformation matrix for reuse
			op.GeoM.Reset()
			op.GeoM.Scale(float64(TileWidth)/float64(tileImage.Bounds().Dx()), float64(TileHeight)/float64(tileImage.Bounds().Dy()))
			op.GeoM.Translate(float64(tileX), float64(tileY))
			screen.DrawImage(tileImage, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) getWinningLine() (startX, startY, endX, endY float32, hasWinner bool) {
	board := g.GameData.Board

	// Check rows
	for row := 0; row < 3; row++ {
		if board[row][0] != TileStateEmpty &&
			board[row][0] == board[row][1] &&
			board[row][1] == board[row][2] {
			// Horizontal line through the middle of the row
			y := float32(HeaderHeight + row*TileHeight + TileHeight/2)
			return 0, y, float32(ScreenWidth), y, true
		}
	}

	// Check columns
	for col := 0; col < 3; col++ {
		if board[0][col] != TileStateEmpty &&
			board[0][col] == board[1][col] &&
			board[1][col] == board[2][col] {
			// Vertical line through the middle of the column
			x := float32(col*TileWidth + TileWidth/2)
			return x, float32(HeaderHeight), x, float32(ScreenHeight), true
		}
	}

	// Check main diagonal (top-left to bottom-right)
	if board[0][0] != TileStateEmpty &&
		board[0][0] == board[1][1] &&
		board[1][1] == board[2][2] {
		return 0, float32(HeaderHeight), float32(ScreenWidth), float32(ScreenHeight), true
	}

	// Check anti-diagonal (top-right to bottom-left)
	if board[0][2] != TileStateEmpty &&
		board[0][2] == board[1][1] &&
		board[1][1] == board[2][0] {
		return float32(ScreenWidth), float32(HeaderHeight), 0, float32(ScreenHeight), true
	}

	return 0, 0, 0, 0, false
}

func (g *Game) drawWinningLine(screen *ebiten.Image) {
	startX, startY, endX, endY, hasWinner := g.getWinningLine()
	if hasWinner {
		var lineColor color.Color
		if g.GameData.Winner == g.PlayerId {
			lineColor = color.RGBA{0, 255, 0, 255}
		} else {
			lineColor = color.RGBA{255, 0, 0, 255}
		}

		vector.StrokeLine(screen, startX, startY, endX, endY, 6, lineColor, false)
	}
}

func (g *Game) SendMessage(msg OutboundMessage) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal: %v", err)
		return
	}
	g.Conn.Call("send", string(msgBytes))
}

func (g *Game) handleTextInput(key rune) {
	if g.clientState != ClientStateEnteringGameID {
		return
	}

	if key == '\b' { // Backspace
		if len(g.inputBuffer) > 0 {
			g.inputBuffer = g.inputBuffer[:len(g.inputBuffer)-1]
		}
	} else if key >= '0' && key <= '9' {
		if len(g.inputBuffer) < 6 { // Limit game ID length to 6 digits
			g.inputBuffer += string(key)
		}
	}
}
