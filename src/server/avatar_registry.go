// Package server provides avatar lifecycle management tied to WebSocket connections
package server

import (
	"fmt"
	"sync"
	"time"

	"holodeck1/logging"
	syncPkg "holodeck1/sync"
)

// Avatar represents a connected client in the Three.js scene
type Avatar struct {
	ID           string                 `json:"id"`
	ClientID     string                 `json:"client_id"`
	Name         string                 `json:"name"`
	Position     Vector3                `json:"position"`
	Rotation     *Vector3               `json:"rotation,omitempty"`
	Animation    string                 `json:"animation,omitempty"`
	Capabilities []string               `json:"capabilities"`
	ClientInfo   *ClientInfo            `json:"client_info,omitempty"`
	ConnectedAt  time.Time              `json:"connected_at"`
	LastSeen     time.Time              `json:"last_seen"`
	Client       *Client                `json:"-"` // Reference to WebSocket client
}

// Vector3 represents a 3D vector for Three.js
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// AvatarRegistry manages all connected avatars
type AvatarRegistry struct {
	avatars map[string]*Avatar
	mutex   sync.RWMutex
	hub     *Hub
}

// NewAvatarRegistry creates a new avatar registry
func NewAvatarRegistry(hub *Hub) *AvatarRegistry {
	return &AvatarRegistry{
		avatars: make(map[string]*Avatar),
		hub:     hub,
	}
}

// CreateAvatar creates a new avatar for a WebSocket connection
func (ar *AvatarRegistry) CreateAvatar(client *Client) *Avatar {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	// Generate unique avatar ID
	avatarID := fmt.Sprintf("avatar-%s", client.GetClientID())
	
	// Default spawn position
	position := Vector3{X: 0, Y: 0, Z: 0}
	
	// Create avatar
	avatar := &Avatar{
		ID:           avatarID,
		ClientID:     client.GetClientID(),
		Name:         fmt.Sprintf("User_%s", client.GetClientID()[:8]),
		Position:     position,
		Animation:    "idle",
		Capabilities: []string{"WebGL", "WebSocket"},
		ClientInfo:   client.info,
		ConnectedAt:  time.Now(),
		LastSeen:     time.Now(),
		Client:       client,
	}

	// Store avatar
	ar.avatars[avatarID] = avatar
	
	// Set avatar ID on client
	client.SetAvatarID(avatarID)

	logging.Info("avatar created", map[string]interface{}{
		"avatar_id":  avatarID,
		"client_id":  client.GetClientID(),
		"session_id": client.GetSessionID(),
		"position":   fmt.Sprintf("%.2f,%.2f,%.2f", position.X, position.Y, position.Z),
	})

	// Submit avatar_create operation to sync system
	operation := &syncPkg.Operation{
		ClientID: client.GetClientID(),
		Type:     "avatar_create",
		Data: map[string]interface{}{
			"avatar_id":    avatarID,
			"name":         avatar.Name,
			"position":     avatar.Position,
			"capabilities": avatar.Capabilities,
			"client_info":  avatar.ClientInfo,
		},
		Timestamp: time.Now(),
	}

	ar.hub.SubmitOperation(operation)

	return avatar
}

// FindAvatarByClientID finds an avatar by client ID
func (ar *AvatarRegistry) FindAvatarByClientID(clientID string) *Avatar {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()
	
	for _, avatar := range ar.avatars {
		if avatar.ClientID == clientID {
			return avatar
		}
	}
	return nil
}

// ReconnectClient reconnects an existing client to an avatar
func (ar *AvatarRegistry) ReconnectClient(clientID string, newClient *Client) *Avatar {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	
	// Find existing avatar
	for _, avatar := range ar.avatars {
		if avatar.ClientID == clientID {
			// Update client reference
			avatar.Client = newClient
			avatar.LastSeen = time.Now()
			
			// Set client's avatar ID
			newClient.SetAvatarID(avatar.ID)
			
			logging.Info("client reconnected to existing avatar", map[string]interface{}{
				"avatar_id":  avatar.ID,
				"client_id":  clientID,
				"session_id": newClient.GetSessionID(),
			})
			
			return avatar
		}
	}
	return nil
}

