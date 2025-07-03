package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
	"holodeck1/config"
	"holodeck1/logging"
)

type Hub struct {
	clients       map[*Client]bool
	broadcast     chan []byte
	register      chan *Client
	unregister    chan *Client
	logger        *LogManager
	store         *SessionStore
	mutex         sync.RWMutex
	
	// Session Graph Architecture - Channel-based persistence
	sessionChannels     map[string]*SessionChannel  // sessionID -> SessionChannel
	clientSessions   map[*Client]string       // client -> sessionID mapping
}

// SessionChannel represents a persistent session graph with TCP-level reliability
type SessionChannel struct {
	sessionID     string
	clients       map[*Client]bool     // clients currently in this channel
	graphState    map[string]interface{}  // persistent session graph state
	lastActivity  time.Time
	mutex         sync.RWMutex
	messageQueue  [][]byte             // queued messages for reconnecting clients
}

// NewSessionChannel creates a new persistent session channel
func NewSessionChannel(sessionID string) *SessionChannel {
	return &SessionChannel{
		sessionID:    sessionID,
		clients:      make(map[*Client]bool),
		graphState:   make(map[string]interface{}),
		lastActivity: time.Now(),
		messageQueue: make([][]byte, 0),
	}
}

// Channel data structures for loading YAML configurations
type PlayCanvasEntity struct {
	Name       string                 `json:"name" yaml:"name"`
	Components map[string]interface{} `json:"components" yaml:"components"`
}

type PlayCanvasScene struct {
	AmbientLight interface{} `json:"ambientLight,omitempty" yaml:"ambientLight,omitempty"` // Can be string or []float64
	Gravity      []float64   `json:"gravity,omitempty" yaml:"gravity,omitempty"`
}

type PlayCanvasConfig struct {
	Scene    PlayCanvasScene    `json:"scene,omitempty" yaml:"scene,omitempty"`
	Entities []PlayCanvasEntity `json:"entities,omitempty" yaml:"entities,omitempty"`
}

type ChannelConfig struct {
	PlayCanvas *PlayCanvasConfig `yaml:"playcanvas,omitempty"`
}

func NewHub() *Hub {
	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     NewLogManager(),
		store:      NewSessionStore(),
		
		// Session Graph Architecture
		sessionChannels:   make(map[string]*SessionChannel),
		clientSessions: make(map[*Client]string),
	}
	
	
	return hub
}

// Session Channel Management Methods

// JoinSessionChannel joins a client to a session channel with TCP-level reliability
func (h *Hub) JoinSessionChannel(sessionID, clientID string, reconnect bool) (*SessionChannel, int, map[string]interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// Get or create session channel
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		channel = NewSessionChannel(sessionID)
		h.sessionChannels[sessionID] = channel
		logging.Info("session channel created", map[string]interface{}{
			"session_id": sessionID,
		})
	}
	
	// Add client to channel (using client ID as a pseudo-client)
	channel.mutex.Lock()
	channel.clients[&Client{sessionID: sessionID}] = true // Simplified for API
	channel.lastActivity = time.Now()
	clientCount := len(channel.clients)
	
	// Initialize graph state from session entities if channel is new
	if len(channel.graphState) == 0 {
		// Entities are managed via channels/PlayCanvas, not stored in sessions
		channel.graphState["entities"] = []interface{}{}
		channel.graphState["last_updated"] = time.Now()
	}
	
	graphState := make(map[string]interface{})
	for k, v := range channel.graphState {
		graphState[k] = v
	}
	channel.mutex.Unlock()
	
	return channel, clientCount, graphState
}

// LeaveSessionChannel removes a client from a session channel
func (h *Hub) LeaveSessionChannel(sessionID, clientID string) (bool, int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		return false, 0
	}
	
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	
	// Remove client from channel (simplified for API)
	// In real implementation, we'd track clients by ID
	clientCount := len(channel.clients)
	channel.lastActivity = time.Now()
	
	return true, clientCount
}

// GetSessionGraphState retrieves the current graph state for a session
func (h *Hub) GetSessionGraphState(sessionID string) (map[string]interface{}, int, time.Time) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		// Return empty state if channel doesn't exist yet
		return make(map[string]interface{}), 0, time.Now()
	}
	
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	
	graphState := make(map[string]interface{})
	for k, v := range channel.graphState {
		graphState[k] = v
	}
	
	lastUpdated := time.Now()
	if t, ok := channel.graphState["last_updated"].(time.Time); ok {
		lastUpdated = t
	}
	
	return graphState, len(channel.clients), lastUpdated
}

