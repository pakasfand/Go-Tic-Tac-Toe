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

func (g *Game) AddPlayer(conn *websocket.Conn, pt PlayerType) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := Player{
		Conn:       conn,
		PlayerType: pt,
		Id:         generateRandomID(),
	}
	g.Players = append(g.Players, player)

	return player.Id
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

func (g *Game) MakeMove(move *GameMove, player *Player) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State.CurrentPlayerId != player.Id {
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
	if player.PlayerType == PlayerTypeCross {
		g.State.Board[move.Row][move.Col] = TileStateCross
	} else {
		g.State.Board[move.Row][move.Col] = TileStateCircle
	}

	// Switch turns
	if player.Id == g.Players[0].Id {
		g.State.CurrentPlayerId = g.Players[1].Id
	} else {
		g.State.CurrentPlayerId = g.Players[0].Id
	}

	move.PlayerType = player.PlayerType

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
	game.State.Board = [3][3]TileState{
		{TileStateEmpty, TileStateEmpty, TileStateEmpty},
		{TileStateEmpty, TileStateEmpty, TileStateEmpty},
		{TileStateEmpty, TileStateEmpty, TileStateEmpty},
	}
	games[gameID] = game
	log.Printf("Number of active games: %d", len(games))
	log.Printf("Created new game with ID: %s", gameID)
	return game
}

func (g *Game) DisconnectGame() bool {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	if _, exists := games[g.ID]; exists {
		delete(games, g.ID)
		log.Printf("Game %s disconnected. Number of active games: %d", g.ID, len(games))
		return true
	}

	return false
}

func FindGame(gameID string) *Game {
	gamesMutex.RLock()
	defer gamesMutex.RUnlock()
	return games[gameID]
}

func (g *Game) FindPlayer(playerId string) *Player {
	gamesMutex.RLock()
	defer gamesMutex.RUnlock()

	for _, player := range g.Players {
		if player.Id == playerId {
			return &player
		}
	}
	return nil
}

func (g *Game) FindPlayerFromType(playerType PlayerType) *Player {
	for _, player := range g.Players {
		if player.PlayerType == playerType {
			return &player
		}
	}
	return nil
}

func FindGameFromConnection(conn *websocket.Conn) *Game {
	for _, game := range games {
		for _, player := range game.Players {
			if player.Conn == conn {
				return game
			}
		}
	}

	return nil
}

func generateGameID() string {
	return strconv.Itoa(rand.Intn(900000) + 100000) // 6-digit number
}

func generateRandomID() string {
	return strconv.Itoa(rand.Intn(900000) + 1000000) // 6-digit number
}
