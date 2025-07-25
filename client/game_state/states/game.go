package states

import (
	"encoding/json"
	"log"

	. "shared_types"
	. "client/game_state"
	. "client/models"
	// . "client/game_state/states"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

var crossImage *ebiten.Image
var circleImage *ebiten.Image

type Game struct {
	Conn             *websocket.Conn
	incomingMessages chan []byte
	gameData         ClientState
	gameID           string
	inputBuffer      string
	prevKeyStates    map[ebiten.Key]bool
	PlayerType       PlayerType
	ServerGameData   GameData
	StateMachine     GameStateMachine
}

func NewGame(conn *websocket.Conn) *Game {
	game := &Game{
		Conn:             conn,
		incomingMessages: make(chan []byte, 10),
		gameData:         ClientStateMenu,
		prevKeyStates:    make(map[ebiten.Key]bool),
		StateMachine: GameStateMachine{},
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

func (g *Game) Cleanup() {
	if g.Conn != nil {
		g.Conn.Close()
	}
}

func (g *Game) ReadMessages() {
	for {
		_, message, err := g.Conn.ReadMessage()
		if err != nil {
			log.Printf("read: %v", err)
			return // Optionally handle reconnection or cleanup
		}
		g.incomingMessages <- message
	}
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
		Type       string      `json:"type"`
		GameID     string      `json:"game_id,omitempty"`
		Data       interface{} `json:"data,omitempty"`
		Move       *GameMove   `json:"move,omitempty"`
		GameData   *GameData   `json:"game_data,omitempty"`
		PlayerType PlayerType  `json:"player_type,omitempty"`
	}

	if err := json.Unmarshal(message, &response); err != nil {
		log.Printf("Error unmarshaling server message: %v", err)
		return
	}

	switch response.Type {
	case "game_created":
		g.gameID = response.GameID
		g.PlayerType = response.PlayerType
		log.Printf("Game created with ID: %s", g.gameID)
	case "game_joined":
		g.gameID = response.GameID
		g.PlayerType = response.PlayerType
		log.Printf("Joined game with ID: %s", g.gameID)
	case "game_start":
		log.Printf("Game is ready to start!")
		g.gameData = ClientStatePlaying
		g.StateMachine.SetState(&PlayState{Game: g})
		g.ServerGameData = *response.GameData
	case "make_move":
		if response.Move != nil {
			if g.ServerGameData.Turn == PlayerTypeCross {
				g.ServerGameData.Board[response.Move.Row][response.Move.Col] = TileStateCross
			} else {
				g.ServerGameData.Board[response.Move.Row][response.Move.Col] = TileStateCircle
			}

			if g.ServerGameData.Turn == PlayerTypeCross {
				g.ServerGameData.Turn = PlayerTypeCircle
			} else {
				g.ServerGameData.Turn = PlayerTypeCross
			}
		}
	case "game_over":
		g.gameData = ClientStateGameOver
		g.StateMachine.SetState(&GameOverState{Game: g})
		g.ServerGameData = *response.GameData
	case "game_reset":
		g.gameData = ClientStatePlaying
		g.StateMachine.SetState(&PlayState{Game: g})
		g.ServerGameData.Board = response.GameData.Board
		g.ServerGameData.Turn = response.GameData.Turn
	case "error":
		log.Printf("Server error: %v", response.Data)
		// g.gameState = ClientStateMenu
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(screen)
}

func (g *Game) TryToSelectTileAtPosition(x, y int) bool {
	tileX := x / TileWidth
	tileY := y / TileHeight

	if tileX < 0 || tileX >= 3 || tileY < 0 || tileY >= 3 {
		return false
	}

	// Only allow move if it's this player's turn
	if g.ServerGameData.Turn != g.PlayerType {
		return false
	}

	if g.ServerGameData.Board[tileY][tileX] != TileStateEmpty {
		return false
	}

	// Send move to server
	move := GameMove{Row: tileY, Col: tileX}
	g.sendMessage(OutboundMessage{Type: MessageTypeMakeMove, GameID: g.gameID, Move: &move})

	return true
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			tileState := g.ServerGameData.Board[row][col]
			if tileState == TileStateEmpty {
				continue
			}

			tileX := col * (ScreenWidth / 3)
			tileY := row * (ScreenHeight / 3)

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

func (g *Game) sendMessage(msg OutboundMessage) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal: %v", err)
		return
	}

	err = g.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Printf("write: %v", err)
	}
}

func (g *Game) handleTextInput(key rune) {
	if g.gameData != ClientStateEnteringGameID {
		return
	}

	if key == '\b' { // Backspace
		if len(g.inputBuffer) > 0 {
			g.inputBuffer = g.inputBuffer[:len(g.inputBuffer)-1]
		}
	} else if key >= '0' && key <= '9' {
		if len(g.inputBuffer) < 6 { // Limit game ID length
			g.inputBuffer += string(key)
		}
	}
}
