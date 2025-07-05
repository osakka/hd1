package entities

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"holodeck1/logging"
	"holodeck1/server"
)

// DeleteEntityHandler handles DELETE /sessions/{sessionId}/entities/{entityId}
func DeleteEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
		logging.Error("invalid URL path for delete entity", map[string]interface{}{
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
	
	// Parse query parameters
	query := r.URL.Query()
	cascadeParam := query.Get("cascade")
	cascade := false
	if cascadeParam != "" {
		if parsedCascade, err := strconv.ParseBool(cascadeParam); err == nil {
			cascade = parsedCascade
		}
	}
	
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
	
	// TODO: Implement PlayCanvas entity deletion
	// For now, just log the operation
	
	logging.Info("entity deleted", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"cascade":    cascade,
	})
	
	// CRITICAL FIX: Use world-based broadcast for multiplayer entity visibility
	h.BroadcastAvatarPositionToChannel(sessionID, "entity_deleted", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"cascade":    cascade,
	})
	
	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}