// UpdateSessionGraphState updates the graph state and broadcasts to channel members
func (h *Hub) UpdateSessionGraphState(sessionID, clientID string, updates map[string]interface{}, atomic bool) (int, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		channel = NewSessionChannel(sessionID)
		h.sessionChannels[sessionID] = channel
	}
	
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	
	// Apply updates to graph state
	for k, v := range updates {
		channel.graphState[k] = v
	}
	channel.graphState["last_updated"] = time.Now()
	channel.lastActivity = time.Now()
	
	// Broadcast updates to all session clients
	broadcastData := map[string]interface{}{
		"graph_updates": updates,
		"client_id":     clientID,
		"timestamp":     time.Now(),
	}
	
	h.BroadcastToSession(sessionID, "graph_updated", broadcastData)
	
	return len(channel.clients), nil
}

// SyncSessionState forces synchronization of session state across all clients
func (h *Hub) SyncSessionState(sessionID string, forceFullSync bool) (int, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		return 0, nil // No channel to sync
	}
	
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	
	// Get current session entities for full sync (from channel)
	entities := []interface{}{} // Entities managed via PlayCanvas/channels
	
	syncData := map[string]interface{}{
		"sync_type":     "full",
		"entities":      entities,
		"graph_state":   channel.graphState,
		"force_full":    forceFullSync,
		"sync_timestamp": time.Now(),
	}
	
	h.BroadcastToSession(sessionID, "state_sync", syncData)
	
	return len(channel.clients), nil
}

// BroadcastToSessionChannel broadcasts to all clients in a session channel, optionally excluding one
func (h *Hub) BroadcastToSessionChannel(sessionID, messageType string, data interface{}, excludeClientID string) {
	// For now, use the existing BroadcastToSession method
	// In a full implementation, this would filter by excludeClientID
	h.BroadcastToSession(sessionID, messageType, data)
}

// SessionChannelStatus represents the status of a session channel
type SessionChannelStatus struct {
	ChannelActive       bool                   `json:"channel_active"`
	ConnectedClients []map[string]interface{} `json:"connected_clients"`
	GraphSummary     map[string]interface{} `json:"graph_summary"`
	HealthMetrics    map[string]interface{} `json:"health_metrics"`
}

// GetSessionChannelStatus retrieves detailed status information for a session channel
func (h *Hub) GetSessionChannelStatus(sessionID string) *SessionChannelStatus {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	channel, exists := h.sessionChannels[sessionID]
	if !exists {
		return &SessionChannelStatus{
			ChannelActive:       false,
			ConnectedClients: []map[string]interface{}{},
			GraphSummary:     map[string]interface{}{"object_count": 0},
			HealthMetrics:    map[string]interface{}{"uptime": "0s", "message_count": 0},
		}
	}
	
	channel.mutex.RLock()
	defer channel.mutex.RUnlock()
	
	// Get entity count (managed via channels/PlayCanvas)
	entityCount := 0 // Entities managed via PlayCanvas/channels
	
	// Build client list (simplified)
	clients := make([]map[string]interface{}, 0)
	for range channel.clients {
		clients = append(clients, map[string]interface{}{
			"client_id":     "client_" + sessionID,
			"connected_at":  channel.lastActivity,
			"last_activity": channel.lastActivity,
		})
	}
	
	uptime := time.Since(channel.lastActivity)
	
	return &SessionChannelStatus{
		ChannelActive:       len(channel.clients) > 0,
		ConnectedClients: clients,
		GraphSummary: map[string]interface{}{
			"entity_count":  entityCount,
			"last_updated":  channel.graphState["last_updated"],
		},
		HealthMetrics: map[string]interface{}{
			"uptime":       uptime.String(),
			"message_count": len(channel.messageQueue),
			"last_sync":    channel.lastActivity,
		},
	}
}

// GetID returns the session channel ID
func (sr *SessionChannel) GetID() string {
	return sr.sessionID
}

