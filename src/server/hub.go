// Package server provides the new TCP-simple WebSocket hub
// Replaces the complex system with sequence number reliability
package server

import (
	"encoding/json"
	stdSync "sync"
	"time"

	"holodeck1/logging"
	"holodeck1/sync"
)

// Hub represents the TCP-simple WebSocket coordination hub
type Hub struct {
	// Core sync system
	sync *sync.ReliableSync
	
	// Client management
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	mutex      stdSync.RWMutex
	
	// Avatar management
	avatarRegistry *AvatarRegistry
	
	// Message routing
	broadcast chan []byte
}

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	ClientID  string                 `json:"client_id,omitempty"`
	Operation *sync.Operation        `json:"operation,omitempty"`
	From      uint64                 `json:"from,omitempty"`
	To        uint64                 `json:"to,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NewHub creates a new TCP-simple WebSocket hub
func NewHub() *Hub {
	hub := &Hub{
		sync:       sync.NewReliableSync(),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 1000),
	}
	
	// Initialize avatar registry
	hub.avatarRegistry = NewAvatarRegistry(hub)
	
	return hub
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	go h.handleOperations()
	
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
			
		case client := <-h.unregister:
			h.unregisterClient(client)
			
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// handleOperations processes sync operations
func (h *Hub) handleOperations() {
	for {
		time.Sleep(10 * time.Millisecond) // Small delay to prevent busy waiting
		// Process any pending operations from sync system
		operations := h.sync.GetPendingOperations()
		for _, op := range operations {
			h.broadcastOperation(op)
		}
	}
}

// registerClient adds a client to the hub and creates an avatar (if not reconnecting)
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.clients[client] = true
	
	// Only create avatar if client doesn't already have one (not a reconnection)
	if client.GetAvatarID() == "" {
		avatar := h.avatarRegistry.CreateAvatar(client)
		
		logging.Info("client registered with new avatar", map[string]interface{}{
			"client_count": len(h.clients),
			"session_id":   client.sessionID,
			"client_id":    client.GetClientID(),
			"avatar_id":    avatar.ID,
			"avatar_count": h.avatarRegistry.GetAvatarCount(),
		})
	} else {
		logging.Info("client registered with existing avatar", map[string]interface{}{
			"client_count": len(h.clients),
			"session_id":   client.sessionID,
			"client_id":    client.GetClientID(),
			"avatar_id":    client.GetAvatarID(),
			"avatar_count": h.avatarRegistry.GetAvatarCount(),
		})
	}
}

// unregisterClient removes a client from the hub and cleans up avatar
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		
		// Remove avatar when client disconnects
		if avatarID := client.GetAvatarID(); avatarID != "" {
			h.avatarRegistry.RemoveAvatar(avatarID)
		}
		
		logging.Info("client unregistered with avatar cleanup", map[string]interface{}{
			"client_count": len(h.clients),
			"session_id":   client.sessionID,
			"client_id":    client.GetClientID(),
			"avatar_id":    client.GetAvatarID(),
			"avatar_count": h.avatarRegistry.GetAvatarCount(),
		})
	}
}

// broadcastMessage sends a message to all clients
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// Client channel is full, disconnect it
			h.unregisterClient(client)
		}
	}
}

// broadcastOperation sends an operation to all clients
func (h *Hub) broadcastOperation(op *sync.Operation) {
	msg := Message{
		Type:      "operation",
		Operation: op,
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		logging.Error("failed to marshal operation", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	h.broadcast <- data
}

// SubmitOperation submits an operation to the sync system
func (h *Hub) SubmitOperation(op *sync.Operation) {
	h.sync.SubmitOperation(op)
	
	logging.Debug("operation submitted", map[string]interface{}{
		"sequence": op.SeqNum,
		"type":     op.Type,
	})
}

// GetSyncStats returns sync system statistics
func (h *Hub) GetSyncStats() map[string]interface{} {
	return h.sync.GetStats()
}

// GetStats returns sync system statistics (alias for compatibility)
func (h *Hub) GetStats() map[string]interface{} {
	return h.sync.GetStats()
}

// GetSync returns the sync system (for handler compatibility)
func (h *Hub) GetSync() *sync.ReliableSync {
	return h.sync
}

// GetFullSync returns all operations for full sync
func (h *Hub) GetFullSync() []*sync.Operation {
	return h.sync.GetAllOperations()
}

// GetMissingOperations returns operations in a range
func (h *Hub) GetMissingOperations(from, to uint64) []*sync.Operation {
	return h.sync.GetOperationsInRange(from, to)
}

// GetAvatarRegistry returns the avatar registry
func (h *Hub) GetAvatarRegistry() *AvatarRegistry {
	return h.avatarRegistry
}