package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const (
	websocketMagicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

// Connection represents a WebSocket connection
type Connection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// Hub manages WebSocket connections
type Hub struct {
	connections map[*Connection]bool
	broadcast   chan []byte
	register    chan *Connection
	unregister  chan *Connection
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		connections: make(map[*Connection]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
			log.Printf("WebSocket connection registered. Total connections: %d", len(h.connections))

		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				conn.Close()
				log.Printf("WebSocket connection unregistered. Total connections: %d", len(h.connections))
			}

		case message := <-h.broadcast:
			for conn := range h.connections {
				if err := conn.WriteMessage(message); err != nil {
					log.Printf("Error writing message: %v", err)
					delete(h.connections, conn)
					conn.Close()
				}
			}
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	select {
	case h.broadcast <- message:
	default:
		log.Println("Broadcast channel is full, dropping message")
	}
}

// Handler creates an HTTP handler for WebSocket connections
func (h *Hub) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is a WebSocket upgrade request
		if !isWebSocketUpgrade(r) {
			http.Error(w, "Expected WebSocket upgrade", http.StatusBadRequest)
			return
		}

		// Hijack the connection
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "WebSocket upgrade not supported", http.StatusInternalServerError)
			return
		}

		conn, bufrw, err := hijacker.Hijack()
		if err != nil {
			log.Printf("Error hijacking connection: %v", err)
			return
		}

		// Perform WebSocket handshake
		if err := performHandshake(bufrw, r); err != nil {
			log.Printf("Error performing handshake: %v", err)
			conn.Close()
			return
		}

		// Create WebSocket connection
		wsConn := &Connection{
			conn:   conn,
			reader: bufrw.Reader,
			writer: bufrw.Writer,
		}

		// Register the connection
		h.register <- wsConn

		// Handle the connection in a goroutine
		go wsConn.handleConnection(h)
	}
}

// isWebSocketUpgrade checks if the request is a WebSocket upgrade request
func isWebSocketUpgrade(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		r.Header.Get("Sec-WebSocket-Key") != ""
}

// performHandshake performs the WebSocket handshake
func performHandshake(bufrw *bufio.ReadWriter, r *http.Request) error {
	key := r.Header.Get("Sec-WebSocket-Key")
	acceptKey := generateAcceptKey(key)

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n"+
			"\r\n",
		acceptKey,
	)

	_, err := bufrw.WriteString(response)
	if err != nil {
		return err
	}

	return bufrw.Flush()
}

// generateAcceptKey generates the WebSocket accept key
func generateAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key + websocketMagicString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// handleConnection handles a WebSocket connection
func (c *Connection) handleConnection(h *Hub) {
	defer func() {
		h.unregister <- c
	}()

	for {
		message, err := c.readMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", string(message))

		// Echo the message back to all clients
		response := fmt.Sprintf("Echo: %s", string(message))
		h.Broadcast([]byte(response))
	}
}

// readMessage reads a WebSocket message (simplified implementation)
func (c *Connection) readMessage() ([]byte, error) {
	// Read the first 2 bytes to get frame info
	header := make([]byte, 2)
	_, err := c.reader.Read(header)
	if err != nil {
		return nil, err
	}

	// Check if this is a close frame
	opcode := header[0] & 0x0F
	if opcode == 0x8 {
		return nil, fmt.Errorf("connection closed by client")
	}

	// Get payload length
	payloadLen := int(header[1] & 0x7F)
	masked := (header[1] & 0x80) != 0

	// Handle extended payload lengths
	if payloadLen == 126 {
		extLen := make([]byte, 2)
		_, err := c.reader.Read(extLen)
		if err != nil {
			return nil, err
		}
		payloadLen = int(binary.BigEndian.Uint16(extLen))
	} else if payloadLen == 127 {
		extLen := make([]byte, 8)
		_, err := c.reader.Read(extLen)
		if err != nil {
			return nil, err
		}
		payloadLen = int(binary.BigEndian.Uint64(extLen))
	}

	// Read masking key if present
	var maskKey []byte
	if masked {
		maskKey = make([]byte, 4)
		_, err := c.reader.Read(maskKey)
		if err != nil {
			return nil, err
		}
	}

	// Read payload
	payload := make([]byte, payloadLen)
	_, err = c.reader.Read(payload)
	if err != nil {
		return nil, err
	}

	// Unmask payload if necessary
	if masked {
		for i := 0; i < len(payload); i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	return payload, nil
}

// WriteMessage writes a WebSocket message
func (c *Connection) WriteMessage(message []byte) error {
	// Create frame header
	frame := make([]byte, 0, len(message)+10)
	
	// First byte: FIN (1) + RSV (000) + Opcode (0001 for text)
	frame = append(frame, 0x81)
	
	// Payload length
	if len(message) < 126 {
		frame = append(frame, byte(len(message)))
	} else if len(message) < 65536 {
		frame = append(frame, 126)
		lengthBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBytes, uint16(len(message)))
		frame = append(frame, lengthBytes...)
	} else {
		frame = append(frame, 127)
		lengthBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(lengthBytes, uint64(len(message)))
		frame = append(frame, lengthBytes...)
	}
	
	// Add payload
	frame = append(frame, message...)
	
	// Write frame
	_, err := c.writer.Write(frame)
	if err != nil {
		return err
	}
	
	return c.writer.Flush()
}

// Close closes the WebSocket connection
func (c *Connection) Close() {
	c.conn.Close()
}
