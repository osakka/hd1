package server

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *LogManager
	store      *SessionStore
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     NewLogManager(),
		store:      NewSessionStore(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected. Total: %d", len(h.clients))
			if h.logger != nil {
				h.logger.Log("info", "SERVER", "Client connected", map[string]interface{}{
					"total_clients": len(h.clients),
				})
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total: %d", len(h.clients))
				if h.logger != nil {
					h.logger.Log("info", "SERVER", "Client disconnected", map[string]interface{}{
						"total_clients": len(h.clients),
					})
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

func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// SessionStore provides persistence for 3D visualization sessions
type SessionStore struct {
	mutex    sync.RWMutex
	sessions map[string]*Session
	objects  map[string]map[string]*Object // sessionId -> objectName -> Object
	worlds   map[string]*World            // sessionId -> World
}

// Session represents a 3D visualization session
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

// Object represents a 3D object in the world
type Object struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Color    string  `json:"color,omitempty"`
	Scale    float64 `json:"scale,omitempty"`
	Rotation string  `json:"rotation,omitempty"`
}

// World represents the 3D world coordinate system
type World struct {
	Size         int     `json:"size"`
	Transparency float64 `json:"transparency"`
	CameraX      float64 `json:"camera_x"`
	CameraY      float64 `json:"camera_y"`
	CameraZ      float64 `json:"camera_z"`
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
		objects:  make(map[string]map[string]*Object),
		worlds:   make(map[string]*World),
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
	}
	
	s.sessions[sessionID] = session
	s.objects[sessionID] = make(map[string]*Object)
	
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
	delete(s.objects, sessionID)
	delete(s.worlds, sessionID)
	
	return true
}

// CreateObject creates a new object in a session
func (s *SessionStore) CreateObject(sessionID, objectName, objectType string, x, y, z float64) (*Object, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Validate session exists
	if _, exists := s.sessions[sessionID]; !exists {
		return nil, &SessionError{Message: "Session not found"}
	}
	
	// Validate coordinates are within world bounds [-12, +12]
	if x < -12 || x > 12 || y < -12 || y > 12 || z < -12 || z > 12 {
		return nil, &CoordinateError{Message: "Coordinates must be within [-12, +12] bounds"}
	}
	
	object := &Object{
		Name: objectName,
		Type: objectType,
		X:    x,
		Y:    y,
		Z:    z,
		Scale: 1.0,
	}
	
	s.objects[sessionID][objectName] = object
	return object, nil
}

// GetObject retrieves an object from a session
func (s *SessionStore) GetObject(sessionID, objectName string) (*Object, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if objects, exists := s.objects[sessionID]; exists {
		object, found := objects[objectName]
		return object, found
	}
	return nil, false
}

// ListObjects returns all objects in a session
func (s *SessionStore) ListObjects(sessionID string) []*Object {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var objects []*Object
	if sessionObjects, exists := s.objects[sessionID]; exists {
		for _, object := range sessionObjects {
			objects = append(objects, object)
		}
	}
	return objects
}

// UpdateObject updates an existing object
func (s *SessionStore) UpdateObject(sessionID, objectName string, updates map[string]interface{}) (*Object, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if objects, exists := s.objects[sessionID]; exists {
		if object, found := objects[objectName]; found {
			// Apply updates with validation
			if x, ok := updates["x"].(float64); ok {
				if x < -12 || x > 12 {
					return nil, &CoordinateError{Message: "X coordinate must be within [-12, +12] bounds"}
				}
				object.X = x
			}
			if y, ok := updates["y"].(float64); ok {
				if y < -12 || y > 12 {
					return nil, &CoordinateError{Message: "Y coordinate must be within [-12, +12] bounds"}
				}
				object.Y = y
			}
			if z, ok := updates["z"].(float64); ok {
				if z < -12 || z > 12 {
					return nil, &CoordinateError{Message: "Z coordinate must be within [-12, +12] bounds"}
				}
				object.Z = z
			}
			if color, ok := updates["color"].(string); ok {
				object.Color = color
			}
			if scale, ok := updates["scale"].(float64); ok {
				object.Scale = scale
			}
			if rotation, ok := updates["rotation"].(string); ok {
				object.Rotation = rotation
			}
			
			return object, nil
		}
	}
	return nil, &ObjectError{Message: "Object not found"}
}

// DeleteObject removes an object from a session
func (s *SessionStore) DeleteObject(sessionID, objectName string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if objects, exists := s.objects[sessionID]; exists {
		if _, found := objects[objectName]; found {
			delete(objects, objectName)
			return true
		}
	}
	return false
}

// InitializeWorld sets up the 3D world for a session
func (s *SessionStore) InitializeWorld(sessionID string, size int, transparency float64, cameraX, cameraY, cameraZ float64) (*World, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Validate session exists
	if _, exists := s.sessions[sessionID]; !exists {
		return nil, &SessionError{Message: "Session not found"}
	}
	
	world := &World{
		Size:         size,
		Transparency: transparency,
		CameraX:      cameraX,
		CameraY:      cameraY,
		CameraZ:      cameraZ,
	}
	
	s.worlds[sessionID] = world
	return world, nil
}

// GetWorld retrieves world configuration for a session
func (s *SessionStore) GetWorld(sessionID string) (*World, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	world, exists := s.worlds[sessionID]
	return world, exists
}

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

type CoordinateError struct {
	Message string
}

func (e *CoordinateError) Error() string {
	return e.Message
}

type ObjectError struct {
	Message string
}

func (e *ObjectError) Error() string {
	return e.Message
}

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