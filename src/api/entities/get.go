package entities

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// GetEntityHandler handles GET /sessions/{sessionId}/entities/{entityId}
func GetEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
		logging.Error("invalid URL path for get entity", map[string]interface{}{
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
	
	// TODO: Implement PlayCanvas entity retrieval
	// For now, return mock entity data for demonstration
	
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
	
	logging.Info("entity retrieved", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
	})
	
	// Return mock entity details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": entityID,
		"name":      "MockEntity",
		"enabled":   true,
		"tags":      []string{"mock", "demo"},
		"transform": map[string]interface{}{
			"position": []float64{0, 0, 0},
			"rotation": []float64{0, 0, 0, 1},
			"scale":    []float64{1, 1, 1},
		},
		"components": []string{},
		"component_details": map[string]interface{}{},
		"children":    []string{},
		"parent":      nil,
		"playcanvas_guid": "pc-guid-mock-" + entityID,
		"created_at":      time.Now().Add(-time.Hour).Format(time.RFC3339),
		"last_updated":    time.Now().Format(time.RFC3339),
	})
}