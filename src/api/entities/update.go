package entities

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// UpdateEntityRequest represents the request body for entity updates
type UpdateEntityRequest struct {
	Name     *string   `json:"name,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Enabled  *bool     `json:"enabled,omitempty"`
	Position []float64 `json:"position,omitempty"`
	Rotation []float64 `json:"rotation,omitempty"`
	Scale    []float64 `json:"scale,omitempty"`
}

// UpdateEntityHandler handles PUT /sessions/{sessionId}/entities/{entityId}
func UpdateEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal server error",
			"message": "Hub interface cast failed",
		})
		return
	}
	// Extract sessionId and entityId from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 4 || pathParts[0] != "sessions" || pathParts[2] != "entities" {
		logging.Error("invalid URL path for update entity", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities/{entityId}",
		})
		return
	}
	
	sessionID := pathParts[1]
	entityID := pathParts[3]
	
	// Validate session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		logging.Warn("session not found", map[string]interface{}{
			"session_id": sessionID,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session not found",
			"message": "Session does not exist or has expired",
		})
		return
	}
	
	// Parse request body
	var req UpdateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("invalid JSON in request", map[string]interface{}{
			"session_id": sessionID,
			"entity_id":  entityID,
			"error":      err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid JSON",
			"message": "Request body must be valid JSON",
		})
		return
	}
	
	// Validate position array if provided
	if req.Position != nil && len(req.Position) != 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid position",
			"message": "Position must be [x, y, z] array",
		})
		return
	}
	
	// Validate rotation array if provided
	if req.Rotation != nil && len(req.Rotation) != 4 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid rotation",
			"message": "Rotation must be [x, y, z, w] quaternion array",
		})
		return
	}
	
	// Validate scale array if provided
	if req.Scale != nil && len(req.Scale) != 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid scale",
			"message": "Scale must be [x, y, z] array",
		})
		return
	}
	
	// Check if entity exists (mock implementation)
	if !strings.HasPrefix(entityID, "entity-") {
		logging.Warn("entity not found", map[string]interface{}{
			"session_id": sessionID,
			"entity_id":  entityID,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Entity not found",
			"message": "Entity does not exist in session",
		})
		return
	}
	
	// TODO: Implement PlayCanvas entity updates
	// For now, just track what would be changed
	changes := []string{}
	if req.Name != nil {
		changes = append(changes, "name")
	}
	if req.Tags != nil {
		changes = append(changes, "tags")
	}
	if req.Enabled != nil {
		changes = append(changes, "enabled")
	}
	if req.Position != nil {
		changes = append(changes, "position")
	}
	if req.Rotation != nil {
		changes = append(changes, "rotation")
	}
	if req.Scale != nil {
		changes = append(changes, "scale")
	}
	
	logging.Info("entity updated", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"changes":    changes,
	})
	
	// Broadcast entity update via WebSocket
	h.BroadcastUpdate("entity_updated", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"changes":    changes,
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"entity_id":  entityID,
		"changes":    changes,
		"updated_at": time.Now().Format(time.RFC3339),
	})
}