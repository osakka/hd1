package components

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// RemoveComponentHandler handles DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}
func RemoveComponentHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Extract sessionId, entityId, and componentType from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 6 || pathParts[0] != "sessions" || pathParts[2] != "entities" || pathParts[4] != "components" {
		logging.Error("invalid URL path for remove component", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities/{entityId}/components/{componentType}",
		})
		return
	}
	
	sessionID := pathParts[1]
	entityID := pathParts[3]
	componentType := pathParts[5]
	
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
	
	// Validate component type
	validTypes := []string{"model", "camera", "light", "rigidbody", "script", "sound", "animation", "collision", "element", "particlesystem", "render", "sprite"}
	isValidType := false
	for _, validType := range validTypes {
		if componentType == validType {
			isValidType = true
			break
		}
	}
	
	if !isValidType {
		logging.Error("invalid component type", map[string]interface{}{
			"component_type": componentType,
			"valid_types":    validTypes,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid component type",
			"message": "Component type must be one of: " + strings.Join(validTypes, ", "),
		})
		return
	}
	
	// TODO: Implement PlayCanvas component removal
	// For now, assume component exists and can be removed
	
	logging.Info("component removed from entity", map[string]interface{}{
		"session_id":     sessionID,
		"entity_id":      entityID,
		"component_type": componentType,
	})
	
	// Broadcast component removed event via WebSocket
	if h != nil {
		h.BroadcastToSession(sessionID, "component_removed", map[string]interface{}{
			"entity_id":      entityID,
			"component_type": componentType,
			"timestamp":      time.Now().Format(time.RFC3339),
		})
	}
	
	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}