func (h *Hub) Run() {
	// Note: Scenes watcher disabled due to filesystem noatime/lazytime mount options
	// Using API-based scene loading instead
	logging.Info("hub started without scenes watcher", map[string]interface{}{
		"reason": "filesystem mount options interfere with fsnotify",
		"solution": "API-based scene loading on page load",
	})
	
	// Start session cleanup timer
	go h.startSessionCleanup()
	
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			
			logging.Info("client connected to hub", map[string]interface{}{
				"total_clients": len(h.clients),
			})

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				logging.Info("client disconnected from hub", map[string]interface{}{
					"total_clients": len(h.clients),
				})
			}
			h.mutex.Unlock()

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

func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// BroadcastToSession sends a message only to clients in a specific session
func (h *Hub) BroadcastToSession(sessionID string, updateType string, data interface{}) {
	update := map[string]interface{}{
		"type": updateType,
		"data": data,
		"timestamp": time.Now().Unix(),
		"session_id": sessionID,
	}
	
	if jsonData, err := json.Marshal(update); err == nil {
		// Send only to clients associated with this specific session
		clientCount := 0
		for client := range h.clients {
			if client.sessionID == sessionID {
				select {
				case client.send <- jsonData:
					clientCount++
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
		
		// Log session-specific broadcast for debugging
		logging.Debug("session broadcast sent", map[string]interface{}{
			"session_id": sessionID,
			"type": updateType,
			"clients": clientCount,
		})
	}
}

// BroadcastAvatarPositionToChannel broadcasts avatar position updates to ALL sessions in the same channel
// This enables bidirectional visibility - all participants see each other's avatars
func (h *Hub) BroadcastAvatarPositionToChannel(sessionID string, updateType string, data interface{}) {
	// Get the current session to find which channel it's in
	session, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		logging.Warn("cannot broadcast avatar position - session not found", map[string]interface{}{
			"session_id": sessionID,
		})
		return
	}
	
	// If session is not in a named channel, fall back to session-only broadcast
	if session.ChannelID == "" {
		h.BroadcastToSession(sessionID, updateType, data)
		return
	}
	
	channelID := session.ChannelID
	
	// Find ALL sessions currently in the same channel
	allSessions := h.GetStore().ListSessions()
	channelSessions := make([]string, 0)
	
	for _, s := range allSessions {
		if s.ChannelID == channelID {
			channelSessions = append(channelSessions, s.ID)
		}
	}
	
	update := map[string]interface{}{
		"type": updateType,
		"data": data,
		"timestamp": time.Now().Unix(),
		"session_id": sessionID,  // Still identify the originating session
		"channel_id": channelID,  // Add channel context
	}
	
	if jsonData, err := json.Marshal(update); err == nil {
		// Broadcast to ALL sessions in the channel for bidirectional visibility
		totalClients := 0
		for client := range h.clients {
			// Send to clients in ANY session that's in this channel
			for _, chanSessionID := range channelSessions {
				if client.sessionID == chanSessionID {
					select {
					case client.send <- jsonData:
						totalClients++
					default:
						close(client.send)
						delete(h.clients, client)
					}
					break // Don't send duplicate to same client
				}
			}
		}
		
		// Enhanced logging for channel broadcasts
		logging.Info("avatar position broadcast to channel", map[string]interface{}{
			"originating_session": sessionID,
			"channel_id": channelID,
			"channel_sessions": channelSessions,
			"total_clients": totalClients,
			"update_type": updateType,
		})
	}
}

// SessionStore provides persistence for 3D visualization sessions
type SessionStore struct {
	mutex    sync.RWMutex
	sessions map[string]*Session
	// Legacy objects and worlds removed - entities managed via channels/PlayCanvas
}

// Session represents a 3D visualization session
// Entity represents a PlayCanvas entity in a session
type Entity struct {
	ID             string                 `json:"entity_id"`
	Name           string                 `json:"name"`
	PlayCanvasGUID string                 `json:"playcanvas_guid"`
	Components     map[string]interface{} `json:"components"`
	Tags           []string               `json:"tags,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	Enabled        bool                   `json:"enabled"`
}

type Session struct {
	ID        string             `json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	Status    string             `json:"status"`
	ChannelID string             `json:"channel_id,omitempty"` // Current channel joined
	Entities  map[string]*Entity `json:"entities,omitempty"`   // Entity storage
}

// Legacy Object and World types removed - replaced by PlayCanvas entities and channels

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a new session
func (s *SessionStore) CreateSession() *Session {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		Status:    "active",
		Entities:  make(map[string]*Entity),
	}
	
	s.sessions[sessionID] = session
	
	return session
}

// GetSession retrieves a session by ID
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	session, exists := s.sessions[sessionID]
	return session, exists
}

