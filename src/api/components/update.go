package components

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// UpdateComponentRequest represents the request body for updating a component
type UpdateComponentRequest struct {
	Properties map[string]interface{} `json:"properties,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
}

// UpdateComponentHandler handles PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}
func UpdateComponentHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
		logging.Error("invalid URL path for update component", map[string]interface{}{
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
	
	// Parse request body
	var req UpdateComponentRequest
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
	
	// TODO: Implement PlayCanvas component property updates
	// For now, collect updated property names from request
	var updatedProperties []string
	for propertyName := range req.Properties {
		updatedProperties = append(updatedProperties, propertyName)
	}
	
	logging.Info("component updated", map[string]interface{}{
		"session_id":         sessionID,
		"entity_id":          entityID,
		"component_type":     componentType,
		"updated_properties": updatedProperties,
		"enabled_changed":    req.Enabled != nil,
		"properties":         req.Properties,
	})
	
	// Broadcast component updated event via WebSocket
	if h != nil {
		h.BroadcastToSession(sessionID, "component_updated", map[string]interface{}{
			"entity_id":          entityID,
			"component_type":     componentType,
			"updated_properties": updatedProperties,
			"properties":         req.Properties,
			"timestamp":          time.Now().Format(time.RFC3339),
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":            true,
		"entity_id":          entityID,
		"component_type":     componentType,
		"updated_properties": updatedProperties,
		"updated_at":         time.Now().Format(time.RFC3339),
	})
}