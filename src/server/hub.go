// Package server provides the new TCP-simple WebSocket hub
// Replaces the complex system with sequence number reliability
package server

import (
	"context"
	stdSync "sync"

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
	
	// Message routing - REMOVED: Using sync system directly
}

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	ClientID  string                 `json:"hd1_id,omitempty"`
	Operation *sync.Operation        `json:"operation,omitempty"`
	From      uint64                 `json:"from,omitempty"`
	To        uint64                 `json:"to,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NewHub creates a new TCP-simple WebSocket hub
func NewHub() *Hub {
	hub := &Hub{
		sync:           sync.NewReliableSync(),
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}
	
	// Initialize avatar registry
	hub.avatarRegistry = NewAvatarRegistry(hub)
	
	return hub
}

// Run starts the hub's main loop with pure in-memory architecture
func (h *Hub) Run(ctx context.Context) {
	logging.Info("HD1 hub started with stateless in-memory architecture", map[string]interface{}{
		"sync_protocol": "TCP-simple reliable",
		"avatar_cleanup": "WebSocket connection-based",
		"stateless": true,
	})
	
	for {
		select {
		case <-ctx.Done():
			logging.Info("hub shutting down", nil)
			return
		case client := <-h.register:
			h.registerClient(client)
			
		case client := <-h.unregister:
			h.unregisterClient(client)
		}
	}
}

// handleOperations - REMOVED: Using sync system directly instead of polling

// registerClient adds a client to the hub and creates an avatar (if not reconnecting)
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.clients[client] = true
	
	// Register client with sync system - SINGLE SOURCE OF TRUTH
	syncChan := h.sync.RegisterClient(client.GetHD1ID())
	client.syncChan = syncChan
	
	// Start sync forwarding goroutine
	go client.forwardSyncOperations()
	
	// Send initial sync for existing operations
	client.sendInitialSync()
	
	// Only create avatar if client doesn't already have one (not a reconnection)
	if client.GetAvatarID() == "" {
		avatar := h.avatarRegistry.CreateAvatar(client)
		
		logging.Info("client registered with new avatar and sync channel", map[string]interface{}{
			"client_count": len(h.clients),
			"hd1_id":       client.GetClientID(),
			"avatar_id":    avatar.ID,
			"avatar_count": h.avatarRegistry.GetAvatarCount(),
		})
	} else {
		logging.Info("client registered with existing avatar and sync channel", map[string]interface{}{
			"client_count": len(h.clients),
			"hd1_id":       client.GetClientID(),
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
		
		// Unregister from sync system - SINGLE SOURCE OF TRUTH
		h.sync.UnregisterClient(client.GetHD1ID())
		
		// Remove avatar when client disconnects
		if avatarID := client.GetAvatarID(); avatarID != "" {
			h.avatarRegistry.RemoveAvatar(avatarID)
		}
		
		logging.Info("client unregistered with avatar cleanup and sync cleanup", map[string]interface{}{
			"client_count": len(h.clients),
			"hd1_id":       client.GetClientID(),
			"avatar_id":    client.GetAvatarID(),
			"avatar_count": h.avatarRegistry.GetAvatarCount(),
		})
	}
}

// broadcastMessage - REMOVED: Using sync system directly instead
// broadcastOperation - REMOVED: Using sync system directly instead

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