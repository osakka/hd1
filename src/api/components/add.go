package components

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// AddComponentRequest represents the request body for adding a component
type AddComponentRequest struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
}

// AddComponentHandler handles POST /sessions/{sessionId}/entities/{entityId}/components
func AddComponentHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
		logging.Error("invalid URL path for add component", map[string]interface{}{
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
	
	// Parse request body
	var req AddComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("failed to parse request body", map[string]interface{}{
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
			"message": "Failed to parse JSON request",
		})
		return
	}
	
	// Validate component type
	validTypes := []string{"model", "camera", "light", "rigidbody", "script", "sound", "animation", "collision", "element", "particlesystem", "render", "sprite"}
	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	
	if !isValidType {
		logging.Error("invalid component type", map[string]interface{}{
			"component_type": req.Type,
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
	
	// Set default enabled state
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	// TODO: Implement PlayCanvas component addition
	// For now, generate mock component ID and response
	componentID := "comp-" + req.Type + "-" + entityID + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	
	logging.Info("component added to entity", map[string]interface{}{
		"session_id":    sessionID,
		"entity_id":     entityID,
		"component_id":  componentID,
		"component_type": req.Type,
		"enabled":       enabled,
		"properties":    req.Properties,
	})
	
	// Broadcast component added event via WebSocket
	if h != nil {
		h.BroadcastToSession(sessionID, "component_added", map[string]interface{}{
			"entity_id":     entityID,
			"component_id":  componentID,
			"component_type": req.Type,
			"enabled":       enabled,
			"timestamp":     time.Now().Format(time.RFC3339),
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"entity_id":    entityID,
		"component_id": componentID,
		"type":         req.Type,
		"enabled":      enabled,
		"created_at":   time.Now().Format(time.RFC3339),
	})
}