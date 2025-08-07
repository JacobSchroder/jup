package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/server/websocket"
)

// WebSocketMessage represents a message structure for WebSocket communication
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NotificationMessage represents a notification to be sent via WebSocket
type NotificationMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Level   string `json:"level"` // "info", "success", "warning", "error"
}

// WebSocketHandler manages WebSocket connections and provides utility functions
type WebSocketHandler struct {
	hub *websocket.Hub
	app *di.App
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(app *di.App) *WebSocketHandler {
	hub := websocket.NewHub()

	// Start the hub in a goroutine
	go hub.Run()

	// Send a welcome message periodically (for demo purposes)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			welcomeMsg := WebSocketMessage{
				Type: "notification",
				Data: NotificationMessage{
					Title:   "System Status",
					Message: "System is running smoothly",
					Level:   "info",
				},
				Timestamp: time.Now(),
			}

			if jsonData, err := json.Marshal(welcomeMsg); err == nil {
				hub.Broadcast(jsonData)
			}
		}
	}()

	return &WebSocketHandler{
		hub: hub,
		app: app,
	}
}

// HandleWebSocket handles WebSocket connection requests
func (wh *WebSocketHandler) HandleWebSocket() http.HandlerFunc {
	return wh.hub.Handler()
}

// BroadcastNotification sends a notification to all connected clients
func (wh *WebSocketHandler) BroadcastNotification(title, message, level string) {
	notification := WebSocketMessage{
		Type: "notification",
		Data: NotificationMessage{
			Title:   title,
			Message: message,
			Level:   level,
		},
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	wh.hub.Broadcast(jsonData)
}

// BroadcastUpdate sends a generic update to all connected clients
func (wh *WebSocketHandler) BroadcastUpdate(updateType string, data any) {
	update := WebSocketMessage{
		Type:      updateType,
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(update)
	if err != nil {
		log.Printf("Error marshaling update: %v", err)
		return
	}

	wh.hub.Broadcast(jsonData)
}

// GetHub returns the WebSocket hub (for advanced usage)
func (wh *WebSocketHandler) GetHub() *websocket.Hub {
	return wh.hub
}
