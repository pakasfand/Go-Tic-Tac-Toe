package models

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"sync"
	
	. "shared_types"

	"github.com/gorilla/websocket"
)

var games = make(map[string]*Game)
var gamesMutex sync.RWMutex

type Game struct {
	ID      string
	Players []Player
	State   GameData
	mu      sync.Mutex
}

func (g *Game) AddPlayer(conn *websocket.Conn, pt PlayerType) {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := Player{
		Conn:       conn,
		PlayerType: pt,
	}
	g.Players = append(g.Players, player)
}

func (g *Game) IsFull() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.Players) >= 2
}

func (g *Game) GetPlayerCount() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.Players)
}

func (g *Game) Broadcast(message OutboundMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for _, player := range g.Players {
		err := player.Conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Printf("Error broadcasting to player: %v", err)
		}
	}
}

func (g *Game) MakeMove(move GameMove, playerType PlayerType) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State.Turn != playerType {
		return false // Not this player's turn
	}

	// Check if the move is valid
	if move.Row < 0 || move.Row >= 3 || move.Col < 0 || move.Col >= 3 {
		return false // Invalid position
	}

	if g.State.Board[move.Row][move.Col] != TileStateEmpty {
		return false // Position already occupied
	}

	// Make the move
	if playerType == PlayerTypeCross {
		g.State.Board[move.Row][move.Col] = TileStateCross
	} else {
		g.State.Board[move.Row][move.Col] = TileStateCircle
	}

	// Switch turns
	if g.State.Turn == PlayerTypeCircle {
		g.State.Turn = PlayerTypeCross
	} else {
		g.State.Turn = PlayerTypeCircle
	}

	return true
}

func CreateNewGame() *Game {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	var gameID string
	for {
		gameID = generateGameID()
		if _, exists := games[gameID]; !exists {
			break
		}
	}

	game := &Game{ID: gameID}
	games[gameID] = game
	log.Printf("Created new game with ID: %s", gameID)
	return game
}

func FindGame(gameID string) *Game {
	gamesMutex.RLock()
	defer gamesMutex.RUnlock()
	return games[gameID]
}

func generateGameID() string {
	return strconv.Itoa(rand.Intn(900000) + 100000) // 6-digit number
}