package components

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// ListEntityComponentsHandler handles GET /sessions/{sessionId}/entities/{entityId}/components
func ListEntityComponentsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	if len(pathParts) < 5 || pathParts[0] != "sessions" || pathParts[2] != "entities" || pathParts[4] != "components" {
		logging.Error("invalid URL path for list entity components", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities/{entityId}/components",
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
	
	// Parse query parameters
	query := r.URL.Query()
	typeFilter := query.Get("type")
	
	// TODO: Implement PlayCanvas component listing for entity
	// Valid component types: model, camera, light, rigidbody, script, sound, animation, collision, element, particlesystem, render, sprite
	components := []map[string]interface{}{}
	
	// Mock response for now - in real implementation, this would query PlayCanvas entity components
	if typeFilter == "" || typeFilter == "model" {
		components = append(components, map[string]interface{}{
			"component_id":     "comp-model-" + entityID + "-1",
			"type":            "model",
			"enabled":         true,
			"playcanvas_id":   "pc-comp-guid-123",
			"properties_count": 5,
			"last_updated":    time.Now().Format(time.RFC3339),
		})
	}
	
	logging.Info("entity components listed", map[string]interface{}{
		"session_id":  sessionID,
		"entity_id":   entityID,
		"type_filter": typeFilter,
		"count":       len(components),
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"entity_id":  entityID,
		"components": components,
		"total":      len(components),
	})
}