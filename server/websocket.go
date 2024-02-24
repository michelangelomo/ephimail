// server/websocket.go
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHub manages all active WebSocket connections
type WebSocketHub struct {
	// Registered clients
	clients map[*WebSocketClient]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *WebSocketClient

	// Unregister requests from clients
	unregister chan *WebSocketClient

	// Mapping of email addresses to clients
	subscriptions map[string][]*WebSocketClient
	subLock       sync.RWMutex
}

// WebSocketClient represents a connected client
type WebSocketClient struct {
	hub *WebSocketHub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Email address this client is subscribed to
	email string
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast:     make(chan []byte, 256),
		register:      make(chan *WebSocketClient),
		unregister:    make(chan *WebSocketClient),
		clients:       make(map[*WebSocketClient]bool),
		subscriptions: make(map[string][]*WebSocketClient),
	}
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from subscriptions
				if client.email != "" {
					h.subLock.Lock()
					clients := h.subscriptions[client.email]
					for i, c := range clients {
						if c == client {
							h.subscriptions[client.email] = append(clients[:i], clients[i+1:]...)
							break
						}
					}
					if len(h.subscriptions[client.email]) == 0 {
						delete(h.subscriptions, client.email)
					}
					h.subLock.Unlock()
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Subscribe subscribes a client to an email address
func (h *WebSocketHub) Subscribe(client *WebSocketClient, email string) {
	h.subLock.Lock()
	defer h.subLock.Unlock()

	// Unsubscribe from previous email if any
	if client.email != "" && client.email != email {
		clients := h.subscriptions[client.email]
		for i, c := range clients {
			if c == client {
				h.subscriptions[client.email] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(h.subscriptions[client.email]) == 0 {
			delete(h.subscriptions, client.email)
		}
	}

	// Subscribe to new email
	client.email = email
	if _, ok := h.subscriptions[email]; !ok {
		h.subscriptions[email] = make([]*WebSocketClient, 0)
	}
	h.subscriptions[email] = append(h.subscriptions[email], client)
}

// NotifyNewEmail notifies all clients subscribed to an email address about a new email
func (h *WebSocketHub) NotifyNewEmail(email, messageID string) {
	h.subLock.RLock()
	defer h.subLock.RUnlock()

	clients, ok := h.subscriptions[email]
	if !ok {
		return
	}

	// Create JSON message
	payload := map[string]string{
		"email": email,
	}

	if messageID != "" {
		payload["message_id"] = messageID
	}

	msgData, err := json.Marshal(WebSocketMessage{
		Type:    "new_email",
		Payload: payload,
	})

	if err != nil {
		log.Printf("Error creating WebSocket message: %v", err)
		return
	}

	// Send message to all subscribed clients
	for _, client := range clients {
		select {
		case client.send <- msgData:
		default:
			// If channel is full, just continue
		}
	}
}

// Upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// ServeWs handles WebSocket requests from clients
func ServeWs(hub *WebSocketHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &WebSocketClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}

// readPump reads messages from the WebSocket connection
func (c *WebSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse the message
		var msg struct {
			Type    string `json:"type"`
			Payload struct {
				Email string `json:"email"`
			} `json:"payload"`
		}

		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing WebSocket message: %v", err)
			continue
		}

		// Handle message based on type
		switch msg.Type {
		case "subscribe":
			if msg.Payload.Email != "" {
				c.hub.Subscribe(c, msg.Payload.Email)
				log.Printf("Client subscribed to %s", msg.Payload.Email)
			}
		case "unsubscribe":
			if msg.Payload.Email != "" && msg.Payload.Email == c.email {
				c.email = ""
				// The client will be removed from subscriptions in the unregister handler
			}
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
