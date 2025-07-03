package avatar

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// SessionAvatar represents the current session avatar state
type SessionAvatar struct {
	AvatarType string                 `json:"avatar_type"`
	EntityID   string                 `json:"entity_id"`
	Created    time.Time              `json:"created"`
	Position   map[string]float64     `json:"position"`
	Rotation   map[string]float64     `json:"rotation"`
	Status     string                 `json:"status"`
	LastUpdate time.Time              `json:"last_update"`
	SessionID  string                 `json:"session_id"`
	UserID     string                 `json:"user_id,omitempty"`
}

// GetSessionAvatarHandler handles GET /sessions/{sessionId}/avatar
func GetSessionAvatarHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logger := logging.GetLogger()
	sessionID := extractSessionID(r.URL.Path)

	logger.Info("Getting session avatar", map[string]interface{}{
		"session_id": sessionID,
	})

	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logger.Error("Hub not available")
		http.Error(w, "Server not available", http.StatusInternalServerError)
		return
	}

	// Get session
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		logger.Warn("Session not found", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Find avatar entity for this session
	entities, err := h.GetStore().GetEntities(sessionID)
	if err != nil {
		logger.Error("Failed to get entities", map[string]interface{}{
			"session_id": sessionID,
			"error":      err,
		})
		http.Error(w, "Failed to get entities", http.StatusInternalServerError)
		return
	}
	
	var avatarEntity *server.Entity
	for _, entity := range entities {
		for _, tag := range entity.Tags {
			if tag == "session-avatar" || tag == "avatar" {
				// Found avatar entity
				avatarEntity = entity
				break
			}
		}
		if avatarEntity != nil {
			break
		}
	}

	if avatarEntity == nil {
		logger.Warn("No avatar found for session", map[string]interface{}{
			"session_id": sessionID,
		})
		response := map[string]interface{}{
			"success": false,
			"message": "No avatar assigned to session",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get avatar position and rotation
	transform := avatarEntity.Components["transform"]
	position := map[string]float64{"x": 0, "y": 0, "z": 0}
	rotation := map[string]float64{"x": 0, "y": 0, "z": 0}

	if transform != nil {
		if transformMap, ok := transform.(map[string]interface{}); ok {
			if pos, ok := transformMap["position"].(map[string]interface{}); ok {
				if x, ok := pos["x"].(float64); ok {
					position["x"] = x
				}
				if y, ok := pos["y"].(float64); ok {
					position["y"] = y
				}
				if z, ok := pos["z"].(float64); ok {
					position["z"] = z
				}
			}
			if rot, ok := transformMap["rotation"].(map[string]interface{}); ok {
				if x, ok := rot["x"].(float64); ok {
					rotation["x"] = x
				}
				if y, ok := rot["y"].(float64); ok {
					rotation["y"] = y
				}
				if z, ok := rot["z"].(float64); ok {
					rotation["z"] = z
				}
			}
		}
	}

	// Determine avatar type from tags
	avatarType := "default"
	for _, tag := range avatarEntity.Tags {
		if tag == "business" {
			avatarType = "business"
			break
		} else if tag == "casual" {
			avatarType = "casual"
			break
		}
	}

	// Build response
	sessionAvatar := SessionAvatar{
		AvatarType: avatarType,
		EntityID:   avatarEntity.ID,
		Created:    avatarEntity.CreatedAt,
		Position:   position,
		Rotation:   rotation,
		Status:     "active",
		LastUpdate: time.Now(),
		SessionID:  sessionID,
	}

	response := map[string]interface{}{
		"success": true,
		"avatar":  sessionAvatar,
		"message": "Session avatar retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode session avatar response", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info("Session avatar retrieved successfully", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  avatarEntity.ID,
	})
}

// extractSessionID extracts session ID from URL path  
func extractSessionID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[1] == "sessions" {
		return parts[2]
	}
	return ""
}