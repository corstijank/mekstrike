package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/corstijank/mekstrike/gamemaster/internal/types"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// In production, this should be more restrictive
		return true
	},
}

type WebSocketService struct {
	// Map of game ID to list of connections
	connections map[string][]*websocket.Conn
	mutex       sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[string][]*websocket.Conn),
	}
}

// HandleWebSocketConnection handles new WebSocket connections for a specific game
func (ws *WebSocketService) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameId := vars["id"]

	if gameId == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add connection to game's connection list
	ws.addConnection(gameId, conn)
	defer ws.removeConnection(gameId, conn)

	log.Printf("New WebSocket connection established for game %s", gameId)

	// Keep connection alive and handle disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for game %s: %v", gameId, err)
			}
			break
		}
	}

	log.Printf("WebSocket connection closed for game %s", gameId)
}

// addConnection adds a WebSocket connection to a game's connection list
func (ws *WebSocketService) addConnection(gameId string, conn *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.connections[gameId] == nil {
		ws.connections[gameId] = make([]*websocket.Conn, 0)
	}
	ws.connections[gameId] = append(ws.connections[gameId], conn)
}

// removeConnection removes a WebSocket connection from a game's connection list
func (ws *WebSocketService) removeConnection(gameId string, conn *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	connections := ws.connections[gameId]
	for i, c := range connections {
		if c == conn {
			// Remove connection from slice
			ws.connections[gameId] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// Clean up empty game entries
	if len(ws.connections[gameId]) == 0 {
		delete(ws.connections, gameId)
	}
}

// BroadcastToGame sends a CloudEvent to all WebSocket connections for a specific game
func (ws *WebSocketService) BroadcastToGame(gameId string, cloudEvent types.CloudEvent) {
	ws.mutex.RLock()
	connections := ws.connections[gameId]
	if len(connections) == 0 {
		ws.mutex.RUnlock()
		return // No connections for this game
	}

	// Make a copy of connections to avoid holding lock during broadcast
	connsCopy := make([]*websocket.Conn, len(connections))
	copy(connsCopy, connections)
	ws.mutex.RUnlock()

	// Prepare message
	message, err := json.Marshal(cloudEvent)
	if err != nil {
		log.Printf("Failed to marshal CloudEvent for WebSocket broadcast: %v", err)
		return
	}

	// Broadcast to all connections for this game
	for _, conn := range connsCopy {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Failed to write WebSocket message for game %s: %v", gameId, err)
			// Remove failed connection
			ws.removeConnection(gameId, conn)
		}
	}

	log.Printf("Broadcasted CloudEvent to %d connections for game %s", len(connsCopy), gameId)
}

// GetConnectionCount returns the number of active connections for a game
func (ws *WebSocketService) GetConnectionCount(gameId string) int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	
	return len(ws.connections[gameId])
}

// GetTotalConnections returns the total number of active connections across all games
func (ws *WebSocketService) GetTotalConnections() int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	
	total := 0
	for _, connections := range ws.connections {
		total += len(connections)
	}
	return total
}