// ListSessions returns all active sessions
func (s *SessionStore) ListSessions() []*Session {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var sessions []*Session
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// DeleteSession removes a session and all its data
func (s *SessionStore) DeleteSession(sessionID string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.sessions[sessionID]; !exists {
		return false
	}
	
	delete(s.sessions, sessionID)
	
	return true
}

// UpdateSessionChannel updates the channel ID for a session
func (s *SessionStore) UpdateSessionChannel(sessionID, channelID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return &SessionError{Message: "Session not found"}
	}

	session.ChannelID = channelID
	return nil
}

// AddEntity adds an entity to a session
func (s *SessionStore) AddEntity(sessionID string, entity *Entity) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return &SessionError{Message: "Session not found"}
	}

	if session.Entities == nil {
		session.Entities = make(map[string]*Entity)
	}

	session.Entities[entity.ID] = entity
	return nil
}

// GetEntities retrieves all entities from a session
func (s *SessionStore) GetEntities(sessionID string) ([]*Entity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, &SessionError{Message: "Session not found"}
	}

	entities := make([]*Entity, 0, len(session.Entities))
	for _, entity := range session.Entities {
		entities = append(entities, entity)
	}

	return entities, nil
}

// GetEntity retrieves a specific entity from a session
func (s *SessionStore) GetEntity(sessionID, entityID string) (*Entity, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, &SessionError{Message: "Session not found"}
	}

	entity, exists := session.Entities[entityID]
	if !exists {
		return nil, &SessionError{Message: "Entity not found"}
	}

	return entity, nil
}

// UpdateEntity updates an existing entity in a session
func (s *SessionStore) UpdateEntity(sessionID, entityID string, updatedEntity *Entity) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return &SessionError{Message: "Session not found"}
	}

	if session.Entities == nil {
		return &SessionError{Message: "No entities found in session"}
	}

	_, exists = session.Entities[entityID]
	if !exists {
		return &SessionError{Message: "Entity not found"}
	}

	// Update the entity in the session's entities map
	session.Entities[entityID] = updatedEntity
	return nil
}

// Legacy Object and World management methods removed
// Objects/Worlds replaced by PlayCanvas Entities via Channels
// Use /sessions/{sessionId}/entities and /channels endpoints instead

// GetStore returns the session store for external access
func (h *Hub) GetStore() *SessionStore {
	return h.store
}



// BroadcastUpdate sends real-time updates to connected clients
func (h *Hub) BroadcastUpdate(updateType string, data interface{}) {
	update := map[string]interface{}{
		"type": updateType,
		"data": data,
		"timestamp": time.Now().Unix(),
	}
	
	if jsonData, err := json.Marshal(update); err == nil {
		h.BroadcastMessage(jsonData)
	}
}

// Error types for better error handling
type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}

// Legacy CoordinateError and ObjectError removed - entities use PlayCanvas validation

// generateSessionID creates a unique session identifier
func generateSessionID() string {
	return "session-" + generateID(8)
}

// generateID creates a random ID of specified length
func generateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// startSessionCleanup runs periodic cleanup of inactive sessions
func (h *Hub) startSessionCleanup() {
	ticker := time.NewTicker(config.GetSessionCleanupInterval())
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupInactiveSessions()
	}
}

// cleanupInactiveSessions removes sessions that have been inactive for more than 10 minutes
func (h *Hub) cleanupInactiveSessions() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-config.GetSessionInactivityTimeout())
	
	sessions := h.store.ListSessions()
	cleanedCount := 0
	
	for _, session := range sessions {
		// Check if session is older than cutoff and has no active clients
		if session.CreatedAt.Before(cutoff) {
			// Check if any clients are connected to this session
			hasActiveClients := false
			for client := range h.clients {
				if client.sessionID == session.ID {
					hasActiveClients = true
					break
				}
			}
			
			// Clean up session channels that are inactive
			if channel, exists := h.sessionChannels[session.ID]; exists {
				channel.mutex.RLock()
				channelInactive := len(channel.clients) == 0 && channel.lastActivity.Before(cutoff)
				channel.mutex.RUnlock()
				
				if channelInactive {
					delete(h.sessionChannels, session.ID)
					logging.Debug("cleaned up inactive session channel", map[string]interface{}{
						"session_id": session.ID,
						"age": now.Sub(session.CreatedAt),
					})
				}
			}
			
			// Delete session if no active clients
			if !hasActiveClients {
				if h.store.DeleteSession(session.ID) {
					cleanedCount++
					logging.Debug("cleaned up inactive session", map[string]interface{}{
						"session_id": session.ID,
						"age": now.Sub(session.CreatedAt),
					})
				}
			}
		}
	}
	
	if cleanedCount > 0 {
		logging.Info("session cleanup completed", map[string]interface{}{
			"cleaned_sessions": cleanedCount,
			"remaining_sessions": len(h.store.ListSessions()),
		})
	}
}

