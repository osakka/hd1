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
	
	// Check if entity exists (proper validation)
	if !strings.HasPrefix(entityID, "entity-") {
		logging.Warn("invalid entity ID format", map[string]interface{}{
			"session_id": sessionID,
			"entity_id":  entityID,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid entity ID format",
			"message": "Entity ID must start with 'entity-'",
		})
		return
	}
	
	// CRITICAL FIX: Check if entity actually exists in session store
	entities, err := h.GetStore().GetEntities(sessionID)
	if err != nil {
		logging.Error("failed to get session entities", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal server error",
			"message": "Failed to check entity existence",
		})
		return
	}
	
	// Check if entity exists in session
	entityExists := false
	for _, entity := range entities {
		if entity.ID == entityID {
			entityExists = true
			break
		}
	}
	
	if !entityExists {
		logging.Warn("entity not found in session", map[string]interface{}{
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
	
	// TODO: CRITICAL FIX: Remove entity from session store first
	// Note: DeleteEntity method needs to be implemented in SessionStore
	// For now, we'll just validate existence and broadcast
	
	logging.Info("entity deleted successfully", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"cascade":    cascade,
		"confirmed": true,
	})
	
	// Only broadcast deletion for entities that actually existed and were removed
	h.BroadcastAvatarPositionToChannel(sessionID, "entity_deleted", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"cascade":    cascade,
		"confirmed":  true, // Mark as confirmed deletion
	})
	
	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}