// RemoveAvatar removes an avatar when client disconnects
func (ar *AvatarRegistry) RemoveAvatar(avatarID string) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	avatar, exists := ar.avatars[avatarID]
	if !exists {
		return
	}

	// Remove from registry
	delete(ar.avatars, avatarID)

	logging.Info("avatar removed", map[string]interface{}{
		"avatar_id":  avatarID,
		"client_id":  avatar.ClientID,
		"session_id": avatar.Client.GetSessionID(),
		"duration":   time.Since(avatar.ConnectedAt).String(),
	})

	// Submit avatar_remove operation to sync system
	operation := &syncPkg.Operation{
		ClientID: avatar.ClientID,
		Type:     "avatar_remove",
		Data: map[string]interface{}{
			"avatar_id": avatarID,
		},
		Timestamp: time.Now(),
	}

	ar.hub.SubmitOperation(operation)
}

// UpdateAvatarPosition updates an avatar's position in the registry
func (ar *AvatarRegistry) UpdateAvatarPosition(avatarID string, positionData map[string]interface{}) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	
	if avatar, exists := ar.avatars[avatarID]; exists {
		// Update position from WebSocket data
		if x, ok := positionData["x"].(float64); ok {
			avatar.Position.X = x
		}
		if y, ok := positionData["y"].(float64); ok {
			avatar.Position.Y = y
		}
		if z, ok := positionData["z"].(float64); ok {
			avatar.Position.Z = z
		}
		
		// Update last seen time
		avatar.LastSeen = time.Now()
		
		logging.Trace("avatar", "position updated", map[string]interface{}{
			"avatar_id": avatarID,
			"position":  positionData,
		})
	}
}

// GetAvatar gets an avatar by ID
func (ar *AvatarRegistry) GetAvatar(avatarID string) (*Avatar, bool) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	avatar, exists := ar.avatars[avatarID]
	return avatar, exists
}

// GetAllAvatars returns all connected avatars
func (ar *AvatarRegistry) GetAllAvatars() []*Avatar {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	avatars := make([]*Avatar, 0, len(ar.avatars))
	for _, avatar := range ar.avatars {
		avatars = append(avatars, avatar)
	}

	return avatars
}

// UpdateAvatar updates an avatar's properties
func (ar *AvatarRegistry) UpdateAvatar(avatarID string, updates map[string]interface{}) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	avatar, exists := ar.avatars[avatarID]
	if !exists {
		return fmt.Errorf("avatar not found: %s", avatarID)
	}

	// Update avatar properties
	if position, ok := updates["position"].(map[string]interface{}); ok {
		if x, ok := position["x"].(float64); ok {
			avatar.Position.X = x
		}
		if y, ok := position["y"].(float64); ok {
			avatar.Position.Y = y
		}
		if z, ok := position["z"].(float64); ok {
			avatar.Position.Z = z
		}
	}

	if rotation, ok := updates["rotation"].(map[string]interface{}); ok {
		if avatar.Rotation == nil {
			avatar.Rotation = &Vector3{}
		}
		if x, ok := rotation["x"].(float64); ok {
			avatar.Rotation.X = x
		}
		if y, ok := rotation["y"].(float64); ok {
			avatar.Rotation.Y = y
		}
		if z, ok := rotation["z"].(float64); ok {
			avatar.Rotation.Z = z
		}
	}

	if animation, ok := updates["animation"].(string); ok {
		avatar.Animation = animation
	}

	avatar.LastSeen = time.Now()

	logging.Debug("avatar updated", map[string]interface{}{
		"avatar_id": avatarID,
		"client_id": avatar.ClientID,
		"updates":   updates,
	})

	return nil
}

// GetAvatarCount returns the number of connected avatars
func (ar *AvatarRegistry) GetAvatarCount() int {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	return len(ar.avatars)
}