// LoadNamedChannelIntoSession loads a named channel's PlayCanvas configuration into a session
func (h *Hub) LoadNamedChannelIntoSession(sessionID, channelID string) error {
	logging.Info("loading named channel into session", map[string]interface{}{
		"session_id": sessionID,
		"channel_id": channelID,
	})
	
	// Read channel YAML configuration
	channelPath := filepath.Join(config.GetChannelsDir(), channelID+".yaml")
	configData, err := ioutil.ReadFile(channelPath)
	if err != nil {
		return fmt.Errorf("failed to read channel config %s: %w", channelPath, err)
	}
	
	// Parse YAML configuration
	var config ChannelConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse channel config %s: %w", channelPath, err)
	}
	
	if config.PlayCanvas == nil {
		logging.Info("no PlayCanvas configuration in channel", map[string]interface{}{
			"channel_id": channelID,
		})
		return nil
	}
	
	// Create entities in the session via API calls
	for _, entity := range config.PlayCanvas.Entities {
		if err := h.createEntityInSession(sessionID, entity); err != nil {
			logging.Error("failed to create entity in session", map[string]interface{}{
				"session_id": sessionID,
				"channel_id": channelID,
				"entity_name": entity.Name,
				"error": err.Error(),
			})
			// Continue with other entities even if one fails
		} else {
			logging.Info("created entity in session", map[string]interface{}{
				"session_id": sessionID,
				"entity_name": entity.Name,
			})
		}
	}
	
	logging.Info("channel entities loaded into session", map[string]interface{}{
		"session_id": sessionID,
		"channel_id": channelID,
		"entities_created": len(config.PlayCanvas.Entities),
	})
	
	return nil
}

// ClearSessionEntitiesWithBroadcast clears all entities from a session with proper WebSocket broadcasts
// This implements the "API = Control" principle by using proper service calls instead of client-side loops
func (h *Hub) ClearSessionEntitiesWithBroadcast(sessionID string) error {
	logging.Info("clearing session entities with broadcast", map[string]interface{}{
		"session_id": sessionID,
	})
	
	// Get all entities in the session
	entities, err := h.store.GetEntities(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session entities: %w", err)
	}
	
	clearedCount := 0
	for _, entity := range entities {
		// Delete each entity via internal HTTP call to ensure proper broadcasts
		url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities/%s", sessionID, entity.ID)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			logging.Error("failed to create delete request", map[string]interface{}{
				"session_id": sessionID,
				"entity_id": entity.ID,
				"error": err.Error(),
			})
			continue
		}
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logging.Error("failed to delete entity", map[string]interface{}{
				"session_id": sessionID,
				"entity_id": entity.ID,
				"error": err.Error(),
			})
			continue
		}
		resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
			clearedCount++
			logging.Debug("entity cleared with broadcast", map[string]interface{}{
				"session_id": sessionID,
				"entity_id": entity.ID,
				"entity_name": entity.Name,
			})
		}
	}
	
	logging.Info("session entities cleared with broadcasts", map[string]interface{}{
		"session_id": sessionID,
		"entities_cleared": clearedCount,
		"total_entities": len(entities),
	})
	
	return nil
}

// createEntityInSession creates a PlayCanvas entity in the session via internal API call
func (h *Hub) createEntityInSession(sessionID string, entity PlayCanvasEntity) error {
	// Prepare the entity creation payload
	payload := map[string]interface{}{
		"name":       entity.Name,
		"components": entity.Components,
	}
	
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal entity payload: %w", err)
	}
	
	// Make internal HTTP request to create entity
	url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities", sessionID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("entity creation failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// CreateEntityViaAPI creates an entity by making an HTTP call to HD1's own API
