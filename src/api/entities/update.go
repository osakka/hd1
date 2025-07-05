package entities

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/memory"
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
	
	// Read request body for both specific and generic updates
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Error("failed to read request body", map[string]interface{}{
			"session_id": sessionID,
			"entity_id":  entityID,
			"error":      err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request",
			"message": "Could not read request body",
		})
		return
	}

	// Parse as specific UpdateEntityRequest first
	var req UpdateEntityRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
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

	// RADICAL OPTIMIZATION: Use pooled map for generic request parsing
	genericReq := memory.GetEntityRequestPool()
	defer memory.PutEntityRequestPool(genericReq)
	json.Unmarshal(bodyBytes, &genericReq)
	
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
	
	// Get the entity from storage to update it
	entity, err := h.GetStore().GetEntity(sessionID, entityID)
	if err != nil {
		logging.Warn("entity not found in storage", map[string]interface{}{
			"session_id": sessionID,
			"entity_id":  entityID,
			"error":      err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Entity not found",
			"message": "Entity does not exist in session storage",
		})
		return
	}

	// Track actual changes made
	changes := []string{}
	
	// Handle generic component updates from request body (for camera position updates)
	if components, ok := genericReq["components"].(map[string]interface{}); ok {
		if entity.Components == nil {
			entity.Components = make(map[string]interface{})
		}
		
		for componentType, componentData := range components {
			entity.Components[componentType] = componentData
			changes = append(changes, "components."+componentType)
		}
	}

	// Apply specific field updates from parsed request
	if req.Name != nil {
		entity.Name = *req.Name
		changes = append(changes, "name")
	}
	if req.Tags != nil {
		entity.Tags = req.Tags
		changes = append(changes, "tags")
	}
	if req.Enabled != nil {
		entity.Enabled = *req.Enabled
		changes = append(changes, "enabled")
	}
	if req.Position != nil {
		if entity.Components == nil {
			entity.Components = make(map[string]interface{})
		}
		if transform, ok := entity.Components["transform"].(map[string]interface{}); ok {
			transform["position"] = req.Position
		} else {
			entity.Components["transform"] = map[string]interface{}{
				"position": req.Position,
			}
		}
		changes = append(changes, "position")
	}
	if req.Rotation != nil {
		if entity.Components == nil {
			entity.Components = make(map[string]interface{})
		}
		if transform, ok := entity.Components["transform"].(map[string]interface{}); ok {
			transform["rotation"] = req.Rotation
		} else {
			entity.Components["transform"] = map[string]interface{}{
				"rotation": req.Rotation,
			}
		}
		changes = append(changes, "rotation")
	}
	if req.Scale != nil {
		if entity.Components == nil {
			entity.Components = make(map[string]interface{})
		}
		if transform, ok := entity.Components["transform"].(map[string]interface{}); ok {
			transform["scale"] = req.Scale
		} else {
			entity.Components["transform"] = map[string]interface{}{
				"scale": req.Scale,
			}
		}
		changes = append(changes, "scale")
	}

	// Save the updated entity back to storage
	if len(changes) > 0 {
		if err := h.GetStore().UpdateEntity(sessionID, entityID, entity); err != nil {
			logging.Error("failed to save updated entity", map[string]interface{}{
				"session_id": sessionID,
				"entity_id":  entityID,
				"error":      err.Error(),
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Failed to save entity updates",
				"message": "Could not persist entity changes",
			})
			return
		}
	}
	
	logging.Info("entity updated", map[string]interface{}{
		"session_id": sessionID,
		"entity_id":  entityID,
		"changes":    changes,
	})
	
	// RADICAL OPTIMIZATION: Use pooled maps for WebSocket broadcast and API response
	broadcastData := memory.GetWebSocketUpdate()
	defer memory.PutWebSocketUpdate(broadcastData)
	
	broadcastData["session_id"] = sessionID
	broadcastData["entity_id"] = entityID
	broadcastData["changes"] = changes
	
	// CRITICAL FIX: Use world-based broadcast for multiplayer entity visibility
	h.BroadcastAvatarPositionToChannel(sessionID, "entity_updated", broadcastData)
	
	// Use pooled map for API response
	responseData := memory.GetWebSocketUpdate()
	defer memory.PutWebSocketUpdate(responseData)
	
	responseData["success"] = true
	responseData["entity_id"] = entityID
	responseData["changes"] = changes
	responseData["updated_at"] = time.Now().Format(time.RFC3339)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}