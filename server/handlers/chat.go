package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for demo
	},
}

type ChatMessage struct {
	Username  string `json:"username"`
	Message   string `json:"chat_message"`
	Timestamp string `json:"timestamp"`
}

type ChatHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

var hub = &ChatHub{
	clients:    make(map[*websocket.Conn]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

func init() {
	go hub.run()
}

func (h *ChatHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			slog.Info("Client connected", "total_clients", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			slog.Info("Client disconnected", "total_clients", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					slog.Error("Error writing message to client", "error", err)
					delete(h.clients, client)
					client.Close()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func HandleChatWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	hub.register <- conn

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", "error", err)
			hub.unregister <- conn
			break
		}

		var chatMsg ChatMessage
		if err := json.Unmarshal(msgBytes, &chatMsg); err != nil {
			slog.Error("Error unmarshaling message", "error", err)
			continue
		}

		slog.Info("Received message", "message", chatMsg)

		// Set timestamp
		chatMsg.Timestamp = time.Now().Format("15:04:05")

		// Default username if not provided
		if chatMsg.Username == "" {
			chatMsg.Username = "Anonymous"
		}

		// Skip empty messages
		if chatMsg.Message == "" {
			continue
		}

		// Escape HTML to prevent XSS
		safeUsername := html.EscapeString(chatMsg.Username)
		safeMessage := html.EscapeString(chatMsg.Message)

		// Create HTML response for HTMX
		htmlMessage := fmt.Sprintf(`
			<div id="messages" hx-swap-oob="beforeend">
				<div class="mb-2 p-3 bg-slate-700 rounded-lg">
					<div class="flex items-center gap-2 mb-1">
						<span class="text-blue-400 font-semibold">%s</span>
						<span class="text-gray-400 text-xs">%s</span>
					</div>
					<div class="text-white">%s</div>
				</div>
			</div>
		`, safeUsername, chatMsg.Timestamp, safeMessage)

		hub.broadcast <- []byte(htmlMessage)
	}
}