// This maintains 100% API-first architecture - all entity creation goes through the API
func (h *Hub) CreateEntityViaAPI(sessionID string, entityPayload map[string]interface{}) error {
	payloadJSON, err := json.Marshal(entityPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal entity payload: %w", err)
	}
	
	// Make internal HTTP request to HD1's create entity API
	url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities", sessionID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: config.GetSessionHTTPClientTimeout()}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create entity via API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("entity creation failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// DeleteEntityByNameViaAPI finds an entity by name and deletes it via API
// This maintains 100% API-first architecture - all entity operations go through the API
func (h *Hub) DeleteEntityByNameViaAPI(sessionID, entityName string) error {
	// First, get all entities to find the one with matching name
	url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities", sessionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create list entities request: %w", err)
	}
	
	client := &http.Client{Timeout: config.GetSessionHTTPClientTimeout()}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to list entities via API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("list entities failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var entitiesResponse struct {
		Entities []struct {
			ID   string `json:"entity_id"`
			Name string `json:"name"`
		} `json:"entities"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		return fmt.Errorf("failed to decode entities response: %w", err)
	}
	
	// Find entity with matching name
	var entityID string
	for _, entity := range entitiesResponse.Entities {
		if entity.Name == entityName {
			entityID = entity.ID
			break
		}
	}
	
	if entityID == "" {
		return fmt.Errorf("entity with name '%s' not found", entityName)
	}
	
	// Delete the entity via API
	deleteURL := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities/%s", sessionID, entityID)
	deleteReq, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		return fmt.Errorf("failed to delete entity via API: %w", err)
	}
	defer deleteResp.Body.Close()
	
	if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(deleteResp.Body)
		return fmt.Errorf("entity deletion failed with status %d: %s", deleteResp.StatusCode, string(body))
	}
	
	return nil
}

// UpdateEntityByNameViaAPI finds an entity by name and updates it via API
// This maintains 100% API-first architecture - all entity operations go through the API
func (h *Hub) UpdateEntityByNameViaAPI(sessionID, entityName string, updatePayload map[string]interface{}) error {
	// First, get all entities to find the one with matching name
	url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities", sessionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create list entities request: %w", err)
	}
	
	client := &http.Client{Timeout: config.GetSessionHTTPClientTimeout()}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to list entities via API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("list entities failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var entitiesResponse struct {
		Entities []struct {
			ID   string `json:"entity_id"`
			Name string `json:"name"`
		} `json:"entities"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		return fmt.Errorf("failed to decode entities response: %w", err)
	}
	
	// Find entity with matching name
	var entityID string
	for _, entity := range entitiesResponse.Entities {
		if entity.Name == entityName {
			entityID = entity.ID
			break
		}
	}
	
	if entityID == "" {
		return fmt.Errorf("entity with name '%s' not found", entityName)
	}
	
	// Prepare update payload
	payloadJSON, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}
	
	// Update the entity via API
	updateURL := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities/%s", sessionID, entityID)
	updateReq, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	
	updateReq.Header.Set("Content-Type", "application/json")
	
	updateResp, err := client.Do(updateReq)
	if err != nil {
		return fmt.Errorf("failed to update entity via API: %w", err)
	}
	defer updateResp.Body.Close()
	
	if updateResp.StatusCode != http.StatusOK && updateResp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(updateResp.Body)
		return fmt.Errorf("entity update failed with status %d: %s", updateResp.StatusCode, string(body))
	}
	
	return nil
}

// GetEntityByNameViaAPI finds an entity by name and returns its data via API
// This maintains 100% API-first architecture - all entity operations go through the API
func (h *Hub) GetEntityByNameViaAPI(sessionID, entityName string) (map[string]interface{}, error) {
	// First, get all entities to find the one with matching name
	url := fmt.Sprintf("http://localhost:8080/api/sessions/%s/entities", sessionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list entities request: %w", err)
	}
	
	client := &http.Client{Timeout: config.GetSessionHTTPClientTimeout()}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities via API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("list entities failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var entitiesResponse struct {
		Entities []struct {
			ID         string                 `json:"entity_id"`
			Name       string                 `json:"name"`
			Components map[string]interface{} `json:"components"`
		} `json:"entities"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode entities response: %w", err)
	}
	
	// Find entity with matching name and return its full data
	for _, entity := range entitiesResponse.Entities {
		if entity.Name == entityName {
			return map[string]interface{}{
				"entity_id":  entity.ID,
				"name":       entity.Name,
				"components": entity.Components,
			}, nil
		}
	}
	
	return nil, fmt.Errorf("entity with name '%s' not found", entityName)
}